package model

// ChangePasswordRequest 表示修改密码请求
type ChangePasswordRequest struct {
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}

// LoginRequest 表示登录请求
type LoginRequest struct {
	UserID   string `json:"user_id"`
	Password string `json:"password"`
	Role     string `json:"role"`
}