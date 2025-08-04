package service

import (
	"testing"

	"github.com/yourusername/student-management-system/internal/model"
	"github.com/yourusername/student-management-system/internal/repository"
)

// MockStudentRepository 模拟学生仓库
type MockStudentRepository struct {
	students map[string]*model.Student
}

func NewMockStudentRepository() *MockStudentRepository {
	return &MockStudentRepository{
		students: make(map[string]*model.Student),
	}
}

func (m *MockStudentRepository) Create(student *model.Student) error {
	m.students[student.ID] = student
	return nil
}

func (m *MockStudentRepository) FindByID(id string) (*model.Student, error) {
	if student, exists := m.students[id]; exists {
		return student, nil
	}
	return nil, repository.ErrNotFound
}

func (m *MockStudentRepository) Update(student *model.Student) error {
	if _, exists := m.students[student.ID]; exists {
		m.students[student.ID] = student
		return nil
	}
	return repository.ErrNotFound
}

func (m *MockStudentRepository) Delete(id string) error {
	if _, exists := m.students[id]; exists {
		delete(m.students, id)
		return nil
	}
	return repository.ErrNotFound
}

func (m *MockStudentRepository) FindAll() ([]*model.Student, error) {
	students := make([]*model.Student, 0, len(m.students))
	for _, student := range m.students {
		students = append(students, student)
	}
	return students, nil
}

func (m *MockStudentRepository) UpdatePassword(id string, hashedPassword string, salt string) error {
	if student, exists := m.students[id]; exists {
		student.Password = hashedPassword
		student.Salt = salt
		return nil
	}
	return repository.ErrNotFound
}

// MockTakesRepository 是TakesRepository的模拟实现
type MockTakesRepository struct{}

func (m *MockTakesRepository) FindByStudentID(studentID string) ([]*model.Takes, error) {
	return nil, nil
}

func (m *MockTakesRepository) FindByStudentAndSection(studentID, sectionID string) (*model.Takes, error) {
	return nil, nil
}

func (m *MockTakesRepository) FindBySection(sectionID string) ([]*model.Takes, error) {
	return nil, nil
}

func (m *MockTakesRepository) FindBySectionID(sectionID string) ([]*model.Takes, error) {
	return nil, nil
}

func (m *MockTakesRepository) Create(takes *model.Takes) error {
	return nil
}

func (m *MockTakesRepository) Delete(studentID, sectionID string) error {
	return nil
}

func (m *MockTakesRepository) UpdateGrade(studentID, sectionID, grade string) error {
	return nil
}

func (m *MockTakesRepository) GetStudentTranscript(studentID string) (*model.Transcript, error) {
	return nil, nil
}

func (m *MockTakesRepository) GetCurrentCourses(studentID string, semester string, year int) ([]*model.Takes, error) {
	return nil, nil
}

func (m *MockTakesRepository) CheckTimeConflict(studentID, sectionID string) (bool, error) {
	return false, nil
}

// MockPrereqRepository 是PrereqRepository的模拟实现
type MockPrereqRepository struct{}

func (m *MockPrereqRepository) FindByID(courseID, prereqID string) (*model.Prereq, error) {
	return nil, nil
}

func (m *MockPrereqRepository) FindAll() ([]*model.Prereq, error) {
	return nil, nil
}

func (m *MockPrereqRepository) FindByCourseID(courseID string) ([]*model.Prereq, error) {
	return nil, nil
}

func (m *MockPrereqRepository) Create(prereq *model.Prereq) error {
	return nil
}

func (m *MockPrereqRepository) Delete(courseID, prereqID string) error {
	return nil
}

func (m *MockPrereqRepository) GetPrereqIDs(courseID string) ([]string, error) {
	return nil, nil
}

func (m *MockPrereqRepository) CheckPrereqsSatisfied(studentID string, courseID string) (bool, error) {
	return true, nil
}

func (m *MockPrereqRepository) HasPrerequisite(courseID string, prereqID string) (bool, error) {
	return false, nil
}

// MockSectionRepository 是SectionRepository的模拟实现
type MockSectionRepository struct{}

func (m *MockSectionRepository) FindByID(id string) (*model.Section, error) {
	return nil, nil
}

func (m *MockSectionRepository) FindByCourseID(courseID string) ([]*model.Section, error) {
	return nil, nil
}

func (m *MockSectionRepository) FindAll() ([]*model.Section, error) {
	return nil, nil
}

func (m *MockSectionRepository) Create(section *model.Section) error {
	return nil
}

func (m *MockSectionRepository) Update(section *model.Section) error {
	return nil
}

func (m *MockSectionRepository) Delete(id string) error {
	return nil
}

func (m *MockSectionRepository) FindByParams(params *model.SectionQueryParams) ([]*model.Section, error) {
	return nil, nil
}

func (m *MockSectionRepository) FindWithDetails(id string) (*model.Section, error) {
	return nil, nil
}

func (m *MockSectionRepository) GetEnrollmentCount(sectionID string) (int, error) {
	return 0, nil
}

func (m *MockSectionRepository) GetSectionClassroom(sectionID string) (*model.Classroom, error) {
	return nil, nil
}

// MockAdvisorRepository 是AdvisorRepository的模拟实现
type MockAdvisorRepository struct{}

// FindByID 根据ID查找导师关系
func (m *MockAdvisorRepository) FindByID(studentID string, instructorID string) (*model.Advisor, error) {
	return nil, nil
}

// FindAll 查找所有导师关系
func (m *MockAdvisorRepository) FindAll() ([]*model.Advisor, error) {
	return nil, nil
}

// FindByStudentID 根据学生ID查找导师关系
func (m *MockAdvisorRepository) FindByStudentID(studentID string) ([]*model.Advisor, error) {
	return nil, nil
}

// FindByInstructorID 根据导师ID查找导师关系
func (m *MockAdvisorRepository) FindByInstructorID(instructorID string) ([]*model.Advisor, error) {
	return nil, nil
}

// Create 创建导师关系
func (m *MockAdvisorRepository) Create(advisor *model.Advisor) error {
	return nil
}

// Delete 删除导师关系
func (m *MockAdvisorRepository) Delete(studentID, instructorID string) error {
	return nil
}

// FindByStudentAndInstructor 根据学生ID和导师ID查找导师关系
func (m *MockAdvisorRepository) FindByStudentAndInstructor(studentID string, instructorID string) (*model.Advisor, error) {
	return nil, nil
}

// Update 更新导师关系
func (m *MockAdvisorRepository) Update(studentID string, instructorID string) error {
	return nil
}

// 测试用例
func TestStudentService_Create(t *testing.T) {
	mockRepo := NewMockStudentRepository()
	// 创建所需的mock仓库
	mockTakesRepo := &MockTakesRepository{}
	mockPrereqRepo := &MockPrereqRepository{}
	mockSectionRepo := &MockSectionRepository{}
	mockAdvisorRepo := &MockAdvisorRepository{}

	service := NewStudentService(mockRepo, mockTakesRepo, mockPrereqRepo, mockSectionRepo, mockAdvisorRepo)

	student := &model.Student{
		ID:   "S001",
		Name: "张三",
		Dept: "计算机科学",
	}

	err := service.Create(student)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// 验证学生是否被创建
	created, err := service.GetByID("S001")
	if err != nil {
		t.Errorf("Expected to find student, got error %v", err)
	}

	if created.Name != "张三" {
		t.Errorf("Expected name '张三', got %s", created.Name)
	}
}

func TestStudentService_GetByID_NotFound(t *testing.T) {
	mockRepo := NewMockStudentRepository()
	// 创建所需的mock仓库
	mockTakesRepo := &MockTakesRepository{}
	mockPrereqRepo := &MockPrereqRepository{}
	mockSectionRepo := &MockSectionRepository{}
	mockAdvisorRepo := &MockAdvisorRepository{}

	service := NewStudentService(mockRepo, mockTakesRepo, mockPrereqRepo, mockSectionRepo, mockAdvisorRepo)

	_, err := service.GetByID("nonexistent")
	if err == nil {
		t.Error("Expected error for non-existent student")
	}
}
