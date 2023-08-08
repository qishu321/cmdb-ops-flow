package api

import (
	"cmdb-ops-flow/models"
	"cmdb-ops-flow/service"
	"cmdb-ops-flow/utils/msg"
	"cmdb-ops-flow/utils/result"
	"github.com/gin-gonic/gin"
	"net/http"
)

func AddUser(c *gin.Context) {
	var data models.User
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusOK, (&result.Result{}).Error(msg.ERROR, err.Error(), msg.GetErrMsg(msg.ERROR)))
		return
	}
	if data.Username == "" || data.Password == "" {
		c.JSON(http.StatusOK, (&result.Result{}).Error(msg.ERROR, nil, msg.GetErrMsg(msg.ERROR_USER_NO_PASSWD)))
		return
	}
	code := models.CheckUser(data.Username)
	if code == msg.SUCCSE {
		service.AddUser(data)
	}
	c.JSON(http.StatusOK, (&result.Result{}).Ok(code, data, msg.GetErrMsg(code)))

}

func DelUser(c *gin.Context) {
	var data models.User
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusOK, (&result.Result{}).Error(msg.ERROR, err.Error(), msg.GetErrMsg(msg.ERROR)))
		return
	}
	code := models.DelUser(data.Userid)

	c.JSON(http.StatusOK, (&result.Result{}).Ok(code, data, msg.GetErrMsg(code)))

}

func GetUser(c *gin.Context) {
	var data models.User
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusOK, (&result.Result{}).Error(msg.ERROR, err.Error(), msg.GetErrMsg(msg.ERROR)))
		return
	}
	list, err := service.GetUserList(data)
	if err != nil {
		c.JSON(http.StatusOK, (&result.Result{}).Error(msg.ERROR, err.Error(), msg.GetErrMsg(msg.ERROR)))
		return
	}
	code := 200
	c.JSON(http.StatusOK, (&result.Result{}).Ok(code, list, msg.GetErrMsg(code)))

}
func EditUser(c *gin.Context) {
	var data models.User
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusOK, (&result.Result{}).Error(msg.ERROR, err.Error(), msg.GetErrMsg(msg.ERROR)))
		return
	}
	list, err := service.EditUser(data)
	if err != nil {
		c.JSON(http.StatusOK, (&result.Result{}).Error(msg.ERROR, err.Error(), msg.GetErrMsg(msg.ERROR)))
		return
	}
	code := 200
	c.JSON(http.StatusOK, (&result.Result{}).Ok(code, list, msg.GetErrMsg(code)))

}
func Login(c *gin.Context) {
	var json models.User
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	data, err := service.Login(json)
	if err != nil {
		c.JSON(http.StatusOK, (&result.Result{}).Error(msg.ERROR_USER_NO_LOGIN, err.Error(), msg.GetErrMsg(msg.ERROR_USER_NO_LOGIN)))
		return
	}
	code := 200
	c.JSON(http.StatusOK, (&result.Result{}).Ok(code, data, msg.GetErrMsg(code)))
}

func Info(c *gin.Context) {
	token := c.MustGet("token").(string)
	data, err := service.Info(token)
	if err != nil {
		c.JSON(http.StatusOK, (&result.Result{}).Error(msg.ERROR_USER_NO_LOGIN, err.Error(), msg.GetErrMsg(msg.ERROR_USER_NO_LOGIN)))
		return
	}
	code := 200
	c.JSON(http.StatusOK, (&result.Result{}).Ok(code, data, msg.GetErrMsg(code)))
}

func Logout(c *gin.Context) {
	token := c.MustGet("token").(string)
	err := service.Logout(token)
	if err != nil {
		c.JSON(http.StatusOK, (&result.Result{}).Error(msg.ERROR_USER_NO_LOGOUT, err.Error(), msg.GetErrMsg(msg.ERROR_USER_NO_LOGOUT)))
		return
	}
	code := 200
	c.JSON(http.StatusOK, (&result.Result{}).Ok(code, "退出登录成功", msg.GetErrMsg(code)))
}
