package model

// Department 表示系部实体
type Department struct {
	DeptName  string  `json:"dept_name"` // 系部名称
	Building  string  `json:"building"`  // 所在建筑
	Budget    float64 `json:"budget"`    // 预算
}

// DepartmentCreateRequest 表示创建院系的请求
type DepartmentCreateRequest struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Building    string  `json:"building"`
	Budget      float64 `json:"budget"`
	Description interface{}
}

// DepartmentUpdateRequest 表示更新院系的请求
type DepartmentUpdateRequest struct {
	Name     string  `json:"name"`
	Building string  `json:"building"`
	Budget   float64 `json:"budget"`
}

// DepartmentStats 表示院系统计信息
type DepartmentStats struct {
	Department      Department `json:"department"`       // 院系信息
	StudentCount    int        `json:"student_count"`    // 学生数量
	InstructorCount int        `json:"instructor_count"` // 教师数量
	CourseCount     int        `json:"course_count"`     // 课程数量
}
