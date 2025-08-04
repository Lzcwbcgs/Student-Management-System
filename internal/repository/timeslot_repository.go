package repository

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"

	"github.com/yourusername/student-management-system/internal/model"
)

// TimeSlotRepository 定义时间段仓库接口
type TimeSlotRepository interface {
	FindByID(id string) (*model.TimeSlot, error)
	FindAll() ([]*model.TimeSlot, error)
	FindByDayOfWeek(dayOfWeek int) ([]*model.TimeSlot, error)
	FindByTimeRange(startTime, endTime string) ([]*model.TimeSlot, error)
	Create(timeSlot *model.TimeSlot) error
	Update(timeSlot *model.TimeSlot) error
	Delete(id string) error
}

// SQLTimeSlotRepository 实现TimeSlotRepository接口
type SQLTimeSlotRepository struct {
	db *sql.DB
}

// NewTimeSlotRepository 创建时间段仓库实例
func NewTimeSlotRepository(db *sql.DB) TimeSlotRepository {
	return &SQLTimeSlotRepository{db: db}
}

// FindByID 根据ID查找时间段
func (r *SQLTimeSlotRepository) FindByID(id string) (*model.TimeSlot, error) {
	query := `SELECT id, start_time, end_time, days FROM time_slot WHERE id = ?`

	var timeSlot model.TimeSlot
	var daysStr string

	err := r.db.QueryRow(query, id).Scan(
		&timeSlot.ID,
		&timeSlot.StartTime,
		&timeSlot.EndTime,
		&daysStr,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("time slot not found")
		}
		return nil, fmt.Errorf("error querying time slot: %w", err)
	}

	// 解析星期几数据
	// 这里需要实现从字符串转换为 []int 的逻辑
	timeSlot.Days = parseDaysString(daysStr)

	return &timeSlot, nil
}

// FindAll 查找所有时间段
func (r *SQLTimeSlotRepository) FindAll() ([]*model.TimeSlot, error) {
	query := `SELECT time_slot_id, day, start_hr, start_min, end_hr, end_min FROM time_slot`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error querying time slots: %w", err)
	}
	defer rows.Close()

	var timeSlots []*model.TimeSlot
	for rows.Next() {
		var timeSlot model.TimeSlot
		var dayStr string
		err := rows.Scan(
			&timeSlot.ID,
			&dayStr,
			&timeSlot.StartHr,
			&timeSlot.StartMin,
			&timeSlot.EndHr,
			&timeSlot.EndMin,
		)
		// 将dayStr转换为Days数组
		timeSlot.Days = []int{}
		if dayStr != "" {
			dayInt, err := strconv.Atoi(dayStr)
			if err == nil {
				timeSlot.Days = append(timeSlot.Days, dayInt)
			}
		}
		if err != nil {
			return nil, fmt.Errorf("error scanning time slot: %w", err)
		}
		timeSlots = append(timeSlots, &timeSlot)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating time slots: %w", err)
	}

	return timeSlots, nil
}

// FindByDayOfWeek 根据星期几查找时间段
func (r *SQLTimeSlotRepository) FindByDayOfWeek(dayOfWeek int) ([]*model.TimeSlot, error) {
	query := `SELECT time_slot_id, day, start_hr, start_min, end_hr, end_min FROM time_slot WHERE ? IN (day)`

	rows, err := r.db.Query(query, dayOfWeek)
	if err != nil {
		return nil, fmt.Errorf("error querying time slots by day of week: %w", err)
	}
	defer rows.Close()

	var timeSlots []*model.TimeSlot
	for rows.Next() {
		var timeSlot model.TimeSlot
		var dayStr string
		err := rows.Scan(
			&timeSlot.ID,
			&dayStr,
			&timeSlot.StartHr,
			&timeSlot.StartMin,
			&timeSlot.EndHr,
			&timeSlot.EndMin,
		)
		// 将dayStr转换为Days数组
		timeSlot.Days = []int{}
		if dayStr != "" {
			dayInt, err := strconv.Atoi(dayStr)
			if err == nil {
				timeSlot.Days = append(timeSlot.Days, dayInt)
			}
		}
		if err != nil {
			return nil, fmt.Errorf("error scanning time slot: %w", err)
		}
		timeSlots = append(timeSlots, &timeSlot)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating time slots: %w", err)
	}

	return timeSlots, nil
}

// FindByTimeRange 根据时间范围查找时间段
func (r *SQLTimeSlotRepository) FindByTimeRange(startTime, endTime string) ([]*model.TimeSlot, error) {
	query := `SELECT time_slot_id, day, start_hr, start_min, end_hr, end_min FROM time_slot WHERE (start_hr * 60 + start_min) >= (CAST(? AS UNSIGNED) * 60 + CAST(? AS UNSIGNED)) AND (end_hr * 60 + end_min) <= (CAST(? AS UNSIGNED) * 60 + CAST(? AS UNSIGNED))`

	rows, err := r.db.Query(query, startTime[:2], startTime[3:], endTime[:2], endTime[3:])
	if err != nil {
		return nil, fmt.Errorf("error querying time slots by time range: %w", err)
	}
	defer rows.Close()

	var timeSlots []*model.TimeSlot
	for rows.Next() {
		var timeSlot model.TimeSlot
		var dayStr string
		err := rows.Scan(
			&timeSlot.ID,
			&dayStr,
			&timeSlot.StartHr,
			&timeSlot.StartMin,
			&timeSlot.EndHr,
			&timeSlot.EndMin,
		)
		// 将dayStr转换为Days数组
		timeSlot.Days = []int{}
		if dayStr != "" {
			dayInt, err := strconv.Atoi(dayStr)
			if err == nil {
				timeSlot.Days = append(timeSlot.Days, dayInt)
			}
		}
		if err != nil {
			return nil, fmt.Errorf("error scanning time slot: %w", err)
		}
		timeSlots = append(timeSlots, &timeSlot)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating time slots: %w", err)
	}

	return timeSlots, nil
}

// Create 创建时间段
func (r *SQLTimeSlotRepository) Create(timeSlot *model.TimeSlot) error {
	query := `INSERT INTO time_slot (time_slot_id, day, start_hr, start_min, end_hr, end_min) VALUES (?, ?, ?, ?, ?, ?)`

	// 使用Days数组的第一个元素作为day字段
	day := 0
	if len(timeSlot.Days) > 0 {
		day = timeSlot.Days[0]
	}

	_, err := r.db.Exec(query,
		timeSlot.ID,
		day,
		timeSlot.StartHr,
		timeSlot.StartMin,
		timeSlot.EndHr,
		timeSlot.EndMin,
	)

	if err != nil {
		return fmt.Errorf("error creating time slot: %w", err)
	}

	return nil
}

// Update 更新时间段
func (r *SQLTimeSlotRepository) Update(timeSlot *model.TimeSlot) error {
	query := `UPDATE time_slot SET day = ?, start_hr = ?, start_min = ?, end_hr = ?, end_min = ? WHERE time_slot_id = ?`

	// 使用Days数组的第一个元素作为day字段
	day := 0
	if len(timeSlot.Days) > 0 {
		day = timeSlot.Days[0]
	}

	result, err := r.db.Exec(query,
		day,
		timeSlot.StartHr,
		timeSlot.StartMin,
		timeSlot.EndHr,
		timeSlot.EndMin,
		timeSlot.ID,
	)

	if err != nil {
		return fmt.Errorf("error updating time slot: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error getting rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("time slot not found")
	}

	return nil
}

// Delete 删除时间段
func (r *SQLTimeSlotRepository) Delete(id string) error {
	query := `DELETE FROM time_slot WHERE time_slot_id = ?`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("error deleting time slot: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error getting rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("time slot not found")
	}

	return nil
}

// parseDaysString 解析星期几字符串，转换为整数切片
func parseDaysString(daysStr string) []int {
	if daysStr == "" {
		return nil
	}

	days := strings.Split(daysStr, ",")
	var result []int
	for _, day := range days {
		day = strings.TrimSpace(day)
		if dayInt, err := strconv.Atoi(day); err == nil {
			result = append(result, dayInt)
		}
	}
	return result
}