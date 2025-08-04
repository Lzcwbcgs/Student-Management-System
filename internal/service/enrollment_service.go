package service

import (
	"errors"
	"fmt"

	"github.com/yourusername/student-management-system/internal/model"
	"github.com/yourusername/student-management-system/internal/repository"
)

// EnrollmentService 定义选课服务接口
type EnrollmentService interface {
	RegisterForCourse(studentID string, sectionID string) error
	DropCourse(studentID string, sectionID string) error
	GetRegisteredCourses(studentID string) ([]*model.Course, error)
	CheckPrerequisites(studentID string, courseID string) (bool, error)
	CheckTimeConflict(studentID string, sectionID string) (bool, error)
	CheckCapacity(sectionID string) (bool, error)
}

// DefaultEnrollmentService 实现EnrollmentService接口
type DefaultEnrollmentService struct {
	takesRepo    repository.TakesRepository
	studentRepo  repository.StudentRepository
	sectionRepo  repository.SectionRepository
	courseRepo   repository.CourseRepository
	prereqRepo   repository.PrereqRepository
	timeSlotRepo repository.TimeSlotRepository
	teachesRepo  repository.TeachesRepository
}

// NewEnrollmentService 创建选课服务实例
func NewEnrollmentService(takesRepo repository.TakesRepository, studentRepo repository.StudentRepository, sectionRepo repository.SectionRepository, courseRepo repository.CourseRepository, prereqRepo *repository.SQLPrereqRepository, timeSlotRepo repository.TimeSlotRepository, teachesRepo repository.TeachesRepository) EnrollmentService {
	return &DefaultEnrollmentService{
		takesRepo:    takesRepo,
		studentRepo:  studentRepo,
		sectionRepo:  sectionRepo,
		courseRepo:   courseRepo,
		prereqRepo:   prereqRepo,
		timeSlotRepo: timeSlotRepo,
		teachesRepo:  teachesRepo,
	}
}

// RegisterForCourse 学生选课
func (s *DefaultEnrollmentService) RegisterForCourse(studentID string, sectionID string) error {
	// 检查学生是否存在
	_, err := s.studentRepo.GetByID(studentID)
	if err != nil {
		return fmt.Errorf("student not found: %w", err)
	}

	// 检查课程段是否存在
	section, err := s.sectionRepo.FindByID(sectionID)
	if err != nil {
		return fmt.Errorf("section not found: %w", err)
	}

	// 检查是否已经选过这门课
	existingTakes, err := s.takesRepo.FindByStudentAndSection(studentID, sectionID)
	if err == nil && existingTakes != nil {
		return errors.New("already registered for this course")
	}

	// 检查先修课程要求
	satisfied, err := s.prereqRepo.CheckPrereqsSatisfied(studentID, section.CourseID)
	if err != nil {
		return fmt.Errorf("error checking prerequisites: %w", err)
	}
	if !satisfied {
		return errors.New("prerequisites not satisfied")
	}

	// 检查时间冲突
	hasConflict, err := s.takesRepo.CheckTimeConflict(studentID, sectionID)
	if err != nil {
		return fmt.Errorf("error checking time conflict: %w", err)
	}
	if hasConflict {
		return errors.New("time conflict with existing courses")
	}

	// 检查容量
	enrollmentCount, err := s.sectionRepo.GetEnrollmentCount(sectionID)
	if err != nil {
		return fmt.Errorf("error getting enrollment count: %w", err)
	}

	classroom, err := s.sectionRepo.GetSectionClassroom(sectionID)
	if err != nil {
		return fmt.Errorf("error getting classroom info: %w", err)
	}

	if enrollmentCount >= classroom.Capacity {
		return errors.New("section is full")
	}

	// 创建选课记录 - 使用正确的字段名
	takes := &model.Takes{
		StudentID: studentID,
		CourseID:  section.CourseID,
		SectionID: sectionID,
		Semester:  section.Semester,
		Year:      section.Year,
		Grade:     "", // 新选课没有成绩
	}

	return s.takesRepo.Create(takes)
}

// DropCourse 学生退课
func (s *DefaultEnrollmentService) DropCourse(studentID string, sectionID string) error {
	// 检查选课记录是否存在
	_, err := s.takesRepo.FindByStudentAndSection(studentID, sectionID)
	if err != nil {
		return fmt.Errorf("enrollment not found: %w", err)
	}

	return s.takesRepo.Delete(studentID, sectionID)
}

// GetRegisteredCourses 获取学生已选课程
func (s *DefaultEnrollmentService) GetRegisteredCourses(studentID string) ([]*model.Course, error) {
	takes, err := s.takesRepo.FindByStudentID(studentID)
	if err != nil {
		return nil, err
	}

	courses := make([]*model.Course, 0, len(takes))
	for _, take := range takes {
		if take.Course != nil {
			courses = append(courses, take.Course)
		}
	}
	return courses, nil
}

// CheckPrerequisites 检查先修课程
func (s *DefaultEnrollmentService) CheckPrerequisites(studentID string, courseID string) (bool, error) {
	return s.prereqRepo.CheckPrereqsSatisfied(studentID, courseID)
}

// CheckTimeConflict 检查时间冲突
func (s *DefaultEnrollmentService) CheckTimeConflict(studentID string, sectionID string) (bool, error) {
	return s.takesRepo.CheckTimeConflict(studentID, sectionID)
}

// CheckCapacity 检查容量
func (s *DefaultEnrollmentService) CheckCapacity(sectionID string) (bool, error) {
	enrollmentCount, err := s.sectionRepo.GetEnrollmentCount(sectionID)
	if err != nil {
		return false, err
	}

	classroom, err := s.sectionRepo.GetSectionClassroom(sectionID)
	if err != nil {
		return false, err
	}

	return enrollmentCount < classroom.Capacity, nil
}
