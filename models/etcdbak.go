package models

import (
	"cmdb-ops-flow/utils/msg"
	"time"
)

type EtcdBak struct {
	ID             int         `json:"id" db:"id" form:"id"`
	Etcdbakid      int64       `gorm:"type:bigint;not null" json:"etcdbakid"`
	Etcdbakname    string      `json:"etcdbakname" db:"etcdbakname" form:"etcdbakname"`
	EtcdGroupname  string      `json:"etcdGroupname" db:"etcdGroupname" form:"etcdGroupname"`
	EtcdGroupnames []EtcdGroup `gorm:"FOREIGNKEY:EtcdGroupnames;ASSOCIATION_FOREIGNKEY:EtcdGroupname"`
	Key            string      `json:"key" db:"key"`
	Value          string      `json:"value" db:"value"`
	CreatedAt      time.Time   `gorm:"default:CURRENT_TIMESTAMP",json:"created_at"`
	Label          string      `json:"label" db:"label" form:"label"`
}

func AddEtcdBak(etcdBak EtcdBak) (interface{}, error) {
	err := db.Create(&etcdBak).Error
	return etcdBak, err
}

//	func EditEtcdBak(etcdBak EtcdBak) (interface{}, error) {
//		err := db.Select("id").Where("etcdbakid = ?", etcdBak.Etcdbakid).First(&etcdBak).Error
//		if err != nil {
//			if errors.Is(err, gorm.ErrRecordNotFound) {
//				return msg.ERROR_SCRIPT_GET_WRONG, err
//			}
//			return msg.InvalidParams, err
//		}
//		updateData := map[string]interface{}{
//			"Etcdbakname": etcdBak.Etcdbakname,
//			"EtcdGroupname": etcdBak.EtcdGroupname,
//			"Key": etcdBak.Key,
//			"Value": etcdBak.Value,
//			"CreatedAt":etcdBak.CreatedAt,
//			"Label":        etcdBak.Label,
//		}
//
//		err = db.Model(&etcdBak).Where("etcdbakid = ?", etcdBak.Etcdbakid).Updates(updateData).Error
//		if err != nil {
//			return etcdBak, err
//		}
//
//		return etcdBak, err
//	}
func DelEtcdBak(name int64) (code int) {
	var etcdBak EtcdBak
	db.Select("id").Where("etcdGroupid = ?", name).First(&etcdBak)
	if etcdBak.ID > 0 {
		err = db.Where("etcdGroupid = ?", name).Delete(&etcdBak).Error
		if err != nil {
			return msg.InvalidParams
		}
		return msg.SUCCSE
	} else {
		return msg.ERROR_job_EDIT_WRONG
	}

}
func GetEtcdBakList(id int) ([]EtcdBak, error) {
	var list []EtcdBak
	if id != 0 {
		res := db.Debug().Where("id = ?", id).Find(&list)
		return list, res.Error
	} else {
		res := db.Debug().Find(&list)
		return list, res.Error
	}
}
func CheckEtcdBak(Etcdbakname string) (code int) {
	var etcdBak EtcdBak
	db.Select("id").Where("etcdGroupname = ?", Etcdbakname).First(&etcdBak)
	if etcdBak.ID > 0 {
		return msg.ERROR_job_GET_INFO
	}
	return msg.SUCCSE
}
