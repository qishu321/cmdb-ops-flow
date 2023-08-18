package service

import (
	"cmdb-ops-flow/models"
	"cmdb-ops-flow/utils/common"
	"errors"
	"time"
)

func AddEtcdbak(etcdbak models.EtcdBak) (data interface{}, err error) {
	if etcdbak.EtcdGroupname == "" || etcdbak.Etcdbakname == "" {
		return nil, errors.New("所有字段都是必填的")
	}
	etcdbak.CreatedAt = time.Now()
	daoetcdbak := models.EtcdBak{
		Etcdbakid:     common.GenerateRandomNumber(),
		Etcdbakname:   etcdbak.Etcdbakname,
		EtcdGroupname: etcdbak.EtcdGroupname,
		Key:           etcdbak.Key,
		Value:         etcdbak.Value,
		CreatedAt:     etcdbak.CreatedAt,
		Label:         etcdbak.Label,
	}

	data, err = models.AddEtcdBak(daoetcdbak)
	return data, err
}

//func EditEtcdbak(jobgroup models.EtcdBak) (data interface{}, err error) {
//
//	daojobgroup := models.EtcdBak{
//		Jobgroupid:   jobgroup.Jobgroupid,
//		Jobgroupname: jobgroup.Jobgroupname,
//		Label:        jobgroup.Label,
//	}
//
//	data, err = models.EditJobGroup(daojobgroup)
//	return data, err
//}

// func GetScriptList(json models.ScriptManager) (data interface{}, err error) {
func GetEtcdbakList(id int) ([]models.EtcdBak, error) {
	list, err := models.GetEtcdBakList(id)
	return list, err
}
