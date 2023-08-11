package api

import (
	"cmdb-ops-flow/models"
	"cmdb-ops-flow/service"
	"cmdb-ops-flow/utils/msg"
	"cmdb-ops-flow/utils/result"
	"github.com/gin-gonic/gin"
	"net/http"
)

func AddScript(c *gin.Context) {
	var data models.ScriptManager
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusOK, (&result.Result{}).Error(msg.ERROR, err.Error(), msg.GetErrMsg(msg.ERROR)))
		return
	}
	code := models.CheckScript(data.Script)
	if code != msg.SUCCSE {
		c.JSON(http.StatusOK, (&result.Result{}).Error(msg.ERROR, "Script name不能重复", msg.GetErrMsg(msg.ERROR)))
		return
	}

	// 调用 service.AddCmdb 执行业务逻辑
	list, err := service.AddScript(data)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, (&result.Result{}).Ok(code, list, msg.GetErrMsg(code)))

}

func GetScript(c *gin.Context) {
	//var data models.ScriptManager
	var data struct {
		ID int `json:"id"`
	}

	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusOK, (&result.Result{}).Error(400, err.Error(), msg.GetErrMsg(msg.ERROR)))
		return
	}
	list, err := service.GetScriptList(data.ID)
	if err != nil {
		c.JSON(http.StatusOK, (&result.Result{}).Error(msg.ERROR, err.Error(), msg.GetErrMsg(msg.ERROR)))
		return
	}
	code := msg.SUCCSE
	c.JSON(http.StatusOK, (&result.Result{}).Ok(code, list, msg.GetErrMsg(code)))

}

func EditScript(c *gin.Context) {
	var data models.ScriptManager
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusOK, (&result.Result{}).Error(msg.ERROR, err.Error(), msg.GetErrMsg(msg.ERROR)))
		return
	}
	list, err := service.EditScript(data)
	if err != nil {
		c.JSON(http.StatusOK, (&result.Result{}).Error(msg.ERROR, err.Error(), msg.GetErrMsg(msg.ERROR)))
		return
	}
	code := msg.SUCCSE
	c.JSON(http.StatusOK, (&result.Result{}).Ok(code, list, msg.GetErrMsg(code)))
}
func DelScript(c *gin.Context) {
	var data struct {
		MachineID int64 `json:"machineid" validate:"required"`
	}
	//var data models.ScriptManager
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusOK, (&result.Result{}).Error(msg.ERROR, err.Error(), msg.GetErrMsg(msg.ERROR)))
		return
	}
	code := models.DelScript(data.MachineID)

	c.JSON(http.StatusOK, (&result.Result{}).Ok(code, msg.SUCCSE, msg.GetErrMsg(code)))

}
