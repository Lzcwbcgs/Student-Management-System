package service

import (
	"errors"
	"fmt"

	"github.com/yourusername/student-management-system/internal/model"
	"github.com/yourusername/student-management-system/internal/repository"
)

// CourseService 定义课程服务接口
type CourseService interface {
	GetCourseByID(id string) (*model.Course, error)
	GetAllCourses() ([]*model.Course, error)
	GetCoursesByDepartment(deptName string) ([]*model.Course, error)
	CreateCourse(req *model.CourseCreateRequest) error
	UpdateCourse(id string, req *model.CourseUpdateRequest) error
	DeleteCourse(id string) error
	GetCourseWithPrereqs(id string) (*model.CourseWithPrereqs, error)
	AddPrerequisite(courseID string, prereqID string) error
	RemovePrerequisite(courseID string, prereqID string) error
	GetCourses(department, title, credits string) ([]*model.Course, error)
}

// DefaultCourseService 实现CourseService接口
type DefaultCourseService struct {
	courseRepo repository.CourseRepository
	prereqRepo repository.PrereqRepository
	deptRepo   repository.DepartmentRepository
}

// NewCourseService 创建课程服务实例
func NewCourseService(courseRepo repository.CourseRepository, prereqRepo repository.PrereqRepository, deptRepo repository.DepartmentRepository) CourseService {
	return &DefaultCourseService{
		courseRepo: courseRepo,
		prereqRepo: prereqRepo,
		deptRepo:   deptRepo,
	}
}

// GetCourseByID 根据ID获取课程信息
func (s *DefaultCourseService) GetCourseByID(id string) (*model.Course, error) {
	return s.courseRepo.FindByID(id)
}

// GetAllCourses 获取所有课程
func (s *DefaultCourseService) GetAllCourses() ([]*model.Course, error) {
	return s.courseRepo.FindAll()
}

// GetCoursesByDepartment 获取指定院系的所有课程
func (s *DefaultCourseService) GetCoursesByDepartment(deptName string) ([]*model.Course, error) {
	// 检查院系是否存在
	_, err := s.deptRepo.FindByID(deptName)
	if err != nil {
		return nil, fmt.Errorf("department not found: %w", err)
	}

	return s.courseRepo.FindByDept(deptName)
}

// CreateCourse 创建课程
func (s *DefaultCourseService) CreateCourse(req *model.CourseCreateRequest) error {
	// 检查院系是否存在
	_, err := s.deptRepo.FindByID(req.Dept)
	if err != nil {
		return fmt.Errorf("department not found: %w", err)
	}

	// 创建课程对象
	course := &model.Course{
		ID:      req.ID,
		Title:   req.Title,
		Dept:    req.Dept,
		Credits: req.Credits, // 使用float64类型，与model.Course一致
		Name:    req.Title,   // 添加Name字段，与model.Course一致
	}

	// 创建课程
	err = s.courseRepo.Create(course)
	if err != nil {
		return err
	}

	// 如果有先修课程，添加先修课程关系
	if len(req.PrereqIDs) > 0 {
		for _, prereqID := range req.PrereqIDs {
			// 检查先修课程是否存在
			_, err := s.courseRepo.FindByID(prereqID)
			if err != nil {
				return fmt.Errorf("prerequisite course not found: %s, %w", prereqID, err)
			}

			// 添加先修课程关系
			prereq := &model.Prereq{
				CourseID: req.ID,
				PrereqID: prereqID,
			}
			err = s.prereqRepo.Create(prereq)
			if err != nil {
				return fmt.Errorf("error adding prerequisite: %w", err)
			}
		}
	}

	return nil
}

// UpdateCourse 更新课程信息
func (s *DefaultCourseService) UpdateCourse(id string, req *model.CourseUpdateRequest) error {
	// 先查询课程是否存在
	course, err := s.courseRepo.FindByID(id)
	if err != nil {
		return fmt.Errorf("course not found: %w", err)
	}
	if course == nil {
		return errors.New("course not found")
	}

	// 检查新院系是否存在
	if req.Dept != "" {
		_, err := s.deptRepo.FindByID(req.Dept)
		if err != nil {
			return fmt.Errorf("department not found: %w", err)
		}
		course.Dept = req.Dept
	}

	// 更新课程信息
	if req.Title != "" {
		course.Title = req.Title
	}
	if req.Credits > 0 {
		course.Credits = req.Credits // 使用float64类型
	}
	course.Name = req.Title // 更新Name字段，与model.Course一致

	return s.courseRepo.Update(course)
}

// DeleteCourse 删除课程
func (s *DefaultCourseService) DeleteCourse(id string) error {
	// 检查课程是否存在
	_, err := s.courseRepo.FindByID(id)
	if err != nil {
		return err
	}

	return s.courseRepo.Delete(id)
}

// GetCourseWithPrereqs 获取课程及其先修课程信息
func (s *DefaultCourseService) GetCourseWithPrereqs(id string) (*model.CourseWithPrereqs, error) {
	course, err := s.courseRepo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("course not found: %w", err)
	}
	if course == nil {
		return nil, errors.New("course not found")
	}

	prereqIDs, err := s.prereqRepo.GetPrereqIDs(id)
	if err != nil {
		return nil, fmt.Errorf("error getting prereq IDs: %w", err)
	}

	courseWithPrereqs := &model.CourseWithPrereqs{}
	courseWithPrereqs.ID = course.ID
	courseWithPrereqs.Title = course.Title
	courseWithPrereqs.Dept = course.Dept
	courseWithPrereqs.Credits = course.Credits
	courseWithPrereqs.Name = course.Name
	courseWithPrereqs.PrereqIDs = prereqIDs
	return courseWithPrereqs, nil
}

// AddPrerequisite 添加先修课程
func (s *DefaultCourseService) AddPrerequisite(courseID string, prereqID string) error {
	// 检查课程是否存在
	_, err := s.courseRepo.FindByID(courseID)
	if err != nil {
		return fmt.Errorf("course not found: %w", err)
	}

	// 检查先修课程是否存在
	_, err = s.courseRepo.FindByID(prereqID)
	if err != nil {
		return fmt.Errorf("prerequisite course not found: %w", err)
	}

	// 检查是否已经是先修课程
	prereqs, err := s.prereqRepo.FindByCourseID(courseID)
	if err != nil {
		return fmt.Errorf("error checking existing prerequisites: %w", err)
	}

	for _, prereq := range prereqs {
		if prereq.PrereqID == prereqID {
			return errors.New("prerequisite relationship already exists")
		}
	}

	// 创建先修课程关系
	prereq := &model.Prereq{
		CourseID: courseID,
		PrereqID: prereqID,
	}

	return s.prereqRepo.Create(prereq)
}

// RemovePrerequisite 移除先修课程
func (s *DefaultCourseService) RemovePrerequisite(courseID string, prereqID string) error {
	// 检查课程是否存在
	_, err := s.courseRepo.FindByID(courseID)
	if err != nil {
		return fmt.Errorf("course not found: %w", err)
	}

	// 检查先修课程是否存在
	_, err = s.courseRepo.FindByID(prereqID)
	if err != nil {
		return fmt.Errorf("prerequisite course not found: %w", err)
	}

	return s.prereqRepo.Delete(courseID, prereqID)
}

// GetCourses 根据参数获取课程列表
func (s *DefaultCourseService) GetCourses(department, title, credits string) ([]*model.Course, error) {
	// 如果有院系参数，使用GetCoursesByDepartment
	if department != "" {
		return s.GetCoursesByDepartment(department)
	}

	// 否则返回所有课程
	return s.GetAllCourses()
}
