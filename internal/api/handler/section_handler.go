package handler

import (
	"net/http"
	"strconv"

	"github.com/yourusername/student-management-system/internal/service"
	"github.com/yourusername/student-management-system/pkg/utils"
)

type SectionHandler struct {
	sectionService service.SectionService
}

func NewSectionHandler(sectionService service.SectionService) *SectionHandler {
	return &SectionHandler{
		sectionService: sectionService,
	}
}

// GetSections 获取课程段列表
func (h *SectionHandler) GetSections(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.WriteErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	// 获取查询参数
	courseID := r.URL.Query().Get("course_id")
	semester := r.URL.Query().Get("semester")
	yearStr := r.URL.Query().Get("year")
	instructorID := r.URL.Query().Get("instructor_id")

	// 转换year为int
	var year int
	if yearStr != "" {
		if y, err := strconv.Atoi(yearStr); err == nil {
			year = y
		}
	}

	sections, err := h.sectionService.GetSections(courseID, semester, year, instructorID)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to get sections")
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, sections)
}
