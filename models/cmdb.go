package models

import (
	"cmdb-ops-flow/utils/msg"
	"errors"
	"github.com/jinzhu/gorm"
)

type Cmdb struct {
	ID         int    `json:"id" db:"id" form:"id"`
	Cmdbid     int64  `gorm:"type:bigint;not null" json:"cmdbid" validate:"required"`
	Cmdbname   string `json:"cmdbname" db:"cmdbname" form:"cmdbname"`
	PublicIP   string `json:"public_ip" db:"public_ip" form:"public_ip"`
	PrivateIP  string `json:"private_ip" db:"private_ip" form:"private_ip"`
	Username   string `json:"username" db:"username" form:"username"`
	Password   string `json:"password" db:"password" form:"password"`
	PrivateKey string `json:"private_key" db:"private_key" form:"private_key"`
	SSHPort    int    `json:"ssh_port" db:"ssh_port" form:"ssh_port"`
	Label      string `json:"label" db:"label" form:"label"`
}

func Addcmdb(cmdb Cmdb) (interface{}, error) {
	err := db.Create(&cmdb).Error
	return cmdb, err
}
func Editcmdb(cmdb Cmdb) (interface{}, error) {
	err := db.Select("id").Where("cmdbid = ?", cmdb.Cmdbid).First(&cmdb).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return msg.ERROR_CMDB_GET_WRONG, err
		}
		return msg.InvalidParams, err
	}
	updateData := map[string]interface{}{
		"Cmdbname":   cmdb.Cmdbname,
		"PublicIP":   cmdb.PublicIP,
		"PrivateIP":  cmdb.PrivateIP,
		"Username":   cmdb.Username,
		"Password":   cmdb.Password,
		"PrivateKey": cmdb.PrivateKey,
		"SSHPort":    cmdb.SSHPort,
		"Label":      cmdb.Label,
	}

	err = db.Model(&cmdb).Where("cmdbid = ?", cmdb.Cmdbid).Updates(updateData).Error
	if err != nil {
		return cmdb, err
	}

	return cmdb, err
}
func Delcmdb(name int64) (code int) {
	var cmdb Cmdb
	db.Select("id").Where("cmdbid = ?", name).First(&cmdb)
	if cmdb.ID > 0 {
		err = db.Where("cmdbid = ?", name).Delete(&cmdb).Error
		if err != nil {
			return msg.InvalidParams
		}
		return msg.SUCCSE
	} else {
		return msg.ERROR_CMDB_EDIT_WRONG
	}

}
func GetcmdbList(id int) ([]Cmdb, error) {
	var list []Cmdb
	if id != 0 {
		res := db.Debug().Where("id = ?", id).Find(&list)
		return list, res.Error
	} else {
		res := db.Debug().Find(&list)
		return list, res.Error
	}
}
func Checkcmdb(name string) (code int) {
	var cmdb Cmdb
	db.Select("id").Where("cmdbname = ?", name).First(&cmdb)
	if cmdb.ID > 0 {
		return msg.ERROR_CMDB_GET_INFO
	}
	return msg.SUCCSE
}
func SearchCmdb(keyword string) ([]Cmdb, int) {
	var cmdb []Cmdb
	db.Where("cmdbname LIKE ? OR public_ip LIKE ? OR private_ip LIKE ? OR label LIKE ?", "%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%").Find(&cmdb)
	if len(cmdb) > 0 {
		return cmdb, msg.SUCCSE
	}
	return nil, msg.ERROR_CMDB_GET_INFO
}
