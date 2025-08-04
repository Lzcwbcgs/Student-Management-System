import request from './request'

// 获取课程列表
export function getCourseList(params) {
  return request({
    url: '/api/courses',
    method: 'get',
    params
  })
}

// 获取课程详情
export function getCourseDetail(id) {
  return request({
    url: `/api/courses/${id}`,
    method: 'get'
  })
}

// 创建课程
export function createCourse(data) {
  return request({
    url: '/api/courses',
    method: 'post',
    data
  })
}

// 更新课程信息
export function updateCourse(id, data) {
  return request({
    url: `/api/courses/${id}`,
    method: 'put',
    data
  })
}

// 删除课程
export function deleteCourse(id) {
  return request({
    url: `/api/courses/${id}`,
    method: 'delete'
  })
}

// 获取课程的先修课程
export function getCoursePrerequisites(id) {
  return request({
    url: `/api/courses/${id}/prerequisites`,
    method: 'get'
  })
}

// 更新课程的先修课程
export function updateCoursePrerequisites(id, data) {
  return request({
    url: `/api/courses/${id}/prerequisites`,
    method: 'put',
    data
  })
}

// 获取教师列表（用于课程教师分配）
export function getInstructorList() {
  return request({
    url: '/api/instructors',
    method: 'get'
  })
}

// 分配教师到课程
export function assignInstructor(courseId, data) {
  return request({
    url: `/api/courses/${courseId}/instructors`,
    method: 'post',
    data
  })
}
