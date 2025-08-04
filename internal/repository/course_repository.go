package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/yourusername/student-management-system/internal/model"
)

// CourseRepository 定义课程仓储接口
type CourseRepository interface {
	FindByID(id string) (*model.Course, error)
	FindAll() ([]*model.Course, error)
	FindByDept(dept string) ([]*model.Course, error)
	Create(course *model.Course) error
	Update(course *model.Course) error
	Delete(id string) error
	FindWithPrereqs(id string) (*model.CourseWithPrereqs, error)
}

// SQLCourseRepository 实现CourseRepository接口
type SQLCourseRepository struct {
	db *sql.DB
}

// NewCourseRepository 创建课程仓储实例
func NewCourseRepository(db *sql.DB) CourseRepository {
	return &SQLCourseRepository{db: db}
}

// FindByID 根据ID查找课程
func (r *SQLCourseRepository) FindByID(id string) (*model.Course, error) {
	query := `SELECT course_id, title, dept_name, credits FROM course WHERE course_id = ?`
	row := r.db.QueryRow(query, id)

	var course model.Course
	err := row.Scan(&course.ID, &course.Title, &course.Dept, &course.Credits)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("course not found: %w", err)
		}
		return nil, fmt.Errorf("error scanning course: %w", err)
	}

	return &course, nil
}

// FindAll 查找所有课程
func (r *SQLCourseRepository) FindAll() ([]*model.Course, error) {
	query := `SELECT course_id, title, dept_name, credits FROM course`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error querying courses: %w", err)
	}
	defer rows.Close()

	var courses []*model.Course
	for rows.Next() {
		var course model.Course
		err := rows.Scan(&course.ID, &course.Title, &course.Dept, &course.Credits)
		if err != nil {
			return nil, fmt.Errorf("error scanning course: %w", err)
		}
		courses = append(courses, &course)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating courses: %w", err)
	}

	return courses, nil
}

// FindByDept 根据院系查找课程
func (r *SQLCourseRepository) FindByDept(dept string) ([]*model.Course, error) {
	query := `SELECT course_id, title, dept_name, credits FROM course WHERE dept_name = ?`
	rows, err := r.db.Query(query, dept)
	if err != nil {
		return nil, fmt.Errorf("error querying courses by dept: %w", err)
	}
	defer rows.Close()

	var courses []*model.Course
	for rows.Next() {
		var course model.Course
		err := rows.Scan(&course.ID, &course.Title, &course.Dept, &course.Credits)
		if err != nil {
			return nil, fmt.Errorf("error scanning course: %w", err)
		}
		courses = append(courses, &course)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating courses: %w", err)
	}

	return courses, nil
}

// Create 创建课程
func (r *SQLCourseRepository) Create(course *model.Course) error {
	query := `INSERT INTO course (course_id, title, dept_name, credits) VALUES (?, ?, ?, ?)`
	_, err := r.db.Exec(query, course.ID, course.Title, course.Dept, course.Credits)
	if err != nil {
		if strings.Contains(err.Error(), "Duplicate entry") {
			return fmt.Errorf("course already exists: %w", err)
		}
		return fmt.Errorf("error creating course: %w", err)
	}

	return nil
}

// Update 更新课程
func (r *SQLCourseRepository) Update(course *model.Course) error {
	query := `UPDATE course SET title = ?, dept_name = ?, credits = ? WHERE course_id = ?`
	result, err := r.db.Exec(query, course.Title, course.Dept, course.Credits, course.ID)
	if err != nil {
		return fmt.Errorf("error updating course: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error getting rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("course not found")
	}

	return nil
}

// Delete 删除课程
func (r *SQLCourseRepository) Delete(id string) error {
	query := `DELETE FROM course WHERE course_id = ?`
	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("error deleting course: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error getting rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("course not found")
	}

	return nil
}

// FindWithPrereqs 查找课程及其先修课程
func (r *SQLCourseRepository) FindWithPrereqs(id string) (*model.CourseWithPrereqs, error) {
	// 查找课程
	course, err := r.FindByID(id)
	if err != nil {
		return nil, err
	}

	// 查找先修课程
	query := `
		SELECT p.prereq_id, c.title, c.dept_name, c.credits 
		FROM prereq p 
		JOIN course c ON p.prereq_id = c.course_id 
		WHERE p.course_id = ?
	`
	rows, err := r.db.Query(query, id)
	if err != nil {
		return nil, fmt.Errorf("error querying prereqs: %w", err)
	}
	defer rows.Close()

	var prereqs []model.Course
	for rows.Next() {
		var prereq model.Course
		err := rows.Scan(&prereq.ID, &prereq.Title, &prereq.Dept, &prereq.Credits)
		if err != nil {
			return nil, fmt.Errorf("error scanning prereq: %w", err)
		}
		prereqs = append(prereqs, prereq)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating prereqs: %w", err)
	}

	return &model.CourseWithPrereqs{
		Course:  *course,
		Prereqs: prereqs,
	}, nil
}