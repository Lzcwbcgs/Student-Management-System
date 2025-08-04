package service

import (
	"errors"
	"fmt"

	"github.com/yourusername/student-management-system/internal/model"
	"github.com/yourusername/student-management-system/internal/repository"
	"github.com/yourusername/student-management-system/pkg/utils"
)

// InstructorService 定义教师服务接口
type InstructorService interface {
	GetInstructorByID(id string) (*model.Instructor, error)
	GetAllInstructors() ([]*model.Instructor, error)
	CreateInstructor(req *model.InstructorCreateRequest) error
	UpdateInstructor(id string, req *model.InstructorUpdateRequest) error
	DeleteInstructor(id string) error
	ChangePassword(id string, req *model.ChangePasswordRequest) error
	GetCurrentTeaching(id string, semester string, year int) ([]*model.Teaches, error)
	AssignGrade(instructorID string, studentID string, sectionID string, grade string) error
	GetAdvisees(instructorID string) ([]*model.Advisor, error)
	AssignTeaching(instructorID string, sectionID string, courseID string, semester string, year int) error
	RemoveTeaching(instructorID string, sectionID string) error
	GetByID(id string) (*model.Instructor, error)
	UpdateProfile(id string, name string) error
	GetTeachingSections(id string) ([]*model.Section, error)
	GetSectionStudents(instructorID string, sectionID string) ([]*model.Student, error)
	UpdateGrade(instructorID string, studentID string, sectionID string, grade string) error
	GetAdviseeInfo(instructorID string, studentID string) (*model.Student, error)
	Authenticate(id string, password string) (string, error)
}

// DefaultInstructorService 实现InstructorService接口
type DefaultInstructorService struct {
	instructorRepo repository.InstructorRepository
	teachesRepo    repository.TeachesRepository
	takesRepo      repository.TakesRepository
	advisorRepo    repository.AdvisorRepository
	sectionRepo    repository.SectionRepository
	studentRepo    repository.StudentRepository
}

// NewInstructorService 创建教师服务实例
func NewInstructorService(instructorRepo repository.InstructorRepository, teachesRepo repository.TeachesRepository, takesRepo repository.TakesRepository, advisorRepo repository.AdvisorRepository, sectionRepo repository.SectionRepository, studentRepo repository.StudentRepository) InstructorService {
	return &DefaultInstructorService{
		instructorRepo: instructorRepo,
		teachesRepo:    teachesRepo,
		takesRepo:      takesRepo,
		advisorRepo:    advisorRepo,
		sectionRepo:    sectionRepo,
		studentRepo:    studentRepo,
	}
}

// GetInstructorByID 根据ID获取教师信息
func (s *DefaultInstructorService) GetInstructorByID(id string) (*model.Instructor, error) {
	return s.instructorRepo.GetByID(id)
}

// GetAllInstructors 获取所有教师
func (s *DefaultInstructorService) GetAllInstructors() ([]*model.Instructor, error) {
	instructors, _, err := s.instructorRepo.List(1, 1000)
	if err != nil {
		return nil, err
	}
	return instructors, nil
}

// CreateInstructor 创建教师
func (s *DefaultInstructorService) CreateInstructor(req *model.InstructorCreateRequest) error {
	// 密码加密
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return fmt.Errorf("error hashing password: %w", err)
	}

	// 创建教师对象
	instructor := &model.Instructor{
		ID:       req.ID,
		Name:     req.Name,
		Dept:     req.Dept,
		Salary:   req.Salary,
		Password: hashedPassword,
	}

	return s.instructorRepo.Create(instructor)
}

// UpdateInstructor 更新教师信息
func (s *DefaultInstructorService) UpdateInstructor(id string, req *model.InstructorUpdateRequest) error {
	// 先查询教师是否存在
	existingInstructor, err := s.instructorRepo.GetByID(id)
	if err != nil {
		return err
	}

	// 更新教师信息
	existingInstructor.Name = req.Name
	existingInstructor.Dept = req.Dept
	existingInstructor.Salary = req.Salary

	return s.instructorRepo.Update(existingInstructor)
}

// DeleteInstructor 删除教师
func (s *DefaultInstructorService) DeleteInstructor(id string) error {
	return s.instructorRepo.Delete(id)
}

// ChangePassword 修改密码
func (s *DefaultInstructorService) ChangePassword(id string, req *model.ChangePasswordRequest) error {
	// 先查询教师是否存在
	instructor, err := s.instructorRepo.GetByID(id)
	if err != nil {
		return err
	}

	// 验证旧密码
	if !utils.CheckPassword(req.OldPassword, instructor.Password) {
		return errors.New("invalid old password")
	}

	// 加密新密码
	hashedPassword, err := utils.HashPassword(req.NewPassword)
	if err != nil {
		return fmt.Errorf("error hashing password: %w", err)
	}

	// 更新密码 - 根据model.Instructor的Salt字段，这里应该传入salt
	return s.instructorRepo.UpdatePassword(id, hashedPassword, instructor.Salt)
}

// GetCurrentTeaching 获取教师当前学期的教学任务
func (s *DefaultInstructorService) GetCurrentTeaching(id string, semester string, year int) ([]*model.Teaches, error) {
	return s.teachesRepo.GetCurrentTeaching(id, semester, year)
}

// AssignGrade 分配成绩
func (s *DefaultInstructorService) AssignGrade(instructorID string, studentID string, sectionID string, grade string) error {
	// 检查教师是否教授这门课
	_, err := s.teachesRepo.FindByInstructorAndSection(instructorID, sectionID)
	if err != nil {
		return fmt.Errorf("instructor not teaching this section: %w", err)
	}

	// 检查学生是否选了这门课
	_, err = s.takesRepo.FindByStudentAndSection(studentID, sectionID)
	if err != nil {
		return fmt.Errorf("student not enrolled in this section: %w", err)
	}

	// 更新成绩
	return s.takesRepo.UpdateGrade(studentID, sectionID, grade)
}

// GetAdvisees 获取导师指导的学生列表
func (s *DefaultInstructorService) GetAdvisees(instructorID string) ([]*model.Advisor, error) {
	return s.advisorRepo.FindByInstructorID(instructorID)
}

// AssignTeaching 分配教学任务
func (s *DefaultInstructorService) AssignTeaching(instructorID string, sectionID string, courseID string, semester string, year int) error {
	// 检查教师是否存在
	_, err := s.instructorRepo.GetByID(instructorID)
	if err != nil {
		return fmt.Errorf("instructor not found: %w", err)
	}

	// 检查课程段是否存在
	_, err = s.sectionRepo.FindByID(sectionID)
	if err != nil {
		return fmt.Errorf("section not found: %w", err)
	}

	// 创建教学关系
	teaches := &model.Teaches{
		InstructorID: instructorID,
		CourseID:     courseID,
		SectionID:    sectionID,
		Semester:     semester,
		Year:         year,
	}

	return s.teachesRepo.Create(teaches)
}

// RemoveTeaching 移除教学任务
func (s *DefaultInstructorService) RemoveTeaching(instructorID string, sectionID string) error {
	// 需要根据teaches表的完整主键来删除
	// 先查询teaches记录获取完整信息
	teaches, err := s.teachesRepo.FindByInstructorAndSection(instructorID, sectionID)
	if err != nil {
		return fmt.Errorf("teaching assignment not found: %w", err)
	}

	return s.teachesRepo.Delete(instructorID, teaches.CourseID, sectionID, teaches.Semester, teaches.Year)
}

// GetByID 根据ID获取教师信息（别名方法）
func (s *DefaultInstructorService) GetByID(id string) (*model.Instructor, error) {
	return s.GetInstructorByID(id)
}

// UpdateProfile 更新教师个人信息
func (s *DefaultInstructorService) UpdateProfile(id string, name string) error {
	instructor, err := s.instructorRepo.GetByID(id)
	if err != nil {
		return err
	}

	instructor.Name = name
	return s.instructorRepo.Update(instructor)
}

// GetTeachingSections 获取教师授课的课程段
func (s *DefaultInstructorService) GetTeachingSections(id string) ([]*model.Section, error) {
	teaches, err := s.teachesRepo.FindByInstructorID(id)
	if err != nil {
		return nil, err
	}

	var sections []*model.Section
	for _, teach := range teaches {
		section, err := s.sectionRepo.FindByID(teach.SectionID)
		if err != nil {
			continue
		}
		sections = append(sections, section)
	}

	return sections, nil
}

// GetSectionStudents 获取课程段的学生名单
func (s *DefaultInstructorService) GetSectionStudents(instructorID string, sectionID string) ([]*model.Student, error) {
	// 检查教师是否教授这门课
	_, err := s.teachesRepo.FindByInstructorAndSection(instructorID, sectionID)
	if err != nil {
		return nil, fmt.Errorf("instructor not teaching this section: %w", err)
	}

	// 获取选课学生
	takes, err := s.takesRepo.FindBySection(sectionID)
	if err != nil {
		return nil, err
	}

	var students []*model.Student
	for _, take := range takes {
		student, err := s.studentRepo.GetByID(take.StudentID)
		if err != nil {
			continue
		}
		students = append(students, student)
	}

	return students, nil
}

// UpdateGrade 更新学生成绩
func (s *DefaultInstructorService) UpdateGrade(instructorID string, studentID string, sectionID string, grade string) error {
	return s.AssignGrade(instructorID, studentID, sectionID, grade)
}

// GetAdviseeInfo 获取指导学生的详细信息
func (s *DefaultInstructorService) GetAdviseeInfo(instructorID string, studentID string) (*model.Student, error) {
	// 检查导师关系 - 使用正确的方法名
	advisors, err := s.advisorRepo.FindByInstructorID(instructorID)
	if err != nil {
		return nil, fmt.Errorf("error getting advisor relationships: %w", err)
	}

	// 检查是否存在导师关系
	found := false
	for _, advisor := range advisors {
		if advisor.StudentID == studentID {
			found = true
			break
		}
	}
	if !found {
		return nil, fmt.Errorf("advisor relationship not found")
	}

	return s.studentRepo.GetByID(studentID)
}

// Authenticate 教师认证
func (s *DefaultInstructorService) Authenticate(id string, password string) (string, error) {
	instructor, err := s.instructorRepo.GetByID(id)
	if err != nil {
		return "", fmt.Errorf("instructor not found: %w", err)
	}

	if !utils.CheckPassword(password, instructor.Password) {
		return "", errors.New("invalid password")
	}

	return instructor.ID, nil
}
