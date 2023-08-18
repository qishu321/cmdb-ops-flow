package api

import (
	"cmdb-ops-flow/models"
	"cmdb-ops-flow/service"
	"cmdb-ops-flow/utils/msg"
	"cmdb-ops-flow/utils/result"
	"context"
	"github.com/gin-gonic/gin"
	clientv3 "go.etcd.io/etcd/client/v3"
	"net/http"
)

func AddEtcd(c *gin.Context) {
	var data models.EtcdGroup
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusOK, (&result.Result{}).Error(msg.ERROR, err.Error(), msg.GetErrMsg(msg.ERROR)))
		return
	}
	code := models.CheckEtcdGroup(data.EtcdGroupname)
	if code != msg.SUCCSE {
		c.JSON(http.StatusOK, (&result.Result{}).Error(msg.ERROR, "etcd name不能重复", msg.GetErrMsg(msg.ERROR)))
		return
	}

	list, err := service.AddEtcd(data)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, (&result.Result{}).Ok(code, list, msg.GetErrMsg(code)))

}
func EditEtcd(c *gin.Context) {
	var data models.EtcdGroup
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusOK, (&result.Result{}).Error(msg.ERROR, err.Error(), msg.GetErrMsg(msg.ERROR)))
		return
	}
	list, err := service.EditEtcd(data)
	if err != nil {
		c.JSON(http.StatusOK, (&result.Result{}).Error(msg.ERROR, err.Error(), msg.GetErrMsg(msg.ERROR)))
		return
	}
	code := msg.SUCCSE
	c.JSON(http.StatusOK, (&result.Result{}).Ok(code, list, msg.GetErrMsg(code)))
}
func GetEtcd(c *gin.Context) {
	//var data models.ScriptManager
	var data struct {
		ID int `json:"id"`
	}

	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusOK, (&result.Result{}).Error(400, err.Error(), msg.GetErrMsg(msg.ERROR)))
		return
	}
	list, err := service.GetEtcdList(data.ID)
	if err != nil {
		c.JSON(http.StatusOK, (&result.Result{}).Error(msg.ERROR, err.Error(), msg.GetErrMsg(msg.ERROR)))
		return
	}
	code := msg.SUCCSE
	c.JSON(http.StatusOK, (&result.Result{}).Ok(code, list, msg.GetErrMsg(code)))

}

func DelEtcd(c *gin.Context) {

	var data models.EtcdGroup
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusOK, (&result.Result{}).Error(msg.ERROR, err.Error(), msg.GetErrMsg(msg.ERROR)))
		return
	}
	code := models.DelEtcdGroup(data.EtcdGroupid)

	c.JSON(http.StatusOK, (&result.Result{}).Ok(code, msg.SUCCSE, msg.GetErrMsg(code)))

}

func EtcdGetall(c *gin.Context) {

	var data models.EtcdGroup

	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusOK, (&result.Result{}).Error(400, err.Error(), msg.GetErrMsg(msg.ERROR)))
		return
	}

	// 初始化 ETCD 客户端
	etcdClient, err := service.Etcdinit(data)
	if err != nil {
		c.JSON(http.StatusOK, (&result.Result{}).Error(msg.ERROR, err.Error(), msg.GetErrMsg(msg.ERROR)))
		return
	}
	defer etcdClient.Close() // 在使用完毕后关闭 ETCD 客户端连接

	// 获取 ETCD 中的数据
	resp, err := etcdClient.Get(context.Background(), "", clientv3.WithPrefix())
	if err != nil {
		c.JSON(http.StatusOK, (&result.Result{}).Error(msg.ERROR, err.Error(), msg.GetErrMsg(msg.ERROR)))
		return
	}

	var list []map[string]string
	for _, kv := range resp.Kvs {
		list = append(list, map[string]string{
			"key":   string(kv.Key),
			"value": string(kv.Value),
		})
	}

	c.JSON(http.StatusOK, (&result.Result{}).Ok(200, list, msg.GetErrMsg(200)))
}
