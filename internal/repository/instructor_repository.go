package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/yourusername/student-management-system/internal/model"
)

// SQLInstructorRepository 实现InstructorRepository接口
type SQLInstructorRepository struct {
	db *sql.DB
}

// NewInstructorRepository 创建教师仓储实例
func NewInstructorRepository(db *sql.DB) InstructorRepository {
	return &SQLInstructorRepository{db: db}
}

// GetByID 根据ID查找教师
func (r *SQLInstructorRepository) GetByID(id string) (*model.Instructor, error) {
	query := `SELECT id, name, dept_name, salary, password, salt FROM instructor WHERE id = ?`
	row := r.db.QueryRow(query, id)

	var instructor model.Instructor
	err := row.Scan(&instructor.ID, &instructor.Name, &instructor.Dept, &instructor.Salary, &instructor.Password, &instructor.Salt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("instructor not found: %w", err)
		}
		return nil, fmt.Errorf("error scanning instructor: %w", err)
	}

	return &instructor, nil
}

// List 查找所有教师
func (r *SQLInstructorRepository) List(page, pageSize int) ([]*model.Instructor, int64, error) {
	var total int64
	countQuery := `SELECT COUNT(*) FROM instructor`
	err := r.db.QueryRow(countQuery).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("error counting instructors: %w", err)
	}

	offset := (page - 1) * pageSize
	query := `SELECT id, name, dept_name, salary FROM instructor LIMIT ? OFFSET ?`
	rows, err := r.db.Query(query, pageSize, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("error querying instructors: %w", err)
	}
	defer rows.Close()

	var instructors []*model.Instructor
	for rows.Next() {
		var instructor model.Instructor
		err := rows.Scan(&instructor.ID, &instructor.Name, &instructor.Dept, &instructor.Salary)
		if err != nil {
			return nil, 0, fmt.Errorf("error scanning instructor: %w", err)
		}
		instructors = append(instructors, &instructor)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("error iterating instructors: %w", err)
	}

	return instructors, total, nil
}

// Create 创建教师
func (r *SQLInstructorRepository) Create(instructor *model.Instructor) error {
	query := `INSERT INTO instructor (id, name, dept_name, salary, password, salt) VALUES (?, ?, ?, ?, ?, ?)`
	_, err := r.db.Exec(query, instructor.ID, instructor.Name, instructor.Dept, instructor.Salary, instructor.Password, instructor.Salt)
	if err != nil {
		if strings.Contains(err.Error(), "Duplicate entry") {
			return fmt.Errorf("instructor already exists: %w", err)
		}
		return fmt.Errorf("error creating instructor: %w", err)
	}

	return nil
}

// Update 更新教师
func (r *SQLInstructorRepository) Update(instructor *model.Instructor) error {
	query := `UPDATE instructor SET name = ?, dept_name = ?, salary = ? WHERE id = ?`
	result, err := r.db.Exec(query, instructor.Name, instructor.Dept, instructor.Salary, instructor.ID)
	if err != nil {
		return fmt.Errorf("error updating instructor: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error getting rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("instructor not found")
	}

	return nil
}

// Delete 删除教师
func (r *SQLInstructorRepository) Delete(id string) error {
	query := `DELETE FROM instructor WHERE id = ?`
	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("error deleting instructor: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error getting rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("instructor not found")
	}

	return nil
}

// ExistsByID 检查指定ID的教师是否存在
func (r *SQLInstructorRepository) ExistsByID(id string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM instructor WHERE id = ?)`
	err := r.db.QueryRow(query, id).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("error checking instructor existence: %w", err)
	}
	return exists, nil
}

// Search 搜索教师
func (r *SQLInstructorRepository) Search(query string) ([]*model.Instructor, error) {
	sqlQuery := `SELECT id, name, dept_name, salary FROM instructor WHERE id LIKE ? OR name LIKE ? OR dept_name LIKE ?`
	pattern := "%" + query + "%"
	rows, err := r.db.Query(sqlQuery, pattern, pattern, pattern)
	if err != nil {
		return nil, fmt.Errorf("error searching instructors: %w", err)
	}
	defer rows.Close()

	var instructors []*model.Instructor
	for rows.Next() {
		var instructor model.Instructor
		err := rows.Scan(&instructor.ID, &instructor.Name, &instructor.Dept, &instructor.Salary)
		if err != nil {
			return nil, fmt.Errorf("error scanning instructor: %w", err)
		}
		instructors = append(instructors, &instructor)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating instructors: %w", err)
	}

	return instructors, nil
}

// UpdatePassword 更新教师密码
func (r *SQLInstructorRepository) UpdatePassword(id, hashedPassword, salt string) error {
	query := `UPDATE instructor SET password = ?, salt = ? WHERE id = ?`
	result, err := r.db.Exec(query, hashedPassword, salt, id)
	if err != nil {
		return fmt.Errorf("error updating instructor password: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error getting rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("instructor not found")
	}

	return nil
}
