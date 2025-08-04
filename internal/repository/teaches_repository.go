package repository

import (
	"database/sql"
	"fmt"

	"github.com/yourusername/student-management-system/internal/model"
)

// TeachesRepository 定义教学关系仓库接口
type TeachesRepository interface {
	FindByInstructorID(instructorID string) ([]*model.Teaches, error)
	FindBySectionID(sectionID string) ([]*model.Teaches, error)
	FindByInstructorAndSection(instructorID, sectionID string) (*model.Teaches, error)
	Create(teaches *model.Teaches) error
	Delete(instructorID, courseID, sectionID, semester string, year int) error
	GetCurrentTeaching(instructorID string, semester string, year int) ([]*model.Teaches, error)
	FindAll() ([]*model.Teaches, error)
}

// SQLTeachesRepository 实现TeachesRepository接口
type SQLTeachesRepository struct {
	db *sql.DB
}

// NewTeachesRepository 创建教学关系仓库实例
func NewTeachesRepository(db *sql.DB) TeachesRepository {
	return &SQLTeachesRepository{db: db}
}

// FindAll 查找所有教学关系
func (r *SQLTeachesRepository) FindAll() ([]*model.Teaches, error) {
	query := `
		SELECT t.id, t.instructor_id, t.course_id, t.sec_id, t.semester, t.year,
		       c.title, c.dept_name, c.credits,
		       s.building, s.room_number, s.time_slot_id
		FROM teaches t
		JOIN course c ON t.course_id = c.course_id
		JOIN section s ON t.sec_id = s.sec_id AND t.semester = s.semester AND t.year = s.year
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error querying teaches: %w", err)
	}
	defer rows.Close()

	var teachesList []*model.Teaches
	for rows.Next() {
		var teaches model.Teaches
		var course model.Course
		var section model.Section

		err := rows.Scan(
			&teaches.ID,
			&teaches.InstructorID,
			&teaches.CourseID,
			&teaches.SectionID,
			&teaches.Semester,
			&teaches.Year,
			&course.Title,
			&course.Dept,
			&course.Credits,
			&section.Building,
			&section.RoomNumber,
			&section.TimeSlotID,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning teaches row: %w", err)
		}

		course.ID = teaches.CourseID
		teaches.Course = &course

		section.ID = teaches.SectionID
		section.CourseID = teaches.CourseID
		section.Semester = teaches.Semester
		section.Year = teaches.Year
		teaches.Section = &section

		teachesList = append(teachesList, &teaches)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating teaches rows: %w", err)
	}

	return teachesList, nil
}

// FindByInstructorID 根据教师ID查找教学关系
func (r *SQLTeachesRepository) FindByInstructorID(instructorID string) ([]*model.Teaches, error) {
	query := `
		SELECT t.id, t.instructor_id, t.course_id, t.sec_id, t.semester, t.year,
		       c.title, c.dept_name, c.credits,
		       s.building, s.room_number, s.time_slot_id
		FROM teaches t
		JOIN course c ON t.course_id = c.course_id
		JOIN section s ON t.sec_id = s.sec_id AND t.semester = s.semester AND t.year = s.year
		WHERE t.instructor_id = ?
	`

	rows, err := r.db.Query(query, instructorID)
	if err != nil {
		return nil, fmt.Errorf("error querying teaches: %w", err)
	}
	defer rows.Close()

	var teachesList []*model.Teaches
	for rows.Next() {
		var teaches model.Teaches
		var course model.Course
		var section model.Section

		err := rows.Scan(
			&teaches.ID,
			&teaches.InstructorID,
			&teaches.CourseID,
			&teaches.SectionID,
			&teaches.Semester,
			&teaches.Year,
			&course.Title,
			&course.Dept,
			&course.Credits,
			&section.Building,
			&section.RoomNumber,
			&section.TimeSlotID,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning teaches: %w", err)
		}

		course.ID = teaches.CourseID
		section.ID = teaches.SectionID
		section.CourseID = teaches.CourseID
		section.Semester = teaches.Semester
		section.Year = teaches.Year

		teaches.Course = &course
		teaches.Section = &section
		teachesList = append(teachesList, &teaches)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating teaches: %w", err)
	}

	return teachesList, nil
}

// FindBySectionID 根据课程段ID查找教学关系
func (r *SQLTeachesRepository) FindBySectionID(sectionID string) ([]*model.Teaches, error) {
	query := `
		SELECT t.id, t.instructor_id, t.course_id, t.sec_id, t.semester, t.year,
		       i.name, i.dept_name, i.salary
		FROM teaches t
		JOIN instructor i ON t.instructor_id = i.id
		WHERE t.sec_id = ?
	`

	rows, err := r.db.Query(query, sectionID)
	if err != nil {
		return nil, fmt.Errorf("error querying teaches: %w", err)
	}
	defer rows.Close()

	var teachesList []*model.Teaches
	for rows.Next() {
		var teaches model.Teaches
		var instructor model.Instructor

		err := rows.Scan(
			&teaches.ID,
			&teaches.InstructorID,
			&teaches.CourseID,
			&teaches.SectionID,
			&teaches.Semester,
			&teaches.Year,
			&instructor.Name,
			&instructor.Dept,
			&instructor.Salary,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning teaches: %w", err)
		}

		instructor.ID = teaches.InstructorID
		teaches.Instructor = &instructor
		teachesList = append(teachesList, &teaches)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating teaches: %w", err)
	}

	return teachesList, nil
}

// FindByInstructorAndSection 根据教师ID和课程段ID查找教学关系
func (r *SQLTeachesRepository) FindByInstructorAndSection(instructorID, sectionID string) (*model.Teaches, error) {
	query := `SELECT id, instructor_id, course_id, sec_id, semester, year FROM teaches WHERE instructor_id = ? AND sec_id = ?`

	var teaches model.Teaches
	err := r.db.QueryRow(query, instructorID, sectionID).Scan(
		&teaches.ID,
		&teaches.InstructorID,
		&teaches.CourseID,
		&teaches.SectionID,
		&teaches.Semester,
		&teaches.Year,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("teaches record not found")
		}
		return nil, fmt.Errorf("error querying teaches: %w", err)
	}

	return &teaches, nil
}

// Create 创建教学关系
func (r *SQLTeachesRepository) Create(teaches *model.Teaches) error {
	query := `INSERT INTO teaches (instructor_id, course_id, sec_id, semester, year) VALUES (?, ?, ?, ?, ?)`

	_, err := r.db.Exec(query,
		teaches.InstructorID,
		teaches.CourseID,
		teaches.SectionID,
		teaches.Semester,
		teaches.Year,
	)

	if err != nil {
		return fmt.Errorf("error creating teaches: %w", err)
	}

	return nil
}

// Delete 删除教学安排
func (r *SQLTeachesRepository) Delete(instructorID, courseID, sectionID, semester string, year int) error {
	query := `DELETE FROM teaches WHERE id = ? AND course_id = ? AND sec_id = ? AND semester = ? AND year = ?`

	result, err := r.db.Exec(query, instructorID, courseID, sectionID, semester, year)
	if err != nil {
		return fmt.Errorf("error deleting teaches: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error getting rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("teaches record not found")
	}

	return nil
}

// GetCurrentTeaching 获取教师当前学期的教学任务
func (r *SQLTeachesRepository) GetCurrentTeaching(instructorID string, semester string, year int) ([]*model.Teaches, error) {
	query := `
		SELECT t.id, t.instructor_id, t.course_id, t.sec_id, t.semester, t.year,
		       c.title, c.dept_name, c.credits,
		       s.building, s.room_number, s.time_slot_id,
		       ts.day, ts.start_hr, ts.start_min, ts.end_hr, ts.end_min,
		       (SELECT COUNT(*) FROM takes tk WHERE tk.sec_id = t.sec_id AND tk.semester = t.semester AND tk.year = t.year) as enrollment
		FROM teaches t
		JOIN course c ON t.course_id = c.course_id
		JOIN section s ON t.sec_id = s.sec_id AND t.semester = s.semester AND t.year = s.year
		JOIN time_slot ts ON s.time_slot_id = ts.time_slot_id
		WHERE t.instructor_id = ? AND t.semester = ? AND t.year = ?
	`
	rows, err := r.db.Query(query, instructorID, semester, year)
	if err != nil {
		return nil, fmt.Errorf("error querying current teaching: %w", err)
	}
	defer rows.Close()

	var teachesList []*model.Teaches
	teachesMap := make(map[string]*model.Teaches)

	for rows.Next() {
		var teaches model.Teaches
		var course model.Course
		var section model.Section
		var dayStr string
		var startHr, startMin, endHr, endMin int
		var enrollment int

		err := rows.Scan(
			&teaches.ID,
			&teaches.InstructorID,
			&teaches.CourseID,
			&teaches.SectionID,
			&teaches.Semester,
			&teaches.Year,
			&course.Title,
			&course.Dept,
			&course.Credits,
			&section.Building,
			&section.RoomNumber,
			&section.TimeSlotID,
			&dayStr,
			&startHr,
			&startMin,
			&endHr,
			&endMin,
			&enrollment,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning teaches: %w", err)
		}

		course.ID = teaches.CourseID
		section.ID = teaches.SectionID
		section.CourseID = teaches.CourseID
		section.Semester = teaches.Semester
		section.Year = teaches.Year
		section.Enrollment = enrollment

		// 将日期字符串转换为整数（1-7表示周一到周日）
		var day int
		switch dayStr {
		case "Monday":
			day = 1
		case "Tuesday":
			day = 2
		case "Wednesday":
			day = 3
		case "Thursday":
			day = 4
		case "Friday":
			day = 5
		case "Saturday":
			day = 6
		case "Sunday":
			day = 7
		}

		// 处理时间段
		timeSlot := &model.TimeSlot{
			ID:       section.TimeSlotID,
			Days:     []int{day},
			StartHr:  startHr,
			StartMin: startMin,
			EndHr:    endHr,
			EndMin:   endMin,
		}

		// 使用map合并相同课程的不同时间段
		if _, ok := teachesMap[teaches.SectionID]; !ok {
			// 新课程
			section.TimeSlot = timeSlot
			section.Course = &course
			teaches.Course = &course
			teaches.Section = &section
			teachesMap[teaches.SectionID] = &teaches
			teachesList = append(teachesList, &teaches)
		}
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating teaches: %w", err)
	}

	return teachesList, nil
}
