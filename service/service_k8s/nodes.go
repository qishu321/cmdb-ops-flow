package service_k8s

import (
	"cmdb-ops-flow/models/k8s"
	"context"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/version"
	"k8s.io/klog"
	"log"
)

func Getnodes(id int) ([]k8s.Node, error) {
	clientSet, err := k8s.GetKubeConfig(id)
	if err != nil {
		log.Fatal(err)
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
		log.Fatal(err)
	}

	version, err := clientSet.ServerVersion()
	if err != nil {
		klog.Error(err)
		return nil, err
	}
	return version, nil
}
