package service

import (
	"cmdb-ops-flow/models"
	"cmdb-ops-flow/utils/common"
	"errors"
)

func AddScript(Script models.ScriptManager) (data interface{}, err error) {
	if Script.Name == "" || Script.Script == "" || Script.Label == "" || Script.Type == "" {
		return nil, errors.New("所有字段都是必填的")
	}

	daoScript := models.ScriptManager{
		MachineID:   common.GenerateRandomNumber(),
		Name:        Script.Name,
		Type:        Script.Type,
		Description: Script.Description,
		Script:      Script.Script,

		Label: Script.Label,
	}

	data, err = models.AddScript(daoScript)
	return data, err
}
func EditScript(Script models.ScriptManager) (data interface{}, err error) {
	daoScript := models.ScriptManager{
		MachineID:   Script.MachineID,
		Name:        Script.Name,
		Type:        Script.Type,
		Description: Script.Description,
		Script:      Script.Script,

		Label: Script.Label,
	}

	data, err = models.EditScript(daoScript)
	return data, err
}

// func GetScriptList(json models.ScriptManager) (data interface{}, err error) {
func GetScriptList(id int) ([]models.ScriptManager, error) {
	list, err := models.GetScriptList(id)
	return list, err
}
