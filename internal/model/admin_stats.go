package model

// AdminStats 管理员统计信息
type AdminStats struct {
	TotalStudents    int     `json:"total_students"`
	TotalInstructors int     `json:"total_instructors"`
	TotalCourses     int     `json:"total_courses"`
	TotalSections    int     `json:"total_sections"`
	TotalDepartments int     `json:"total_departments"`
	TotalClassrooms  int     `json:"total_classrooms"`
	TotalBudget      float64 `json:"total_budget"`
} 