package model

// Classroom 表示教室实体
type Classroom struct {
	Building   string `json:"building"`    // 教学楼
	RoomNumber string `json:"room_number"` // 教室号
	Capacity   int    `json:"capacity"`    // 容量
}

// ClassroomCreateRequest 表示创建教室的请求
type ClassroomCreateRequest struct {
	Building   string `json:"building"`
	RoomNumber string `json:"room_number"`
	Capacity   int    `json:"capacity"`
	Facilities interface{}
}

// ClassroomUpdateRequest 表示更新教室的请求
type ClassroomUpdateRequest struct {
	Capacity int `json:"capacity"`
}
