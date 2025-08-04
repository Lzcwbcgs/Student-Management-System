package handler

import (
	"encoding/json"
	"github.com/yourusername/student-management-system/internal/service"
	"github.com/yourusername/student-management-system/pkg/utils"
	"net/http"
)

type InstructorHandler struct {
	instructorService service.InstructorService
}

func NewInstructorHandler(instructorService service.InstructorService) *InstructorHandler {
	return &InstructorHandler{
		instructorService: instructorService,
	}
}

// GetProfile 获取教师个人信息
func (h *InstructorHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.WriteErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	instructorID := r.Context().Value("userID").(string)

	instructor, err := h.instructorService.GetByID(instructorID)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to get instructor profile")
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, instructor)
}

// UpdateProfile 更新教师个人信息
func (h *InstructorHandler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
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

	instructorID := r.Context().Value("userID").(string)

	err := h.instructorService.UpdateProfile(instructorID, updateData.Name)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to update profile")
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, map[string]string{"message": "Profile updated successfully"})
}

// GetSections 获取教师授课列表
func (h *InstructorHandler) GetSections(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.WriteErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	instructorID := r.Context().Value("userID").(string)

	sections, err := h.instructorService.GetTeachingSections(instructorID)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to get sections")
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, sections)
}

// GetSectionStudents 获取课程学生名单
func (h *InstructorHandler) GetSectionStudents(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.WriteErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	sectionID := r.URL.Query().Get("section_id")
	if sectionID == "" {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Section ID is required")
		return
	}

	instructorID := r.Context().Value("userID").(string)

	students, err := h.instructorService.GetSectionStudents(instructorID, sectionID)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to get section students")
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, students)
}

// UpdateGrade 更新学生成绩
func (h *InstructorHandler) UpdateGrade(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		utils.WriteErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var gradeData struct {
		StudentID string `json:"student_id"`
		SectionID string `json:"section_id"`
		Grade     string `json:"grade"`
	}

	if err := json.NewDecoder(r.Body).Decode(&gradeData); err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	instructorID := r.Context().Value("userID").(string)

	err := h.instructorService.UpdateGrade(instructorID, gradeData.StudentID, gradeData.SectionID, gradeData.Grade)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to update grade")
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, map[string]string{"message": "Grade updated successfully"})
}

// GetAdvisees 获取导师指导的学生列表
func (h *InstructorHandler) GetAdvisees(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.WriteErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	instructorID := r.Context().Value("userID").(string)

	advisees, err := h.instructorService.GetAdvisees(instructorID)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to get advisees")
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, advisees)
}

// GetAdviseeInfo 获取指导学生的详细信息
func (h *InstructorHandler) GetAdviseeInfo(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.WriteErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	studentID := r.URL.Query().Get("student_id")
	if studentID == "" {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Student ID is required")
		return
	}

	instructorID := r.Context().Value("userID").(string)

	info, err := h.instructorService.GetAdviseeInfo(instructorID, studentID)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to get advisee info")
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, info)
}
