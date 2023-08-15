package service

import (
	"cmdb-ops-flow/models"
	"cmdb-ops-flow/utils/common"
	"errors"
)

func AddJob(job models.Job) (data interface{}, err error) {
	if job.Jobname == "" || job.Label == "" || job.Jobgroup == "" {
		return nil, errors.New("所有字段都是必填的")
	}
	daojob := models.Job{
		Jobid:       common.GenerateRandomNumber(),
		Jobname:     job.Jobname,
		Jobleve:     job.Jobleve,
		Jobgroup:    job.Jobgroup,
		Params:      job.Params,
		Scriptname:  job.Scriptname,
		Type:        job.Type,
		Jobcmdbname: job.Jobcmdbname,
		Label:       job.Label,
	}

	data, err = models.AddJob(daojob)
	return data, err
}
func EditJob(job models.Job) (data interface{}, err error) {

	daojob := models.Job{
		Jobid:    job.Jobid,
		Jobname:  job.Jobname,
		Jobleve:  job.Jobleve,
		Jobgroup: job.Jobgroup,

		Params:      job.Params,
		Scriptname:  job.Scriptname,
		Type:        job.Type,
		Jobcmdbname: job.Jobcmdbname,

		Label: job.Label,
	}
	data, err = models.EditJob(daojob)
	return data, err
}

// func GetScriptList(json models.ScriptManager) (data interface{}, err error) {
func GetJobList(id int) ([]models.Job, error) {
	list, err := models.GetJobList(id)
	return list, err
}
