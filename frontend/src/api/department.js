// department.js
import request from '@/api/request'

export function getDepartmentList(params) {
  return request({
    url: '/admin/department/list',
    method: 'get',
    params
  })
}

export function addDepartment(data) {
  return request({
    url: '/admin/department',
    method: 'post',
    data
  })
}

export function updateDepartment(data) {
  return request({
    url: '/admin/department',
    method: 'put',
    data
  })
}

export function deleteDepartment(id) {
  return request({
    url: `/admin/department/${id}`,
    method: 'delete'
  })
}