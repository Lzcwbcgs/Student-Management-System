package service

import (
	"fmt"

	"github.com/yourusername/student-management-system/internal/model"
	"github.com/yourusername/student-management-system/internal/repository"
)

// DepartmentService 定义院系服务接口
type DepartmentService interface {
	GetDepartmentByID(id string) (*model.Department, error)
	GetAllDepartments() ([]*model.Department, error)
	CreateDepartment(deptName string, building string, budget float64) error
	UpdateDepartment(deptName string, building string, budget float64) error
	DeleteDepartment(deptName string) error
	GetStudentCount(departmentID string) (int, error)
	GetInstructorCount(departmentID string) (int, error)
	GetCourseCount(departmentID string) (int, error)
	GetDepartmentStats() ([]*model.DepartmentStats, error)
}

// DefaultDepartmentService 实现DepartmentService接口
type DefaultDepartmentService struct {
	departmentRepo repository.DepartmentRepository
	studentRepo    repository.StudentRepository
	instructorRepo repository.InstructorRepository
	courseRepo     repository.CourseRepository
}

// NewDepartmentService 创建院系服务实例
func NewDepartmentService(
	departmentRepo repository.DepartmentRepository,
	studentRepo repository.StudentRepository,
	instructorRepo repository.InstructorRepository,
	courseRepo repository.CourseRepository,
) DepartmentService {
	return &DefaultDepartmentService{
		departmentRepo: departmentRepo,
		studentRepo:    studentRepo,
		instructorRepo: instructorRepo,
		courseRepo:     courseRepo,
	}
}

// GetDepartmentByID 根据ID获取院系信息
func (s *DefaultDepartmentService) GetDepartmentByID(id string) (*model.Department, error) {
	return s.departmentRepo.FindByID(id)
}

// GetAllDepartments 获取所有院系
func (s *DefaultDepartmentService) GetAllDepartments() ([]*model.Department, error) {
	return s.departmentRepo.FindAll()
}

// CreateDepartment 创建院系 - 使用正确的字段名
func (s *DefaultDepartmentService) CreateDepartment(deptName string, building string, budget float64) error {
	// 检查院系ID是否已存在
	_, err := s.departmentRepo.FindByID(deptName)
	if err == nil {
		return fmt.Errorf("department with name %s already exists", deptName)
	}

	// 创建院系对象 - 使用正确的字段名
	department := &model.Department{
		DeptName: deptName,  // 使用DeptName字段
		Building: building,  // 使用Building字段
		Budget:   budget,    // 使用Budget字段
	}

	return s.departmentRepo.Create(department)
}

// UpdateDepartment 更新院系信息 - 使用正确的字段名
func (s *DefaultDepartmentService) UpdateDepartment(deptName string, building string, budget float64) error {
	// 先查询院系是否存在
	existingDepartment, err := s.departmentRepo.FindByID(deptName)
	if err != nil {
		return err
	}

	// 更新院系信息 - 使用正确的字段名
	existingDepartment.Building = building
	existingDepartment.Budget = budget

	return s.departmentRepo.Update(existingDepartment)
}

// DeleteDepartment 删除院系
func (s *DefaultDepartmentService) DeleteDepartment(deptName string) error {
	// 检查院系是否存在
	_, err := s.departmentRepo.FindByID(deptName)
	if err != nil {
		return err
	}

	// 检查院系是否有关联的学生
	studentCount, err := s.departmentRepo.GetStudentCount(deptName)
	if err != nil {
		return fmt.Errorf("error checking student count: %w", err)
	}

	if studentCount > 0 {
		return fmt.Errorf("cannot delete department: it has %d associated students", studentCount)
	}

	// 检查院系是否有关联的教师
	instructorCount, err := s.departmentRepo.GetInstructorCount(deptName)
	if err != nil {
		return fmt.Errorf("error checking instructor count: %w", err)
	}

	if instructorCount > 0 {
		return fmt.Errorf("cannot delete department: it has %d associated instructors", instructorCount)
	}

	// 检查院系是否有关联的课程
	courseCount, err := s.departmentRepo.GetCourseCount(deptName)
	if err != nil {
		return fmt.Errorf("error checking course count: %w", err)
	}

	if courseCount > 0 {
		return fmt.Errorf("cannot delete department: it has %d associated courses", courseCount)
	}

	return s.departmentRepo.Delete(deptName)
}

// GetStudentCount 获取院系学生数量
func (s *DefaultDepartmentService) GetStudentCount(departmentID string) (int, error) {
	// 检查院系是否存在
	_, err := s.departmentRepo.FindByID(departmentID)
	if err != nil {
		return 0, err
	}

	return s.departmentRepo.GetStudentCount(departmentID)
}

// GetInstructorCount 获取院系教师数量
func (s *DefaultDepartmentService) GetInstructorCount(departmentID string) (int, error) {
	// 检查院系是否存在
	_, err := s.departmentRepo.FindByID(departmentID)
	if err != nil {
		return 0, err
	}

	return s.departmentRepo.GetInstructorCount(departmentID)
}

// GetCourseCount 获取院系课程数量
func (s *DefaultDepartmentService) GetCourseCount(departmentID string) (int, error) {
	// 检查院系是否存在
	_, err := s.departmentRepo.FindByID(departmentID)
	if err != nil {
		return 0, err
	}

	return s.departmentRepo.GetCourseCount(departmentID)
}

// GetDepartmentStats 获取所有院系统计信息
func (s *DefaultDepartmentService) GetDepartmentStats() ([]*model.DepartmentStats, error) {
	return s.departmentRepo.GetDepartmentStats()
}