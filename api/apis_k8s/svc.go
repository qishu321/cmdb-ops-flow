package apis_k8s

import (
	"cmdb-ops-flow/models/k8s"
	"cmdb-ops-flow/service/service_k8s"
	"cmdb-ops-flow/utils/msg"
	"cmdb-ops-flow/utils/result"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetallSvc(c *gin.Context) {
	var data k8s.Svc
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusOK, (&result.Result{}).Error(400, err.Error(), msg.GetErrMsg(msg.ERROR)))
		return
	}
	list, err := service_k8s.Getsvs(data.ID, data.Namespace)
	if err != nil {
		c.JSON(http.StatusOK, (&result.Result{}).Error(msg.ERROR, err.Error(), msg.GetErrMsg(msg.ERROR)))
		return
	}
	code := msg.SUCCSE
	c.JSON(http.StatusOK, (&result.Result{}).Ok(code, list, msg.GetErrMsg(code)))
}

func AddSvc(c *gin.Context) {

	var data k8s.Svc
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusOK, (&result.Result{}).Error(5001, err.Error(), msg.GetErrMsg(msg.ERROR)))
		return
	}
	list, err := service_k8s.AddSvc(data.ID, data)
	if err != nil {
		c.JSON(http.StatusOK, (&result.Result{}).Error(msg.ERROR, err.Error(), msg.GetErrMsg(msg.ERROR)))
		return
	}
	code := msg.SUCCSE
	c.JSON(http.StatusOK, (&result.Result{}).Ok(code, list, msg.GetErrMsg(code)))
}
func EditSvc(c *gin.Context) {

	var data k8s.Svc
	if err := c.ShouldBindJSON(&data); err != nil {
		fmt.Println(err)
		c.JSON(http.StatusOK, (&result.Result{}).Error(5001, err.Error(), msg.GetErrMsg(msg.ERROR)))
		return
	}
	list, err := service_k8s.EditSvc(data.ID, data)
	if err != nil {
		c.JSON(http.StatusOK, (&result.Result{}).Error(msg.ERROR, err.Error(), msg.GetErrMsg(msg.ERROR)))
		return
	}
	code := msg.SUCCSE
	c.JSON(http.StatusOK, (&result.Result{}).Ok(code, list, msg.GetErrMsg(code)))
}

func DelSvc(c *gin.Context) {

	var data k8s.Svc
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusOK, (&result.Result{}).Error(5001, err.Error(), msg.GetErrMsg(msg.ERROR)))
		return
	}
	code, err := service_k8s.Delsvs(data.ID, data.Namespace, data.Name)
	if err != nil {
		c.JSON(http.StatusOK, (&result.Result{}).Error(msg.ERROR, err.Error(), msg.GetErrMsg(msg.ERROR)))
		return
	}
	c.JSON(http.StatusOK, (&result.Result{}).Ok(code, "删除成功", msg.GetErrMsg(code)))
}
