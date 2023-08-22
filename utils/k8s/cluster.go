package main

import (
	"fmt"
	"io/ioutil"
	"k8s.io/client-go/tools/clientcmd"
)

const fixedConfigPath = "./../../conf/kube"

func main() {
	fileInfoList, err := ioutil.ReadDir(fixedConfigPath)
	if err != nil {
		fmt.Println("Error reading directory:", err)
		return
	}

	var list []string

	for _, info := range fileInfoList {
		name := info.Name()
		list = append(list, name)
		KubeConfig, _ := clientcmd.BuildConfigFromFlags("", fixedConfigPath+"/"+name+"/"+"config")
		fmt.Println(KubeConfig)
	}

	fmt.Println(list)
}

//func main() {
//	// 从当前目录开始向上找 conf/kube/config 目录
//	currentDir, err := filepath.Abs(".")
//	if err != nil {
//		fmt.Println("Error getting current directory:", err)
//		return
//	}
//	fmt.Println("curr",currentDir)
//	// 递归向上查找 kubeconfig 文件
//	for {
//		kubeConfigPath := filepath.Join(currentDir, "conf", "kube", "config")
//		_, err := ioutil.ReadFile(kubeConfigPath)
//		if err == nil {
//			// 找到 kubeconfig 文件，现在可以使用它
//			fmt.Println("Found kubeconfig at:", kubeConfigPath)
//
//			// 使用 kubeConfigPath 构建 kubeconfig
//			config, err := clientcmd.BuildConfigFromFlags("", kubeConfigPath)
//			if err != nil {
//				fmt.Println("Error building kubeconfig:", err)
//				return
//			}
//			fmt.Println(config)
//
//			// 使用 config 连接到 Kubernetes 集群，执行你的操作
//
//			break
//		}
//		// 向上一级目录继续查找
//		parentDir := filepath.Dir(currentDir)
//		if parentDir == currentDir {
//			fmt.Println("Kubeconfig not found in parent directories.")
//			break
//		}
//
//		currentDir = parentDir
//	}
//
//	// 如果代码运行到这里，表示未找到 kubeconfig 文件
//	fmt.Println("Kubeconfig not found.")
//}
