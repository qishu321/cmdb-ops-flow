package service

import (
	"cmdb-ops-flow/models"
	"cmdb-ops-flow/utils/common"
	"context"
	"errors"
	clientv3 "go.etcd.io/etcd/client/v3"
	"time"
)

func AddEtcd(etcd models.EtcdGroup) (data interface{}, err error) {
	if etcd.EtcdGroupname == "" || etcd.EtcdEndpoints == "" {
		return nil, errors.New("所有字段都是必填的")
	}
	daoetcd := models.EtcdGroup{
		EtcdGroupid:   common.GenerateRandomNumber(),
		EtcdGroupname: etcd.EtcdGroupname,
		EtcdEndpoints: etcd.EtcdEndpoints,
		Label:         etcd.Label,
	}

	data, err = models.AddEtcdGroup(daoetcd)
	return data, err
}
func EditEtcd(etcd models.EtcdGroup) (data interface{}, err error) {

	daoetcd := models.EtcdGroup{
		EtcdGroupid:   etcd.EtcdGroupid,
		EtcdGroupname: etcd.EtcdGroupname,
		EtcdEndpoints: etcd.EtcdEndpoints,
		Label:         etcd.Label,
	}

	data, err = models.EditEtcdGroup(daoetcd)
	return data, err
}

// func GetScriptList(json models.ScriptManager) (data interface{}, err error) {
func GetEtcdList(id int) ([]models.EtcdGroup, error) {
	list, err := models.GetEtcdGroupList(id)
	return list, err
}

var (
	EtcdClient *clientv3.Client
	opTimeout  = 5 * time.Second
)

func Etcdinit(etcd models.EtcdGroup) (*clientv3.Client, error) {
	var err error
	EtcdClient, err = clientv3.New(clientv3.Config{
		Endpoints:   []string{etcd.EtcdEndpoints},
		DialTimeout: opTimeout,
	})
	if err != nil {
		return nil, err // 返回错误
	}

	// 在函数结束之前返回 etcdClient
	return EtcdClient, nil
}

func EtcdPut(key string, value string, etcdEndpoints string) (*clientv3.PutResponse, error) {
	var err error
	EtcdClient, err = clientv3.New(clientv3.Config{
		Endpoints:   []string{etcdEndpoints},
		DialTimeout: opTimeout,
	})
	if err != nil {
		return nil, err // 返回错误
	}

	ctx, cancel := context.WithTimeout(context.Background(), opTimeout)
	defer cancel()

	kvc := clientv3.NewKV(EtcdClient)

	resp, err := kvc.Put(ctx, key, value)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
