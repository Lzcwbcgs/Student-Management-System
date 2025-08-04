package repository

import (
	"database/sql"
	"fmt"

	"github.com/yourusername/student-management-system/internal/model"
)

// TakesRepository 定义学生选课仓库接口
type TakesRepository interface {
	FindByStudentID(studentID string) ([]*model.Takes, error)
	FindByStudentAndSection(studentID, sectionID string) (*model.Takes, error)
	FindBySection(sectionID string) ([]*model.Takes, error)
	FindBySectionID(sectionID string) ([]*model.Takes, error)
	Create(takes *model.Takes) error
	Delete(studentID, sectionID string) error
	UpdateGrade(studentID, sectionID, grade string) error
	GetStudentTranscript(studentID string) (*model.Transcript, error)
	GetCurrentCourses(studentID string, semester string, year int) ([]*model.Takes, error)
	CheckTimeConflict(studentID, sectionID string) (bool, error)
}

// SQLTakesRepository 实现TakesRepository接口
type SQLTakesRepository struct {
	db *sql.DB
}

// NewTakesRepository 创建学生选课仓库实例
func NewTakesRepository(db *sql.DB) TakesRepository {
	return &SQLTakesRepository{db: db}
}

// FindByStudentID 根据学生ID查找选课记录
func (r *SQLTakesRepository) FindByStudentID(studentID string) ([]*model.Takes, error) {
	query := `
		SELECT t.student_id, t.course_id, t.sec_id, t.semester, t.year, t.grade,
		       c.title, c.dept_name, c.credits,
		       s.building, s.room_number, s.time_slot_id
		FROM takes t
		JOIN course c ON t.course_id = c.course_id
		JOIN section s ON t.sec_id = s.sec_id AND t.semester = s.semester AND t.year = s.year
		WHERE t.student_id = ?
	`

	rows, err := r.db.Query(query, studentID)
	if err != nil {
		return nil, fmt.Errorf("error querying takes: %w", err)
	}
	defer rows.Close()

	var takesList []*model.Takes
	for rows.Next() {
		var takes model.Takes
		var course model.Course
		var section model.Section

		err := rows.Scan(
			&takes.StudentID,
			&takes.CourseID,
			&takes.SectionID,
			&takes.Semester,
			&takes.Year,
			&takes.Grade,
			&course.Title,
			&course.Dept,
			&course.Credits,
			&section.Building,
			&section.RoomNumber,
			&section.TimeSlotID,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning takes: %w", err)
		}

		course.ID = takes.CourseID
		section.ID = takes.SectionID
		section.CourseID = takes.CourseID
		section.Semester = takes.Semester
		section.Year = takes.Year

		takes.Course = &course
		takes.Section = &section
		takesList = append(takesList, &takes)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating takes: %w", err)
	}

	return takesList, nil
}

// FindByStudentAndSection 根据学生ID和课程段ID查找选课记录
func (r *SQLTakesRepository) FindByStudentAndSection(studentID, sectionID string) (*model.Takes, error) {
	query := `SELECT student_id, course_id, sec_id, semester, year, grade FROM takes WHERE student_id = ? AND sec_id = ?`

	var takes model.Takes
	err := r.db.QueryRow(query, studentID, sectionID).Scan(
		&takes.StudentID,
		&takes.CourseID,
		&takes.SectionID,
		&takes.Semester,
		&takes.Year,
		&takes.Grade,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("takes record not found")
		}
		return nil, fmt.Errorf("error querying takes: %w", err)
	}

	return &takes, nil
}

// FindBySection 根据课程段ID查找所有选课记录
func (r *SQLTakesRepository) FindBySection(sectionID string) ([]*model.Takes, error) {
	query := `
		SELECT t.student_id, t.course_id, t.sec_id, t.semester, t.year, t.grade,
		       s.name, s.dept_name, s.tot_cred
		FROM takes t
		JOIN student s ON t.student_id = s.id
		WHERE t.sec_id = ?
	`

	rows, err := r.db.Query(query, sectionID)
	if err != nil {
		return nil, fmt.Errorf("error querying takes: %w", err)
	}
	defer rows.Close()

	var takesList []*model.Takes
	for rows.Next() {
		var takes model.Takes
		var student model.Student

		err := rows.Scan(
			&takes.StudentID,
			&takes.CourseID,
			&takes.SectionID,
			&takes.Semester,
			&takes.Year,
			&takes.Grade,
			&student.Name,
			&student.Dept,
			&student.TotCred,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning takes: %w", err)
		}

		student.ID = takes.StudentID
		takes.Student = &student
		takesList = append(takesList, &takes)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating takes: %w", err)
	}

	return takesList, nil
}

// Create 创建选课记录
func (r *SQLTakesRepository) Create(takes *model.Takes) error {
	query := `INSERT INTO takes (student_id, course_id, sec_id, semester, year, grade) VALUES (?, ?, ?, ?, ?, ?)`

	_, err := r.db.Exec(query,
		takes.StudentID,
		takes.CourseID,
		takes.SectionID,
		takes.Semester,
		takes.Year,
		takes.Grade,
	)

	if err != nil {
		return fmt.Errorf("error creating takes: %w", err)
	}

	return nil
}

// Delete 删除选课记录
func (r *SQLTakesRepository) Delete(studentID, sectionID string) error {
	query := `DELETE FROM takes WHERE student_id = ? AND sec_id = ?`

	result, err := r.db.Exec(query, studentID, sectionID)
	if err != nil {
		return fmt.Errorf("error deleting takes: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error getting rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("takes record not found")
	}

	return nil
}

// UpdateGrade 更新成绩
func (r *SQLTakesRepository) UpdateGrade(studentID, sectionID, grade string) error {
	query := `UPDATE takes SET grade = ? WHERE student_id = ? AND sec_id = ?`

	result, err := r.db.Exec(query, grade, studentID, sectionID)
	if err != nil {
		return fmt.Errorf("error updating grade: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error getting rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("takes record not found")
	}

	return nil
}

// GetStudentTranscript 获取学生成绩单
func (r *SQLTakesRepository) GetStudentTranscript(studentID string) (*model.Transcript, error) {
	// 首先获取学生信息
	studentQuery := `SELECT id, name, dept_name, tot_cred FROM student WHERE id = ?`
	var student model.Student
	err := r.db.QueryRow(studentQuery, studentID).Scan(
		&student.ID,
		&student.Name,
		&student.Dept,
		&student.TotCred,
	)
	if err != nil {
		return nil, fmt.Errorf("error getting student info: %w", err)
	}

	// 获取学生所有课程成绩
	query := `
		SELECT t.course_id, t.sec_id, t.semester, t.year, t.grade,
		       c.title, c.dept_name, c.credits
		FROM takes t
		JOIN course c ON t.course_id = c.course_id
		WHERE t.student_id = ?
		ORDER BY t.year DESC, t.semester DESC
	`

	rows, err := r.db.Query(query, studentID)
	if err != nil {
		return nil, fmt.Errorf("error querying transcript: %w", err)
	}
	defer rows.Close()

	var transcript model.Transcript
	transcript.Student = *student.ToDTO() // 使用现有的StudentDTO
	var totalCredits float64
	var totalGradePoints float64

	for rows.Next() {
		var courseGrade model.CourseGrade
		var semester, year string
		var grade sql.NullString

		err := rows.Scan(
			&courseGrade.CourseID,
			&semester,
			&year,
			&grade,
			&courseGrade.Title,
			&courseGrade.Semester,
			&courseGrade.Year,
			&courseGrade.Credits,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning transcript: %w", err)
		}

		gradeStr := ""
		if grade.Valid {
			gradeStr = grade.String
		}
		courseGrade.Grade = gradeStr

		// 计算绩点
		courseGrade.GradePoint = calculateGradePoint(gradeStr)

		transcript.Courses = append(transcript.Courses, courseGrade)

		if gradeStr != "" && gradeStr != "F" {
			totalCredits += courseGrade.Credits
			totalGradePoints += courseGrade.Credits * courseGrade.GradePoint
		}
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating transcript: %w", err)
	}

	transcript.TotalCred = totalCredits
	if totalCredits > 0 {
		transcript.GPA = totalGradePoints / totalCredits
	}

	return &transcript, nil
}

// GetCurrentCourses 获取学生当前学期的课程
func (r *SQLTakesRepository) GetCurrentCourses(studentID string, semester string, year int) ([]*model.Takes, error) {
	query := `
		SELECT t.student_id, t.course_id, t.sec_id, t.semester, t.year, t.grade,
		       c.title, c.dept_name, c.credits,
		       s.building, s.room_number, s.time_slot_id,
		       ts.day, ts.start_hr, ts.start_min, ts.end_hr, ts.end_min
		FROM takes t
		JOIN course c ON t.course_id = c.course_id
		JOIN section s ON t.sec_id = s.sec_id AND t.semester = s.semester AND t.year = s.year
		JOIN time_slot ts ON s.time_slot_id = ts.time_slot_id
		WHERE t.student_id = ? AND t.semester = ? AND t.year = ?
	`

	rows, err := r.db.Query(query, studentID, semester, year)
	if err != nil {
		return nil, fmt.Errorf("error querying current courses: %w", err)
	}
	defer rows.Close()

	var takesList []*model.Takes
	takesMap := make(map[string]*model.Takes)

	for rows.Next() {
		var takes model.Takes
		var course model.Course
		var section model.Section
		var day string
		var startHr, startMin, endHr, endMin int

		err := rows.Scan(
			&takes.StudentID,
			&takes.CourseID,
			&takes.SectionID,
			&takes.Semester,
			&takes.Year,
			&takes.Grade,
			&course.Title,
			&course.Dept,
			&course.Credits,
			&section.Building,
			&section.RoomNumber,
			&section.TimeSlotID,
			&day,
			&startHr,
			&startMin,
			&endHr,
			&endMin,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning takes: %w", err)
		}

		course.ID = takes.CourseID
		section.ID = takes.SectionID
		section.CourseID = takes.CourseID
		section.Semester = takes.Semester
		section.Year = takes.Year

		// 创建时间段
		dayNum := convertDayToInt(day)
		timeSlot := &model.TimeSlot{
			ID:       section.TimeSlotID,
			Days:     []int{dayNum},
			StartHr:  startHr,
			StartMin: startMin,
			EndHr:    endHr,
			EndMin:   endMin,
		}

		// 使用map合并相同课程的不同时间段
		if _, ok := takesMap[takes.SectionID]; !ok {
			// 新课程
			section.TimeSlot = timeSlot
			section.Course = &course
			takes.Course = &course
			takes.Section = &section
			takesMap[takes.SectionID] = &takes
			takesList = append(takesList, &takes)
		}
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating takes: %w", err)
	}

	return takesList, nil
}

// FindBySectionID 根据课程段ID查找选课记录
func (r *SQLTakesRepository) FindBySectionID(sectionID string) ([]*model.Takes, error) {
	return r.FindBySection(sectionID)
}

// CheckTimeConflict 检查时间冲突
func (r *SQLTakesRepository) CheckTimeConflict(studentID, sectionID string) (bool, error) {
	// 获取要选的课程的时间段
	timeSlotQuery := `
		SELECT ts.time_slot_id, ts.day, ts.start_hr, ts.start_min, ts.end_hr, ts.end_min
		FROM section s
		JOIN time_slot ts ON s.time_slot_id = ts.time_slot_id
		WHERE s.sec_id = ?
	`
	timeSlotRows, err := r.db.Query(timeSlotQuery, sectionID)
	if err != nil {
		return false, fmt.Errorf("error querying time slot: %w", err)
	}
	defer timeSlotRows.Close()

	type TimeSlotInfo struct {
		Day      string
		StartHr  int
		StartMin int
		EndHr    int
		EndMin   int
	}

	var newTimeSlots []TimeSlotInfo
	for timeSlotRows.Next() {
		var timeSlotID, day string
		var startHr, startMin, endHr, endMin int
		err := timeSlotRows.Scan(&timeSlotID, &day, &startHr, &startMin, &endHr, &endMin)
		if err != nil {
			return false, fmt.Errorf("error scanning time slot: %w", err)
		}
		newTimeSlots = append(newTimeSlots, TimeSlotInfo{Day: day, StartHr: startHr, StartMin: startMin, EndHr: endHr, EndMin: endMin})
	}

	if err := timeSlotRows.Err(); err != nil {
		return false, fmt.Errorf("error iterating time slots: %w", err)
	}

	// 获取学生当前学期已选课程的时间段
	currentCoursesQuery := `
		SELECT ts.day, ts.start_hr, ts.start_min, ts.end_hr, ts.end_min
		FROM takes t
		JOIN section s ON t.sec_id = s.sec_id AND t.semester = s.semester AND t.year = s.year
		JOIN time_slot ts ON s.time_slot_id = ts.time_slot_id
		WHERE t.student_id = ? AND t.sec_id <> ?
	`
	currentCoursesRows, err := r.db.Query(currentCoursesQuery, studentID, sectionID)
	if err != nil {
		return false, fmt.Errorf("error querying current courses: %w", err)
	}
	defer currentCoursesRows.Close()

	var currentTimeSlots []TimeSlotInfo
	for currentCoursesRows.Next() {
		var day string
		var startHr, startMin, endHr, endMin int
		err := currentCoursesRows.Scan(&day, &startHr, &startMin, &endHr, &endMin)
		if err != nil {
			return false, fmt.Errorf("error scanning current course: %w", err)
		}
		currentTimeSlots = append(currentTimeSlots, TimeSlotInfo{Day: day, StartHr: startHr, StartMin: startMin, EndHr: endHr, EndMin: endMin})
	}

	if err := currentCoursesRows.Err(); err != nil {
		return false, fmt.Errorf("error iterating current courses: %w", err)
	}

	// 检查时间冲突
	for _, newSlot := range newTimeSlots {
		for _, currentSlot := range currentTimeSlots {
			if newSlot.Day == currentSlot.Day {
				// 检查时间重叠
				newStart := newSlot.StartHr*60 + newSlot.StartMin
				newEnd := newSlot.EndHr*60 + newSlot.EndMin
				currentStart := currentSlot.StartHr*60 + currentSlot.StartMin
				currentEnd := currentSlot.EndHr*60 + currentSlot.EndMin

				if (newStart < currentEnd) && (newEnd > currentStart) {
					return true, nil // 发现时间冲突
				}
			}
		}
	}

	return false, nil // 没有时间冲突
}

// convertDayToInt 将星期转换为数字
func convertDayToInt(day string) int {
	switch day {
	case "Monday":
		return 1
	case "Tuesday":
		return 2
	case "Wednesday":
		return 3
	case "Thursday":
		return 4
	case "Friday":
		return 5
	case "Saturday":
		return 6
	case "Sunday":
		return 7
	default:
		return 0
	}
}

// calculateGradePoint 计算绩点
func calculateGradePoint(grade string) float64 {
	switch grade {
	case "A":
		return 4.0
	case "A-":
		return 3.7
	case "B+":
		return 3.3
	case "B":
		return 3.0
	case "B-":
		return 2.7
	case "C+":
		return 2.3
	case "C":
		return 2.0
	case "C-":
		return 1.7
	case "D+":
		return 1.3
	case "D":
		return 1.0
	case "F":
		return 0.0
	default:
		return 0.0
	}
}
