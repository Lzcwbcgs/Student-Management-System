package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/yourusername/student-management-system/internal/model"
)

// ClassroomRepository 定义教室仓储接口
type ClassroomRepository interface {
	FindByID(building, roomNumber string) (*model.Classroom, error)
	FindAll() ([]*model.Classroom, error)
	FindByBuildingAndRoom(building, roomNumber string) (*model.Classroom, error)
	FindByBuilding(building string) ([]*model.Classroom, error)
	Create(classroom *model.Classroom) error
	Update(classroom *model.Classroom) error
	Delete(building, roomNumber string) error
	FindAvailable(capacity int, semester string, year int, timeSlotID string) ([]*model.Classroom, error)
}

// SQLClassroomRepository 实现ClassroomRepository接口
type SQLClassroomRepository struct {
	db *sql.DB
}

// NewClassroomRepository 创建教室仓储实例
func NewClassroomRepository(db *sql.DB) ClassroomRepository {
	return &SQLClassroomRepository{db: db}
}

// FindByBuilding 根据教学楼查找教室
func (r *SQLClassroomRepository) FindByBuilding(building string) ([]*model.Classroom, error) {
	query := `SELECT building, room_number, capacity FROM classroom WHERE building = ?`

	rows, err := r.db.Query(query, building)
	if err != nil {
		return nil, fmt.Errorf("error querying classrooms by building: %w", err)
	}
	defer rows.Close()

	var classrooms []*model.Classroom
	for rows.Next() {
		var classroom model.Classroom
		err := rows.Scan(
			&classroom.Building,
			&classroom.RoomNumber,
			&classroom.Capacity,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning classroom: %w", err)
		}
		classrooms = append(classrooms, &classroom)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating classrooms: %w", err)
	}

	return classrooms, nil
}

// FindByID 根据教学楼和教室号查找教室
func (r *SQLClassroomRepository) FindByID(building, roomNumber string) (*model.Classroom, error) {
	query := `SELECT building, room_number, capacity FROM classroom WHERE building = ? AND room_number = ?`
	row := r.db.QueryRow(query, building, roomNumber)

	var classroom model.Classroom
	err := row.Scan(&classroom.Building, &classroom.RoomNumber, &classroom.Capacity)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("classroom not found: %w", err)
		}
		return nil, fmt.Errorf("error scanning classroom: %w", err)
	}

	return &classroom, nil
}

// FindAll 查找所有教室
func (r *SQLClassroomRepository) FindAll() ([]*model.Classroom, error) {
	query := `SELECT building, room_number, capacity FROM classroom`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error querying classrooms: %w", err)
	}
	defer rows.Close()

	var classrooms []*model.Classroom
	for rows.Next() {
		var classroom model.Classroom
		err := rows.Scan(&classroom.Building, &classroom.RoomNumber, &classroom.Capacity)
		if err != nil {
			return nil, fmt.Errorf("error scanning classroom: %w", err)
		}
		classrooms = append(classrooms, &classroom)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating classrooms: %w", err)
	}

	return classrooms, nil
}

// FindByBuilding 根据教学楼查找教室
//func (r *SQLClassroomRepository) FindByBuilding(building string) ([]*model.Classroom, error) {
//	query := `SELECT building, room_number, capacity FROM classroom WHERE building = ? ORDER BY room_number`
//	rows, err := r.db.Query(query, building)
//	if err != nil {
//		return nil, fmt.Errorf("error querying classrooms: %w", err)
//	}
//	defer rows.Close()
//
//	var classrooms []*model.Classroom
//	for rows.Next() {
//		var classroom model.Classroom
//		err := rows.Scan(&classroom.Building, &classroom.RoomNumber, &classroom.Capacity)
//		if err != nil {
//			return nil, fmt.Errorf("error scanning classroom: %w", err)
//		}
//		classrooms = append(classrooms, &classroom)
//	}
//
//	if err := rows.Err(); err != nil {
//		return nil, fmt.Errorf("error iterating classrooms: %w", err)
//	}
//
//	return classrooms, nil
//}

// Create 创建教室
func (r *SQLClassroomRepository) Create(classroom *model.Classroom) error {
	query := `INSERT INTO classroom (building, room_number, capacity) VALUES (?, ?, ?)`
	_, err := r.db.Exec(query, classroom.Building, classroom.RoomNumber, classroom.Capacity)
	if err != nil {
		if strings.Contains(err.Error(), "Duplicate entry") {
			return fmt.Errorf("classroom already exists: %w", err)
		}
		return fmt.Errorf("error creating classroom: %w", err)
	}

	return nil
}

// Update 更新教室
func (r *SQLClassroomRepository) Update(classroom *model.Classroom) error {
	query := `UPDATE classroom SET capacity = ? WHERE building = ? AND room_number = ?`
	result, err := r.db.Exec(query, classroom.Capacity, classroom.Building, classroom.RoomNumber)
	if err != nil {
		return fmt.Errorf("error updating classroom: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error getting rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("classroom not found")
	}

	return nil
}

// Delete 删除教室
func (r *SQLClassroomRepository) Delete(building, roomNumber string) error {
	// 检查教室是否被使用
	checkQuery := `
		SELECT COUNT(*) FROM section 
		WHERE building = ? AND room_number = ?
	`
	var count int
	err := r.db.QueryRow(checkQuery, building, roomNumber).Scan(&count)
	if err != nil {
		return fmt.Errorf("error checking classroom usage: %w", err)
	}

	if count > 0 {
		return fmt.Errorf("cannot delete classroom: it is being used by %d sections", count)
	}

	// 删除教室
	query := `DELETE FROM classroom WHERE building = ? AND room_number = ?`
	result, err := r.db.Exec(query, building, roomNumber)
	if err != nil {
		return fmt.Errorf("error deleting classroom: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error getting rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("classroom not found")
	}

	return nil
}

// FindAvailable 查找可用教室
func (r *SQLClassroomRepository) FindAvailable(capacity int, semester string, year int, timeSlotID string) ([]*model.Classroom, error) {
	// 查询在指定时间段没有被占用且容量满足要求的教室
	query := `
		SELECT c.building, c.room_number, c.capacity
		FROM classroom c
		WHERE c.capacity >= ?
		AND NOT EXISTS (
			SELECT 1 FROM section s
			WHERE s.building = c.building
			AND s.room_number = c.room_number
			AND s.semester = ?
			AND s.year = ?
			AND s.time_slot_id = ?
		)
		ORDER BY c.capacity
	`
	rows, err := r.db.Query(query, capacity, semester, year, timeSlotID)
	if err != nil {
		return nil, fmt.Errorf("error querying available classrooms: %w", err)
	}
	defer rows.Close()

	var classrooms []*model.Classroom
	for rows.Next() {
		var classroom model.Classroom
		err := rows.Scan(&classroom.Building, &classroom.RoomNumber, &classroom.Capacity)
		if err != nil {
			return nil, fmt.Errorf("error scanning classroom: %w", err)
		}
		classrooms = append(classrooms, &classroom)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating classrooms: %w", err)
	}

	return classrooms, nil
}

// FindByBuildingAndRoom 根据建筑和房间号查找教室
func (r *SQLClassroomRepository) FindByBuildingAndRoom(building string, roomNumber string) (*model.Classroom, error) {
	query := `SELECT building, room_number, capacity FROM classroom WHERE building = ? AND room_number = ?`

	var classroom model.Classroom
	err := r.db.QueryRow(query, building, roomNumber).Scan(
		&classroom.Building,
		&classroom.RoomNumber,
		&classroom.Capacity,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("classroom not found")
		}
		return nil, fmt.Errorf("error querying classroom: %w", err)
	}

	return &classroom, nil
}
