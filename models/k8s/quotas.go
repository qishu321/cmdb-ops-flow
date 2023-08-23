package k8s

type Quota struct {
	LimitsCpu   string `json:"limits_cpu"`
	LimitsMem   string `json:"limits_mem"`
	RequestsCpu string `json:"requests_cpu"`
	RequestsMem string `json:"requests_mem"`
}
