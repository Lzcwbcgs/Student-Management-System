import request from './request'

export function login({ user_id, password, role }) {
  return request({
    url: '/api/login',
    method: 'post',
    data: { user_id, password, role }
  })
}

export function register(data) {
  return request({
    url: '/api/register',
    method: 'post',
    data
  })
}

export function logout() {
  return request({
    url: '/api/logout',
    method: 'post'
  })
}

export function getUserInfo() {
  return request({
    url: '/api/user-info',
    method: 'get'
  })
}
