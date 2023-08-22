package main

import (
	"cmdb-ops-flow/models"
	"context"
	"fmt"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"os"
)

func main() {
	models.InitDb()

	// 调用 Kubeinit 函数初始化 KubeConfig
	kubeConfigs, err := Kubeinit(1)
	if err != nil {
		fmt.Printf("Error initializing KubeConfig: %v\n", err)
		return
	}

	// 打印解析后的 KubeConfig 数据
	fmt.Println("Parsed KubeConfig:")
	for _, kubeConfig := range kubeConfigs {
		fmt.Println("KubeConfig Name:", kubeConfig.Kubeconfigname)
		fmt.Println("KubeConfig Kubeconfigdata:", kubeConfig.Kubeconfigdata)

		KubeConfig, err := clientcmd.NewClientConfigFromBytes([]byte(kubeConfig.Kubeconfigdata))
		if err != nil {
			fmt.Printf("Error creating Kubernetes client config: %v\n", err)
			return
		}

		clientConfig, err := KubeConfig.ClientConfig()
		if err != nil {
			fmt.Printf("Error creating Kubernetes client: %v\n", err)
			return
		}

		clientset, err := kubernetes.NewForConfig(clientConfig)
		if err != nil {
			fmt.Printf("Error creating Kubernetes client: %v\n", err)
			os.Exit(1)
		}

		pods, err := clientset.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			fmt.Printf("Error listing pods: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Pods in the cluster:\n")
		for _, pod := range pods.Items {
			fmt.Printf(" - %s\n", pod.GetName())
		}
	}
}

func Kubeinit(id int) ([]models.KubeConfig, error) {
	// 调用 GetKubeConfigList 获取 KubeConfig 数据
	kubeConfigs, err := models.GetKubeConfigList(id)
	if err != nil {
		return nil, err
	}

	// 创建一个切片，用于存储解析后的 KubeConfig 对象
	var parsedConfigs []models.KubeConfig

	// 遍历每个 KubeConfig 数据并解析成 KubeConfig 对象
	for _, kubeConfig := range kubeConfigs {
		// 创建一个新的 KubeConfig 对象并填充它的字段
		parsedConfig := models.KubeConfig{
			ID:             kubeConfig.ID,
			Kubeconfigid:   kubeConfig.Kubeconfigid,
			Kubeconfigname: kubeConfig.Kubeconfigname,
			Kubeconfigdata: kubeConfig.Kubeconfigdata,
			Label:          kubeConfig.Label,
		}

		// 将解析后的 KubeConfig 对象添加到切片中
		parsedConfigs = append(parsedConfigs, parsedConfig)
	}

	// 返回解析后的 KubeConfig 对象切片和可能的错误
	return parsedConfigs, nil
}
