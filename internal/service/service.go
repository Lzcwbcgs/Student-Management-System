package service

// AuthService 定义认证相关的业务逻辑
type AuthService interface {
	// Login 用户登录
	Login(userType string, id string, password string) (string, error)
	// ValidateToken 验证token
	ValidateToken(token string) (string, error)
}
