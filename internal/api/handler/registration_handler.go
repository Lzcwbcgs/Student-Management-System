package handler

import (
	"encoding/json"
	"net/http"

	"github.com/yourusername/student-management-system/internal/service"
	"github.com/yourusername/student-management-system/pkg/utils"
)

type RegistrationHandler struct {
	enrollmentService service.EnrollmentService
}

func NewRegistrationHandler(enrollmentService service.EnrollmentService) *RegistrationHandler {
	return &RegistrationHandler{
		enrollmentService: enrollmentService,
	}
}

// RegisterCourse 学生选课
func (h *RegistrationHandler) RegisterCourse(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.WriteErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var registrationData struct {
		SectionID string `json:"section_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&registrationData); err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	studentID := r.Context().Value("userID").(string)

	err := h.enrollmentService.RegisterForCourse(studentID, registrationData.SectionID)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to register course")
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, map[string]string{"message": "Course registered successfully"})
}

// DropCourse 学生退课
func (h *RegistrationHandler) DropCourse(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		utils.WriteErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var dropData struct {
		SectionID string `json:"section_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&dropData); err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	studentID := r.Context().Value("userID").(string)

	err := h.enrollmentService.DropCourse(studentID, dropData.SectionID)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to drop course")
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, map[string]string{"message": "Course dropped successfully"})
}

// GetRegisteredCourses 获取学生已选课程
func (h *RegistrationHandler) GetRegisteredCourses(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.WriteErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	studentID := r.Context().Value("userID").(string)

	courses, err := h.enrollmentService.GetRegisteredCourses(studentID)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to get registered courses")
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, courses)
}
