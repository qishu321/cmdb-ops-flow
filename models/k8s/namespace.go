package k8s

import "time"

type NameSpace struct {
	ID int `json:"id"`

	Name              string      `json:"name"`
	Status            string      `json:"status"`
	Labels            Labels      `json:"labels"`
	Annotations       Annotations `json:"annotations"`
	CreationTimestamp time.Time   `json:"creation_timestamp"`
}
