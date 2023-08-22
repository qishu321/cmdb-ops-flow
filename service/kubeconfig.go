package service

import (
	"cmdb-ops-flow/models"
	"cmdb-ops-flow/utils/common"
	"encoding/base64"
	"gopkg.in/yaml.v3"
)

func AddKubeConfig(configData []byte, kubeconfigname string) (data interface{}, err error) {
	var kubeConfig models.ClusterConfig
	err = yaml.Unmarshal(configData, &kubeConfig)
	if err != nil {
		return nil, err
	}

	// Serialize kubeConfig back to YAML
	yamlData, err := yaml.Marshal(&kubeConfig)
	if err != nil {
		return nil, err
	}
	base64Data := base64.StdEncoding.EncodeToString(yamlData)

	daoKubeConfig := models.KubeConfig{
		Kubeconfigid:   common.GenerateRandomNumber(),
		Kubeconfigname: kubeconfigname,
		Kubeconfigdata: base64Data, // Store the YAML representation
		Label:          kubeconfigname,
	}

	// Call your model's AddKubeConfig function
	data, err = models.AddKubeConfig(daoKubeConfig)
	if err != nil {
		return nil, err
	}

	// Return the serialized kubeConfig as the response data
	return data, err
}

//func AddKubeConfig(configData []byte) (interface{}, error) {
//	var kubeConfig models.ClusterConfig
//	err := yaml.Unmarshal(configData, &kubeConfig)
//	if err != nil {
//		return nil, err
//	}
//
//	// Serialize kubeConfig back to YAML
//	yamlData, err := yaml.Marshal(&kubeConfig)
//	if err != nil {
//		return nil, err
//	}
//	base64Data := base64.StdEncoding.EncodeToString(yamlData)
//
//	daoKubeConfig := models.KubeConfig{
//		Kubeconfigid:   common.GenerateRandomNumber(),
//		Kubeconfigname: "123",
//		Kubeconfigdata: base64Data, // Store the YAML representation
//		Label:         "123",
//	}
//
//	// Call your model's AddKubeConfig function
//	_, err = models.AddKubeConfig(daoKubeConfig)
//	if err != nil {
//		return nil, err
//	}
//
//	// Return the serialized kubeConfig as the response data
//	return string(yamlData), nil
//}

func EditKubeConfig(config models.KubeConfig) (data interface{}, err error) {

	daoKubeConfig := models.KubeConfig{
		Kubeconfigid:   config.Kubeconfigid,
		Kubeconfigname: config.Kubeconfigname,
		Kubeconfigdata: config.Kubeconfigdata,
		Label:          config.Label,
	}

	data, err = models.EditKubeConfig(daoKubeConfig)
	return data, err
}

// func GetScriptList(json models.ScriptManager) (data interface{}, err error) {
func GetKubeConfigList(id int) ([]models.KubeConfig, error) {
	list, err := models.GetKubeConfigList(id)
	return list, err
}
