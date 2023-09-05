package prometheus

import (
	"bytes"
	"cmdb-ops-flow/models/prometheus"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

func SendMessage(notification prometheus.Notification, defaultRobot string) string {
	var buffer bytes.Buffer

	// 获取本地时区
	loc, err := time.LoadLocation("Local")
	if err != nil {
		log.Println("LoadLocation failed,", err)
		loc = time.UTC // 如果加载本地时区失败，则默认使用UTC时区
	}

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

		// 使用指定时区格式化触发时间
		buffer.WriteString(fmt.Sprintf("\n>触发时间: %s\n", alert.StartsAt.In(loc).Format("2006-01-02 15:04:05")))
		// 如果告警已恢复，同时显示恢复时间
		if notification.Status == "resolved" {
			buffer.WriteString(fmt.Sprintf("\n>恢复时间: %s\n", alert.EndsAt.In(loc).Format("2006-01-02 15:04:05")))
		}
		buffer.WriteString(fmt.Sprintf(`<@%s>`, "guomengfei"))

	}

	var m prometheus.Message
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
		return m.MsgType
	}
	resp := string(jsons)
	client := &http.Client{}

	req, err := http.NewRequest("POST", defaultRobot, strings.NewReader(resp))
	if err != nil {
		log.Println("SendMessage http NewRequest failed,", err)
		return resp
	}
	req.Header.Set("Content-Type", "application/json")
	r, err := client.Do(req)
	if err != nil {
		log.Println("SendMessage client Do failed", err)
		return resp
	}
	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println("SendMessage ReadAll Body failed", err)
		return resp
	}
	log.Println("SendMessage success,body:", string(body))
	return resp
}
