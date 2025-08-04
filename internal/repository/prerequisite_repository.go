package repository

import (
	"database/sql"
	"fmt"

	"github.com/yourusername/student-management-system/internal/model"
)

// PrerequisiteRepository 定义先修课程仓库接口
type PrerequisiteRepository interface {
	FindByCourseID(courseID string) ([]*model.Prereq, error)
	Create(prereq *model.Prereq) error
	Delete(courseID, prereqID string) error
}

// SQLPrerequisiteRepository 实现PrerequisiteRepository接口
type SQLPrerequisiteRepository struct {
	db *sql.DB
}

// NewPrerequisiteRepository 创建先修课程仓库实例
func NewPrerequisiteRepository(db *sql.DB) PrerequisiteRepository {
	return &SQLPrerequisiteRepository{db: db}
}

// FindByCourseID 根据课程ID查找先修课程
func (r *SQLPrerequisiteRepository) FindByCourseID(courseID string) ([]*model.Prereq, error) {
	query := `
		SELECT p.course_id, p.prereq_id, c.title, c.dept_name, c.credits
		FROM prereq p
		JOIN course c ON p.prereq_id = c.course_id
		WHERE p.course_id = ?
	`

	rows, err := r.db.Query(query, courseID)
	if err != nil {
		return nil, fmt.Errorf("error querying prerequisites: %w", err)
	}
	defer rows.Close()

	var prereqs []*model.Prereq
	for rows.Next() {
		var prereq model.Prereq
		var course model.Course

		err := rows.Scan(
			&prereq.CourseID,
			&prereq.PrereqID,
			&course.Title,
			&course.Dept,
			&course.Credits,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning prerequisite row: %w", err)
		}

		course.ID = prereq.PrereqID
		prereq.PrereqInfo = &course

		prereqs = append(prereqs, &prereq)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating prerequisite rows: %w", err)
	}

	return prereqs, nil
}

// Create 创建先修课程关系
func (r *SQLPrerequisiteRepository) Create(prereq *model.Prereq) error {
	query := `INSERT INTO prereq (course_id, prereq_id) VALUES (?, ?)`

	_, err := r.db.Exec(query, prereq.CourseID, prereq.PrereqID)
	if err != nil {
		return fmt.Errorf("error creating prerequisite: %w", err)
	}

	return nil
}

// Delete 删除先修课程关系
func (r *SQLPrerequisiteRepository) Delete(courseID, prereqID string) error {
	query := `DELETE FROM prereq WHERE course_id = ? AND prereq_id = ?`

	_, err := r.db.Exec(query, courseID, prereqID)
	if err != nil {
		return fmt.Errorf("error deleting prerequisite: %w", err)
	}

	return nil
}