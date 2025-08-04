package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/yourusername/student-management-system/internal/model"
)

// SQLStudentRepository 实现StudentRepository接口
type SQLStudentRepository struct {
	db *sql.DB
}

// NewStudentRepository 创建学生仓储实例
func NewStudentRepository(db *sql.DB) StudentRepository {
	return &SQLStudentRepository{db: db}
}

// GetByID 根据ID查找学生
func (r *SQLStudentRepository) GetByID(id string) (*model.Student, error) {
	query := `SELECT id, name, dept_name, tot_cred, password, salt FROM student WHERE id = ?`
	row := r.db.QueryRow(query, id)

	var student model.Student
	err := row.Scan(&student.ID, &student.Name, &student.Dept, &student.TotCred, &student.Password, &student.Salt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("student not found: %w", err)
		}
		return nil, fmt.Errorf("error scanning student: %w", err)
	}

	return &student, nil
}

// List 查找所有学生
func (r *SQLStudentRepository) List(page, pageSize int) ([]*model.Student, int64, error) {
	var total int64
	countQuery := `SELECT COUNT(*) FROM student`
	err := r.db.QueryRow(countQuery).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("error counting students: %w", err)
	}

	offset := (page - 1) * pageSize
	query := `SELECT id, name, dept_name, tot_cred FROM student LIMIT ? OFFSET ?`
	rows, err := r.db.Query(query, pageSize, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("error querying students: %w", err)
	}
	defer rows.Close()

	var students []*model.Student
	for rows.Next() {
		var student model.Student
		err := rows.Scan(&student.ID, &student.Name, &student.Dept, &student.TotCred)
		if err != nil {
			return nil, 0, fmt.Errorf("error scanning student: %w", err)
		}
		students = append(students, &student)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("error iterating students: %w", err)
	}

	return students, total, nil
}

// Create 创建学生
func (r *SQLStudentRepository) Create(student *model.Student) error {
	query := `INSERT INTO student (id, name, dept_name, tot_cred, password, salt) VALUES (?, ?, ?, ?, ?, ?)`
	_, err := r.db.Exec(query, student.ID, student.Name, student.Dept, student.TotCred, student.Password, student.Salt)
	if err != nil {
		if strings.Contains(err.Error(), "Duplicate entry") {
			return fmt.Errorf("student already exists: %w", err)
		}
		return fmt.Errorf("error creating student: %w", err)
	}

	return nil
}

// Update 更新学生
func (r *SQLStudentRepository) Update(student *model.Student) error {
	query := `UPDATE student SET name = ?, dept_name = ?, tot_cred = ? WHERE id = ?`
	result, err := r.db.Exec(query, student.Name, student.Dept, student.TotCred, student.ID)
	if err != nil {
		return fmt.Errorf("error updating student: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error getting rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("student not found")
	}

	return nil
}

// Delete 删除学生
func (r *SQLStudentRepository) Delete(id string) error {
	query := `DELETE FROM student WHERE id = ?`
	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("error deleting student: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error getting rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("student not found")
	}

	return nil
}

// ExistsByID 检查指定ID的学生是否存在
func (r *SQLStudentRepository) ExistsByID(id string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM student WHERE id = ?)`
	err := r.db.QueryRow(query, id).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("error checking student existence: %w", err)
	}
	return exists, nil
}

// Search 搜索学生
func (r *SQLStudentRepository) Search(query string) ([]*model.Student, error) {
	sqlQuery := `SELECT id, name, dept_name, tot_cred FROM student WHERE id LIKE ? OR name LIKE ? OR dept_name LIKE ?`
	pattern := "%" + query + "%"
	rows, err := r.db.Query(sqlQuery, pattern, pattern, pattern)
	if err != nil {
		return nil, fmt.Errorf("error searching students: %w", err)
	}
	defer rows.Close()

	var students []*model.Student
	for rows.Next() {
		var student model.Student
		err := rows.Scan(&student.ID, &student.Name, &student.Dept, &student.TotCred)
		if err != nil {
			return nil, fmt.Errorf("error scanning student: %w", err)
		}
		students = append(students, &student)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating students: %w", err)
	}

	return students, nil
}

// UpdatePassword 更新学生密码
func (r *SQLStudentRepository) UpdatePassword(id, hashedPassword, salt string) error {
	query := `UPDATE student SET password = ?, salt = ? WHERE id = ?`
	result, err := r.db.Exec(query, hashedPassword, salt, id)
	if err != nil {
		return fmt.Errorf("error updating student password: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error getting rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("student not found")
	}

	return nil
}
