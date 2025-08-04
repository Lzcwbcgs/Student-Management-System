package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/yourusername/student-management-system/internal/model"
	"github.com/yourusername/student-management-system/internal/service"
	"github.com/yourusername/student-management-system/pkg/utils"
)

type AdminHandler struct {
	adminService service.AdminService
}

func NewAdminHandler(adminService *service.DefaultAdminService) *AdminHandler {
	return &AdminHandler{
		adminService: adminService,
	}
}

// GetStudents 获取学生列表
func (h *AdminHandler) GetStudents(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.WriteErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	students, err := h.adminService.GetAllStudents()
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to get students")
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, students)
}

// CreateStudent 创建学生
func (h *AdminHandler) CreateStudent(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.WriteErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var studentData struct {
		ID   string `json:"id"`
		Name string `json:"name"`
		Dept string `json:"dept_name"`
	}

	if err := json.NewDecoder(r.Body).Decode(&studentData); err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	err := h.adminService.CreateStudent(studentData.ID, studentData.Name, studentData.Dept)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.WriteJSONResponse(w, http.StatusCreated, map[string]string{"message": "Student created successfully"})
}

// UpdateStudent 更新学生信息
func (h *AdminHandler) UpdateStudent(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		utils.WriteErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var studentData struct {
		ID   string `json:"id"`
		Name string `json:"name"`
		Dept string `json:"dept_name"`
	}

	if err := json.NewDecoder(r.Body).Decode(&studentData); err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	err := h.adminService.UpdateStudent(studentData.ID, studentData.Name, studentData.Dept)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, map[string]string{"message": "Student updated successfully"})
}

// DeleteStudent 删除学生
func (h *AdminHandler) DeleteStudent(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		utils.WriteErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	studentID := r.URL.Query().Get("id")
	if studentID == "" {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Student ID is required")
		return
	}

	err := h.adminService.DeleteStudent(studentID)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, map[string]string{"message": "Student deleted successfully"})
}

// GetInstructors 获取教师列表
func (h *AdminHandler) GetInstructors(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.WriteErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	instructors, err := h.adminService.GetAllInstructors()
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to get instructors")
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, instructors)
}

// CreateInstructor 创建教师
func (h *AdminHandler) CreateInstructor(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.WriteErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var instructorData struct {
		ID     string  `json:"id"`
		Name   string  `json:"name"`
		Dept   string  `json:"dept_name"`
		Salary float64 `json:"salary"`
	}

	if err := json.NewDecoder(r.Body).Decode(&instructorData); err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	err := h.adminService.CreateInstructor(instructorData.ID, instructorData.Name, instructorData.Dept, instructorData.Salary)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.WriteJSONResponse(w, http.StatusCreated, map[string]string{"message": "Instructor created successfully"})
}

// UpdateInstructor 更新教师信息
func (h *AdminHandler) UpdateInstructor(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		utils.WriteErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var instructorData struct {
		ID     string  `json:"id"`
		Name   string  `json:"name"`
		Dept   string  `json:"dept_name"`
		Salary float64 `json:"salary"`
	}

	if err := json.NewDecoder(r.Body).Decode(&instructorData); err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	err := h.adminService.UpdateInstructor(instructorData.ID, instructorData.Name, instructorData.Dept, instructorData.Salary)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, map[string]string{"message": "Instructor updated successfully"})
}

// DeleteInstructor 删除教师
func (h *AdminHandler) DeleteInstructor(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		utils.WriteErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	instructorID := r.URL.Query().Get("id")
	if instructorID == "" {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Instructor ID is required")
		return
	}

	err := h.adminService.DeleteInstructor(instructorID)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, map[string]string{"message": "Instructor deleted successfully"})
}

// GetDepartments 获取院系列表
func (h *AdminHandler) GetDepartments(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.WriteErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	departments, err := h.adminService.GetAllDepartments()
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to get departments")
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, departments)
}

// CreateDepartment 创建院系
func (h *AdminHandler) CreateDepartment(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.WriteErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var deptData struct {
		Name     string  `json:"dept_name"`
		Building string  `json:"building"`
		Budget   float64 `json:"budget"`
	}

	if err := json.NewDecoder(r.Body).Decode(&deptData); err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	err := h.adminService.CreateDepartment(deptData.Name, deptData.Building, deptData.Budget)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.WriteJSONResponse(w, http.StatusCreated, map[string]string{"message": "Department created successfully"})
}

// UpdateDepartment 更新院系信息
func (h *AdminHandler) UpdateDepartment(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		utils.WriteErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var deptData struct {
		Name     string  `json:"dept_name"`
		Building string  `json:"building"`
		Budget   float64 `json:"budget"`
	}

	if err := json.NewDecoder(r.Body).Decode(&deptData); err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	err := h.adminService.UpdateDepartment(deptData.Name, deptData.Building, deptData.Budget)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, map[string]string{"message": "Department updated successfully"})
}

// DeleteDepartment 删除院系
func (h *AdminHandler) DeleteDepartment(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		utils.WriteErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	deptName := r.URL.Query().Get("name")
	if deptName == "" {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Department name is required")
		return
	}

	err := h.adminService.DeleteDepartment(deptName)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, map[string]string{"message": "Department deleted successfully"})
}

// GetCourses 获取课程列表
func (h *AdminHandler) GetCourses(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.WriteErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	courses, err := h.adminService.GetAllCourses()
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to get courses")
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, courses)
}

// CreateCourse 创建课程
func (h *AdminHandler) CreateCourse(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.WriteErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var courseData struct {
		ID      string `json:"course_id"`
		Title   string `json:"title"`
		Dept    string `json:"dept_name"`
		Credits int    `json:"credits"`
	}

	if err := json.NewDecoder(r.Body).Decode(&courseData); err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	err := h.adminService.CreateCourse(courseData.ID, courseData.Title, courseData.Dept, courseData.Credits)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.WriteJSONResponse(w, http.StatusCreated, map[string]string{"message": "Course created successfully"})
}

// UpdateCourse 更新课程信息
func (h *AdminHandler) UpdateCourse(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		utils.WriteErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var courseData struct {
		ID      string `json:"course_id"`
		Title   string `json:"title"`
		Dept    string `json:"dept_name"`
		Credits int    `json:"credits"`
	}

	if err := json.NewDecoder(r.Body).Decode(&courseData); err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	err := h.adminService.UpdateCourse(courseData.ID, courseData.Title, courseData.Dept, courseData.Credits)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, map[string]string{"message": "Course updated successfully"})
}

// DeleteCourse 删除课程
func (h *AdminHandler) DeleteCourse(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		utils.WriteErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	courseID := r.URL.Query().Get("id")
	if courseID == "" {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Course ID is required")
		return
	}

	err := h.adminService.DeleteCourse(courseID)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, map[string]string{"message": "Course deleted successfully"})
}

// GetPrereqs 获取先修课程列表
func (h *AdminHandler) GetPrereqs(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.WriteErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	courseID := r.URL.Query().Get("course_id")
	if courseID == "" {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Course ID is required")
		return
	}

	prereqs, err := h.adminService.GetPrereqs(courseID)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to get prerequisites")
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, prereqs)
}

// CreatePrereq 创建先修课程关系
func (h *AdminHandler) CreatePrereq(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.WriteErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var prereqData struct {
		CourseID string `json:"course_id"`
		PrereqID string `json:"prereq_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&prereqData); err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	err := h.adminService.CreatePrereq(prereqData.CourseID, prereqData.PrereqID)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.WriteJSONResponse(w, http.StatusCreated, map[string]string{"message": "Prerequisite created successfully"})
}

// DeletePrereq 删除先修课程关系
func (h *AdminHandler) DeletePrereq(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		utils.WriteErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	courseID := r.URL.Query().Get("course_id")
	prereqID := r.URL.Query().Get("prereq_id")

	if courseID == "" || prereqID == "" {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Course ID and Prerequisite ID are required")
		return
	}

	err := h.adminService.DeletePrereq(courseID, prereqID)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, map[string]string{"message": "Prerequisite deleted successfully"})
}

// GetClassrooms 获取教室列表
func (h *AdminHandler) GetClassrooms(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.WriteErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	classrooms, err := h.adminService.GetAllClassrooms()
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to get classrooms")
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, classrooms)
}

// CreateClassroom 创建教室
func (h *AdminHandler) CreateClassroom(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.WriteErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var classroomData struct {
		Building string `json:"building"`
		Room     string `json:"room_number"`
		Capacity int    `json:"capacity"`
	}

	if err := json.NewDecoder(r.Body).Decode(&classroomData); err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	err := h.adminService.CreateClassroom(classroomData.Building, classroomData.Room, classroomData.Capacity)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.WriteJSONResponse(w, http.StatusCreated, map[string]string{"message": "Classroom created successfully"})
}

// UpdateClassroom 更新教室信息
func (h *AdminHandler) UpdateClassroom(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		utils.WriteErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var classroomData struct {
		Building string `json:"building"`
		Room     string `json:"room_number"`
		Capacity int    `json:"capacity"`
	}

	if err := json.NewDecoder(r.Body).Decode(&classroomData); err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	err := h.adminService.UpdateClassroom(classroomData.Building, classroomData.Room, classroomData.Capacity)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, map[string]string{"message": "Classroom updated successfully"})
}

// DeleteClassroom 删除教室
func (h *AdminHandler) DeleteClassroom(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		utils.WriteErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	building := r.URL.Query().Get("building")
	room := r.URL.Query().Get("room")

	if building == "" || room == "" {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Building and Room are required")
		return
	}

	err := h.adminService.DeleteClassroom(building, room)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, map[string]string{"message": "Classroom deleted successfully"})
}

// GetSections 获取课程段列表
func (h *AdminHandler) GetSections(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.WriteErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	sections, err := h.adminService.GetAllSections()
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to get sections")
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, sections)
}

// CreateSection 创建课程段
func (h *AdminHandler) CreateSection(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.WriteErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var sectionData struct {
		CourseID   string `json:"course_id"`
		SecID      string `json:"sec_id"`
		Semester   string `json:"semester"`
		Year       int    `json:"year"`
		Building   string `json:"building"`
		Room       string `json:"room_number"`
		TimeSlotID string `json:"time_slot_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&sectionData); err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	req := &model.SectionCreateRequest{
		ID:         sectionData.SecID,
		CourseID:   sectionData.CourseID,
		Semester:   sectionData.Semester,
		Year:       sectionData.Year,
		Building:   sectionData.Building,
		RoomNumber: sectionData.Room,
		TimeSlotID: sectionData.TimeSlotID,
	}

	err := h.adminService.CreateSection(req)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.WriteJSONResponse(w, http.StatusCreated, map[string]string{"message": "Section created successfully"})
}

// UpdateSection 更新课程段信息
func (h *AdminHandler) UpdateSection(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		utils.WriteErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var sectionData struct {
		CourseID   string `json:"course_id"`
		SecID      string `json:"sec_id"`
		Semester   string `json:"semester"`
		Year       int    `json:"year"`
		Building   string `json:"building"`
		Room       string `json:"room_number"`
		TimeSlotID string `json:"time_slot_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&sectionData); err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	req := &model.SectionUpdateRequest{
		Semester:   sectionData.Semester,
		Year:       sectionData.Year,
		Building:   sectionData.Building,
		RoomNumber: sectionData.Room,
		TimeSlotID: sectionData.TimeSlotID,
	}

	// 使用课程ID和章节ID作为标识符
	sectionID := sectionData.SecID
	err := h.adminService.UpdateSection(sectionID, req)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, map[string]string{"message": "Section updated successfully"})
}

// DeleteSection 删除课程段
func (h *AdminHandler) DeleteSection(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		utils.WriteErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	courseID := r.URL.Query().Get("course_id")
	secID := r.URL.Query().Get("sec_id")
	semester := r.URL.Query().Get("semester")
	yearStr := r.URL.Query().Get("year")

	if courseID == "" || secID == "" || semester == "" || yearStr == "" {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Course ID, Section ID, Semester, and Year are required")
		return
	}

	year, err := strconv.Atoi(yearStr)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Invalid year format")
		return
	}

	err = h.adminService.DeleteSection(courseID, secID, semester, year)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, map[string]string{"message": "Section deleted successfully"})
}

// GetTeaches 获取教学关系列表
func (h *AdminHandler) GetTeaches(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.WriteErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	teaches, err := h.adminService.GetAllTeaches()
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to get teaches")
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, teaches)
}

// CreateTeaches 创建教学关系
func (h *AdminHandler) CreateTeaches(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.WriteErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var teachesData struct {
		InstructorID string `json:"instructor_id"`
		CourseID     string `json:"course_id"`
		SecID        string `json:"sec_id"`
		Semester     string `json:"semester"`
		Year         int    `json:"year"`
	}

	if err := json.NewDecoder(r.Body).Decode(&teachesData); err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	err := h.adminService.CreateTeaches(teachesData.InstructorID, teachesData.CourseID, teachesData.SecID, teachesData.Semester, teachesData.Year)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.WriteJSONResponse(w, http.StatusCreated, map[string]string{"message": "Teaching assignment created successfully"})
}

// DeleteTeaches 删除教学关系
func (h *AdminHandler) DeleteTeaches(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		utils.WriteErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	instructorID := r.URL.Query().Get("instructor_id")
	courseID := r.URL.Query().Get("course_id")
	secID := r.URL.Query().Get("sec_id")
	semester := r.URL.Query().Get("semester")
	yearStr := r.URL.Query().Get("year")

	if instructorID == "" || courseID == "" || secID == "" || semester == "" || yearStr == "" {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "All parameters are required")
		return
	}

	year, err := strconv.Atoi(yearStr)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Invalid year format")
		return
	}

	err = h.adminService.DeleteTeaches(instructorID, courseID, secID, semester, year)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, map[string]string{"message": "Teaching assignment deleted successfully"})
}

// GetAdvisors 获取导师关系列表
func (h *AdminHandler) GetAdvisors(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.WriteErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	advisors, err := h.adminService.GetAllAdvisors()
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to get advisors")
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, advisors)
}

// CreateAdvisor 创建导师关系
func (h *AdminHandler) CreateAdvisor(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.WriteErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var advisorData struct {
		StudentID    string `json:"student_id"`
		InstructorID string `json:"instructor_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&advisorData); err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	err := h.adminService.CreateAdvisor(advisorData.StudentID, advisorData.InstructorID)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.WriteJSONResponse(w, http.StatusCreated, map[string]string{"message": "Advisor relationship created successfully"})
}

// DeleteAdvisor 删除导师关系
func (h *AdminHandler) DeleteAdvisor(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		utils.WriteErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	studentID := r.URL.Query().Get("student_id")
	instructorID := r.URL.Query().Get("instructor_id")

	if studentID == "" || instructorID == "" {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Student ID and Instructor ID are required")
		return
	}

	err := h.adminService.DeleteAdvisor(studentID, instructorID)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, map[string]string{"message": "Advisor relationship deleted successfully"})
}

// GetStats 获取系统统计信息
func (h *AdminHandler) GetStats(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.WriteErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	stats, err := h.adminService.GetSystemStats()
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to get system stats")
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, stats)
}
