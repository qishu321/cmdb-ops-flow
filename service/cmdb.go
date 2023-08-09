package service

import (
	"cmdb-ops-flow/conf"
	"cmdb-ops-flow/models"
	"cmdb-ops-flow/utils/common"
	"errors"
	"fmt"
)

func AddCmdb(cmdb models.Cmdb) (data interface{}, err error) {
	if cmdb.Cmdbname == "" || cmdb.PublicIP == "" || cmdb.PrivateIP == "" ||
		cmdb.Username == "" || cmdb.SSHPort == 0 {
		return nil, errors.New("所有字段都是必填的")
	}
	// 如果 Password 和 PrivateKey 都为空，则返回错误
	if cmdb.Password == "" && cmdb.PrivateKey == "" {
		return nil, errors.New("Password 和 PrivateKey 至少要填写一个")
	}

	key := []byte(conf.Encryptkey)
	fmt.Println(key)
	//fmt.Println(cmdb.Password)
	passsword, err := common.Encrypt(key, cmdb.Password)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(passsword)
	daocmdb := models.Cmdb{
		Cmdbid:     common.GenerateRandomNumber(),
		Cmdbname:   cmdb.Cmdbname,
		PublicIP:   cmdb.PublicIP,
		PrivateIP:  cmdb.PrivateIP,
		Username:   cmdb.Username,
		Password:   passsword,
		PrivateKey: cmdb.PrivateKey,
		SSHPort:    cmdb.SSHPort,
		Label:      cmdb.Label,
	}

	data, err = models.Addcmdb(daocmdb)
	return data, err
}
func EditCmdb(cmdb models.Cmdb) (data interface{}, err error) {
	key := []byte(conf.Encryptkey)
	passsword, _ := common.Encrypt(key, cmdb.Password)
	daocmdb := models.Cmdb{
		Cmdbid:     cmdb.Cmdbid,
		Cmdbname:   cmdb.Cmdbname,
		PublicIP:   cmdb.PublicIP,
		PrivateIP:  cmdb.PrivateIP,
		Username:   cmdb.Username,
		Password:   passsword,
		PrivateKey: cmdb.PrivateKey,
		SSHPort:    cmdb.SSHPort,
		Label:      cmdb.Label,
	}

	data, err = models.Editcmdb(daocmdb)
	return data, err
}
func GetCmdbList(json models.Cmdb) (data interface{}, err error) {
	list, err := models.GetcmdbList(json.ID)
	return list, err
}
