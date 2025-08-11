package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/yourusername/student-management-system/internal/model"
	"github.com/yourusername/student-management-system/internal/service"
	"github.com/yourusername/student-management-system/pkg/utils"
)

type AuthHandler struct {
	studentService    service.StudentService
	instructorService service.InstructorService
}

func NewAuthHandler(studentService service.StudentService, instructorService service.InstructorService) *AuthHandler {
	return &AuthHandler{
		studentService:    studentService,
		instructorService: instructorService,
	}
}

// Register 用户注册
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.WriteErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var registerData struct {
		ID         string `json:"id"`
		Name       string `json:"name"`
		Password   string `json:"password"`
		Type       string `json:"type"` // "student" or "instructor"
		Department string `json:"department"`
	}

	if err := json.NewDecoder(r.Body).Decode(&registerData); err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	var err error
	switch registerData.Type {
	case "student":
		// 创建学生注册请求
		req := &model.StudentCreateRequest{
			ID:       registerData.ID,
			Name:     registerData.Name,
			Password: registerData.Password,
			Dept:     registerData.Department,
		}
		err = h.studentService.CreateStudent(req)
	case "instructor":
		// 创建教师注册请求
		req := &model.InstructorCreateRequest{
			ID:       registerData.ID,
			Name:     registerData.Name,
			Password: registerData.Password,
			Dept:     registerData.Department,
		}
		err = h.instructorService.CreateInstructor(req)
	default:
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Invalid user type")
		return
	}

	if err != nil {
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to create user")
		return
	}

	response := map[string]string{
		"message": "User registered successfully",
	}

	utils.WriteJSONResponse(w, http.StatusOK, response)
}

// Login 用户登录
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.WriteErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var loginData struct {
		UserID   string `json:"user_id"`
		Password string `json:"password"`
		Role     string `json:"role"` // "student", "instructor", "admin"
	}

	if err := json.NewDecoder(r.Body).Decode(&loginData); err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// 验证用户身份
	var userID string
	var err error

	switch loginData.Role {
	case "student":
		userID, err = h.studentService.Authenticate(loginData.UserID, loginData.Password)
	case "instructor":
		userID, err = h.instructorService.Authenticate(loginData.UserID, loginData.Password)
	case "admin":
		// 管理员验证逻辑（这里简化处理）
		if loginData.UserID == "admin" && loginData.Password == "admin123" {
			userID = "admin"
		} else {
			err = utils.NewError("Invalid admin credentials")
		}
	default:
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Invalid role")
		return
	}

	if err != nil {
		log.Printf("Login failed for user %s (role: %s): %v", loginData.UserID, loginData.Role, err)
		utils.WriteErrorResponse(w, http.StatusUnauthorized, "Invalid credentials")
		return
	}

	// 生成JWT token
	token, err := utils.GenerateJWT(userID, loginData.UserID, loginData.Role) // 使用 loginData.UserID 替代 userData.Username
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to generate token")
		return
	}

	response := map[string]interface{}{
		"token":   token,
		"user_id": userID,
		"role":    loginData.Role,
	}

	utils.WriteJSONResponse(w, http.StatusOK, response)
}
