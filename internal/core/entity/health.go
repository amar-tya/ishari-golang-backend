package entity

// HealthStatus represents the health status of a service
type HealthStatus struct {
	Status   string `json:"status"`
	Database string `json:"database,omitempty"`
}
