package model

// SystemStats 系统统计信息
type SystemStats struct {
	ActiveUsers     int     `json:"active_users"`
	ServerUptime    float64 `json:"server_uptime"`
	DatabaseSize    int64   `json:"database_size"`
	MemoryUsage     float64 `json:"memory_usage"`
	CPUUsage        float64 `json:"cpu_usage"`
	TotalOperations int64   `json:"total_operations"`
	ErrorRate       float64 `json:"error_rate"`
}
