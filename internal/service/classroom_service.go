package service

import (
	"fmt"

	"github.com/yourusername/student-management-system/internal/model"
	"github.com/yourusername/student-management-system/internal/repository"
)

// ClassroomService 定义教室服务接口
type ClassroomService interface {
	GetClassroomByID(building, roomNumber string) (*model.Classroom, error)
	GetAllClassrooms() ([]*model.Classroom, error)
	GetClassroomsByBuilding(building string) ([]*model.Classroom, error)
	CreateClassroom(building, roomNumber string, capacity int) error
	UpdateClassroom(building, roomNumber string, capacity int) error
	DeleteClassroom(building, roomNumber string) error
	GetAvailableClassrooms(capacity int, semester string, year int, timeSlotID string) ([]*model.Classroom, error)
	GetClassroomUsage(building, roomNumber string, semester string, year int) ([]*model.Section, error)
}

// DefaultClassroomService 实现ClassroomService接口
type DefaultClassroomService struct {
	classroomRepo repository.ClassroomRepository
	sectionRepo   repository.SectionRepository
}

// NewClassroomService 创建教室服务实例
func NewClassroomService(
	classroomRepo repository.ClassroomRepository,
	sectionRepo repository.SectionRepository,
) ClassroomService {
	return &DefaultClassroomService{
		classroomRepo: classroomRepo,
		sectionRepo:   sectionRepo,
	}
}

// GetClassroomByID 根据ID获取教室信息
func (s *DefaultClassroomService) GetClassroomByID(building, roomNumber string) (*model.Classroom, error) {
	return s.classroomRepo.FindByBuildingAndRoom(building, roomNumber)
}

// GetAllClassrooms 获取所有教室
func (s *DefaultClassroomService) GetAllClassrooms() ([]*model.Classroom, error) {
	return s.classroomRepo.FindAll()
}

// GetClassroomsByBuilding 获取指定教学楼的所有教室
func (s *DefaultClassroomService) GetClassroomsByBuilding(building string) ([]*model.Classroom, error) {
	return s.classroomRepo.FindByBuilding(building)
}

// CreateClassroom 创建教室 - 使用正确的字段名
func (s *DefaultClassroomService) CreateClassroom(building, roomNumber string, capacity int) error {
	// 检查教室是否已存在
	_, err := s.classroomRepo.FindByBuildingAndRoom(building, roomNumber)
	if err == nil {
		return fmt.Errorf("classroom already exists: %s-%s", building, roomNumber)
	}

	// 创建教室对象 - 使用正确的字段名
	classroom := &model.Classroom{
		Building:   building,   // 使用Building字段
		RoomNumber: roomNumber, // 使用RoomNumber字段
		Capacity:   capacity,   // 使用Capacity字段
	}

	return s.classroomRepo.Create(classroom)
}

// UpdateClassroom 更新教室信息 - 使用正确的字段名
func (s *DefaultClassroomService) UpdateClassroom(building, roomNumber string, capacity int) error {
	// 先查询教室是否存在
	existingClassroom, err := s.classroomRepo.FindByBuildingAndRoom(building, roomNumber)
	if err != nil {
		return err
	}

	// 更新教室信息 - 使用正确的字段名
	existingClassroom.Capacity = capacity

	return s.classroomRepo.Update(existingClassroom)
}

// DeleteClassroom 删除教室
func (s *DefaultClassroomService) DeleteClassroom(building, roomNumber string) error {
	// 检查教室是否存在
	_, err := s.classroomRepo.FindByBuildingAndRoom(building, roomNumber)
	if err != nil {
		return err
	}

	// 检查教室是否被课程章节使用
	params := &model.SectionQueryParams{
		Building:   building,
		RoomNumber: roomNumber,
	}

	sections, err := s.sectionRepo.FindByParams(params)
	if err != nil {
		return fmt.Errorf("error checking classroom usage: %w", err)
	}

	if len(sections) > 0 {
		return fmt.Errorf("cannot delete classroom: it is being used by %d course sections", len(sections))
	}

	return s.classroomRepo.Delete(building, roomNumber)
}

// GetAvailableClassrooms 获取可用教室
func (s *DefaultClassroomService) GetAvailableClassrooms(capacity int, semester string, year int, timeSlotID string) ([]*model.Classroom, error) {
	return s.classroomRepo.FindAvailable(capacity, semester, year, timeSlotID)
}

// GetClassroomUsage 获取教室使用情况
func (s *DefaultClassroomService) GetClassroomUsage(building, roomNumber string, semester string, year int) ([]*model.Section, error) {
	// 检查教室是否存在
	_, err := s.classroomRepo.FindByBuildingAndRoom(building, roomNumber)
	if err != nil {
		return nil, err
	}

	// 查询使用该教室的所有课程章节
	params := &model.SectionQueryParams{
		Building:   building,
		RoomNumber: roomNumber,
		Semester:   semester,
		Year:       year,
	}

	return s.sectionRepo.FindByParams(params)
}