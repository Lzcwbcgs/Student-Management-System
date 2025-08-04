package repository

import "github.com/yourusername/student-management-system/internal/model"

// BaseRepository 定义基础仓储接口
type BaseRepository[T any] interface {
	GetByID(id string) (*T, error)
	Create(entity *T) error
	Update(entity *T) error
	Delete(id string) error
	List(page, pageSize int) ([]*T, int64, error)
	ExistsByID(id string) (bool, error)
	Search(query string) ([]*T, error)
}

// AuthRepository 定义认证相关的仓储接口
type AuthRepository interface {
	UpdatePassword(id, hashedPassword, salt string) error
}

// StudentRepository 定义学生仓储接口
type StudentRepository interface {
	BaseRepository[model.Student]
	AuthRepository
}

// InstructorRepository 定义教师仓储接口
type InstructorRepository interface {
	BaseRepository[model.Instructor]
	AuthRepository
}
