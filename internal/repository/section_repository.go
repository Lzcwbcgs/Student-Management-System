package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"github.com/yourusername/student-management-system/internal/model"
)

// SectionRepository 定义课程章节仓库接口
type SectionRepository interface {
	FindByID(id string) (*model.Section, error)
	FindAll() ([]*model.Section, error)
	FindByCourseID(courseID string) ([]*model.Section, error)
	FindByParams(params *model.SectionQueryParams) ([]*model.Section, error)
	Create(section *model.Section) error
	Update(section *model.Section) error
	Delete(id string) error
	GetEnrollmentCount(sectionID string) (int, error)
	FindWithDetails(id string) (*model.Section, error)
	GetSectionClassroom(sectionID string) (*model.Classroom, error)
}

// SQLSectionRepository 实现SectionRepository接口
type SQLSectionRepository struct {
	db *sql.DB
}

// NewSectionRepository 创建课程章节仓库实例
func NewSectionRepository(db *sql.DB) SectionRepository {
	return &SQLSectionRepository{db: db}
}

// FindByID 根据ID查找课程章节
func (r *SQLSectionRepository) FindByID(id string) (*model.Section, error) {
	query := `SELECT course_id, sec_id, semester, year, building, room_number, time_slot_id FROM section WHERE sec_id = ?`

	var section model.Section
	err := r.db.QueryRow(query, id).Scan(
		&section.CourseID,
		&section.ID,
		&section.Semester,
		&section.Year,
		&section.Building,
		&section.RoomNumber,
		&section.TimeSlotID,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("section not found")
		}
		return nil, fmt.Errorf("error querying section: %w", err)
	}

	return &section, nil // 返回指针类型
}

// GetSectionClassroom 获取课程章节的教室信息
func (r *SQLSectionRepository) GetSectionClassroom(sectionID string) (*model.Classroom, error) {
	query := `
		SELECT c.building, c.room_number, c.capacity 
		FROM section s 
		JOIN classroom c ON s.building = c.building AND s.room_number = c.room_number 
		WHERE s.sec_id = ?
	`

	var classroom model.Classroom
	err := r.db.QueryRow(query, sectionID).Scan(
		&classroom.Building,
		&classroom.RoomNumber,
		&classroom.Capacity,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("classroom not found for section: %w", err)
		}
		return nil, fmt.Errorf("error getting classroom info: %w", err)
	}

	return &classroom, nil
}

// FindAll 查找所有课程章节
func (r *SQLSectionRepository) FindAll() ([]*model.Section, error) {
	query := `SELECT course_id, sec_id, semester, year, building, room_number, time_slot_id FROM section`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error querying sections: %w", err)
	}
	defer rows.Close()

	var sections []*model.Section
	for rows.Next() {
		var section model.Section
		err := rows.Scan(
			&section.CourseID,
			&section.ID,
			&section.Semester,
			&section.Year,
			&section.Building,
			&section.RoomNumber,
			&section.TimeSlotID,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning section: %w", err)
		}
		sections = append(sections, &section)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating sections: %w", err)
	}

	return sections, nil
}

// FindByCourseID 根据课程ID查找课程章节
func (r *SQLSectionRepository) FindByCourseID(courseID string) ([]*model.Section, error) {
	query := `SELECT course_id, sec_id, semester, year, building, room_number, time_slot_id FROM section WHERE course_id = ?`

	rows, err := r.db.Query(query, courseID)
	if err != nil {
		return nil, fmt.Errorf("error querying sections: %w", err)
	}
	defer rows.Close()

	var sections []*model.Section
	for rows.Next() {
		var section model.Section
		err := rows.Scan(
			&section.CourseID,
			&section.ID,
			&section.Semester,
			&section.Year,
			&section.Building,
			&section.RoomNumber,
			&section.TimeSlotID,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning section: %w", err)
		}
		sections = append(sections, &section)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating sections: %w", err)
	}

	return sections, nil
}

// FindByParams 根据参数查找课程章节
func (r *SQLSectionRepository) FindByParams(params *model.SectionQueryParams) ([]*model.Section, error) {
	query := `SELECT s.course_id, s.sec_id, s.semester, s.year, s.building, s.room_number, s.time_slot_id 
			  FROM section s 
			  JOIN course c ON s.course_id = c.course_id 
			  WHERE 1=1`

	var args []interface{}

	if params.CourseID != "" {
		query += ` AND s.course_id = ?`
		args = append(args, params.CourseID)
	}

	if params.Semester != "" {
		query += ` AND s.semester = ?`
		args = append(args, params.Semester)
	}

	if params.Year != 0 {
		query += ` AND s.year = ?`
		args = append(args, params.Year)
	}

	if params.Dept != "" {
		query += ` AND c.dept_name = ?`
		args = append(args, params.Dept)
	}

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("error querying sections: %w", err)
	}
	defer rows.Close()

	var sections []*model.Section
	for rows.Next() {
		var section model.Section
		err := rows.Scan(
			&section.CourseID,
			&section.ID,
			&section.Semester,
			&section.Year,
			&section.Building,
			&section.RoomNumber,
			&section.TimeSlotID,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning section: %w", err)
		}
		sections = append(sections, &section)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating sections: %w", err)
	}

	return sections, nil
}

// Create 创建课程章节
func (r *SQLSectionRepository) Create(section *model.Section) error {
	query := `INSERT INTO section (course_id, sec_id, semester, year, building, room_number, time_slot_id) VALUES (?, ?, ?, ?, ?, ?, ?)`

	_, err := r.db.Exec(query,
		section.CourseID,
		section.ID,
		section.Semester,
		section.Year,
		section.Building,
		section.RoomNumber,
		section.TimeSlotID,
	)

	if err != nil {
		return fmt.Errorf("error creating section: %w", err)
	}

	return nil
}

// Update 更新课程章节
func (r *SQLSectionRepository) Update(section *model.Section) error {
	query := `UPDATE section SET semester = ?, year = ?, building = ?, room_number = ?, time_slot_id = ? WHERE course_id = ? AND sec_id = ?`

	result, err := r.db.Exec(query,
		section.Semester,
		section.Year,
		section.Building,
		section.RoomNumber,
		section.TimeSlotID,
		section.CourseID,
		section.ID,
	)

	if err != nil {
		return fmt.Errorf("error updating section: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error getting rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("section not found")
	}

	return nil
}

// Delete 删除课程章节
func (r *SQLSectionRepository) Delete(id string) error {
	query := `DELETE FROM section WHERE sec_id = ?`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("error deleting section: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error getting rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("section not found")
	}

	return nil
}

// GetEnrollmentCount 获取课程章节的选课人数
func (r *SQLSectionRepository) GetEnrollmentCount(sectionID string) (int, error) {
	query := `SELECT COUNT(*) FROM takes WHERE sec_id = ?`

	var count int
	err := r.db.QueryRow(query, sectionID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("error getting enrollment count: %w", err)
	}

	return count, nil
}

// FindWithDetails 查找课程章节详细信息
func (r *SQLSectionRepository) FindWithDetails(id string) (*model.Section, error) {
	// 首先查找基本信息
	section, err := r.FindByID(id)
	if err != nil {
		return nil, err
	}

	// 查找课程信息
	courseQuery := `SELECT course_id, title, dept_name, credits FROM course WHERE course_id = ?`
	var course model.Course
	err = r.db.QueryRow(courseQuery, section.CourseID).Scan(
		&course.ID,
		&course.Title,
		&course.Dept,
		&course.Credits,
	)
	if err != nil {
		return nil, fmt.Errorf("error getting course info: %w", err)
	}
	section.Course = &course

	// 查找教室信息
	classroomQuery := `SELECT building, room_number, capacity FROM classroom WHERE building = ? AND room_number = ?`
	var classroom model.Classroom
	err = r.db.QueryRow(classroomQuery, section.Building, section.RoomNumber).Scan(
		&classroom.Building,
		&classroom.RoomNumber,
		&classroom.Capacity,
	)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("error getting classroom info: %w", err)
	}
	if !errors.Is(err, sql.ErrNoRows) {
		section.Classroom = &classroom
	}

	// 查找时间段信息
	timeSlotQuery := `SELECT time_slot_id, day, start_hr, start_min, end_hr, end_min FROM time_slot WHERE time_slot_id = ?`
	var timeSlot model.TimeSlot
	var dayStr string
	err = r.db.QueryRow(timeSlotQuery, section.TimeSlotID).Scan(
		&timeSlot.ID,
		&dayStr,
		&timeSlot.StartHr,
		&timeSlot.StartMin,
		&timeSlot.EndHr,
		&timeSlot.EndMin,
	)
	
	// 将dayStr转换为int并添加到Days数组
	if dayStr != "" {
		day, err := strconv.Atoi(dayStr)
		if err == nil {
			timeSlot.Days = append(timeSlot.Days, day)
		}
	}
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("error getting time slot info: %w", err)
	}
	if !errors.Is(err, sql.ErrNoRows) {
		section.TimeSlot = &timeSlot
	}

	// 查找授课教师
	instructorQuery := `
		SELECT i.id, i.name, i.dept_name, i.salary 
		FROM teaches t 
		JOIN instructor i ON t.id = i.id 
		WHERE t.sec_id = ? AND t.semester = ? AND t.year = ?
	`
	instructorRows, err := r.db.Query(instructorQuery, section.ID, section.Semester, section.Year)
	if err != nil {
		return nil, fmt.Errorf("error querying instructors: %w", err)
	}
	defer instructorRows.Close()

	var instructors []model.Instructor
	for instructorRows.Next() {
		var instructor model.Instructor
		err := instructorRows.Scan(
			&instructor.ID,
			&instructor.Name,
			&instructor.Dept,
			&instructor.Salary,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning instructor: %w", err)
		}
		instructors = append(instructors, instructor)
	}

	if err := instructorRows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating instructors: %w", err)
	}

	section.Instructors = instructors

	return section, nil
}
