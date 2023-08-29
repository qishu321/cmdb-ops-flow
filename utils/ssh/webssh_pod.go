package ssh

import (
	"cmdb-ops-flow/models"
	"cmdb-ops-flow/models/k8s"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/remotecommand"
	"sync"

	"net/http"
)

// 包级变量,升级器
var upgrader = websocket.Upgrader{}

func init() {
	//初始化
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		// 解决跨域问题
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
}

// websocket消息
type WsMessage struct {
	MessageType int
	Data        []byte
}

// 封装websocket连接
type WsConnection struct {
	conn      *websocket.Conn // 底层websocket
	inChan    chan *WsMessage // 读取队列
	outChan   chan *WsMessage // 发送队列
	mutex     sync.Mutex      // 避免重复关闭管道
	isClosed  bool
	closeChan chan byte // 关闭通知
}

// 读取协程
func (wsConn *WsConnection) wsReadLoop() {
	for {
		// 读一条message
		msgType, data, err := wsConn.conn.ReadMessage()
		if err != nil {
			break
		}

		// 放入请求队列
		wsConn.inChan <- &WsMessage{
			msgType,
			data,
		}

		//select {
		//case wsConn.inChan <- msg:
		//case <- wsConn.closeChan:
		//
		//}
	}
}

// 发送协程
func (wsConn *WsConnection) wsWriteLoop() {
	// 服务端返回给页面的数据

	for {
		var msg *WsMessage
		select {
		// 取一个应答
		case msg = <-wsConn.outChan:
			// 写给web  websocket
			if err := wsConn.conn.WriteMessage(msg.MessageType, msg.Data); err != nil {
				break
			}
		case <-wsConn.closeChan:
			wsConn.WsClose()
		}
	}
}

func InitWebsocket(resp http.ResponseWriter, req *http.Request) (wsConn *WsConnection, err error) {
	// 应答客户端告知升级连接为websocket
	conn, err := upgrader.Upgrade(resp, req, nil)
	if err != nil {
		return
	}
	wsConn = &WsConnection{
		conn:      conn,
		inChan:    make(chan *WsMessage, 1000),
		outChan:   make(chan *WsMessage, 1000),
		closeChan: make(chan byte),
		isClosed:  false,
	}

	// 页面读入输入 协程
	go wsConn.wsReadLoop()
	// 服务端返回数据 协程
	go wsConn.wsWriteLoop()

	return
}

// 发送返回消息到协程
func (wsConn *WsConnection) WsWrite(messageType int, data []byte) (err error) {
	select {
	case wsConn.outChan <- &WsMessage{messageType, data}:

	case <-wsConn.closeChan:
		err = errors.New("WsWrite websocket closed")
		break
	}
	return
}

// 读取 页面消息到协程
func (wsConn *WsConnection) WsRead() (msg *WsMessage, err error) {
	select {
	case msg = <-wsConn.inChan:
		return
	case <-wsConn.closeChan:
		err = errors.New("WsRead websocket closed")
		break
	}
	return
}

// 关闭连接
func (wsConn *WsConnection) WsClose() {
	wsConn.conn.Close()
	wsConn.mutex.Lock()
	defer wsConn.mutex.Unlock()
	if !wsConn.isClosed {
		wsConn.isClosed = true
		close(wsConn.closeChan)
	}
}

// ssh流式处理器
type streamHandler struct {
	wsConn      *WsConnection
	resizeEvent chan remotecommand.TerminalSize
}

// web终端发来的包
type xtermMessage struct {
	MsgType string `json:"type"`  // 类型:resize客户端调整终端, input客户端输入
	Input   string `json:"input"` // msgtype=input情况下使用
	Rows    uint16 `json:"rows"`  // msgtype=resize情况下使用
	Cols    uint16 `json:"cols"`  // msgtype=resize情况下使用
}

// executor回调获取web是否resize
func (handler *streamHandler) Next() (size *remotecommand.TerminalSize) {
	ret := <-handler.resizeEvent
	size = &ret
	return
}

// executor回调读取web端的输入
func (handler *streamHandler) Read(p []byte) (size int, err error) {

	// 读web发来的输入
	var msg *WsMessage
	if msg, err = handler.wsConn.WsRead(); err != nil {
		handler.wsConn.WsClose()
		return
	}
	// 解析客户端请求
	//if err = json.Unmarshal([]byte(msg.Data), &xtermMsg); err != nil {
	//	return
	//}

	xtermMsg := &xtermMessage{
		//MsgType: string(msg.MessageType),
		Input: string(msg.Data),
	}
	// 放到channel里，等remotecommand executor调用我们的Next取走
	handler.resizeEvent <- remotecommand.TerminalSize{Width: xtermMsg.Cols, Height: xtermMsg.Rows}
	size = len(xtermMsg.Input)
	copy(p, xtermMsg.Input)
	return

}

// executor回调向web端输出
func (handler *streamHandler) Write(p []byte) (size int, err error) {
	// 产生副本
	copyData := make([]byte, len(p))
	copy(copyData, p)
	size = len(p)
	err = handler.wsConn.WsWrite(websocket.TextMessage, copyData)
	return
}

func StartProcess(wsConn *WsConnection, id int, podName string, namespace string, container string, command string) error {
	// URL:
	// https://172.16.0.143:6443/api/v1/namespaces/default/pods/nginx-deployment-5cbd8757f-d5qvx/exec?command=sh&container=nginx&stderr=true&stdin=true&stdout=true&tty=true
	clientSet, err := k8s.GetKubeConfig(id)
	if err != nil {
		return nil
	}

	req := clientSet.CoreV1().RESTClient().Post().
		Resource("pods").
		Name(podName).
		Namespace(namespace).
		SubResource("exec").
		VersionedParams(&v1.PodExecOptions{
			Container: container,
			Command: []string{
				command,
			},
			Stdin:  true,
			Stdout: true,
			Stderr: true,
			TTY:    true,
		}, scheme.ParameterCodec)

	// 创建到容器的连接
	kubeconfig, err := models.GetKubeConfigList(id)
	if err != nil {
		return nil
	}

	// 解码Base64编码的kubeconfig数据
	kubeconfigBytes, err := base64.StdEncoding.DecodeString(kubeconfig[0].Kubeconfigdata)
	if err != nil {
		return fmt.Errorf("invalid Base64 data")
	}

	// 使用kubeconfigBytes来构建Kubernetes配置
	config, _ := clientcmd.RESTConfigFromKubeConfig(kubeconfigBytes)
	if err != nil {
		fmt.Println("使用kubeconfigBytes来构建Kubernetes配置:", err) // 添加错误日志输出

		return nil
	}
	executor, err := remotecommand.NewSPDYExecutor(config, "POST", req.URL())
	if err != nil {
		return err
	}

	// 配置与容器之间的数据流处理回调
	handler := &streamHandler{wsConn: wsConn, resizeEvent: make(chan remotecommand.TerminalSize)}
	if err = executor.Stream(remotecommand.StreamOptions{
		Stdin:             handler,
		Stdout:            handler,
		Stderr:            handler,
		TerminalSizeQueue: handler,
		Tty:               true,
	}); err != nil {
		fmt.Println("handler", err)
		return err

	}
	return err

}
