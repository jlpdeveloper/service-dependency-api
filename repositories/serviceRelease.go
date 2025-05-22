package repositories

type ServiceReleaseInfo struct {
	ServiceName string `json:"service_name"`
	ServiceType string `json:"service_type"`
	Release
}
