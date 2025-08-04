package repository

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/yourusername/student-management-system/internal/model"
)

// PrereqRepository 定义先修课程仓储接口
type PrereqRepository interface {
	FindByID(courseID, prereqID string) (*model.Prereq, error)
	FindAll() ([]*model.Prereq, error)
	FindByCourseID(courseID string) ([]*model.Prereq, error)
	Create(prereq *model.Prereq) error
	Delete(courseID, prereqID string) error
	GetPrereqIDs(courseID string) ([]string, error)
	CheckPrereqsSatisfied(studentID string, courseID string) (bool, error)
	HasPrerequisite(courseID string, prereqID string) (bool, error)
}

// SQLPrereqRepository 实现PrereqRepository接口
type SQLPrereqRepository struct {
	db *sql.DB
}

// NewPrereqRepository 创建先修课程仓储实例
func NewPrereqRepository(db *sql.DB) *SQLPrereqRepository {
	return &SQLPrereqRepository{db: db}
}

// FindByCourseID 根据课程ID查找先修课程
func (r *SQLPrereqRepository) FindByCourseID(courseID string) ([]*model.Prereq, error) {
	query := `
		SELECT p.course_id, p.prereq_id, c.title, c.dept_name, c.credits 
		FROM prereq p 
		JOIN course c ON p.prereq_id = c.course_id 
		WHERE p.course_id = ?
	`
	rows, err := r.db.Query(query, courseID)
	if err != nil {
		return nil, fmt.Errorf("error querying prereqs: %w", err)
	}
	defer rows.Close()

	var prereqs []*model.Prereq
	for rows.Next() {
		var prereq model.Prereq
		var course model.Course
		err := rows.Scan(&prereq.CourseID, &prereq.PrereqID, &course.Title, &course.Dept, &course.Credits)
		if err != nil {
			return nil, fmt.Errorf("error scanning prereq: %w", err)
		}
		course.ID = prereq.PrereqID
		prereq.PrereqInfo = &course
		prereqs = append(prereqs, &prereq)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating prereqs: %w", err)
	}

	return prereqs, nil
}

// Create 创建先修课程关系
func (r *SQLPrereqRepository) Create(prereq *model.Prereq) error {
	query := `INSERT INTO prereq (course_id, prereq_id) VALUES (?, ?)`
	_, err := r.db.Exec(query, prereq.CourseID, prereq.PrereqID)
	if err != nil {
		if strings.Contains(err.Error(), "Duplicate entry") {
			return fmt.Errorf("prereq relationship already exists: %w", err)
		}
		return fmt.Errorf("error creating prereq: %w", err)
	}

	return nil
}

// Delete 删除先修课程关系
func (r *SQLPrereqRepository) Delete(courseID, prereqID string) error {
	query := `DELETE FROM prereq WHERE course_id = ? AND prereq_id = ?`
	result, err := r.db.Exec(query, courseID, prereqID)
	if err != nil {
		return fmt.Errorf("error deleting prereq: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error getting rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("prereq relationship not found")
	}

	return nil
}

// CheckPrereqsSatisfied 检查学生是否满足课程的所有先修课程要求
func (r *SQLPrereqRepository) CheckPrereqsSatisfied(studentID, courseID string) (bool, error) {
	// 查询课程的所有先修课程
	query := `
		SELECT prereq_id 
		FROM prereq 
		WHERE course_id = ?
	`
	rows, err := r.db.Query(query, courseID)
	if err != nil {
		return false, fmt.Errorf("error querying prereqs: %w", err)
	}
	defer rows.Close()

	var prereqIDs []string
	for rows.Next() {
		var prereqID string
		err := rows.Scan(&prereqID)
		if err != nil {
			return false, fmt.Errorf("error scanning prereq: %w", err)
		}
		prereqIDs = append(prereqIDs, prereqID)
	}

	if err := rows.Err(); err != nil {
		return false, fmt.Errorf("error iterating prereqs: %w", err)
	}

	// 如果没有先修课程，则满足条件
	if len(prereqIDs) == 0 {
		return true, nil
	}

	// 检查学生是否已经通过所有先修课程
	for _, prereqID := range prereqIDs {
		// 查询学生是否通过了该先修课程
		query := `
			SELECT COUNT(*) 
			FROM takes 
			WHERE student_id = ? AND course_id = ? AND grade IS NOT NULL AND grade <> 'F'
		`
		var count int
		err := r.db.QueryRow(query, studentID, prereqID).Scan(&count)
		if err != nil {
			return false, fmt.Errorf("error checking prereq completion: %w", err)
		}

		// 如果学生没有通过该先修课程，则不满足条件
		if count == 0 {
			return false, nil
		}
	}

	// 学生已通过所有先修课程
	return true, nil
}

func (r *SQLPrereqRepository) GetPrereqIDs(courseID string) ([]string, error) {
	query := `SELECT prereq_id FROM prereq WHERE course_id = ?`

	rows, err := r.db.Query(query, courseID)
	if err != nil {
		return nil, fmt.Errorf("error querying prereqs: %w", err)
	}
	defer rows.Close()

	var prereqIDs []string
	for rows.Next() {
		var prereqID string
		err := rows.Scan(&prereqID)
		if err != nil {
			return nil, fmt.Errorf("error scanning prereq_id: %w", err)
		}
		prereqIDs = append(prereqIDs, prereqID)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating prereqs: %w", err)
	}

	return prereqIDs, nil
}

func (r *SQLPrereqRepository) FindAll() ([]*model.Prereq, error) {
	query := `SELECT course_id, prereq_id FROM prereq`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error querying prereqs: %w", err)
	}
	defer rows.Close()

	var prereqs []*model.Prereq
	for rows.Next() {
		var prereq model.Prereq
		err := rows.Scan(
			&prereq.CourseID,
			&prereq.PrereqID,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning prereq: %w", err)
		}
		prereqs = append(prereqs, &prereq)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating prereqs: %w", err)
	}

	return prereqs, nil
}

// FindByID 根据课程ID和前置课程ID查找前置课程关系
func (r *SQLPrereqRepository) FindByID(courseID string, prereqID string) (*model.Prereq, error) {
	query := `SELECT course_id, prereq_id FROM prereq WHERE course_id = ? AND prereq_id = ?`

	var prereq model.Prereq
	err := r.db.QueryRow(query, courseID, prereqID).Scan(
		&prereq.CourseID,
		&prereq.PrereqID,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("prerequisite relationship not found")
		}
		return nil, fmt.Errorf("error querying prerequisite: %w", err)
	}

	return &prereq, nil
}

// HasPrerequisite 检查是否存在先修课程关系
func (r *SQLPrereqRepository) HasPrerequisite(courseID string, prereqID string) (bool, error) {
	query := `SELECT 1 FROM prereq WHERE course_id = ? AND prereq_id = ?`

	var exists int
	err := r.db.QueryRow(query, courseID, prereqID).Scan(&exists)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, fmt.Errorf("error checking prerequisite: %w", err)
	}

	return true, nil
}
