package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/yourusername/student-management-system/internal/model"
)

// AdvisorRepository 定义导师关系仓储接口
type AdvisorRepository interface {
	FindByID(studentID, instructorID string) (*model.Advisor, error)
	FindAll() ([]*model.Advisor, error)
	FindByStudentID(studentID string) ([]*model.Advisor, error)
	FindByInstructorID(instructorID string) ([]*model.Advisor, error)
	Create(advisor *model.Advisor) error
	Update(studentID string, instructorID string) error
	Delete(studentID, instructorID string) error
	FindByStudentAndInstructor(studentID string, instructorID string) (*model.Advisor, error)
}

// SQLAdvisorRepository 实现AdvisorRepository接口
type SQLAdvisorRepository struct {
	db *sql.DB
}

// NewAdvisorRepository 创建导师关系仓储实例
func NewAdvisorRepository(db *sql.DB) *SQLAdvisorRepository {
	return &SQLAdvisorRepository{db: db}
}

// FindByID 根据学生ID和导师ID查找导师关系
func (r *SQLAdvisorRepository) FindByID(studentID, instructorID string) (*model.Advisor, error) {
	query := `
		SELECT a.student_id, a.instructor_id, s.name as student_name, s.dept_name as student_dept,
		       i.name as instructor_name, i.dept_name as instructor_dept
		FROM advisor a
		JOIN student s ON a.student_id = s.id
		JOIN instructor i ON a.instructor_id = i.id
		WHERE a.student_id = ? AND a.instructor_id = ?
	`
	row := r.db.QueryRow(query, studentID, instructorID)

	var advisor model.Advisor
	var student model.Student
	var instructor model.Instructor

	err := row.Scan(
		&advisor.StudentID,
		&advisor.InstructorID,
		&student.Name,
		&student.Dept,
		&instructor.Name,
		&instructor.Dept,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("advisor relationship not found")
		}
		return nil, fmt.Errorf("error querying advisor: %w", err)
	}

	advisor.Student = &student
	advisor.Instructor = &instructor

	return &advisor, nil
}

// FindByStudentID 根据学生ID查找导师关系
func (r *SQLAdvisorRepository) FindByStudentID(studentID string) ([]*model.Advisor, error) {
	query := `
		SELECT a.student_id, a.instructor_id, s.name as student_name, s.dept_name as student_dept,
		       i.name as instructor_name, i.dept_name as instructor_dept
		FROM advisor a
		JOIN student s ON a.student_id = s.id
		JOIN instructor i ON a.instructor_id = i.id
		WHERE a.student_id = ?
	`
	rows, err := r.db.Query(query, studentID)
	if err != nil {
		return nil, fmt.Errorf("error querying advisors: %w", err)
	}
	defer rows.Close()

	var advisors []*model.Advisor
	for rows.Next() {
		var advisor model.Advisor
		var student model.Student
		var instructor model.Instructor

		err := rows.Scan(
			&advisor.StudentID,
			&advisor.InstructorID,
			&student.Name,
			&student.Dept,
			&instructor.Name,
			&instructor.Dept,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning advisor: %w", err)
		}

		student.ID = advisor.StudentID
		instructor.ID = advisor.InstructorID
		advisor.Student = &student
		advisor.Instructor = &instructor
		advisors = append(advisors, &advisor)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating advisors: %w", err)
	}

	return advisors, nil
}

// FindByInstructorID 根据教师ID查找所有指导的学生
func (r *SQLAdvisorRepository) FindByInstructorID(instructorID string) ([]*model.Advisor, error) {
	query := `
		SELECT a.student_id, a.instructor_id, s.name as student_name, s.dept_name as student_dept, s.tot_cred
		FROM advisor a
		JOIN student s ON a.student_id = s.id
		WHERE a.instructor_id = ?
		ORDER BY s.name
	`
	rows, err := r.db.Query(query, instructorID)
	if err != nil {
		return nil, fmt.Errorf("error querying advisors: %w", err)
	}
	defer rows.Close()

	var advisorList []*model.Advisor
	for rows.Next() {
		var advisor model.Advisor
		var student model.Student

		err := rows.Scan(
			&advisor.StudentID,
			&advisor.InstructorID,
			&student.Name,
			&student.Dept,
			&student.TotCred,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning advisor: %w", err)
		}

		student.ID = advisor.StudentID
		advisor.Student = &student
		advisorList = append(advisorList, &advisor)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating advisors: %w", err)
	}

	return advisorList, nil
}

// Create 创建导师关系
func (r *SQLAdvisorRepository) Create(advisor *model.Advisor) error {
	// 检查学生是否已有导师
	checkQuery := `SELECT instructor_id FROM advisor WHERE student_id = ?`
	var existingInstructorID string
	err := r.db.QueryRow(checkQuery, advisor.StudentID).Scan(&existingInstructorID)
	if err == nil {
		// 学生已有导师
		return fmt.Errorf("student already has an advisor (instructor_id: %s)", existingInstructorID)
	} else if !errors.Is(err, sql.ErrNoRows) {
		// 查询出错
		return fmt.Errorf("error checking existing advisor: %w", err)
	}

	// 创建新的导师关系
	query := `INSERT INTO advisor (student_id, instructor_id) VALUES (?, ?)`
	_, err = r.db.Exec(query, advisor.StudentID, advisor.InstructorID)
	if err != nil {
		if strings.Contains(err.Error(), "Duplicate entry") {
			return fmt.Errorf("advisor relationship already exists: %w", err)
		}
		return fmt.Errorf("error creating advisor relationship: %w", err)
	}

	return nil
}

// Update 更新学生的导师
func (r *SQLAdvisorRepository) Update(studentID string, instructorID string) error {
	// 检查学生是否存在导师关系
	checkQuery := `SELECT instructor_id FROM advisor WHERE student_id = ?`
	var existingInstructorID string
	err := r.db.QueryRow(checkQuery, studentID).Scan(&existingInstructorID)
	if errors.Is(err, sql.ErrNoRows) {
		// 学生没有导师，创建新关系
		return r.Create(&model.Advisor{StudentID: studentID, InstructorID: instructorID})
	} else if err != nil {
		// 查询出错
		return fmt.Errorf("error checking existing advisor: %w", err)
	}

	// 更新导师关系
	query := `UPDATE advisor SET instructor_id = ? WHERE student_id = ?`
	result, err := r.db.Exec(query, instructorID, studentID)
	if err != nil {
		return fmt.Errorf("error updating advisor: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error getting rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("advisor relationship not found")
	}

	return nil
}

// FindByStudentAndInstructor 根据学生和导师ID查找导师关系
func (r *SQLAdvisorRepository) FindByStudentAndInstructor(studentID string, instructorID string) (*model.Advisor, error) {
	query := `
		SELECT a.student_id, a.instructor_id, s.name as student_name, s.dept_name as student_dept,
		       i.name as instructor_name, i.dept_name as instructor_dept
		FROM advisor a
		JOIN student s ON a.student_id = s.id
		JOIN instructor i ON a.instructor_id = i.id
		WHERE a.student_id = ? AND a.instructor_id = ?
	`
	row := r.db.QueryRow(query, studentID, instructorID)

	var advisor model.Advisor
	var student model.Student
	var instructor model.Instructor

	err := row.Scan(
		&advisor.StudentID,
		&advisor.InstructorID,
		&student.Name,
		&student.Dept,
		&instructor.Name,
		&instructor.Dept,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("advisor relationship not found: %w", err)
		}
		return nil, fmt.Errorf("error scanning advisor: %w", err)
	}

	student.ID = advisor.StudentID
	instructor.ID = advisor.InstructorID
	advisor.Student = &student
	advisor.Instructor = &instructor

	return &advisor, nil
}

// Delete 删除导师关系
func (r *SQLAdvisorRepository) Delete(studentID, instructorID string) error {
	query := `DELETE FROM advisor WHERE s_id = ? AND i_id = ?`

	result, err := r.db.Exec(query, studentID, instructorID)
	if err != nil {
		return fmt.Errorf("error deleting advisor: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error getting rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("advisor relationship not found")
	}

	return nil
}

// 注意：AdvisorRepository接口已在文件顶部定义，这里不需要重复定义

// 在 SQLAdvisorRepository 中实现 FindAll 方法
func (r *SQLAdvisorRepository) FindAll() ([]*model.Advisor, error) {
	query := `SELECT s_id, i_id FROM advisor`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error querying advisors: %w", err)
	}
	defer rows.Close()

	var advisors []*model.Advisor
	for rows.Next() {
		var advisor model.Advisor
		err := rows.Scan(
			&advisor.StudentID,
			&advisor.InstructorID,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning advisor: %w", err)
		}
		advisors = append(advisors, &advisor)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating advisors: %w", err)
	}

	return advisors, nil
}
