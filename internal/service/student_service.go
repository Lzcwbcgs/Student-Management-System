package service

import (
	"errors"
	"fmt"

	"github.com/yourusername/student-management-system/internal/model"
	"github.com/yourusername/student-management-system/internal/repository"
	"github.com/yourusername/student-management-system/pkg/utils"
)

// StudentService 定义学生服务接口
type StudentService interface {
	GetStudentByID(id string) (*model.Student, error)
	GetAllStudents() ([]*model.Student, error)
	CreateStudent(req *model.StudentCreateRequest) error
	UpdateStudent(id string, req *model.StudentUpdateRequest) error
	DeleteStudent(id string) error
	ChangePassword(id string, req *model.ChangePasswordRequest) error
	GetStudentTranscript(id string) (*model.Transcript, error)
	GetCurrentCourses(id string, semester string, year int) ([]*model.Takes, error)
	RegisterForCourse(studentID string, sectionID string, courseID string, semester string, year int) error
	DropCourse(studentID string, sectionID string) error
	GetByID(id string) (*model.Student, error)
	UpdateProfile(id string, name string) error
	GetAdvisor(id string) (*model.Advisor, error)
	GetEnrolledCourses(id string) ([]*model.Takes, error)
	GetTranscript(id string) (*model.Transcript, error)
	Authenticate(id string, password string) (string, error)
	Create(student *model.Student) error
}

// DefaultStudentService 实现StudentService接口
type DefaultStudentService struct {
	studentRepo repository.StudentRepository
	takesRepo   repository.TakesRepository
	prereqRepo  repository.PrereqRepository
	sectionRepo repository.SectionRepository
	advisorRepo repository.AdvisorRepository
}

// NewStudentService 创建学生服务实例
func NewStudentService(studentRepo repository.StudentRepository, takesRepo repository.TakesRepository, prereqRepo repository.PrereqRepository, sectionRepo repository.SectionRepository, advisorRepo repository.AdvisorRepository) StudentService {
	return &DefaultStudentService{
		studentRepo: studentRepo,
		takesRepo:   takesRepo,
		prereqRepo:  prereqRepo,
		sectionRepo: sectionRepo,
		advisorRepo: advisorRepo,
	}
}

// GetStudentByID 根据ID获取学生信息
func (s *DefaultStudentService) GetStudentByID(id string) (*model.Student, error) {
	return s.studentRepo.GetByID(id)
}

// GetAllStudents 获取所有学生
func (s *DefaultStudentService) GetAllStudents() ([]*model.Student, error) {
	students, _, err := s.studentRepo.List(1, 1000) // 设置一个较大的页大小来获取所有学生
	if err != nil {
		return nil, err
	}
	return students, nil
}

// CreateStudent 创建学生
func (s *DefaultStudentService) CreateStudent(req *model.StudentCreateRequest) error {
	// 密码加密
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return fmt.Errorf("error hashing password: %w", err)
	}

	// 创建学生对象
	student := &model.Student{
		ID:       req.ID,
		Name:     req.Name,
		Dept:     req.Dept,
		Password: hashedPassword,
		TotCred:  0.0, // 新学生总学分为0，使用float64类型
	}

	return s.studentRepo.Create(student)
}

// UpdateStudent 更新学生信息
func (s *DefaultStudentService) UpdateStudent(id string, req *model.StudentUpdateRequest) error {
	// 先查询学生是否存在
	existingStudent, err := s.studentRepo.GetByID(id)
	if err != nil {
		return err
	}

	// 更新学生信息
	existingStudent.Name = req.Name
	existingStudent.Dept = req.Dept

	return s.studentRepo.Update(existingStudent)
}

// DeleteStudent 删除学生
func (s *DefaultStudentService) DeleteStudent(id string) error {
	return s.studentRepo.Delete(id)
}

// Create 创建学生
func (s *DefaultStudentService) Create(student *model.Student) error {
	return s.studentRepo.Create(student)
}

// ChangePassword 修改密码
func (s *DefaultStudentService) ChangePassword(id string, req *model.ChangePasswordRequest) error {
	// 先查询学生是否存在
	student, err := s.studentRepo.GetByID(id)
	if err != nil {
		return err
	}

	// 验证旧密码
	if !utils.CheckPassword(req.OldPassword, student.Password) {
		return errors.New("invalid old password")
	}

	// 加密新密码
	hashedPassword, err := utils.HashPassword(req.NewPassword)
	if err != nil {
		return fmt.Errorf("error hashing password: %w", err)
	}

	// 更新密码 - 根据model.Student的Salt字段，这里应该传入salt
	return s.studentRepo.UpdatePassword(id, hashedPassword, student.Salt)
}

// GetStudentTranscript 获取学生成绩单
func (s *DefaultStudentService) GetStudentTranscript(id string) (*model.Transcript, error) {
	return s.takesRepo.GetStudentTranscript(id)
}

// GetCurrentCourses 获取学生当前学期的课程
func (s *DefaultStudentService) GetCurrentCourses(id string, semester string, year int) ([]*model.Takes, error) {
	return s.takesRepo.GetCurrentCourses(id, semester, year)
}

// RegisterForCourse 学生选课
func (s *DefaultStudentService) RegisterForCourse(studentID string, sectionID string, courseID string, semester string, year int) error {
	// 检查学生是否存在
	_, err := s.studentRepo.GetByID(studentID)
	if err != nil {
		return fmt.Errorf("student not found: %w", err)
	}

	// 检查课程段是否存在
	_, err = s.sectionRepo.FindByID(sectionID)
	if err != nil {
		return fmt.Errorf("section not found: %w", err)
	}

	// 检查是否已经选过这门课
	existingTakes, err := s.takesRepo.FindByStudentAndSection(studentID, sectionID)
	if err == nil && existingTakes != nil {
		return errors.New("already registered for this course")
	}

	// 检查先修课程要求
	satisfied, err := s.prereqRepo.CheckPrereqsSatisfied(studentID, courseID)
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

	// 创建选课记录
	takes := &model.Takes{
		StudentID: studentID,
		CourseID:  courseID,
		SectionID: sectionID,
		Semester:  semester,
		Year:      year,
		Grade:     "", // 新选课没有成绩
	}

	return s.takesRepo.Create(takes)
}

// DropCourse 学生退课
func (s *DefaultStudentService) DropCourse(studentID string, sectionID string) error {
	// 检查选课记录是否存在
	_, err := s.takesRepo.FindByStudentAndSection(studentID, sectionID)
	if err != nil {
		return fmt.Errorf("enrollment not found: %w", err)
	}

	return s.takesRepo.Delete(studentID, sectionID)
}

// GetByID 根据ID获取学生信息（别名方法）
func (s *DefaultStudentService) GetByID(id string) (*model.Student, error) {
	return s.GetStudentByID(id)
}

// UpdateProfile 更新学生个人信息
func (s *DefaultStudentService) UpdateProfile(id string, name string) error {
	student, err := s.studentRepo.GetByID(id)
	if err != nil {
		return err
	}

	student.Name = name
	return s.studentRepo.Update(student)
}

// GetAdvisor 获取学生导师信息
func (s *DefaultStudentService) GetAdvisor(id string) (*model.Advisor, error) {
	advisors, err := s.advisorRepo.FindByStudentID(id)
	if err != nil {
		return nil, err
	}
	if len(advisors) == 0 {
		return nil, errors.New("no advisor found")
	}
	return advisors[0], nil
}

// GetEnrolledCourses 获取学生已选课程
func (s *DefaultStudentService) GetEnrolledCourses(id string) ([]*model.Takes, error) {
	return s.takesRepo.FindByStudentID(id)
}

// GetTranscript 获取学生成绩单（别名方法）
func (s *DefaultStudentService) GetTranscript(id string) (*model.Transcript, error) {
	return s.GetStudentTranscript(id)
}

// Authenticate 学生认证
func (s *DefaultStudentService) Authenticate(id string, password string) (string, error) {
	student, err := s.studentRepo.GetByID(id)
	if err != nil {
		return "", fmt.Errorf("student not found: %w", err)
	}

	if !utils.CheckPassword(password, student.Password) {
		return "", errors.New("invalid password")
	}

	return student.ID, nil
}
