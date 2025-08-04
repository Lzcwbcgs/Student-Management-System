package model

// StudentCreateRequest 创建学生的请求模型
type StudentCreateRequest struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Dept     string `json:"dept"`
	Password string `json:"password"`
}

// StudentUpdateRequest 更新学生信息的请求模型
type StudentUpdateRequest struct {
	Name string `json:"name"`
	Dept string `json:"dept"`
}
