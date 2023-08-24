package service_k8s

import (
	"cmdb-ops-flow/models/k8s"
	"context"
	"fmt"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func GetallPods(id int) ([]k8s.Pod, error) {
	clientSet, err := k8s.GetKubeConfig(id)
	if err != nil {
		fmt.Println(err)
	}
	podList, _ := clientSet.CoreV1().Pods("").List(context.Background(), metav1.ListOptions{})

	pods := make([]k8s.Pod, 0, 20)
	for _, item := range podList.Items {
		containers := make([]k8s.Container, 0, len(item.Status.ContainerStatuses))
		for _, containerStatus := range item.Status.ContainerStatuses {
			container := k8s.Container{Name: containerStatus.Name, Ready: containerStatus.Ready, RestartCount: int(containerStatus.RestartCount),
				Image: containerStatus.Image, ImageId: containerStatus.ImageID, ContainerId: containerStatus.ContainerID}
			containers = append(containers, container)
		}
		pod := k8s.Pod{Name: item.Name, Namespace: item.Namespace, Status: string(item.Status.Phase), CreationTimestamp: item.CreationTimestamp.Time,
			Containers: containers, Labels: item.Labels, Annotations: item.Annotations, PodIp: item.Status.PodIP, NodeName: item.Spec.NodeName}
		pods = append(pods, pod)
	}

	return pods, err
}

func GetPods(id int, ns string) ([]k8s.Pod, error) {
	clientSet, err := k8s.GetKubeConfig(id)
	if err != nil {
		fmt.Println(err)
	}
	podList, _ := clientSet.CoreV1().Pods(ns).List(context.Background(), metav1.ListOptions{})

	pods := make([]k8s.Pod, 0, 20)
	for _, item := range podList.Items {
		containers := make([]k8s.Container, 0, len(item.Status.ContainerStatuses))
		for _, containerStatus := range item.Status.ContainerStatuses {
			container := k8s.Container{Name: containerStatus.Name, Ready: containerStatus.Ready, RestartCount: int(containerStatus.RestartCount),
				Image: containerStatus.Image, ImageId: containerStatus.ImageID, ContainerId: containerStatus.ContainerID}
			containers = append(containers, container)
		}
		pod := k8s.Pod{Name: item.Name, Namespace: item.Namespace, Status: string(item.Status.Phase), CreationTimestamp: item.CreationTimestamp.Time,
			Containers: containers, Labels: item.Labels, Annotations: item.Annotations, PodIp: item.Status.PodIP, NodeName: item.Spec.NodeName}
		pods = append(pods, pod)
	}

	return pods, err
}
