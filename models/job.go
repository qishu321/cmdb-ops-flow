package models

import (
	"cmdb-ops-flow/utils/msg"
	"errors"
	"github.com/jinzhu/gorm"
)

//先新建jobgroup组，然后组里就都是相同的jobgroup,不可修改。
//jobleve 优先级判断，判断jobgroup组里的所有数据，然后根据jobleve的从小往大去排序，创建的时候要进行判断，不能有相同的jobleve

type Job struct {
	ID               int    `json:"id" db:"id" form:"id"`
	Jobid            int64  `gorm:"type:bigint;not null" json:"jobid" validate:"required"`
	Jobleve          int64  `json:"jobleve" db:"jobleve" for:"jobleve" column:jobleve"`
	Jobname          string `json:"jobname" db:"jobname" form:"jobname"`
	Jobgroupid       int64  `json:"jobgroupid" db:"jobgroupid" for:"jobgroupid" column:jobgroupid"`
	Params           string `json:"params" db:"params" form:"params"`
	Machineid_Script int64  `json:"machineid_script" db:"machineid_script" for:"machineid_script" column:machineid_script"`
	Type             string `json:"type" binding:"required" form:"type"`
	Label            string `json:"label" db:"label" form:"label"`
}

func AddJob(Script Job) (interface{}, error) {
	err := db.Create(&Script).Error
	return Script, err
}
func EditJob(Script Job) (interface{}, error) {
	err := db.Select("id").Where("jobid = ?", Script.Jobid).First(&Script).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return msg.ERROR_SCRIPT_GET_WRONG, err
		}
		return msg.InvalidParams, err
	}
	updateData := map[string]interface{}{
		"Jobleve":          Script.Jobleve,
		"Jobname":          Script.Jobname,
		"Params":           Script.Params,
		"Machineid_Script": Script.Machineid_Script,
		"Type":             Script.Type,
		"Label":            Script.Label,
	}

	err = db.Model(&Script).Where("jobid = ?", Script.Jobid).Updates(updateData).Error
	if err != nil {
		return Script, err
	}

	return Script, err
}
func DelJob(name int64) (code int) {
	var job Job
	db.Select("id").Where("jobid = ?", name).First(&job)
	if job.ID > 0 {
		err = db.Where("jobid = ?", name).Delete(&job).Error
		if err != nil {
			return msg.InvalidParams
		}
		return msg.SUCCSE
	} else {
		return msg.ERROR_job_EDIT_WRONG
	}

}
func GetJobList(id int) ([]Job, error) {
	var list []Job
	if id != 0 {
		res := db.Debug().Where("id = ?", id).Find(&list)
		return list, res.Error
	} else {
		res := db.Debug().Find(&list)
		return list, res.Error
	}
}
func CheckJob(name string) (code int) {
	var job Job
	db.Select("id").Where("name = ?", name).First(&job)
	if job.ID > 0 {
		return msg.ERROR_job_GET_INFO
	}
	return msg.SUCCSE
}
