package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/yourusername/student-management-system/internal/model"
)

// DepartmentRepository 定义院系仓储接口
type DepartmentRepository interface {
	FindByID(id string) (*model.Department, error)
	FindAll() ([]*model.Department, error)
	FindByDepartment(deptName string) (*model.Department, error)
	Create(department *model.Department) error
	Update(department *model.Department) error
	Delete(deptName string) error
	GetStudentCount(deptName string) (int, error)
	GetInstructorCount(deptName string) (int, error)
	GetCourseCount(deptName string) (int, error)
	GetDepartmentStats() ([]*model.DepartmentStats, error)
}

// SQLDepartmentRepository 实现DepartmentRepository接口
type SQLDepartmentRepository struct {
	db *sql.DB
}

// NewDepartmentRepository 创建院系仓储实例
func NewDepartmentRepository(db *sql.DB) DepartmentRepository {
	return &SQLDepartmentRepository{db: db}
}

// FindByID 根据院系名称查找院系
func (r *SQLDepartmentRepository) FindByID(deptName string) (*model.Department, error) {
	query := `SELECT dept_name, building, budget FROM department WHERE dept_name = ?`
	row := r.db.QueryRow(query, deptName)

	var department model.Department
	err := row.Scan(&department.DeptName, &department.Building, &department.Budget)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("department not found: %w", err)
		}
		return nil, fmt.Errorf("error scanning department: %w", err)
	}

	return &department, nil
}

// FindAll 查找所有院系
func (r *SQLDepartmentRepository) FindAll() ([]*model.Department, error) {
	query := `SELECT dept_name, building, budget FROM department ORDER BY dept_name`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error querying departments: %w", err)
	}
	defer rows.Close()

	var departments []*model.Department
	for rows.Next() {
		var department model.Department
		err := rows.Scan(&department.DeptName, &department.Building, &department.Budget)
		if err != nil {
			return nil, fmt.Errorf("error scanning department: %w", err)
		}
		departments = append(departments, &department)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating departments: %w", err)
	}

	return departments, nil
}

// Create 创建院系
func (r *SQLDepartmentRepository) Create(department *model.Department) error {
	query := `INSERT INTO department (dept_name, building, budget) VALUES (?, ?, ?)`
	_, err := r.db.Exec(query, department.DeptName, department.Building, department.Budget)
	if err != nil {
		if strings.Contains(err.Error(), "Duplicate entry") {
			return fmt.Errorf("department already exists: %w", err)
		}
		return fmt.Errorf("error creating department: %w", err)
	}

	return nil
}

// Update 更新院系
func (r *SQLDepartmentRepository) Update(department *model.Department) error {
	query := `UPDATE department SET building = ?, budget = ? WHERE dept_name = ?`
	result, err := r.db.Exec(query, department.Building, department.Budget, department.DeptName)
	if err != nil {
		return fmt.Errorf("error updating department: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error getting rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("department not found")
	}

	return nil
}

// Delete 删除院系
func (r *SQLDepartmentRepository) Delete(deptName string) error {
	// 检查院系是否有关联的学生
	studentQuery := `SELECT COUNT(*) FROM student WHERE dept_name = ?`
	var studentCount int
	err := r.db.QueryRow(studentQuery, deptName).Scan(&studentCount)
	if err != nil {
		return fmt.Errorf("error checking student count: %w", err)
	}

	if studentCount > 0 {
		return fmt.Errorf("cannot delete department: it has %d associated students", studentCount)
	}

	// 检查院系是否有关联的教师
	instructorQuery := `SELECT COUNT(*) FROM instructor WHERE dept_name = ?`
	var instructorCount int
	err = r.db.QueryRow(instructorQuery, deptName).Scan(&instructorCount)
	if err != nil {
		return fmt.Errorf("error checking instructor count: %w", err)
	}

	if instructorCount > 0 {
		return fmt.Errorf("cannot delete department: it has %d associated instructors", instructorCount)
	}

	// 检查院系是否有关联的课程
	courseQuery := `SELECT COUNT(*) FROM course WHERE dept_name = ?`
	var courseCount int
	err = r.db.QueryRow(courseQuery, deptName).Scan(&courseCount)
	if err != nil {
		return fmt.Errorf("error checking course count: %w", err)
	}

	if courseCount > 0 {
		return fmt.Errorf("cannot delete department: it has %d associated courses", courseCount)
	}

	// 删除院系
	query := `DELETE FROM department WHERE dept_name = ?`
	result, err := r.db.Exec(query, deptName)
	if err != nil {
		return fmt.Errorf("error deleting department: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error getting rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("department not found")
	}

	return nil
}

// GetStudentCount 获取院系学生数量
func (r *SQLDepartmentRepository) GetStudentCount(deptName string) (int, error) {
	query := `SELECT COUNT(*) FROM student WHERE dept_name = ?`
	var count int
	err := r.db.QueryRow(query, deptName).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("error getting student count: %w", err)
	}

	return count, nil
}

// GetInstructorCount 获取院系教师数量
func (r *SQLDepartmentRepository) GetInstructorCount(deptName string) (int, error) {
	query := `SELECT COUNT(*) FROM instructor WHERE dept_name = ?`
	var count int
	err := r.db.QueryRow(query, deptName).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("error getting instructor count: %w", err)
	}

	return count, nil
}

// GetCourseCount 获取院系课程数量
func (r *SQLDepartmentRepository) GetCourseCount(deptName string) (int, error) {
	query := `SELECT COUNT(*) FROM course WHERE dept_name = ?`
	var count int
	err := r.db.QueryRow(query, deptName).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("error getting course count: %w", err)
	}

	return count, nil
}

// GetDepartmentStats 获取所有院系统计信息
func (r *SQLDepartmentRepository) GetDepartmentStats() ([]*model.DepartmentStats, error) {
	query := `
		SELECT 
			d.dept_name, 
			d.building, 
			d.budget,
			(SELECT COUNT(*) FROM student s WHERE s.dept_name = d.dept_name) as student_count,
			(SELECT COUNT(*) FROM instructor i WHERE i.dept_name = d.dept_name) as instructor_count,
			(SELECT COUNT(*) FROM course c WHERE c.dept_name = d.dept_name) as course_count
		FROM department d
		ORDER BY d.dept_name
	`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error querying department stats: %w", err)
	}
	defer rows.Close()

	var stats []*model.DepartmentStats
	for rows.Next() {
		var stat model.DepartmentStats
		var department model.Department

		err := rows.Scan(
			&department.DeptName,
			&department.Building,
			&department.Budget,
			&stat.StudentCount,
			&stat.InstructorCount,
			&stat.CourseCount,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning department stats: %w", err)
		}

		stat.Department = department
		stats = append(stats, &stat)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating department stats: %w", err)
	}

	return stats, nil
}

// 以下是额外的方法实现

// 在 SQLDepartmentRepository 中实现缺少的方法
func (r *SQLDepartmentRepository) FindByDepartment(deptName string) (*model.Department, error) {
	return r.FindByID(deptName)
}

// FindByDepartment方法已通过FindByID实现