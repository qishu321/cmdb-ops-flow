package k8s

type Container struct {
	Name         string `json:"name"`
	Ready        bool   `json:"ready"`
	Image        string `json:"image"`
	ImageId      string `json:"image_id"`
	ContainerId  string `json:"container_id"`
	RestartCount int    `json:"restart_count"`
}
