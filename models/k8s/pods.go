package k8s

import (
	"time"
)

type Pod struct {
	ID int `json:"id" form:"id"`

	Name              string      `json:"name" form:"name"`
	PodIp             string      `json:"pod_ip"`
	Status            string      `json:"status"`
	Labels            Labels      `json:"labels"`
	NodeName          string      `json:"node_name"`
	Namespace         string      `json:"namespace" form:"namespace"`
	Containers        []Container `json:"containers"`
	ContainerName     string      `json:"container_name" form:"container_name"`
	Annotations       Annotations `json:"annotations"`
	CreationTimestamp time.Time   `json:"creation_timestamp"`
}
