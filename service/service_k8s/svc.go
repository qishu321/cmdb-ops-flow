package service_k8s

import (
	"cmdb-ops-flow/models/k8s"
	"context"
	"fmt"
	v1 "k8s.io/api/core/v1"
	"strconv"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func Getsvs(id int, ns string) ([]k8s.Svc, error) {
	clientSet, err := k8s.GetKubeConfig(id)
	if err != nil {
		return nil, err
	}

	serviceList, err := clientSet.CoreV1().Services(ns).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	var svcs []k8s.Svc

	for _, service := range serviceList.Items {
		prots := make([]string, 0, len(service.Spec.Ports))
		for _, servicePort := range service.Spec.Ports {
			portStr := fmt.Sprintf("%d/%s", servicePort.Port, servicePort.Protocol)
			prots = append(prots, portStr)
		}

		if service.Spec.SessionAffinity == "ClientIP" {
			svc := k8s.Svc{
				ID:              id,
				Name:            service.Name,
				Namespace:       service.Namespace,
				Type:            string(service.Spec.Type),
				ClusterIp:       service.Spec.ClusterIP,
				Ports:           prots,
				Selector:        service.Spec.Selector,
				Labels:          service.Labels,
				Annotations:     service.Annotations,
				NodePort:        service.Spec.Ports[0].NodePort, // Assuming there's at least one port defined
				Protocol:        string(service.Spec.Ports[0].Protocol),
				SessionAffinity: string(service.Spec.SessionAffinity),
				SessionAffinityConfig: k8s.SessionAffinityConfig{
					ClientIP: k8s.ClientIP{
						TimeoutSeconds: *service.Spec.SessionAffinityConfig.ClientIP.TimeoutSeconds,
					},
				},
			}
			svcs = append(svcs, svc)
		} else {
			svc := k8s.Svc{
				ID:              id,
				Name:            service.Name,
				Namespace:       service.Namespace,
				Type:            string(service.Spec.Type),
				ClusterIp:       service.Spec.ClusterIP,
				Ports:           prots,
				Selector:        service.Spec.Selector,
				Labels:          service.Labels,
				Annotations:     service.Annotations,
				NodePort:        service.Spec.Ports[0].NodePort, // Assuming there's at least one port defined
				Protocol:        string(service.Spec.Ports[0].Protocol),
				SessionAffinity: string(service.Spec.SessionAffinity),
			}
			//svc := k8s.Svc{
			//	ID:          id,
			//	Name:        service.Name,
			//	Namespace:   service.Namespace,
			//	Type:        string(service.Spec.Type),
			//	ClusterIp:   service.Spec.ClusterIP,
			//	Ports:       prots,
			//	Selector:    service.Spec.Selector,
			//	Labels:      service.Labels,
			//	Annotations: service.Annotations,
			//	NodePort:    service.Spec.Ports[0].NodePort,         // Assuming there's at least one port defined
			//	Protocol:    string(service.Spec.Ports[0].Protocol),
			//	SessionAffinity: string(service.Spec.SessionAffinity),
			//	SessionAffinityConfig: k8s.SessionAffinityConfig{
			//		ClientIP:k8s.ClientIP{
			//			TimeoutSeconds: *service.Spec.SessionAffinityConfig.ClientIP.TimeoutSeconds,
			//		},
			//	},
			//
			//
			//}
			svcs = append(svcs, svc)
		}
	}

	return svcs, nil
}

func EditSvc(id int, updatedSvc k8s.Svc) (*v1.Service, error) {
	clientSet, err := k8s.GetKubeConfig(id)
	if err != nil {
		return nil, err
	}

	// 根据 updatedSvc 的 Namespace 和 Name 获取现有的 Service 资源
	existingSvc, err := clientSet.CoreV1().Services(updatedSvc.Namespace).Get(context.TODO(), updatedSvc.Name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	// 更新现有的 Service 资源
	existingSvc.Labels = updatedSvc.Labels
	existingSvc.Annotations = updatedSvc.Annotations
	existingSvc.Spec.Type = v1.ServiceType(updatedSvc.Type)

	// 更新端口信息
	var updatedServicePorts []v1.ServicePort
	for _, portStr := range updatedSvc.Ports {
		portNumber, err := strconv.Atoi(portStr)
		if err != nil {
			return nil, err
		}

		protocol := v1.ProtocolTCP
		if updatedSvc.Protocol == "udp" {
			protocol = v1.ProtocolUDP
		}

		servicePort := v1.ServicePort{
			Name:     portStr,
			Port:     int32(portNumber),
			Protocol: protocol,
		}
		updatedServicePorts = append(updatedServicePorts, servicePort)
	}
	existingSvc.Spec.Ports = updatedServicePorts
	existingSvc.Spec.Ports[0].NodePort = updatedSvc.NodePort

	existingSvc.Spec.SessionAffinity = v1.ServiceAffinity(updatedSvc.SessionAffinity)
	existingSvc.Spec.SessionAffinityConfig = &v1.SessionAffinityConfig{
		ClientIP: &v1.ClientIPConfig{
			TimeoutSeconds: &updatedSvc.SessionAffinityConfig.ClientIP.TimeoutSeconds,
		},
	}

	// 其他字段根据你的需求进行更新

	// 更新 Service 资源
	updatedService, err := clientSet.CoreV1().Services(updatedSvc.Namespace).Update(context.TODO(), existingSvc, metav1.UpdateOptions{})
	if err != nil {
		return nil, err
	}

	return updatedService, nil
}

func AddSvc(id int, data k8s.Svc) (*v1.Service, error) {
	clientSet, err := k8s.GetKubeConfig(id)
	if err != nil {
		fmt.Println(err)
	}
	var servicePorts []v1.ServicePort

	for _, portStr := range data.Ports {
		portNumber, err := strconv.Atoi(portStr)
		if err != nil {
			fmt.Println(err)
			// 处理端口转换错误
		}

		servicePort := v1.ServicePort{
			Name:     portStr,
			Port:     int32(portNumber),
			Protocol: v1.ProtocolTCP, // 默认使用 TCP 协议
			// Other fields like TargetPort can also be set here
		}

		// 如果协议为 "udp"，则设置为 UDP 协议
		if data.Protocol == "udp" {
			servicePort.Protocol = v1.ProtocolUDP
		}

		servicePorts = append(servicePorts, servicePort)
	}
	//sessionAffinityConfig := &k8s.SessionAffinityConfig{
	//	ClientIP: k8s.ClientIP{
	//		TimeoutSeconds: data.SessionAffinityConfig.ClientIP.TimeoutSeconds,
	//	},
	//}

	serviceType := v1.ServiceTypeClusterIP
	if data.Type == "NodePort" {
		serviceType = v1.ServiceTypeNodePort
	}

	service := &v1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:   data.Name,
			Labels: data.Labels,
		},
		Spec: v1.ServiceSpec{
			Selector:        data.Selector,
			Type:            serviceType,
			Ports:           servicePorts,
			SessionAffinity: v1.ServiceAffinity(data.SessionAffinity),
			SessionAffinityConfig: &v1.SessionAffinityConfig{
				ClientIP: &v1.ClientIPConfig{
					TimeoutSeconds: &data.SessionAffinityConfig.ClientIP.TimeoutSeconds,
				},
			},
		},
	}

	serviceList, err := clientSet.CoreV1().Services(data.Namespace).Create(context.TODO(), service, metav1.CreateOptions{})

	return serviceList, err
}

func Delsvs(id int, nsName string, svcName string) (code int, err error) {
	clientSet, err := k8s.GetKubeConfig(id)
	if err != nil {
		fmt.Println(err)
	}
	err = clientSet.CoreV1().Services(nsName).Delete(context.Background(), svcName, metav1.DeleteOptions{})
	if err != nil {
		return 400, err
	}

	return 200, err
}
