package k8s

import (
	v1 "k8s.io/api/core/v1"
	"time"
)

type Labels map[string]string
type Annotations map[string]string

//type LimitRange map[string]string
//type ResourceQuota map[string]string

type Node struct {
	Name                    string      `json:"name"`
	Status                  string      `json:"status"`
	Taints                  []v1.Taint  `json:"taints"`
	Labels                  Labels      `json:"labels"`
	OsImage                 string      `json:"os_image"`
	InternalIp              string      `json:"internal_ip"`
	Annotations             Annotations `json:"annotations"`
	KernelVersion           string      `json:"kernel_version"`
	KubeletVersion          string      `json:"kubelet_version"`
	CreationTimestamp       time.Time   `json:"creation_timestamp"`
	ContainerRuntimeVersion string      `json:"container_runtime_version"`
}

type NodeMetrics struct {
	ID           int     `json:"id"`
	Name         string  `json:"name"`
	UsedCpu      float64 `json:"used_cpu"`
	TotalCpu     float64 `json:"total_cpu"`
	UsedMemory   float64 `json:"used_memory"`
	TotalMemory  float64 `json:"total_memory"`
	ReadyNodeNum int     `json:"readyNodeNum"`
	TotalNodeNum int     `json:"totalNodeNum"`
}
