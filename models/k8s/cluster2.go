package k8s

//func GetKubeConfigList(id int) ([]models.KubeConfig, error) {
//	// 假设您的models包含一个用于从数据库中检索kubeconfig数据的函数
//	kubeconfig, err := models.GetKubeConfigList(id)
//	if err != nil {
//		return nil, err
//	}
//
//	// 解码Base64编码的kubeconfig数据
//	kubeconfigBytes, err := base64.StdEncoding.DecodeString(kubeconfig[0].Kubeconfigdata)
//	if err != nil {
//		return nil, err
//	}
//
//	// 使用kubeconfigBytes来构建Kubernetes配置
//	config, err := clientcmd.RESTConfigFromKubeConfig(kubeconfigBytes)
//	if err != nil {
//		return nil, err
//	}
//
//	// 创建Kubernetes客户端
//	clientset, err := kubernetes.NewForConfig(config)
//	if err != nil {
//		return nil, err
//	}
//
//	nodeList, err := clientset.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{}) // 获取集群中所有节点的信息列表
//	if err != nil {
//		fmt.Println(err)
//	}
//
//	for _, node := range nodeList.Items {
//		fmt.Printf("%s\t%s\t%s\t%s\t%s\t%s\t%s\n",
//			node.Name,
//			node.Status.Phase,
//			node.Status.NodeInfo.OSImage,
//			node.Status.NodeInfo.KubeletVersion,
//			node.Status.NodeInfo.OperatingSystem,
//			node.Status.NodeInfo.Architecture,
//			node.CreationTimestamp,
//		)
//	}
//
//	// 返回Kubeconfig列表和任何可能的错误
//	return kubeconfig, nil
//}
