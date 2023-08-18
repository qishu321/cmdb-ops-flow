package api

import (
	"cmdb-ops-flow/models"
	"cmdb-ops-flow/service"
	"cmdb-ops-flow/utils/msg"
	"cmdb-ops-flow/utils/result"
	"github.com/gin-gonic/gin"
	"net/http"
)

func AddEtcdbak(c *gin.Context) {
	var data models.EtcdBak
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusOK, (&result.Result{}).Error(msg.ERROR, err.Error(), msg.GetErrMsg(msg.ERROR)))
		return
	}
	code := models.CheckEtcdBak(data.Etcdbakname)
	if code != msg.SUCCSE {
		c.JSON(http.StatusOK, (&result.Result{}).Error(msg.ERROR, "etcd name不能重复", msg.GetErrMsg(msg.ERROR)))
		return
	}

	list, err := service.AddEtcdbak(data)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, (&result.Result{}).Ok(code, list, msg.GetErrMsg(code)))

}

func GetEtcdbak(c *gin.Context) {
	//var data models.ScriptManager
	var data struct {
		ID int `json:"id"`
	}

	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusOK, (&result.Result{}).Error(400, err.Error(), msg.GetErrMsg(msg.ERROR)))
		return
	}
	list, err := service.GetEtcdbakList(data.ID)
	if err != nil {
		c.JSON(http.StatusOK, (&result.Result{}).Error(msg.ERROR, err.Error(), msg.GetErrMsg(msg.ERROR)))
		return
	}
	code := msg.SUCCSE
	c.JSON(http.StatusOK, (&result.Result{}).Ok(code, list, msg.GetErrMsg(code)))

}

func DelEtcdbak(c *gin.Context) {

	var data models.EtcdBak
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusOK, (&result.Result{}).Error(msg.ERROR, err.Error(), msg.GetErrMsg(msg.ERROR)))
		return
	}
	code := models.DelEtcdBak(data.Etcdbakid)

	c.JSON(http.StatusOK, (&result.Result{}).Ok(code, msg.SUCCSE, msg.GetErrMsg(code)))

}
