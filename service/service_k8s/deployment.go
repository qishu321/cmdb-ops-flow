package service_k8s

import (
	"cmdb-ops-flow/models/k8s"
	"context"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
)

func GetDeployments(id int, nsName string) ([]k8s.Deployment, error) {
	clientSet, err := k8s.GetKubeConfig(id)
	if err != nil {
		return nil, err
	}

	deploymentList, err := clientSet.AppsV1().Deployments(nsName).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	var deployments []k8s.Deployment

	for _, deployment := range deploymentList.Items {

		podList, err := clientSet.CoreV1().Pods(nsName).List(context.TODO(), metav1.ListOptions{
			LabelSelector: labels.SelectorFromSet(deployment.Spec.Selector.MatchLabels).String(),
		})
		if err != nil {
			return nil, err
		}

		// 提取所需的 Pod 信息
		var pods []k8s.Pod
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

		// 将 Kubernetes 的 Deployment 转换为自定义的 Deployment 结构
		deploy := k8s.Deployment{
			Name:        deployment.ObjectMeta.Name,
			Namespace:   deployment.ObjectMeta.Namespace,
			Replicas:    *deployment.Spec.Replicas, // 示例，假设 Replicas 是一个指针
			Labels:      deployment.ObjectMeta.Labels,
			Annotations: deployment.ObjectMeta.Annotations,
			Strategy: k8s.Strategy{ // 使用 k8s.Strategy 而不是直接的 Strategy
				Type: string(deployment.Spec.Strategy.Type), // 将 DeploymentStrategyType 转为字符串
				RollingUpdate: k8s.RollingUpdateStrategy{ // 同样使用 k8s.RollingUpdateStrategy
					MaxSurge:       deployment.Spec.Strategy.RollingUpdate.MaxSurge.IntValue(),
					MaxUnavailable: deployment.Spec.Strategy.RollingUpdate.MaxUnavailable.IntValue(),
				},
			},
			Selector:          k8s.Selector(deployment.Spec.Selector.MatchLabels), // 转换为你自己定义的 Selector 类型
			Pods:              pods,
			CreationTimestamp: deployment.ObjectMeta.CreationTimestamp,
		}

		deployments = append(deployments, deploy)
	}

	return deployments, nil
}
