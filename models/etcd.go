package models

import (
	"cmdb-ops-flow/utils/msg"
	"errors"
	"github.com/jinzhu/gorm"
)

type EtcdGroup struct {
	ID            int    `json:"id" db:"id" form:"id"`
	EtcdGroupid   int64  `gorm:"type:bigint;not null" json:"etcdGroupid"`
	EtcdGroupname string `json:"etcdGroupname" db:"etcdGroupname" form:"etcdGroupname"`
	EtcdEndpoints string `json:"etcdEndpoints" db:"etcdEndpoints" form:"etcdEndpoints"`
	Label         string `json:"label" db:"label" form:"label"`
}

func AddEtcdGroup(etcdgroup EtcdGroup) (interface{}, error) {
	err := db.Create(&etcdgroup).Error
	return etcdgroup, err
}
func EditEtcdGroup(etcdgroup EtcdGroup) (interface{}, error) {
	err := db.Select("id").Where("etcdGroupid = ?", etcdgroup.EtcdGroupid).First(&etcdgroup).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return msg.ERROR_SCRIPT_GET_WRONG, err
		}
		return msg.InvalidParams, err
	}
	updateData := map[string]interface{}{
		"EtcdGroupname": etcdgroup.EtcdGroupname,
		"EtcdEndpoints": etcdgroup.EtcdEndpoints,

		"Label": etcdgroup.Label,
	}

	err = db.Model(&etcdgroup).Where("etcdGroupid = ?", etcdgroup.EtcdGroupid).Updates(updateData).Error
	if err != nil {
		return etcdgroup, err
	}

	return etcdgroup, err
}
func DelEtcdGroup(name int64) (code int) {
	var etcdgroup EtcdGroup
	db.Select("id").Where("etcdGroupid = ?", name).First(&etcdgroup)
	if etcdgroup.ID > 0 {
		err = db.Where("etcdGroupid = ?", name).Delete(&etcdgroup).Error
		if err != nil {
			return msg.InvalidParams
		}
		return msg.SUCCSE
	} else {
		return msg.ERROR_job_EDIT_WRONG
	}

}
func GetEtcdGroupList(id int) ([]EtcdGroup, error) {
	var list []EtcdGroup
	if id != 0 {
		res := db.Debug().Where("id = ?", id).Find(&list)
		return list, res.Error
	} else {
		res := db.Debug().Find(&list)
		return list, res.Error
	}
}
func CheckEtcdGroup(etcdGroupname string) (code int) {
	var jobgroup EtcdGroup
	db.Select("id").Where("etcdGroupname = ?", etcdGroupname).First(&jobgroup)
	if jobgroup.ID > 0 {
		return msg.ERROR_job_GET_INFO
	}
	return msg.SUCCSE
}
