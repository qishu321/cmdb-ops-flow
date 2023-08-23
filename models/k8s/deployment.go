package k8s

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

type Deployment struct {
	ID                int         `json:"id"`
	Name              string      `json:"name"`
	Namespace         string      `json:"namespace"`
	Replicas          int32       `json:"replicas"`
	Labels            Labels      `json:"labels"`
	Annotations       Annotations `json:"annotations"`
	Strategy          Strategy    `json:"strategy"`
	Selector          Selector    `json:"selector"`
	Pods              []Pod       `json:"template"`
	CreationTimestamp metav1.Time `json:"creation_timestamp"`
}

type Strategy struct {
	Type          string                `json:"type"`
	RollingUpdate RollingUpdateStrategy `json:"rollingUpdate"`
}

type RollingUpdateStrategy struct {
	MaxSurge       int `json:"maxSurge"`
	MaxUnavailable int `json:"maxUnavailable"`
}
