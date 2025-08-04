package utils

import "errors"

// NewError 创建新的错误
func NewError(message string) error {
	return errors.New(message)
}

// ErrorResponse 错误响应结构
type ErrorResponse struct {
	Error   string `json:"error"`
	Code    int    `json:"code,omitempty"`
	Details string `json:"details,omitempty"`
} 