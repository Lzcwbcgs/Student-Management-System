package model

// Takes 表示学生选课记录
type Takes struct {
	ID        string  `json:"id"`        // 记录ID
	StudentID string  `json:"student_id"` // 学生ID
	CourseID  string  `json:"course_id"`  // 课程ID
	SectionID string  `json:"section_id"` // 章节ID
	Semester  string  `json:"semester"`   // 学期
	Year      int     `json:"year"`       // 年份
	Grade     string  `json:"grade"`      // 成绩
	
	// 关联信息
	Student   *Student `json:"student,omitempty"`   // 学生信息
	Course    *Course  `json:"course,omitempty"`    // 课程信息
	Section   *Section `json:"section,omitempty"`   // 章节信息
}

// TakesCreateRequest 表示创建选课记录的请求
type TakesCreateRequest struct {
	StudentID string `json:"student_id"`
	SectionID string `json:"section_id"`
}

// GradeUpdateRequest 表示更新成绩的请求
type GradeUpdateRequest struct {
	StudentID string `json:"student_id"`
	SectionID string `json:"section_id"`
	Grade     string `json:"grade"`
}

// Transcript 表示成绩单
type Transcript struct {
	Student    StudentDTO `json:"student"`
	Courses    []CourseGrade `json:"courses"`
	TotalCred  float64 `json:"total_cred"`
	GPA        float64 `json:"gpa"`
}

// CourseGrade 表示课程成绩
type CourseGrade struct {
	CourseID   string  `json:"course_id"`
	Title      string  `json:"title"`
	Semester   string  `json:"semester"`
	Year       int     `json:"year"`
	Credits    float64 `json:"credits"`
	Grade      string  `json:"grade"`
	GradePoint float64 `json:"grade_point"`
}