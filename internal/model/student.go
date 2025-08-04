package model

// Student 表示学生实体
type Student struct {
	ID       string  `json:"id"`       // 学生ID
	Name     string  `json:"name"`     // 学生姓名
	Dept     string  `json:"dept"`     // 所属院系
	TotCred  float64 `json:"tot_cred"` // 总学分
	Password string  `json:"password,omitempty"` // 密码（哈希后）
	Salt     string  `json:"salt,omitempty"`     // 密码盐值
}

// StudentDTO 表示学生数据传输对象（不包含敏感信息）
type StudentDTO struct {
	ID      string  `json:"id"`
	Name    string  `json:"name"`
	Dept    string  `json:"dept"`
	TotCred float64 `json:"tot_cred"`
}

// ToDTO 将Student转换为StudentDTO
func (s *Student) ToDTO() *StudentDTO {
	return &StudentDTO{
		ID:      s.ID,
		Name:    s.Name,
		Dept:    s.Dept,
		TotCred: s.TotCred,
	}
}

// StudentProfileUpdateRequest 表示学生个人信息更新请求
type StudentProfileUpdateRequest struct {
	Name     string `json:"name"`
	Password string `json:"password,omitempty"`
}

// StudentLoginRequest 表示学生登录请求
type StudentLoginRequest struct {
	ID       string `json:"id"`
	Password string `json:"password"`
}