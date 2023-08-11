package service

import (
	"cmdb-ops-flow/models"
	"cmdb-ops-flow/utils/common"
	"errors"
)

func AddJobGroup(jobgroup models.JobGroup) (data interface{}, err error) {
	if jobgroup.Jobgroupname == "" || jobgroup.Label == "" {
		return nil, errors.New("所有字段都是必填的")
	}
	daojobgroup := models.JobGroup{
		Jobgroupid:   common.GenerateRandomNumber(),
		Jobgroupname: jobgroup.Jobgroupname,
		Label:        jobgroup.Label,
	}

	data, err = models.AddJobGroup(daojobgroup)
	return data, err
}
func EditJobGroup(jobgroup models.JobGroup) (data interface{}, err error) {

	daojobgroup := models.JobGroup{
		Jobgroupid:   jobgroup.Jobgroupid,
		Jobgroupname: jobgroup.Jobgroupname,
		Label:        jobgroup.Label,
	}

	data, err = models.EditJobGroup(daojobgroup)
	return data, err
}

// func GetScriptList(json models.ScriptManager) (data interface{}, err error) {
func GetJobGroupList(id int) ([]models.JobGroup, error) {
	list, err := models.GetJobGroupList(id)
	return list, err
}
