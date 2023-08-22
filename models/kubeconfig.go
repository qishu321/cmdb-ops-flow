package models

import (
	"cmdb-ops-flow/utils/msg"
	"errors"
	"github.com/jinzhu/gorm"
)

type KubeConfig struct {
	ID             int    `json:"id" db:"id" form:"id"`
	Kubeconfigid   int64  `gorm:"type:bigint;not null" json:"kubeconfigid"`
	Kubeconfigname string `json:"kubeconfigname" db:"kubeconfigname" form:"kubeconfigname"`
	Kubeconfigdata string `json:"kubeconfigdata" db:"kubeconfigdata" form:"kubeconfigdata" gorm:"type:text"`
	Label          string `json:"label" db:"label" form:"label"`
}
type ClusterConfig struct {
	APIVersion string `yaml:"apiVersion"`
	Clusters   []struct {
		Name    string `yaml:"name"`
		Cluster struct {
			CertificateAuthorityData string `yaml:"certificate-authority-data"`
			Server                   string `yaml:"server"`
		} `yaml:"cluster"`
	} `yaml:"clusters"`
	Contexts []struct {
		Name    string `yaml:"name"`
		Context struct {
			Cluster string `yaml:"cluster"`
			User    string `yaml:"user"`
		} `yaml:"context"`
	} `yaml:"contexts"`
	CurrentContext string   `yaml:"current-context"`
	Kind           string   `yaml:"kind"`
	Preferences    struct{} `yaml:"preferences"`
	Users          []struct {
		Name string `yaml:"name"`
		User struct {
			ClientCertificateData string `yaml:"client-certificate-data"`
			ClientKeyData         string `yaml:"client-key-data"`
		} `yaml:"user"`
	} `yaml:"users"`
}

func AddKubeConfig(config KubeConfig) (interface{}, error) {
	err := db.Create(&config).Error
	return config, err
}
func EditKubeConfig(config KubeConfig) (interface{}, error) {
	err := db.Select("id").Where("kubeconfigid = ?", config.Kubeconfigid).First(&config).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return msg.ERROR_SCRIPT_GET_WRONG, err
		}
		return msg.InvalidParams, err
	}
	updateData := map[string]interface{}{
		"Kubeconfigname": config.Kubeconfigname,
		"Kubeconfigdata": config.Kubeconfigdata,
		"Label":          config.Label,
	}

	err = db.Model(&config).Where("kubeconfigid = ?", config.Kubeconfigid).Updates(updateData).Error
	if err != nil {
		return config, err
	}

	return config, err
}
func DelKubeConfig(name int64) (code int) {
	var config KubeConfig
	db.Select("id").Where("kubeconfigname = ?", name).First(&config)
	if config.ID > 0 {
		err = db.Where("kubeconfigid = ?", name).Delete(&config).Error
		if err != nil {
			return msg.InvalidParams
		}
		return msg.SUCCSE
	} else {
		return msg.ERROR_job_EDIT_WRONG
	}

}
func GetKubeConfigList(id int) ([]KubeConfig, error) {
	var list []KubeConfig
	if id != 0 {
		res := db.Debug().Where("id = ?", id).Find(&list)
		return list, res.Error
	} else {
		res := db.Debug().Find(&list)
		return list, res.Error
	}
}
func CheckKubeConfig(KubeConfigname string) (code int) {
	var config KubeConfig
	db.Select("id").Where("kubeconfigname = ?", KubeConfigname).First(&config)
	if config.ID > 0 {
		return msg.ERROR_job_GET_INFO
	}
	return msg.SUCCSE
}
