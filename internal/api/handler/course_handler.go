package handler

import (
	"net/http"

	"github.com/yourusername/student-management-system/internal/service"
	"github.com/yourusername/student-management-system/pkg/utils"
	"github.com/yourusername/student-management-system/internal/model"
)

type CourseHandler struct {
	courseService service.CourseService
}

func NewCourseHandler(courseService service.CourseService) *CourseHandler {
	return &CourseHandler{
		courseService: courseService,
	}
}

// GetCourses 获取课程列表
func (h *CourseHandler) GetCourses(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.WriteErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	// 获取查询参数
	department := r.URL.Query().Get("department")
	title := r.URL.Query().Get("title")

	// 根据参数获取课程列表
	var courses []*model.Course
	var err error
	
	if department != "" || title != "" {
		// 如果有查询参数，使用GetCoursesByDepartment或类似方法
		courses, err = h.courseService.GetAllCourses()
	} else {
		courses, err = h.courseService.GetAllCourses()
	}
	
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to get courses")
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, courses)
}