package k8s

import (
	"cmdb-ops-flow/models"
	"encoding/base64"
	"fmt"
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
		return nil, fmt.Errorf("invalid Base64 data")
	}

	// 使用kubeconfigBytes来构建Kubernetes配置
	config, err := clientcmd.RESTConfigFromKubeConfig(kubeconfigBytes)
	if err != nil {
		fmt.Println("使用kubeconfigBytes来构建Kubernetes配置:", err) // 添加错误日志输出

		return nil, err
	}

	// 创建Kubernetes客户端
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		fmt.Println("创建Kubernetes客户端", err) // 添加错误日志输出
		return nil, err
	}

	return clientset, nil
}
