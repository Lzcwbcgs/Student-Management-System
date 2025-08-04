package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/yourusername/student-management-system/pkg/utils"
)

type AuthMiddleware struct{}

func NewAuthMiddleware() *AuthMiddleware {
	return &AuthMiddleware{}
}

// Authenticate 验证JWT token
func (m *AuthMiddleware) Authenticate(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			utils.WriteErrorResponse(w, http.StatusUnauthorized, "Authorization header required")
			return
		}

		// 检查Bearer token格式
		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			utils.WriteErrorResponse(w, http.StatusUnauthorized, "Invalid authorization header format")
			return
		}

		token := tokenParts[1]

		// 验证JWT token
		claims, err := utils.ValidateJWT(token)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusUnauthorized, "Invalid token")
			return
		}

		// 将用户信息添加到请求上下文
		ctx := context.WithValue(r.Context(), "userID", claims.UserID)
		ctx = context.WithValue(ctx, "role", claims.Role)

		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

// AuthorizeStudent 授权学生访问
func (m *AuthMiddleware) AuthorizeStudent(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		role := r.Context().Value("role").(string)
		if role != "student" && role != "admin" {
			utils.WriteErrorResponse(w, http.StatusForbidden, "Access denied")
			return
		}
		next.ServeHTTP(w, r)
	}
}

// AuthorizeInstructor 授权教师访问
func (m *AuthMiddleware) AuthorizeInstructor(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		role := r.Context().Value("role").(string)
		if role != "instructor" && role != "admin" {
			utils.WriteErrorResponse(w, http.StatusForbidden, "Access denied")
			return
		}
		next.ServeHTTP(w, r)
	}
}

// AuthorizeAdmin 授权管理员访问
func (m *AuthMiddleware) AuthorizeAdmin(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		role := r.Context().Value("role").(string)
		if role != "admin" {
			utils.WriteErrorResponse(w, http.StatusForbidden, "Access denied")
			return
		}
		next.ServeHTTP(w, r)
	}
} 