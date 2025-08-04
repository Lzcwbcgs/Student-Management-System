package model

// Instructor 表示教师实体
type Instructor struct {
	ID       string  `json:"id"`       // 教师ID
	Name     string  `json:"name"`     // 教师姓名
	Dept     string  `json:"dept"`     // 所属院系
	Salary   float64 `json:"salary"`   // 薪水
	Password string  `json:"password,omitempty"` // 密码（哈希后）
	Salt     string  `json:"salt,omitempty"`     // 密码盐值
}

// InstructorDTO 表示教师数据传输对象（不包含敏感信息）
type InstructorDTO struct {
	ID     string  `json:"id"`
	Name   string  `json:"name"`
	Dept   string  `json:"dept"`
	Salary float64 `json:"salary,omitempty"` // 可能对普通用户隐藏
}

// ToDTO 将Instructor转换为InstructorDTO
func (i *Instructor) ToDTO() *InstructorDTO {
	return &InstructorDTO{
		ID:     i.ID,
		Name:   i.Name,
		Dept:   i.Dept,
		Salary: i.Salary,
	}
}

// InstructorProfileUpdateRequest 表示教师个人信息更新请求
type InstructorProfileUpdateRequest struct {
	Name     string `json:"name"`
	Password string `json:"password,omitempty"`
}

// InstructorLoginRequest 表示教师登录请求
type InstructorLoginRequest struct {
	ID       string `json:"id"`
	Password string `json:"password"`
}