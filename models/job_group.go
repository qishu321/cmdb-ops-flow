package models

import (
	"cmdb-ops-flow/utils/msg"
	"errors"
	"github.com/jinzhu/gorm"
)

//先新建jobgroup组，然后组里就都是相同的jobgroup,不可修改。
//jobleve 优先级判断，判断jobgroup组里的所有数据，然后根据jobleve的从小往大去排序，创建的时候要进行判断，不能有相同的jobleve

type JobGroup struct {
	ID           int    `json:"id" db:"id" form:"id"`
	Jobgroupid   int64  `gorm:"type:bigint;not null" json:"jobgroupid" validate:"required"`
	Jobgroupname string `json:"jobgroupname" db:"jobgroupname" form:"jobgroupname"`
	Label        string `json:"label" db:"label" form:"label"`
}

func AddJobGroup(jobgroup JobGroup) (interface{}, error) {
	err := db.Create(&jobgroup).Error
	return jobgroup, err
}
func EditJobGroup(jobgroup JobGroup) (interface{}, error) {
	err := db.Select("id").Where("jobgroupid = ?", jobgroup.Jobgroupid).First(&jobgroup).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return msg.ERROR_SCRIPT_GET_WRONG, err
		}
		return msg.InvalidParams, err
	}
	updateData := map[string]interface{}{
		"Jobgroup": jobgroup.Jobgroupname,
		"Label":    jobgroup.Label,
	}

	err = db.Model(&jobgroup).Where("jobgroupid = ?", jobgroup.Jobgroupid).Updates(updateData).Error
	if err != nil {
		return jobgroup, err
	}

	return jobgroup, err
}
func DelJobGroup(name int64) (code int) {
	var jobgroup JobGroup
	db.Select("id").Where("jobgroupid = ?", name).First(&jobgroup)
	if jobgroup.ID > 0 {
		err = db.Where("jobgroupid = ?", name).Delete(&jobgroup).Error
		if err != nil {
			return msg.InvalidParams
		}
		return msg.SUCCSE
	} else {
		return msg.ERROR_job_EDIT_WRONG
	}

}
func GetJobGroupList(id int) ([]JobGroup, error) {
	var list []JobGroup
	if id != 0 {
		res := db.Debug().Where("id = ?", id).Find(&list)
		return list, res.Error
	} else {
		res := db.Debug().Find(&list)
		return list, res.Error
	}
}
func CheckJobGroup(jobgroupname string) (code int) {
	var jobgroup JobGroup
	db.Select("id").Where("jobgroupname = ?", jobgroupname).First(&jobgroup)
	if jobgroup.ID > 0 {
		return msg.ERROR_job_GET_INFO
	}
	return msg.SUCCSE
}
