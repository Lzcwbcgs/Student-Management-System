package model

// InstructorCreateRequest 创建教师的请求模型
type InstructorCreateRequest struct {
	ID       string  `json:"id"`
	Name     string  `json:"name"`
	Dept     string  `json:"dept"`
	Salary   float64 `json:"salary"`
	Password string  `json:"password"`
}

// InstructorUpdateRequest 更新教师信息的请求模型
type InstructorUpdateRequest struct {
	Name   string  `json:"name"`
	Dept   string  `json:"dept"`
	Salary float64 `json:"salary"`
}
