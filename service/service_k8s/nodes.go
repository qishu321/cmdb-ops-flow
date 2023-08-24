package service_k8s

import (
	"cmdb-ops-flow/models/k8s"
	"context"
	"fmt"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/version"
	"k8s.io/klog"
)

func GetNodeMetricsInfo(id int) ([]k8s.NodeMetrics, error) {
	clientSet, err := k8s.GetKubeConfig(id)
	if err != nil {
		fmt.Println(err)
	}

	ctx := context.Background()
	nodeList, err := clientSet.CoreV1().Nodes().List(ctx, metav1.ListOptions{})
	if err != nil {
		klog.Error(err)
		return nil, err
	}

	var nodeMetricsList []k8s.NodeMetrics
	var totalNodeMetrics k8s.NodeMetrics

	for _, node := range nodeList.Items {
		var nodeMetrics k8s.NodeMetrics

		nodeMetrics.Name = node.Name
		nodeMetrics.TotalNodeNum = len(nodeList.Items)
		for _, condition := range node.Status.Conditions {
			if condition.Type == "Ready" {
				if condition.Status == "True" {
					nodeMetrics.ReadyNodeNum += 1
					totalNodeMetrics.ReadyNodeNum += 1
				}
			}
		}
		cpu := node.Status.Allocatable.Cpu().AsApproximateFloat64()
		nodeMetrics.TotalCpu += cpu
		totalNodeMetrics.TotalCpu += cpu
		memory := node.Status.Allocatable.Memory().AsApproximateFloat64()
		nodeMetrics.TotalMemory += memory
		totalNodeMetrics.TotalMemory += memory

		podsList, _ := clientSet.CoreV1().Pods("").List(ctx, metav1.ListOptions{FieldSelector: fmt.Sprintf("spec.nodeName=%s", node.Name)})
		pods := podsList.Items
		for _, pod := range pods {
			for _, container := range pod.Spec.Containers {
				cpu := container.Resources.Requests.Cpu().AsApproximateFloat64()
				nodeMetrics.UsedCpu += cpu
				totalNodeMetrics.UsedCpu += cpu
				memory := container.Resources.Requests.Memory().AsApproximateFloat64()
				nodeMetrics.UsedMemory += memory
				totalNodeMetrics.UsedMemory += memory
			}
		}

		nodeMetricsList = append(nodeMetricsList, nodeMetrics)
	}

	// Create a flag to track if there's an empty name node
	var emptyNameNodeFound bool

	for i, nodeMetrics := range nodeMetricsList {
		if nodeMetrics.Name == "" {
			nodeMetricsList[i].Name = "汇总值"
			emptyNameNodeFound = true
			break
		}
	}

	// If no empty name node was found, add the summary node at the beginning
	if !emptyNameNodeFound {
		nodeMetricsList = append([]k8s.NodeMetrics{
			{
				ID:           0,
				Name:         "汇总值",
				UsedCpu:      totalNodeMetrics.UsedCpu,
				TotalCpu:     totalNodeMetrics.TotalCpu,
				UsedMemory:   totalNodeMetrics.UsedMemory,
				TotalMemory:  totalNodeMetrics.TotalMemory,
				ReadyNodeNum: totalNodeMetrics.ReadyNodeNum,
				TotalNodeNum: totalNodeMetrics.TotalNodeNum,
			},
		}, nodeMetricsList...)
	}

	return nodeMetricsList, nil
}

func Getnodes(id int) ([]k8s.Node, error) {
	clientSet, err := k8s.GetKubeConfig(id)
	if err != nil {
		fmt.Println(err)
	}

	ctx := context.Background()
	nodeList, err := clientSet.CoreV1().Nodes().List(ctx, metav1.ListOptions{})
	if err != nil {
		klog.Error(err)
		return nil, err
	}

	nodes := make([]k8s.Node, 0, 5)
	for _, item := range nodeList.Items {
		node := k8s.Node{Name: item.Name, Labels: item.Labels, Annotations: item.Annotations, CreationTimestamp: item.CreationTimestamp.Time,
			Taints: item.Spec.Taints, Status: getReadyStatus(item.Status.Conditions), InternalIp: getInternalIp(item.Status.Addresses),
			KernelVersion: item.Status.NodeInfo.KernelVersion, KubeletVersion: item.Status.NodeInfo.KubeletVersion,
			ContainerRuntimeVersion: item.Status.NodeInfo.ContainerRuntimeVersion, OsImage: item.Status.NodeInfo.OSImage}
		nodes = append(nodes, node)
	}
	return nodes, nil
}

func getReadyStatus(conditions []v1.NodeCondition) string {
	for _, condition := range conditions {
		if condition.Type == v1.NodeReady {
			return string(condition.Status)
		}
	}
	return "notfound"
}
func getInternalIp(addresses []v1.NodeAddress) string {
	for _, address := range addresses {
		if address.Type == v1.NodeInternalIP {
			return address.Address
		}
	}
	return "notfound"
}

func Version(id int) (*version.Info, error) {
	clientSet, err := k8s.GetKubeConfig(id)
	if err != nil {
		fmt.Println(err)
	}

	version, err := clientSet.ServerVersion()
	if err != nil {
		klog.Error(err)
		return nil, err
	}
	return version, nil
}
