package service_k8s

import (
	"cmdb-ops-flow/models/k8s"
	"context"
	"fmt"
	"io"
	v1 "k8s.io/api/core/v1"
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

//func GetPodLogs(id int, ns string) ([]k8s.Pod, error) {
//	clientSet, err := k8s.GetKubeConfig(id)
//	if err != nil {
//		fmt.Println(err)
//	}
//	podList, _ := clientSet.CoreV1().Pods(ns).List(context.Background(), metav1.ListOptions{})
//
//	pods := make([]k8s.Pod, 0, 20)
//	for _, item := range podList.Items {
//		containers := make([]k8s.Container, 0, len(item.Status.ContainerStatuses))
//		for _, containerStatus := range item.Status.ContainerStatuses {
//			container := k8s.Container{Name: containerStatus.Name, Ready: containerStatus.Ready, RestartCount: int(containerStatus.RestartCount),
//				Image: containerStatus.Image, ImageId: containerStatus.ImageID, ContainerId: containerStatus.ContainerID}
//			containers = append(containers, container)
//		}
//		pod := k8s.Pod{Name: item.Name, Namespace: item.Namespace, Status: string(item.Status.Phase), CreationTimestamp: item.CreationTimestamp.Time,
//			Containers: containers, Labels: item.Labels, Annotations: item.Annotations, PodIp: item.Status.PodIP, NodeName: item.Spec.NodeName}
//		pods = append(pods, pod)
//	}
//	req,err := clientSet.CoreV1().Pods(ns).GetLogs(pods[0].Name, &v1.PodLogOptions{Follow: true, Container: pods[0].Containers[0].Name}).Stream(context.TODO()) //hello 为container name
//	if err != nil {
//		fmt.Println(err)
//	}
//
//	return req, err
//}

func GetPodLogs(id int, ns string, podName string, containerName string) (io.ReadCloser, error) {
	clientset, err := k8s.GetKubeConfig(id)
	if err != nil {
		return nil, err
	}
	tailLines := int64(25) // 将100转换为int64类型

	req := clientset.CoreV1().Pods(ns).GetLogs(podName, &v1.PodLogOptions{
		TailLines: &tailLines, // 使用指向int64类型的指针
		Container: containerName,
	})
	fmt.Println(req)
	logsStream, err := req.Stream(context.TODO())
	if err != nil {
		return nil, fmt.Errorf("error opening log stream for container %s: %v", containerName, err)
	}

	return logsStream, nil
}
