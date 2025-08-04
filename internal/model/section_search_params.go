package model

// SectionSearchParams 定义查询课程章节的参数
type SectionSearchParams struct {
	CourseID   string `json:"course_id,omitempty"`
	Semester   string `json:"semester,omitempty"`
	Year       int    `json:"year,omitempty"`
	Building   string `json:"building,omitempty"`
	RoomNumber string `json:"room_number,omitempty"`
	TimeSlotID string `json:"time_slot_id,omitempty"`
}
