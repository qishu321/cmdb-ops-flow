package main

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.etcd.io/etcd/client/v3"
)

var (
	etcdClient *clientv3.Client
	opTimeout  = 5 * time.Second
)

func main() {
	r := gin.Default()

	var err error
	etcdClient, err = clientv3.New(clientv3.Config{
		Endpoints:   []string{"http://192.168.2.81:2379"},
		DialTimeout: opTimeout,
	})
	if err != nil {
		log.Fatal(err)
	}
	defer etcdClient.Close()

	r.GET("/get/", handleGet)
	r.GET("/list", handleList)

	r.POST("/put", handlePut)
	r.DELETE("/delete/:key", handleDelete)
	r.GET("/backup", handleBackup)
	r.POST("/restore", handleRestore)

	if err := r.Run(":8080"); err != nil {
		fmt.Println(err)
	}
}
func handleList(c *gin.Context) {
	resp, err := etcdClient.Get(context.Background(), "", clientv3.WithPrefix())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var data []map[string]string
	for _, kv := range resp.Kvs {
		data = append(data, map[string]string{
			"key":   string(kv.Key),
			"value": string(kv.Value),
		})
	}

	c.JSON(http.StatusOK, data)
}

func handleGet(c *gin.Context) {
	key := c.DefaultQuery("key", "") // 获取名为 "key" 的查询参数，默认为空字符串
	if key == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing key parameter"})
		return
	}

	resp, err := etcdClient.Get(context.Background(), key)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if len(resp.Kvs) > 0 {
		value := string(resp.Kvs[0].Value)
		c.JSON(http.StatusOK, gin.H{"key": key, "value": value})
	} else {
		c.JSON(http.StatusOK, gin.H{"key": key, "value": "not found"})
	}
}

func handlePut(c *gin.Context) {
	var data struct {
		Key   string `json:"key"`
		Value string `json:"value"`
	}
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err := etcdPut(data.Key, data.Value)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Successfully put key-value pair"})
}

func handleDelete(c *gin.Context) {
	key := c.Param("key")
	_, err := etcdDelete(key)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Successfully deleted key"})
}

func handleBackup(c *gin.Context) {
	backupData, err := etcdBackup()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Data(http.StatusOK, "application/octet-stream", backupData)
}

func handleRestore(c *gin.Context) {
	restoreFile, _, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	defer restoreFile.Close()

	restoreData, err := ioutil.ReadAll(restoreFile)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	err = etcdRestore(restoreData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Restore completed"})
}

func etcdPut(key, value string) (*clientv3.PutResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), opTimeout)
	defer cancel()

	kvc := clientv3.NewKV(etcdClient)

	resp, err := kvc.Put(ctx, key, value)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func etcdDelete(key string) (*clientv3.DeleteResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), opTimeout)
	defer cancel()

	kvc := clientv3.NewKV(etcdClient)

	resp, err := kvc.Delete(ctx, key)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func etcdBackup() ([]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), opTimeout)
	defer cancel()

	resp, err := etcdClient.Get(ctx, "", clientv3.WithSerializable())
	if err != nil {
		return nil, err
	}

	buf := new(bytes.Buffer)
	for _, kv := range resp.Kvs {
		fmt.Fprintf(buf, "%s=%s\n", kv.Key, kv.Value)
	}
	return buf.Bytes(), nil
}

func etcdRestore(data []byte) error {
	ctx, cancel := context.WithTimeout(context.Background(), opTimeout)
	defer cancel()

	kvc := clientv3.NewKV(etcdClient)

	lines := bytes.Split(data, []byte("\n"))
	for _, line := range lines {
		parts := bytes.SplitN(line, []byte("="), 2)
		if len(parts) != 2 {
			continue
		}
		_, err := kvc.Put(ctx, string(parts[0]), string(parts[1]))
		if err != nil {
			return err
		}
	}
	return nil
}
