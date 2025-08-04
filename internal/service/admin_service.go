package service

import (
	"errors"
	"fmt"

	"github.com/yourusername/student-management-system/internal/model"
	"github.com/yourusername/student-management-system/internal/repository"
)

// AdminService 定义管理员服务接口
type AdminService interface {
	// 学生管理
	GetAllStudents() ([]*model.Student, error)
	CreateStudent(id string, name string, dept string) error
	UpdateStudent(id string, name string, dept string) error
	DeleteStudent(id string) error

	// 教师管理
	GetAllInstructors() ([]*model.Instructor, error)
	CreateInstructor(id string, name string, dept string, salary float64) error
	UpdateInstructor(id string, name string, dept string, salary float64) error
	DeleteInstructor(id string) error

	// 课程管理
	GetAllCourses() ([]*model.Course, error)
	CreateCourse(id string, title string, dept string, credits int) error
	UpdateCourse(id string, title string, dept string, credits int) error
	DeleteCourse(id string) error

	// 章节管理
	GetAllSections() ([]*model.Section, error)
	CreateSection(req *model.SectionCreateRequest) error
	UpdateSection(id string, req *model.SectionUpdateRequest) error
	DeleteSection(id string, secID string, semester string, year int) error

	// 系部管理
	GetAllDepartments() ([]*model.Department, error)
	CreateDepartment(deptName string, building string, budget float64) error
	UpdateDepartment(deptName string, building string, budget float64) error
	DeleteDepartment(deptName string) error

	// 教室管理
	GetAllClassrooms() ([]*model.Classroom, error)
	CreateClassroom(building string, roomNumber string, capacity int) error
	UpdateClassroom(building string, roomNumber string, capacity int) error
	DeleteClassroom(building string, roomNumber string) error

	// 先修课程管理
	GetAllPrereqs() ([]*model.Prereq, error)
	CreatePrereq(courseID string, prereqID string) error
	DeletePrereq(courseID string, prereqID string) error

	// 教学安排管理
	GetAllTeaches() ([]*model.Teaches, error)
	CreateTeaches(instructorID string, courseID string, sectionID string, semester string, year int) error
	DeleteTeaches(instructorID string, courseID string, sectionID string, semester string, year int) error

	// 导师关系管理
	GetAllAdvisors() ([]*model.Advisor, error)
	CreateAdvisor(studentID string, instructorID string) error
	DeleteAdvisor(studentID string, instructorID string) error

	// 统计信息
	GetStats() (*model.AdminStats, error)
	GetSystemStats() (*model.SystemStats, error)
	GetPrereqs(id string) ([]*model.Prereq, error)
}

// DefaultAdminService 实现AdminService接口
type DefaultAdminService struct {
	studentRepo    repository.StudentRepository
	instructorRepo repository.InstructorRepository
	courseRepo     repository.CourseRepository
	sectionRepo    repository.SectionRepository
	departmentRepo repository.DepartmentRepository
	classroomRepo  repository.ClassroomRepository
	timeSlotRepo   repository.TimeSlotRepository
	teachesRepo    repository.TeachesRepository
	advisorRepo    repository.AdvisorRepository
	prereqRepo     repository.PrereqRepository
}

func (s *DefaultAdminService) DeleteSection(id string, secID string, semester string, year int) error {
	//TODO implement me
	panic("implement me")
}

func (s *DefaultAdminService) GetSystemStats() (*model.SystemStats, error) {
	// TODO: 实现系统统计信息的收集
	stats := &model.SystemStats{
		ActiveUsers:     0, // 需要实现实际的统计逻辑
		ServerUptime:    0,
		DatabaseSize:    0,
		MemoryUsage:     0,
		CPUUsage:        0,
		TotalOperations: 0,
		ErrorRate:       0,
	}
	return stats, nil
}

func (s *DefaultAdminService) GetPrereqs(id string) ([]*model.Prereq, error) {
	return s.prereqRepo.FindByCourseID(id)
}

// NewAdminService 创建新的AdminService实例
func NewAdminService(studentRepo repository.StudentRepository, instructorRepo repository.InstructorRepository, courseRepo repository.CourseRepository, sectionRepo repository.SectionRepository, departmentRepo repository.DepartmentRepository, classroomRepo repository.ClassroomRepository, timeSlotRepo repository.TimeSlotRepository, teachesRepo repository.TeachesRepository, advisorRepo repository.AdvisorRepository, prereqRepo repository.PrereqRepository) *DefaultAdminService {
	return &DefaultAdminService{
		studentRepo:    studentRepo,
		instructorRepo: instructorRepo,
		courseRepo:     courseRepo,
		sectionRepo:    sectionRepo,
		departmentRepo: departmentRepo,
		classroomRepo:  classroomRepo,
		timeSlotRepo:   timeSlotRepo,
		teachesRepo:    teachesRepo,
		advisorRepo:    advisorRepo,
		prereqRepo:     prereqRepo,
	}
}

// GetAllStudents 获取所有学生
func (s *DefaultAdminService) GetAllStudents() ([]*model.Student, error) {
	students, _, err := s.studentRepo.List(1, 1000) // 设置一个较大的页大小来获取所有学生
	if err != nil {
		return nil, err
	}
	return students, nil
}

// CreateStudent 创建学生
func (s *DefaultAdminService) CreateStudent(id string, name string, dept string) error {
	student := &model.Student{
		ID:      id,
		Name:    name,
		Dept:    dept,
		TotCred: 0.0, // 使用float64类型
	}
	return s.studentRepo.Create(student)
}

// UpdateStudent 更新学生信息
func (s *DefaultAdminService) UpdateStudent(id string, name string, dept string) error {
	student, err := s.studentRepo.GetByID(id)
	if err != nil {
		return fmt.Errorf("student not found: %w", err)
	}
	if student == nil {
		return errors.New("student not found")
	}

	student.Name = name
	student.Dept = dept
	return s.studentRepo.Update(student)
}

// DeleteStudent 删除学生
func (s *DefaultAdminService) DeleteStudent(id string) error {
	return s.studentRepo.Delete(id)
}

// GetAllInstructors 获取所有教师
func (s *DefaultAdminService) GetAllInstructors() ([]*model.Instructor, error) {
	instructors, _, err := s.instructorRepo.List(1, 1000) // 设置一个较大的页大小来获取所有教师
	if err != nil {
		return nil, err
	}
	return instructors, nil
}

// CreateInstructor 创建教师
func (s *DefaultAdminService) CreateInstructor(id string, name string, dept string, salary float64) error {
	instructor := &model.Instructor{
		ID:     id,
		Name:   name,
		Dept:   dept,
		Salary: salary,
	}
	return s.instructorRepo.Create(instructor)
}

// UpdateInstructor 更新教师信息
func (s *DefaultAdminService) UpdateInstructor(id string, name string, dept string, salary float64) error {
	instructor, err := s.instructorRepo.GetByID(id)
	if err != nil {
		return fmt.Errorf("instructor not found: %w", err)
	}
	if instructor == nil {
		return errors.New("instructor not found")
	}

	instructor.Name = name
	instructor.Dept = dept
	instructor.Salary = salary
	return s.instructorRepo.Update(instructor)
}

// DeleteInstructor 删除教师
func (s *DefaultAdminService) DeleteInstructor(id string) error {
	return s.instructorRepo.Delete(id)
}

// GetAllCourses 获取所有课程
func (s *DefaultAdminService) GetAllCourses() ([]*model.Course, error) {
	return s.courseRepo.FindAll()
}

// CreateCourse 创建课程
func (s *DefaultAdminService) CreateCourse(id string, title string, dept string, credits int) error {
	course := &model.Course{
		ID:      id,
		Title:   title,
		Dept:    dept,
		Credits: float64(credits), // 转换为float64类型
		Name:    title,            // 添加Name字段
	}
	return s.courseRepo.Create(course)
}

// UpdateCourse 更新课程
func (s *DefaultAdminService) UpdateCourse(id string, title string, dept string, credits int) error {
	course, err := s.courseRepo.FindByID(id)
	if err != nil {
		return fmt.Errorf("course not found: %w", err)
	}
	if course == nil {
		return errors.New("course not found")
	}

	course.Title = title
	course.Dept = dept
	course.Credits = float64(credits) // 转换为float64类型
	course.Name = title               // 更新Name字段
	return s.courseRepo.Update(course)
}

// DeleteCourse 删除课程
func (s *DefaultAdminService) DeleteCourse(id string) error {
	return s.courseRepo.Delete(id)
}

// GetAllSections 获取所有章节
func (s *DefaultAdminService) GetAllSections() ([]*model.Section, error) {
	return s.sectionRepo.FindAll()
}

// CreateSection 创建章节
func (s *DefaultAdminService) CreateSection(req *model.SectionCreateRequest) error {
	section := &model.Section{
		ID:         req.ID,
		CourseID:   req.CourseID,
		Semester:   req.Semester,
		Year:       req.Year,
		Building:   req.Building,
		RoomNumber: req.RoomNumber,
		TimeSlotID: req.TimeSlotID,
		Enrollment: 0,
	}
	return s.sectionRepo.Create(section)
}

// UpdateSection 更新章节
func (s *DefaultAdminService) UpdateSection(id string, req *model.SectionUpdateRequest) error {
	section, err := s.sectionRepo.FindByID(id)
	if err != nil {
		return fmt.Errorf("section not found: %w", err)
	}
	if section == nil {
		return errors.New("section not found")
	}

	if req.Semester != "" {
		section.Semester = req.Semester
	}
	if req.Year != 0 {
		section.Year = req.Year
	}
	if req.Building != "" {
		section.Building = req.Building
	}
	if req.RoomNumber != "" {
		section.RoomNumber = req.RoomNumber
	}
	if req.TimeSlotID != "" {
		section.TimeSlotID = req.TimeSlotID
	}

	return s.sectionRepo.Update(section)
}

// GetAllDepartments 获取所有系部
func (s *DefaultAdminService) GetAllDepartments() ([]*model.Department, error) {
	return s.departmentRepo.FindAll()
}

// CreateDepartment 创建系部
func (s *DefaultAdminService) CreateDepartment(deptName string, building string, budget float64) error {
	dept := &model.Department{
		DeptName: deptName, // 使用正确的字段名
		Building: building,
		Budget:   budget,
	}
	return s.departmentRepo.Create(dept)
}

// UpdateDepartment 更新系部
func (s *DefaultAdminService) UpdateDepartment(deptName string, building string, budget float64) error {
	dept, err := s.departmentRepo.FindByID(deptName)
	if err != nil {
		return fmt.Errorf("department not found: %w", err)
	}
	if dept == nil {
		return errors.New("department not found")
	}

	dept.Building = building
	dept.Budget = budget
	return s.departmentRepo.Update(dept)
}

// DeleteDepartment 删除系部
func (s *DefaultAdminService) DeleteDepartment(deptName string) error {
	return s.departmentRepo.Delete(deptName)
}

// GetAllClassrooms 获取所有教室
func (s *DefaultAdminService) GetAllClassrooms() ([]*model.Classroom, error) {
	return s.classroomRepo.FindAll()
}

// CreateClassroom 创建教室
func (s *DefaultAdminService) CreateClassroom(building string, roomNumber string, capacity int) error {
	classroom := &model.Classroom{
		Building:   building,
		RoomNumber: roomNumber,
		Capacity:   capacity,
	}
	return s.classroomRepo.Create(classroom)
}

// UpdateClassroom 更新教室
func (s *DefaultAdminService) UpdateClassroom(building string, roomNumber string, capacity int) error {
	classroom, err := s.classroomRepo.FindByBuildingAndRoom(building, roomNumber)
	if err != nil {
		return fmt.Errorf("classroom not found: %w", err)
	}
	if classroom == nil {
		return errors.New("classroom not found")
	}

	classroom.Capacity = capacity
	return s.classroomRepo.Update(classroom)
}

// DeleteClassroom 删除教室
func (s *DefaultAdminService) DeleteClassroom(building string, roomNumber string) error {
	return s.classroomRepo.Delete(building, roomNumber)
}

// GetAllPrereqs 获取所有先修课程
func (s *DefaultAdminService) GetAllPrereqs() ([]*model.Prereq, error) {
	return s.prereqRepo.FindAll()
}

// CreatePrereq 创建先修课程关系
func (s *DefaultAdminService) CreatePrereq(courseID string, prereqID string) error {
	prereq := &model.Prereq{
		CourseID: courseID,
		PrereqID: prereqID,
	}
	return s.prereqRepo.Create(prereq)
}

// DeletePrereq 删除先修课程关系
func (s *DefaultAdminService) DeletePrereq(courseID string, prereqID string) error {
	return s.prereqRepo.Delete(courseID, prereqID)
}

// GetAllTeaches 获取所有教学安排
func (s *DefaultAdminService) GetAllTeaches() ([]*model.Teaches, error) {
	return s.teachesRepo.FindAll()
}

// CreateTeaches 创建教学安排
func (s *DefaultAdminService) CreateTeaches(instructorID string, courseID string, sectionID string, semester string, year int) error {
	teaches := &model.Teaches{
		InstructorID: instructorID,
		CourseID:     courseID,
		SectionID:    sectionID,
		Semester:     semester,
		Year:         year,
	}
	return s.teachesRepo.Create(teaches)
}

// DeleteTeaches 删除教学安排
func (s *DefaultAdminService) DeleteTeaches(instructorID string, courseID string, sectionID string, semester string, year int) error {
	return s.teachesRepo.Delete(instructorID, courseID, sectionID, semester, year)
}

// GetAllAdvisors 获取所有导师关系
func (s *DefaultAdminService) GetAllAdvisors() ([]*model.Advisor, error) {
	return s.advisorRepo.FindAll()
}

// CreateAdvisor 创建导师关系
func (s *DefaultAdminService) CreateAdvisor(studentID string, instructorID string) error {
	advisor := &model.Advisor{
		StudentID:    studentID,
		InstructorID: instructorID,
	}
	return s.advisorRepo.Create(advisor)
}

// DeleteAdvisor 删除导师关系
func (s *DefaultAdminService) DeleteAdvisor(studentID string, instructorID string) error {
	return s.advisorRepo.Delete(studentID, instructorID)
}

// GetStats 获取统计信息
func (s *DefaultAdminService) GetStats() (*model.AdminStats, error) {
	// 这里需要实现统计逻辑
	// 暂时返回空结构
	return &model.AdminStats{}, nil
}
