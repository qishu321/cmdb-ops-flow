package api

import (
	"cmdb-ops-flow/conf"
	"cmdb-ops-flow/utils/common"
	"cmdb-ops-flow/utils/ssh"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	// 解决跨域问题
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

//type Query struct {
//	UserName  string `form:"username" binding:"required"`
//	Password  string `form:"password"`
//	IP        string `form:"ip" binding:"required"`
//	Port      int    `form:"port" binding:"required"`
//	Command   string `form:"command" binding:"required"`
//	AuthModel string `form:"authmodel" binding:"required"`
//}

// VisitorWebsocketServer https://github.com/widaT/webssh  websocket连接实现webssh
func VisitorWebsocketServer(c *gin.Context) {
	ip := c.Query("ip")
	port := c.Query("port")
	username := c.Query("username")
	Password := c.Query("password")
	authmodel := c.Query("authmodel")

	wsConn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("Failed to upgrade connection: %v", err)
		return
	}
	defer wsConn.Close()

	port1, _ := strconv.Atoi(port)
	key := []byte(conf.Encryptkey)

	password, err := common.Decrypt(key, Password)
	if err != nil {
		fmt.Println("解密失败:", err)
		return
	}
	fmt.Println(password)
	config := &ssh.SSHClientConfig{

		Timeout:   time.Second * 5,
		IP:        ip,
		Port:      port1,
		UserName:  username,
		Password:  password,
		AuthModel: authmodel,
	}
	sshClient, err := ssh.NewSSHClient(config)
	if err != nil {
		wsConn.WriteControl(websocket.CloseMessage,
			[]byte(err.Error()), time.Now().Add(time.Second))
		return
	}
	defer sshClient.Close()

	turn, err := ssh.NewTurn(wsConn, sshClient)
	if err != nil {
		fmt.Println("NewTurn," + err.Error())
		wsConn.WriteControl(websocket.CloseMessage,
			[]byte(err.Error()), time.Now().Add(time.Second))
		return
	}
	defer turn.Close()

	ctx, cancel := context.WithCancel(context.Background())
	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() {
		defer wg.Done()
		err := turn.LoopRead(ctx)
		if err != nil {
			fmt.Printf("%#v", err)
		}
	}()
	go func() {
		defer wg.Done()
		err := turn.SessionWait()
		if err != nil {
			fmt.Printf("%#v", err)
		}
		cancel()
	}()
	wg.Wait()
}
