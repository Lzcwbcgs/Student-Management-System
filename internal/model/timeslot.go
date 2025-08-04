package model

// TimeSlot 时间段模型
type TimeSlot struct {
	ID        string   `json:"id"`         // 时间段ID
	Days      []int    `json:"days"`       // 星期几列表 (1-7 表示周一到周日)
	StartTime string   `json:"start_time"` // 开始时间 (HH:MM 格式)
	EndTime   string   `json:"end_time"`   // 结束时间 (HH:MM 格式)
	StartHr   int      `json:"start_hr"`   // 开始小时
	StartMin  int      `json:"start_min"`  // 开始分钟
	EndHr     int      `json:"end_hr"`     // 结束小时
	EndMin    int      `json:"end_min"`    // 结束分钟
}

// TimeSlotCreateRequest 创建时间段请求
type TimeSlotCreateRequest struct {
	ID        string `json:"id"`
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
	Days      []int  `json:"days"`
}

// TimeSlotUpdateRequest 更新时间段请求
type TimeSlotUpdateRequest struct {
	StartTime string `json:"start_time,omitempty"`
	EndTime   string `json:"end_time,omitempty"`
	Days      []int  `json:"days,omitempty"`
	StartHr   int    `json:"start_hr,omitempty"`
	StartMin  int    `json:"start_min,omitempty"`
	EndHr     int    `json:"end_hr,omitempty"`
	EndMin    int    `json:"end_min,omitempty"`
}
// TimeDuration 表示时间段详情
type TimeDuration struct {
	StartTime string `json:"start_time"` // 开始时间 (HH:MM 格式)
	EndTime   string `json:"end_time"`   // 结束时间 (HH:MM 格式)
	Days      []int  `json:"days"`       // 上课日期列表 (1-7 表示周一到周日)
}

// TimeSlotGroup 表示时间段组（支持多个时间段）
type TimeSlotGroup struct {
	ID        string     `json:"id"`         // 时间段组ID
	TimeSlots []TimeSlot `json:"time_slots"` // 时间段列表
}

// TimeSlotRequest 表示时间段请求的基础字段
type TimeSlotRequest struct {
	StartHr   int    `json:"start_hr,omitempty"`
	StartMin  int    `json:"start_min,omitempty"`
	EndHr     int    `json:"end_hr,omitempty"`
	EndMin    int    `json:"end_min,omitempty"`
	StartTime string `json:"start_time,omitempty"`
	EndTime   string `json:"end_time,omitempty"`
}

// TimeSlotResponse 表示时间段的响应
type TimeSlotResponse struct {
	ID        string   `json:"id"`
	Days      []int    `json:"days"`
	StartTime string   `json:"start_time"`
	EndTime   string   `json:"end_time"`
	StartHr   int      `json:"start_hr"`
	StartMin  int      `json:"start_min"`
	EndHr     int      `json:"end_hr"`
	EndMin    int      `json:"end_min"`
}
