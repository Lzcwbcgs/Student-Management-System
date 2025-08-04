package handler

import (
	"encoding/json"
	"github.com/yourusername/student-management-system/internal/service"
	"github.com/yourusername/student-management-system/pkg/utils"
	"net/http"
)

type StudentHandler struct {
	studentService service.StudentService
}

func NewStudentHandler(studentService service.StudentService) *StudentHandler {
	return &StudentHandler{
		studentService: studentService,
	}
}

// GetProfile 获取学生个人信息
func (h *StudentHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.WriteErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	// 从JWT中获取学生ID
	studentID := r.Context().Value("userID").(string)

	student, err := h.studentService.GetByID(studentID)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to get student profile")
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, student)
}

// UpdateProfile 更新学生个人信息
func (h *StudentHandler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		utils.WriteErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var updateData struct {
		Name string `json:"name"`
	}

	if err := json.NewDecoder(r.Body).Decode(&updateData); err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	studentID := r.Context().Value("userID").(string)

	err := h.studentService.UpdateProfile(studentID, updateData.Name)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to update profile")
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, map[string]string{"message": "Profile updated successfully"})
}

// GetAdvisor 获取学生导师信息
func (h *StudentHandler) GetAdvisor(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.WriteErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	studentID := r.Context().Value("userID").(string)

	advisor, err := h.studentService.GetAdvisor(studentID)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to get advisor")
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, advisor)
}

// GetCourses 获取学生已选课程
func (h *StudentHandler) GetCourses(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.WriteErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	studentID := r.Context().Value("userID").(string)

	courses, err := h.studentService.GetEnrolledCourses(studentID)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to get courses")
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, courses)
}

// GetTranscript 获取学生成绩单
func (h *StudentHandler) GetTranscript(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.WriteErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	studentID := r.Context().Value("userID").(string)

	transcript, err := h.studentService.GetTranscript(studentID)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to get transcript")
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, transcript)
}
