package model

// Advisor 表示学生导师关系
type Advisor struct {
	StudentID    string      `json:"student_id"`    // 学生ID
	InstructorID string      `json:"instructor_id"` // 导师ID
	
	// 关联信息
	Student      *Student    `json:"student,omitempty"`    // 学生信息
	Instructor   *Instructor `json:"instructor,omitempty"` // 导师信息
}

// AdvisorCreateRequest 表示创建导师关系的请求
type AdvisorCreateRequest struct {
	StudentID    string `json:"student_id"`
	InstructorID string `json:"instructor_id"`
}