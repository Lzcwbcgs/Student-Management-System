import request from './request'

// 获取系统统计数据
export function getSystemStats() {
  return request({
    url: '/api/admin/stats',
    method: 'get'
  })
}

// 获取学生统计数据
export function getStudentStats() {
  return request({
    url: '/api/student/stats',
    method: 'get'
  })
}

// 获取教师统计数据
export function getInstructorStats() {
  return request({
    url: '/api/instructor/stats',
    method: 'get'
  })
}

// 获取课程统计数据
export function getCourseStats() {
  return request({
    url: '/api/course/stats',
    method: 'get'
  })
}
