package conf

import (
	"fmt"
	"github.com/go-ini/ini"
	"time"
)

var (
	AppMode      string
	HttpPort     string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration

	Md5Key     string
	Encryptkey string
	Wxhookkey  string
	DbHost     string
	DbPort     string
	DbUser     string
	DbPassWord string
	DbName     string
)

func init() {
	file, err := ini.Load("conf/config.ini")
	if err != nil {
		fmt.Println("配置文件读取失败，请检查文件路径:", err)
	}
	LoadServer(file)
	LoadData(file)
	LoadApp(file)
}
func LoadServer(file *ini.File) {
	AppMode = file.Section("server").Key("AppMode").String()
	HttpPort = file.Section("server").Key("HttpPort").String()
	ReadTimeout = time.Duration(file.Section("server").Key("ReadTimeout").MustInt(60)) * time.Second
	WriteTimeout = time.Duration(file.Section("server").Key("WriteTimeout").MustInt(60)) * time.Second
}
func LoadApp(file *ini.File) {
	Md5Key = file.Section("app").Key("Md5Key").String()
	Encryptkey = file.Section("app").Key("Encryptkey").String()
	Wxhookkey = file.Section("app").Key("Wxhookkey").String()

}

func LoadData(file *ini.File) {
	DbHost = file.Section("database").Key("DbHost").String()
	DbPort = file.Section("database").Key("DbPort").String()
	DbUser = file.Section("database").Key("DbUser").String()
	DbPassWord = file.Section("database").Key("DbPassWord").String()
	DbName = file.Section("database").Key("DbName").String()
}
