package api_prom

import (
	"cmdb-ops-flow/conf"
	"cmdb-ops-flow/models/prometheus"
	"cmdb-ops-flow/utils/msg"
	"cmdb-ops-flow/utils/result"
	"github.com/gin-gonic/gin"

	svc_prome "cmdb-ops-flow/service/prometheus"
	"net/http"
)

var defaultRobot = "https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key="

func Alter(c *gin.Context) {
	aaa := defaultRobot + conf.Wxhookkey
	var notification prometheus.Notification
	err := c.BindJSON(&notification)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	data := svc_prome.SendMessage(notification, aaa)
	code := msg.SUCCSE
	c.JSON(http.StatusOK, (&result.Result{}).Ok(code, data, msg.GetErrMsg(code)))
}
