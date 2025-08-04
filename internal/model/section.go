package model

// Section 表示课程章节实体
type Section struct {
	ID         string `json:"id"`           // 章节ID
	CourseID   string `json:"course_id"`    // 课程ID
	Semester   string `json:"semester"`     // 学期
	Year       int    `json:"year"`         // 年份
	Building   string `json:"building"`     // 教学楼
	RoomNumber string `json:"room_number"`  // 教室号
	TimeSlotID string `json:"time_slot_id"` // 时间段ID
	Enrollment int    `json:"enrollment"`   // 选课人数

	// 关联信息
	Course      *Course      `json:"course,omitempty"`      // 课程信息
	TimeSlot    *TimeSlot    `json:"time_slot,omitempty"`   // 时间段信息
	Classroom   *Classroom   `json:"classroom,omitempty"`   // 教室信息
	Instructors []Instructor `json:"instructors,omitempty"` // 授课教师
}

// SectionDTO 表示课程章节数据传输对象
type SectionDTO struct {
	ID          string          `json:"id"`
	CourseID    string          `json:"course_id"`
	Semester    string          `json:"semester"`
	Year        int             `json:"year"`
	Building    string          `json:"building"`
	RoomNumber  string          `json:"room_number"`
	TimeSlotID  string          `json:"time_slot_id"`
	Enrollment  int             `json:"enrollment"`
	Course      *Course         `json:"course,omitempty"`
	TimeSlot    *TimeSlot       `json:"time_slot,omitempty"`
	Classroom   *Classroom      `json:"classroom,omitempty"`
	Instructors []InstructorDTO `json:"instructors,omitempty"`
}

// SectionCreateRequest 表示创建课程章节的请求
type SectionCreateRequest struct {
	ID         string `json:"id"`
	CourseID   string `json:"course_id"`
	Semester   string `json:"semester"`
	Year       int    `json:"year"`
	Building   string `json:"building"`
	RoomNumber string `json:"room_number"`
	TimeSlotID string `json:"time_slot_id"`
}

// SectionUpdateRequest 表示更新课程章节的请求
type SectionUpdateRequest struct {
	Semester   string `json:"semester"`
	Year       int    `json:"year"`
	Building   string `json:"building"`
	RoomNumber string `json:"room_number"`
	TimeSlotID string `json:"time_slot_id"`
}

// SectionQueryParams 表示查询课程章节的参数
type SectionQueryParams struct {
	CourseID     string `json:"course_id"`
	Semester     string `json:"semester"`
	Year         int    `json:"year"`
	Dept         string `json:"dept"`
	InstructorID string
	TimeSlotID   string
	Building     string `json:"building"`
	RoomNumber   string `json:"room_number"`
}
