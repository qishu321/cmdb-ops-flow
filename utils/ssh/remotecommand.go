package ssh

import (
	"cmdb-ops-flow/conf"
	"cmdb-ops-flow/models"
	"cmdb-ops-flow/utils/common"
	"fmt"
	"time"
)

func CreateRemoteFile(script *models.ScriptManager, cmdb *models.Cmdb) (string, error) {
	// 构造请求参数
	key := []byte(conf.Encryptkey)
	password, err := common.Decrypt(key, cmdb.Password)
	if err != nil {
		fmt.Println("解密失败:", err)
		return password, err
	}
	r := SSHClientConfig{
		Timeout:   time.Second * 5,
		IP:        cmdb.PrivateIP,
		Port:      cmdb.SSHPort,
		UserName:  cmdb.Username,
		Password:  password,
		AuthModel: "PASSWORD",
	}

	// 在远程服务器上创建文件

	code, err := CreateFileOnRemoteServer(&r, script.Name, script.Script)
	if err != nil {
		return "创建文件失败", err
	}
	result := fmt.Sprintf("文件创建成功: %d", code)
	return result, nil
}

func ExecuteRemoteCommand(script *models.ScriptManager, cmdb *models.Cmdb) (string, error) {
	// 构造请求参数
	key := []byte(conf.Encryptkey)
	password, error := common.Decrypt(key, cmdb.Password)
	if error != nil {
		fmt.Println("解密失败:", error)
		return "解密失败:", error
	}
	r := SSHClientConfig{
		Timeout:   time.Second * 5,
		IP:        cmdb.PrivateIP,
		Port:      cmdb.SSHPort,
		UserName:  cmdb.Username,
		Password:  password,
		AuthModel: "PASSWORD",
	}

	output, err := SshCommand(&r, "bash /tmp/"+script.Name)
	if err != nil {
		return "", err
	}

	result := fmt.Sprintf("命令执行成功:\n%s", output)
	return result, nil

}
