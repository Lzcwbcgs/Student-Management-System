package model

// EnrollmentStatus 表示选课状态
type EnrollmentStatus string

const (
	EnrollmentStatusActive   EnrollmentStatus = "active"   // 正常选课状态
	EnrollmentStatusDropped  EnrollmentStatus = "dropped"  // 已退课
	EnrollmentStatusFailed   EnrollmentStatus = "failed"   // 不及格
	EnrollmentStatusPassed   EnrollmentStatus = "passed"   // 及格
	EnrollmentStatusWaiting  EnrollmentStatus = "waiting"  // 等待审核
	EnrollmentStatusRejected EnrollmentStatus = "rejected" // 被拒绝
)

// Enrollment 表示选课记录
type Enrollment struct {
	ID        string `json:"id"`
	StudentID string `json:"student_id"`
	SectionID string `json:"section_id"`
	Grade     string `json:"grade"`

	// 关联信息
	Student *Student `json:"student,omitempty"`
	Section *Section `json:"section,omitempty"`
}

// EnrollmentQueryParams 表示查询选课记录的参数
type EnrollmentQueryParams struct {
	StudentID string `json:"student_id"`
	SectionID string `json:"section_id"`
	Semester  string `json:"semester"`
	Year      int    `json:"year"`
}
