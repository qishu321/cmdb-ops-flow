package k8s

import "time"

type Pod struct {
	Name              string      `json:"name"`
	PodIp             string      `json:"pod_ip"`
	Status            string      `json:"status"`
	Labels            Labels      `json:"labels"`
	NodeName          string      `json:"node_name"`
	Namespace         string      `json:"namespace"`
	Containers        []Container `json:"containers"`
	Annotations       Annotations `json:"annotations"`
	CreationTimestamp time.Time   `json:"creation_timestamp"`
}
