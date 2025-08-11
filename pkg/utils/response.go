package utils

import (
	"encoding/json"
	"net/http"
)

// Response 定义API响应的标准格式
type Response struct {
	Code    int         `json:"code"`    // 状态码
	Message string      `json:"message"` // 响应消息
	Data    interface{} `json:"data"`    // 响应数据
}

// ApiResponse 统一响应结构（与原来的Response结构保持一致）
type ApiResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// NewResponse 创建一个新的响应对象
func NewResponse(code int, message string, data interface{}) *Response {
	return &Response{
		Code:    code,
		Message: message,
		Data:    data,
	}
}

// JSON 将响应以JSON格式写入HTTP响应
func (r *Response) JSON(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(r.Code)

	if err := json.NewEncoder(w).Encode(r); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// Success 返回成功响应
func Success(w http.ResponseWriter, message string, data interface{}) {
	resp := NewResponse(http.StatusOK, message, data)
	resp.JSON(w)
}

// Created 返回资源创建成功响应
func Created(w http.ResponseWriter, message string, data interface{}) {
	resp := NewResponse(http.StatusCreated, message, data)
	resp.JSON(w)
}

// BadRequest 返回请求错误响应
func BadRequest(w http.ResponseWriter, message string) {
	resp := NewResponse(http.StatusBadRequest, message, nil)
	resp.JSON(w)
}

// Unauthorized 返回未授权响应
func Unauthorized(w http.ResponseWriter, message string) {
	resp := NewResponse(http.StatusUnauthorized, message, nil)
	resp.JSON(w)
}

// Forbidden 返回禁止访问响应
func Forbidden(w http.ResponseWriter, message string) {
	resp := NewResponse(http.StatusForbidden, message, nil)
	resp.JSON(w)
}

// NotFound 返回资源未找到响应
func NotFound(w http.ResponseWriter, message string) {
	resp := NewResponse(http.StatusNotFound, message, nil)
	resp.JSON(w)
}

// InternalServerError 返回服务器内部错误响应
func InternalServerError(w http.ResponseWriter, message string) {
	resp := NewResponse(http.StatusInternalServerError, message, nil)
	resp.JSON(w)
}

// ValidationError 返回数据验证错误响应
func ValidationError(w http.ResponseWriter, message string, errors interface{}) {
	resp := NewResponse(http.StatusUnprocessableEntity, message, errors)
	resp.JSON(w)
}

// WriteJSONResponse 写入JSON响应（修改为统一格式）
func WriteJSONResponse(w http.ResponseWriter, statusCode int, data interface{}) {
	response := ApiResponse{
		Code:    statusCode,
		Message: "success",
		Data:    data,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

// WriteErrorResponse 写入错误响应（修改为统一格式）
func WriteErrorResponse(w http.ResponseWriter, statusCode int, message string) {
	response := ApiResponse{
		Code:    statusCode,
		Message: message,
		Data:    nil,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(response)
}

// WriteSuccessResponse 写入成功响应（修改为统一格式）
func WriteSuccessResponse(w http.ResponseWriter, message string) {
	response := ApiResponse{
		Code:    http.StatusOK,
		Message: message,
		Data:    nil,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
