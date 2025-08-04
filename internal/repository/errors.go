package repository

import "errors"

// 定义仓储层通用错误
var (
	// ErrNotFound 表示未找到请求的资源
	ErrNotFound = errors.New("resource not found")
	
	// ErrDuplicate 表示资源已存在
	ErrDuplicate = errors.New("resource already exists")
	
	// ErrDatabase 表示数据库操作错误
	ErrDatabase = errors.New("database operation failed")
	
	// ErrInvalidInput 表示输入参数无效
	ErrInvalidInput = errors.New("invalid input parameters")
)