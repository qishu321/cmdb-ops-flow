package models

import (
	"cmdb-ops-flow/utils/msg"
	"errors"
	"github.com/jinzhu/gorm"
	"strings"
)

//先新建jobgroup组，然后组里就都是相同的jobgroup,不可修改。
//jobleve 优先级判断，判断jobgroup组里的所有数据，然后根据jobleve的从小往大去排序，创建的时候要进行判断，不能有相同的jobleve

type Job struct {
	ID      int   `json:"id" db:"id" form:"id"`
	Jobid   int64 `gorm:"type:bigint;not null" json:"jobid" validate:"required"`
	Jobleve int64 `json:"jobleve" db:"jobleve" form:"jobleve" column:"jobleve"`

	Jobname   string     `json:"jobname" db:"jobname" form:"jobname"`
	Jobgroup  string     `json:"jobgroup" db:"jobgroup" form:"jobgroup"`
	Jobgroups []JobGroup `gorm:"FOREIGNKEY:Jobgroupname;ASSOCIATION_FOREIGNKEY:Jobgroup"`
	//Jobgroupid       int64  `json:"jobgroupid" db:"jobgroupid" for:"jobgroupid" column:jobgroupid"`
	Params      string          `json:"params" db:"params" form:"params"`
	Scriptname  string          `json:"scriptname" db:"scriptname" form:"scriptname"`
	Scriptnames []ScriptManager `gorm:"FOREIGNKEY:Name;ASSOCIATION_FOREIGNKEY:Scriptname"`
	Type        string          `json:"type"  form:"type"`
	Jobcmdbname string          `json:"jobcmdbname" db:"jobcmdbname" form:"jobcmdbname"`
	Cmdbnames   []Cmdb          `gorm:"FOREIGNKEY:Cmdbname;ASSOCIATION_FOREIGNKEY:Jobcmdbname"`
	Label       string          `json:"label" db:"label" form:"label"`
}

func GetJobList(id int) ([]Job, error) {
	var list []Job
	if id != 0 {
		res := db.Debug().Where("id = ?", id).Find(&list)
		return list, res.Error
	} else {
		res := db.Debug().Preload("Jobgroups").Preload("Scriptnames").Find(&list)
		for i, job := range list {
			cmdbNames := strings.Split(job.Jobcmdbname, ",")
			var cmdbList []Cmdb
			db.Where("cmdbname IN (?)", cmdbNames).Find(&cmdbList)
			list[i].Cmdbnames = cmdbList
		}
		return list, res.Error

		//res := db.Debug().Preload("Jobgroups").Preload("Scriptnames").Preload("Cmdbnames").Find(&list)
		//return list, res.Error
	}
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
		"Jobleve":     Script.Jobleve,
		"Jobname":     Script.Jobname,
		"Jobcmdbname": Script.Jobcmdbname,

		"Params": Script.Params,
		"Type":   Script.Type,
		"Label":  Script.Label,
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

func CheckJob(name string) (code int) {
	var job Job
	db.Select("id").Where("name = ?", name).First(&job)
	if job.ID > 0 {
		return msg.ERROR_job_GET_INFO
	}
	return msg.SUCCSE
}
