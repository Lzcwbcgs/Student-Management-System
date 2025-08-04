package service

import (
	"fmt"

	"github.com/yourusername/student-management-system/internal/model"
	"github.com/yourusername/student-management-system/internal/repository"
)

// TimeSlotService 定义时间段服务接口
type TimeSlotService interface {
	GetTimeSlotByID(id string) (*model.TimeSlot, error)
	GetAllTimeSlots() ([]*model.TimeSlot, error)
	GetTimeSlotsByDayOfWeek(dayOfWeek int) ([]*model.TimeSlot, error)
	GetTimeSlotsByTimeRange(startTime, endTime string) ([]*model.TimeSlot, error)
	CreateTimeSlot(req *model.TimeSlotCreateRequest) error
	UpdateTimeSlot(id string, req *model.TimeSlotUpdateRequest) error
	DeleteTimeSlot(id string) error
	GetTimeSlotUsage(id string, semester string, year int) ([]*model.Section, error)
}

// DefaultTimeSlotService 实现TimeSlotService接口
type DefaultTimeSlotService struct {
	timeslotRepo repository.TimeSlotRepository
	sectionRepo  repository.SectionRepository
}

// NewTimeSlotService 创建时间段服务实例
func NewTimeSlotService(
	timeslotRepo repository.TimeSlotRepository,
	sectionRepo repository.SectionRepository,
) TimeSlotService {
	return &DefaultTimeSlotService{
		timeslotRepo: timeslotRepo,
		sectionRepo:  sectionRepo,
	}
}

// GetTimeSlotByID 根据ID获取时间段信息
func (s *DefaultTimeSlotService) GetTimeSlotByID(id string) (*model.TimeSlot, error) {
	return s.timeslotRepo.FindByID(id)
}

// GetAllTimeSlots 获取所有时间段
func (s *DefaultTimeSlotService) GetAllTimeSlots() ([]*model.TimeSlot, error) {
	return s.timeslotRepo.FindAll()
}

// GetTimeSlotsByDayOfWeek 获取指定星期几的所有时间段
func (s *DefaultTimeSlotService) GetTimeSlotsByDayOfWeek(dayOfWeek int) ([]*model.TimeSlot, error) {
	if dayOfWeek < 1 || dayOfWeek > 7 {
		return nil, fmt.Errorf("invalid day of week: %d, must be between 1 and 7", dayOfWeek)
	}

	return s.timeslotRepo.FindByDayOfWeek(dayOfWeek)
}

// GetTimeSlotsByTimeRange 获取指定时间范围内的所有时间段
func (s *DefaultTimeSlotService) GetTimeSlotsByTimeRange(startTime, endTime string) ([]*model.TimeSlot, error) {
	return s.timeslotRepo.FindByTimeRange(startTime, endTime)
}

// CreateTimeSlot 创建时间段
func (s *DefaultTimeSlotService) CreateTimeSlot(req *model.TimeSlotCreateRequest) error {
	// 验证时间段数据
	if len(req.Days) == 0 {
		return fmt.Errorf("at least one day must be specified")
	}

	for _, day := range req.Days {
		if day < 1 || day > 7 {
			return fmt.Errorf("invalid day of week: %d, must be between 1 and 7", day)
		}
	}

	// 创建时间段对象
	timeSlot := &model.TimeSlot{
		ID:        req.ID,
		StartTime: req.StartTime,
		EndTime:   req.EndTime,
		Days:      req.Days,
	}

	return s.timeslotRepo.Create(timeSlot)
}

// UpdateTimeSlot 更新时间段信息
func (s *DefaultTimeSlotService) UpdateTimeSlot(id string, req *model.TimeSlotUpdateRequest) error {
	// 先查询时间段是否存在
	existingTimeSlot, err := s.timeslotRepo.FindByID(id)
	if err != nil {
		return err
	}

	// 更新时间段信息
	if req.StartTime != "" {
		existingTimeSlot.StartTime = req.StartTime
	}

	if req.EndTime != "" {
		existingTimeSlot.EndTime = req.EndTime
	}

	if req.Days != nil && len(req.Days) > 0 {
		// 验证星期几数据
		for _, day := range req.Days {
			if day < 1 || day > 7 {
				return fmt.Errorf("invalid day of week: %d, must be between 1 and 7", day)
			}
		}
		existingTimeSlot.Days = req.Days
	}

	// 检查时间段是否被课程章节使用
	if req.StartTime != "" || req.EndTime != "" || (req.Days != nil && len(req.Days) > 0) {
		params := &model.SectionQueryParams{
			TimeSlotID: id,
		}

		sections, err := s.sectionRepo.FindByParams(params)
		if err != nil {
			return fmt.Errorf("error checking time slot usage: %w", err)
		}

		if len(sections) > 0 {
			return fmt.Errorf("cannot update time slot: it is being used by %d course sections", len(sections))
		}
	}

	return s.timeslotRepo.Update(existingTimeSlot)
}

// DeleteTimeSlot 删除时间段
func (s *DefaultTimeSlotService) DeleteTimeSlot(id string) error {
	// 检查时间段是否存在
	_, err := s.timeslotRepo.FindByID(id)
	if err != nil {
		return err
	}

	// 检查时间段是否被课程章节使用
	params := &model.SectionQueryParams{
		TimeSlotID: id,
	}

	sections, err := s.sectionRepo.FindByParams(params)
	if err != nil {
		return fmt.Errorf("error checking time slot usage: %w", err)
	}

	if len(sections) > 0 {
		return fmt.Errorf("cannot delete time slot: it is being used by %d course sections", len(sections))
	}

	return s.timeslotRepo.Delete(id)
}

// GetTimeSlotUsage 获取时间段使用情况
func (s *DefaultTimeSlotService) GetTimeSlotUsage(id string, semester string, year int) ([]*model.Section, error) {
	// 检查时间段是否存在
	_, err := s.timeslotRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	// 查询使用该时间段的所有课程章节
	params := &model.SectionQueryParams{
		TimeSlotID: id,
		Semester:   semester,
		Year:       year,
	}

	return s.sectionRepo.FindByParams(params)
}
