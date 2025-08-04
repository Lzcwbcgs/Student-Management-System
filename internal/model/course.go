package model

// Course 表示课程实体
type Course struct {
	ID      string   `json:"id"`                // 课程ID
	Title   string   `json:"title"`             // 课程名称
	Dept    string   `json:"dept"`              // 所属院系
	Credits float64  `json:"credits"`           // 学分
	Name    string   `json:"name"`              // 课程名称，用于显示
	Prereqs []Prereq `json:"prereqs,omitempty"` // 先修课程
}

// Prereq 表示先修课程关系
type Prereq struct {
	CourseID   string  `json:"course_id"`             // 课程ID
	PrereqID   string  `json:"prereq_id"`             // 先修课程ID
	PrereqInfo *Course `json:"prereq_info,omitempty"` // 先修课程详细信息
}

// CourseWithPrereqs 表示包含先修课程信息的课程
type CourseWithPrereqs struct {
	Course
	Prereqs   []Course `json:"prereqs,omitempty"`
	PrereqIDs []string `json:"prereq_ids,omitempty"`
}

// CourseCreateRequest 表示创建课程的请求
type CourseCreateRequest struct {
	ID        string   `json:"id"`
	Title     string   `json:"title"`
	Dept      string   `json:"dept"`
	Credits   float64  `json:"credits"`
	PrereqIDs []string `json:"prereq_ids,omitempty"`
}

// CourseUpdateRequest 表示更新课程的请求
type CourseUpdateRequest struct {
	Title   string  `json:"title"`
	Dept    string  `json:"dept"`
	Credits float64 `json:"credits"`
}
