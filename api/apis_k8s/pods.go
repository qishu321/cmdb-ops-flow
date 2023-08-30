package apis_k8s

import (
	"cmdb-ops-flow/models/k8s"
	"cmdb-ops-flow/service/service_k8s"
	"cmdb-ops-flow/utils/msg"
	"cmdb-ops-flow/utils/result"
	"cmdb-ops-flow/utils/ssh"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"gopkg.in/yaml.v3"
	"io"
	"k8s.io/apimachinery/pkg/runtime"
	"log"
	"net/http"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

const logBufferSize = 1024 // 适当设置缓冲区大小

func GetPodLogs(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("WebSocket 连接失败：%v", err)
		return
	}
	defer conn.Close()

	var data k8s.Pod
	if err := c.ShouldBindQuery(&data); err != nil {
		log.Printf("请求数据绑定失败：%v", err)
		return
	}
	fmt.Println("123", data)
	logsStream, err := service_k8s.GetPodLogs(data.ID, data.Namespace, data.Name, data.ContainerName)
	if err != nil {
		log.Printf("获取日志时出错：%v", err)
		return
	}

	logBuffer := make([]byte, logBufferSize)
	for {
		n, err := logsStream.Read(logBuffer)
		if err != nil && err != io.EOF {
			log.Printf("读取日志时出错：%v", err)
			return
		}

		if n > 0 {
			logData := logBuffer[:n]
			if err := conn.WriteMessage(websocket.TextMessage, logData); err != nil {
				log.Printf("WebSocket 发送消息失败：%v", err)
				return
			}
		}

		if err == io.EOF {
			break
		}
	}
}

func GetAllPods(c *gin.Context) {
	var data struct {
		ID int `json:"id"`
	}
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusOK, (&result.Result{}).Error(400, err.Error(), msg.GetErrMsg(msg.ERROR)))
		return
	}
	list, err := service_k8s.GetallPods(data.ID)
	if err != nil {
		c.JSON(http.StatusOK, (&result.Result{}).Error(msg.ERROR, err.Error(), msg.GetErrMsg(msg.ERROR)))
		return
	}
	code := msg.SUCCSE
	c.JSON(http.StatusOK, (&result.Result{}).Ok(code, list, msg.GetErrMsg(code)))
}
func GetPods(c *gin.Context) {
	var data k8s.Pod
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusOK, (&result.Result{}).Error(400, err.Error(), msg.GetErrMsg(msg.ERROR)))
		return
	}
	list, err := service_k8s.GetPods(data.ID, data.Namespace)
	if err != nil {
		c.JSON(http.StatusOK, (&result.Result{}).Error(msg.ERROR, err.Error(), msg.GetErrMsg(msg.ERROR)))
		return
	}
	code := msg.SUCCSE
	c.JSON(http.StatusOK, (&result.Result{}).Ok(code, list, msg.GetErrMsg(code)))
}
func GetPodsYaml(c *gin.Context) {
	var data k8s.Pod
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusOK, (&result.Result{}).Error(400, err.Error(), msg.GetErrMsg(msg.ERROR)))
		return
	}
	pod, err := service_k8s.GetPodsYaml(data.ID, data.Namespace, data.Name)
	if err != nil {
		c.JSON(http.StatusOK, (&result.Result{}).Error(msg.ERROR, err.Error(), msg.GetErrMsg(msg.ERROR)))
		return
	}
	podMap, err := runtime.DefaultUnstructuredConverter.ToUnstructured(pod)
	if err != nil {
		c.JSON(http.StatusOK, (&result.Result{}).Error(msg.ERROR, err.Error(), msg.GetErrMsg(msg.ERROR)))
		return
	}
	podYAMLBytes, err := yaml.Marshal(podMap)
	if err != nil {
		c.JSON(http.StatusOK, (&result.Result{}).Error(msg.ERROR, err.Error(), msg.GetErrMsg(msg.ERROR)))
		return
	}

	podYaml := string(podYAMLBytes)

	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to convert Pod to YAML"})
		return
	}

	c.JSON(http.StatusOK, (&result.Result{}).Ok(200, podYaml, msg.GetErrMsg(200)))

	//c.String(http.StatusOK,  podYaml)
}

//var (
//	logBufferSize = 10 * 1024 * 1024 // 10M
//)

//func GetPodLogs(c *gin.Context) {
//	var data k8s.Pod
//	if err := c.ShouldBindJSON(&data); err != nil {
//		c.JSON(http.StatusBadRequest, (&result.Result{}).Error(http.StatusBadRequest, err.Error(), msg.GetErrMsg(msg.ERROR)))
//		return
//	}
//
//	logsStream, err := service_k8s.GetPodLogs(data.ID, data.Namespace, data.Name, data.ContainerName)
//	if err != nil {
//		log.Printf("Error getting logs: %v", err)
//
//		c.JSON(http.StatusInternalServerError, (&result.Result{}).Error(http.StatusInternalServerError, "Error getting logs", msg.GetErrMsg(msg.ERROR)))
//		return
//	}
//
//	logBuffer := make([]byte, logBufferSize)
//	n, err := logsStream.Read(logBuffer)
//	if err != nil && err != io.EOF {
//		c.JSON(http.StatusInternalServerError, (&result.Result{}).Error(http.StatusInternalServerError, "Error reading logs", msg.GetErrMsg(msg.ERROR)))
//		return
//	}
//	c.JSON(http.StatusOK, (&result.Result{}).Ok(200, string(logBuffer[:n]), msg.GetErrMsg(200)))
//
//}

type query struct {
	Id            int    `form:"id" binding:"required"`
	Namespace     string `form:"namespace" binding:"required"`
	PodName       string `form:"pod_name" binding:"required"`
	ContainerName string `form:"container_name" binding:"required"`
	Command       string `form:"command" binding:"required"`
}

func SshPod(c *gin.Context) {
	wsConn, err := ssh.InitWebsocket(c.Writer, c.Request)
	if err != nil {
		fmt.Println("InitWebsocket err", err)
		wsConn.WsClose()
		return
	}
	var r query
	if err := c.ShouldBindQuery(&r); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"err": err.Error(),
		})
		return
	}

	wsConn.WsWrite(websocket.TextMessage, []byte("你已进入 命名空间："+r.Namespace+" 容器组："+r.ContainerName+" 容器名："+r.ContainerName+"的终端"))

	if err := ssh.StartProcess(wsConn, r.Id, r.PodName, r.Namespace, r.ContainerName, r.Command); err != nil {
		fmt.Println("StartProcess err", err)
		wsConn.WsClose()
		return
	}
}
