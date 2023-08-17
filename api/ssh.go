package api

import (
	"cmdb-ops-flow/conf"
	"cmdb-ops-flow/utils/common"
	"cmdb-ops-flow/utils/msg"
	"cmdb-ops-flow/utils/result"
	"cmdb-ops-flow/utils/ssh"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	"time"
)

// Query 查询参数

func SshcreateFile(c *gin.Context) {
	var r ssh.SSHClientConfig
	// 绑定并校验请求参数
	if err := c.ShouldBind(&r); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"err": err.Error(),
		})
		return
	}
	key := []byte(conf.Encryptkey)
	password, err := common.Decrypt(key, r.Password)
	if err != nil {
		fmt.Println("解密失败:", err)
		return
	}
	config := &ssh.SSHClientConfig{
		Timeout:   time.Second * 5,
		IP:        r.IP,
		Port:      r.Port,
		UserName:  r.UserName,
		Password:  password,
		AuthModel: "PASSWORD",
	}

	// 开始处理 SSH 会话

	code, err := ssh.CreateFileOnRemoteServer(config, r.Filename, r.Content)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, (&result.Result{}).Ok(code, "文件创建成功", msg.GetErrMsg(code)))

}

func SshCommand(c *gin.Context) {
	var r ssh.SSHClientConfig
	// 绑定并校验请求参数
	if err := c.ShouldBind(&r); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"err": err.Error(),
		})
		return
	}
	key := []byte(conf.Encryptkey)
	password, err := common.Decrypt(key, r.Password)
	if err != nil {
		fmt.Println("解密失败:", err)
		return
	}
	config := &ssh.SSHClientConfig{
		Timeout:   time.Second * 5,
		IP:        r.IP,
		Port:      r.Port,
		UserName:  r.UserName,
		Password:  password,
		AuthModel: "PASSWORD",
	}

	// 开始处理 SSH 会话
	output, err := ssh.SshCommand(config, r.Command)
	if err != nil {
		c.JSON(http.StatusOK, (&result.Result{}).Error(msg.ERROR, err.Error(), msg.GetErrMsg(msg.ERROR)))
		return
	}
	formattedOutput := strings.ReplaceAll(output, "\n", "<br>")

	c.JSON(http.StatusOK, (&result.Result{}).Ok(200, formattedOutput, msg.GetErrMsg(200)))

}
