package service

import (
	"errors"
	"fmt"

	"github.com/yourusername/student-management-system/internal/model"
	"github.com/yourusername/student-management-system/internal/repository"
)

// SectionService 定义课程章节服务接口
type SectionService interface {
	GetSectionByID(id string) (*model.Section, error)
	GetAllSections() ([]*model.Section, error)
	GetSectionsByCourseID(courseID string) ([]*model.Section, error)
	GetSectionsByParams(params *model.SectionQueryParams) ([]*model.Section, error)
	CreateSection(req *model.SectionCreateRequest) error
	UpdateSection(id string, req *model.SectionUpdateRequest) error
	DeleteSection(id string) error
	GetSectionWithDetails(id string) (*model.Section, error)
	GetSections(courseID, semester string, year int, instructorID string) ([]*model.Section, error)
}

// DefaultSectionService 实现SectionService接口
type DefaultSectionService struct {
	sectionRepo   repository.SectionRepository
	courseRepo    repository.CourseRepository
	classroomRepo repository.ClassroomRepository
	timeSlotRepo  repository.TimeSlotRepository
}

// NewSectionService 创建新的SectionService实例
func NewSectionService(sectionRepo repository.SectionRepository, courseRepo repository.CourseRepository, classroomRepo repository.ClassroomRepository, timeSlotRepo repository.TimeSlotRepository) SectionService {
	return &DefaultSectionService{
		sectionRepo:   sectionRepo,
		courseRepo:    courseRepo,
		classroomRepo: classroomRepo,
		timeSlotRepo:  timeSlotRepo,
	}
}

// GetSectionByID 根据ID获取课程章节
func (s *DefaultSectionService) GetSectionByID(id string) (*model.Section, error) {
	return s.sectionRepo.FindByID(id)
}

// GetAllSections 获取所有课程章节
func (s *DefaultSectionService) GetAllSections() ([]*model.Section, error) {
	return s.sectionRepo.FindAll()
}

// GetSectionsByCourseID 根据课程ID获取章节列表
func (s *DefaultSectionService) GetSectionsByCourseID(courseID string) ([]*model.Section, error) {
	return s.sectionRepo.FindByCourseID(courseID)
}

// GetSectionsByParams 根据参数获取章节列表
func (s *DefaultSectionService) GetSectionsByParams(params *model.SectionQueryParams) ([]*model.Section, error) {
	return s.sectionRepo.FindByParams(params)
}

// GetSections 根据参数获取章节列表
func (s *DefaultSectionService) GetSections(courseID, semester string, year int, instructorID string) ([]*model.Section, error) {
	params := &model.SectionQueryParams{
		CourseID:     courseID,
		Semester:     semester,
		Year:         year,
		InstructorID: instructorID,
	}
	return s.sectionRepo.FindByParams(params)
}

// CreateSection 创建课程章节
func (s *DefaultSectionService) CreateSection(req *model.SectionCreateRequest) error {
	// 验证课程是否存在
	course, err := s.courseRepo.FindByID(req.CourseID)
	if err != nil {
		return fmt.Errorf("course not found: %w", err)
	}
	if course == nil {
		return errors.New("course not found")
	}

	// 验证教室是否存在
	classroom, err := s.classroomRepo.FindByBuildingAndRoom(req.Building, req.RoomNumber)
	if err != nil {
		return fmt.Errorf("classroom not found: %w", err)
	}
	if classroom == nil {
		return errors.New("classroom not found")
	}

	// 验证时间段是否存在
	timeSlot, err := s.timeSlotRepo.FindByID(req.TimeSlotID)
	if err != nil {
		return fmt.Errorf("time slot not found: %w", err)
	}
	if timeSlot == nil {
		return errors.New("time slot not found")
	}

	// 创建Section - 使用正确的字段名
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

// UpdateSection 更新课程章节
func (s *DefaultSectionService) UpdateSection(id string, req *model.SectionUpdateRequest) error {
	section, err := s.sectionRepo.FindByID(id)
	if err != nil {
		return fmt.Errorf("section not found: %w", err)
	}
	if section == nil {
		return errors.New("section not found")
	}

	// 更新字段
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
		// 验证时间段是否存在
		timeSlot, err := s.timeSlotRepo.FindByID(req.TimeSlotID)
		if err != nil {
			return fmt.Errorf("time slot not found: %w", err)
		}
		if timeSlot == nil {
			return errors.New("time slot not found")
		}
		section.TimeSlotID = req.TimeSlotID
	}

	return s.sectionRepo.Update(section)
}

// DeleteSection 删除课程章节
func (s *DefaultSectionService) DeleteSection(id string) error {
	return s.sectionRepo.Delete(id)
}

// GetSectionWithDetails 获取章节详细信息
func (s *DefaultSectionService) GetSectionWithDetails(id string) (*model.Section, error) {
	return s.sectionRepo.FindWithDetails(id)
}
