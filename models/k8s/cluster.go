package k8s

import (
	"cmdb-ops-flow/models"
	"encoding/base64"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func GetKubeConfig(id int) (*kubernetes.Clientset, error) {
	// 假设您的models包含一个用于从数据库中检索kubeconfig数据的函数
	kubeconfig, err := models.GetKubeConfigList(id)
	if err != nil {
		return nil, err
	}

	// 解码Base64编码的kubeconfig数据
	kubeconfigBytes, err := base64.StdEncoding.DecodeString(kubeconfig[0].Kubeconfigdata)
	if err != nil {
		return nil, err
	}

	// 使用kubeconfigBytes来构建Kubernetes配置
	config, err := clientcmd.RESTConfigFromKubeConfig(kubeconfigBytes)
	if err != nil {
		return nil, err
	}

	// 创建Kubernetes客户端
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return clientset, nil
}
