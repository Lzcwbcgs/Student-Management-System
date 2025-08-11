import axios from 'axios'
import router from '@/router' 
import NProgress from 'nprogress'
import 'nprogress/nprogress.css'

// 创建 axios 实例
const request = axios.create({
  baseURL: 'http://localhost:8080',  // 替换为您的后端 API 地址
  timeout: 15000
})

// 请求拦截器
request.interceptors.request.use(
  config => {
    NProgress.start()
    const token = localStorage.getItem('token')
    if (token) {
      config.headers['Authorization'] = `Bearer ${token}`
    }
    return config
  },
  error => {
    NProgress.done()
    return Promise.reject(error)
  }
)

// 响应拦截器
request.interceptors.response.use(
  response => {
    NProgress.done()
    return response.data
  },
  error => {
    NProgress.done()
    if (error.response.status === 401) {
      // 未授权，跳转到登录页
      router.push('/login')
    }
    return Promise.reject(error)
  }
)

export default request
