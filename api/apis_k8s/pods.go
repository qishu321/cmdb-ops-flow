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
	"net/http"
)

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
