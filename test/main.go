package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type Message struct {
	MsgType  string `json:"msgtype"`
	Markdown struct {
		Content string `json:"content"`
	} `json:"markdown"`
}

type Alert struct {
	Labels      map[string]string `json:"labels"`
	Annotations map[string]string `json:"annotations"`
	StartsAt    time.Time         `json:"startsAt"`
	EndsAt      time.Time         `json:"endsAt"`
}

type Notification struct {
	Version           string            `json:"version"`
	GroupKey          string            `json:"groupKey"`
	Status            string            `json:"status"`
	Receiver          string            `json:"receiver"`
	GroupLabels       map[string]string `json:"groupLabels"`
	CommonLabels      map[string]string `json:"commonLabels"`
	CommonAnnotations map[string]string `json:"commonAnnotations"`
	ExternalURL       string            `json:"externalURL"`
	Alerts            []Alert           `json:"alerts"`
}

var defaultRobot = "https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key="

func SendMessage(notification Notification, defaultRobot string) {
	var buffer bytes.Buffer

	// 构建告警消息
	for _, alert := range notification.Alerts {
		// 设置告警产生时级别状态颜色为黄色
		severityColor := "warning"
		if notification.Status == "resolved" {
			// 如果告警恢复，将字体颜色设置为绿色
			severityColor = "info"
		}
		buffer.WriteString(fmt.Sprintf("\n# <font color=\"%s\">级别状态: %s</font> (%s)\n", severityColor, alert.Labels["severity"], notification.Status))
		buffer.WriteString(fmt.Sprintf("\n>告警主题: %s\n", alert.Annotations["summary"]))
		buffer.WriteString(fmt.Sprintf("\n>告警类型: %s\n", alert.Labels["alertname"]))
		buffer.WriteString(fmt.Sprintf("\n>故障主机: %s\n", alert.Labels["instance"]))
		buffer.WriteString(fmt.Sprintf("\n>告警详情: %s\n", alert.Annotations["description"]))
		buffer.WriteString(fmt.Sprintf("\n>触发时间: %s\n", alert.StartsAt.Format("2006-01-02 15:04:05")))
		buffer.WriteString(fmt.Sprintf(`<@%s>`, "guomengfei"))
	}

	var m Message
	m.MsgType = "markdown"

	if notification.Status == "resolved" {
		// 如果告警恢复，将字体颜色设置为绿色
		m.Markdown.Content = "<font color=\"info\">告警已恢复</font>" + buffer.String()
	} else if notification.Status == "firing" {
		// 如果告警仍在触发中，展示原始告警信息
		m.Markdown.Content = buffer.String()
	}

	jsons, err := json.Marshal(m)
	if err != nil {
		log.Println("SendMessage Marshal failed,", err)
		return
	}
	resp := string(jsons)
	client := &http.Client{}

	req, err := http.NewRequest("POST", defaultRobot, strings.NewReader(resp))
	if err != nil {
		log.Println("SendMessage http NewRequest failed,", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")
	r, err := client.Do(req)
	if err != nil {
		log.Println("SendMessage client Do failed", err)
		return
	}
	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println("SendMessage ReadAll Body failed", err)
		return
	}
	log.Println("SendMessage success,body:", string(body))
}

func Alter(c *gin.Context) {
	var notification Notification
	err := c.BindJSON(&notification)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	SendMessage(notification, defaultRobot)
}

func main() {
	t := gin.Default()
	t.POST("/Alter", Alter)
	t.Run(":8090")
}
