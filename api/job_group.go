package api

import (
	"cmdb-ops-flow/models"
	"cmdb-ops-flow/service"
	"cmdb-ops-flow/utils/msg"
	"cmdb-ops-flow/utils/result"
	"github.com/gin-gonic/gin"
	"net/http"
)

func AddJobGroup(c *gin.Context) {
	var data models.JobGroup
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusOK, (&result.Result{}).Error(msg.ERROR, err.Error(), msg.GetErrMsg(msg.ERROR)))
		return
	}
	code := models.CheckJobGroup(data.Jobgroupname)
	if code != msg.SUCCSE {
		c.JSON(http.StatusOK, (&result.Result{}).Error(msg.ERROR, "JobGroup name不能重复", msg.GetErrMsg(msg.ERROR)))
		return
	}

	// 调用 service.AddCmdb 执行业务逻辑
	list, err := service.AddJobGroup(data)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, (&result.Result{}).Ok(code, list, msg.GetErrMsg(code)))

}

func GetJobGroup(c *gin.Context) {
	//var data models.ScriptManager
	var data models.JobGroup

	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusOK, (&result.Result{}).Error(400, err.Error(), msg.GetErrMsg(msg.ERROR)))
		return
	}
	list, err := service.GetJobGroupList(data.ID)
	if err != nil {
		c.JSON(http.StatusOK, (&result.Result{}).Error(msg.ERROR, err.Error(), msg.GetErrMsg(msg.ERROR)))
		return
	}
	code := msg.SUCCSE
	c.JSON(http.StatusOK, (&result.Result{}).Ok(code, list, msg.GetErrMsg(code)))

}
func GetSearchJobGroup(c *gin.Context) {
	var data struct {
		Keyword string `form:"keyword" binding:"required"`
	}
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusOK, (&result.Result{}).Error(msg.ERROR, err.Error(), msg.GetErrMsg(msg.ERROR)))
		return
	}
	list, code := models.SearchJobGroup(data.Keyword)
	if code != msg.SUCCSE {
		c.JSON(http.StatusOK, (&result.Result{}).Error(msg.ERROR, nil, msg.GetErrMsg(msg.ERROR)))
		return
	}
	c.JSON(http.StatusOK, (&result.Result{}).Ok(code, list, msg.GetErrMsg(code)))
}

func EditJobGroup(c *gin.Context) {
	var data models.JobGroup
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusOK, (&result.Result{}).Error(msg.ERROR, err.Error(), msg.GetErrMsg(msg.ERROR)))
		return
	}
	list, err := service.EditJobGroup(data)
	if err != nil {
		c.JSON(http.StatusOK, (&result.Result{}).Error(msg.ERROR, err.Error(), msg.GetErrMsg(msg.ERROR)))
		return
	}
	code := msg.SUCCSE
	c.JSON(http.StatusOK, (&result.Result{}).Ok(code, list, msg.GetErrMsg(code)))
}
func DelJobGroup(c *gin.Context) {

	var data models.JobGroup
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusOK, (&result.Result{}).Error(msg.ERROR, err.Error(), msg.GetErrMsg(msg.ERROR)))
		return
	}
	code := models.DelJobGroup(data.Jobgroupid)

	c.JSON(http.StatusOK, (&result.Result{}).Ok(code, msg.SUCCSE, msg.GetErrMsg(code)))

}
