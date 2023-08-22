package api

import (
	"cmdb-ops-flow/models"
	"cmdb-ops-flow/service"
	"cmdb-ops-flow/utils/msg"
	"cmdb-ops-flow/utils/result"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
)

func AddKubeConfig(c *gin.Context) {
	file, err := c.FormFile("yamlFile")
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error1": err.Error()})
		return
	}
	kubeconfigname := c.PostForm("kubeconfigname")

	fileContent, err := file.Open()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error2": err.Error()})
		return
	}
	defer fileContent.Close()
	yamlData, err := ioutil.ReadAll(fileContent)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error3": err.Error()})
		return
	}

	data, err := service.AddKubeConfig(yamlData, kubeconfigname)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, (&result.Result{}).Ok(200, data, msg.GetErrMsg(200)))

}
func EditKubeConfig(c *gin.Context) {
	var data models.KubeConfig
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusOK, (&result.Result{}).Error(msg.ERROR, err.Error(), msg.GetErrMsg(msg.ERROR)))
		return
	}
	list, err := service.EditKubeConfig(data)
	if err != nil {
		c.JSON(http.StatusOK, (&result.Result{}).Error(msg.ERROR, err.Error(), msg.GetErrMsg(msg.ERROR)))
		return
	}
	code := msg.SUCCSE
	c.JSON(http.StatusOK, (&result.Result{}).Ok(code, list, msg.GetErrMsg(code)))
}
func GetKubeConfig(c *gin.Context) {
	//var data models.ScriptManager
	var data struct {
		ID int `json:"id"`
	}

	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusOK, (&result.Result{}).Error(400, err.Error(), msg.GetErrMsg(msg.ERROR)))
		return
	}
	list, err := service.GetKubeConfigList(data.ID)
	if err != nil {
		c.JSON(http.StatusOK, (&result.Result{}).Error(msg.ERROR, err.Error(), msg.GetErrMsg(msg.ERROR)))
		return
	}
	code := msg.SUCCSE
	c.JSON(http.StatusOK, (&result.Result{}).Ok(code, list, msg.GetErrMsg(code)))

}

func DelKubeConfig(c *gin.Context) {

	var data models.KubeConfig
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusOK, (&result.Result{}).Error(msg.ERROR, err.Error(), msg.GetErrMsg(msg.ERROR)))
		return
	}
	code := models.DelKubeConfig(data.Kubeconfigid)

	c.JSON(http.StatusOK, (&result.Result{}).Ok(code, msg.SUCCSE, msg.GetErrMsg(code)))

}
