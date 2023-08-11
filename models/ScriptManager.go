package models

import (
	"cmdb-ops-flow/utils/msg"
	"errors"
	"github.com/jinzhu/gorm"
)

type ScriptManager struct {
	ID        int    `json:"id" db:"id" form:"id"`
	Name      string `json:"name" binding:"required" form:"name"`
	MachineID int64  `gorm:"type:bigint;not null;column:machineid" json:"machineid" validate:"required"`

	Type        string `json:"type" binding:"required" form:"type"`
	Description string `json:"description"  form:"description"`
	Script      string `json:"script" binding:"required" form:"script" gorm:"type:text"`

	Label string `json:"label" db:"label" form:"label"`
}

func AddScript(Script ScriptManager) (interface{}, error) {
	err := db.Create(&Script).Error
	return Script, err
}
func EditScript(Script ScriptManager) (interface{}, error) {
	err := db.Select("id").Where("machineid = ?", Script.MachineID).First(&Script).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return msg.ERROR_SCRIPT_GET_WRONG, err
		}
		return msg.InvalidParams, err
	}
	updateData := map[string]interface{}{
		"Name":        Script.Name,
		"Type":        Script.Type,
		"Description": Script.Description,
		"Script":      Script.Script,
		"Label":       Script.Label,
	}

	err = db.Model(&Script).Where("machineid = ?", Script.MachineID).Updates(updateData).Error
	if err != nil {
		return Script, err
	}

	return Script, err
}
func DelScript(name int64) (code int) {
	var Script ScriptManager
	db.Select("id").Where("machineid = ?", name).First(&Script)
	if Script.ID > 0 {
		err = db.Where("machineid = ?", name).Delete(&Script).Error
		if err != nil {
			return msg.InvalidParams
		}
		return msg.SUCCSE
	} else {
		return msg.ERROR_SCRIPT_EDIT_WRONG
	}

}
func GetScriptList(id int) ([]ScriptManager, error) {
	var list []ScriptManager
	if id != 0 {
		res := db.Debug().Where("id = ?", id).Find(&list)
		return list, res.Error
	} else {
		res := db.Debug().Find(&list)
		return list, res.Error
	}
}
func CheckScript(name string) (code int) {
	var Script ScriptManager
	db.Select("id").Where("name = ?", name).First(&Script)
	if Script.ID > 0 {
		return msg.ERROR_SCRIPT_GET_INFO
	}
	return msg.SUCCSE
}
