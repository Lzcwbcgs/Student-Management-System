package model

// Teaches 表示教师教授课程的关系
type Teaches struct {
	ID           string `json:"id"`           // 记录ID
	InstructorID string `json:"instructor_id"` // 教师ID
	CourseID     string `json:"course_id"`     // 课程ID
	SectionID    string `json:"section_id"`    // 章节ID
	Semester     string `json:"semester"`      // 学期
	Year         int    `json:"year"`          // 年份
	
	// 关联信息
	Instructor   *Instructor `json:"instructor,omitempty"` // 教师信息
	Course       *Course     `json:"course,omitempty"`     // 课程信息
	Section      *Section    `json:"section,omitempty"`    // 章节信息
}

// TeachesCreateRequest 表示创建教学关系的请求
type TeachesCreateRequest struct {
	InstructorID string `json:"instructor_id"`
	SectionID    string `json:"section_id"`
}