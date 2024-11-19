import axios from 'axios'
import { ElMessage } from 'element-plus'
import { getToken } from '@/utils/auth'

// 创建axios实例
const service = axios.create({
  baseURL: (import.meta.env.VITE_APP_API_BASE_URL || '') as string,
  timeout: 5000,
  headers: {
    'Content-Type': 'application/json'
  }
})

// 请求拦截器
service.interceptors.request.use(
  (config) => {
    const token = getToken()
    console.log('Current token:', token) // 调试日志
    
    if (token) {
      config.headers['Token'] = token
      // 同时也设置 Authorization 头，以防万一
      config.headers['Authorization'] = `Bearer ${token}`
    }
    
    console.log('Request headers:', config.headers) // 调试日志
    
    if (config.method) {
      config.method = config.method.toLowerCase()
    }

    return config
  },
  (error) => {
    console.error('Request error:', error)
    return Promise.reject(error)
  }
)

// 响应拦截器
service.interceptors.response.use(
  (response) => {
    const res = response.data
    console.log('Response:', response.config.url, res) // 调试日志

    if (res.code === 401) {
      ElMessage.error(res.message || '请先登录')
      // 可以在这里处理登录跳转
      return Promise.reject(new Error(res.message || '请先登录'))
    }
    
    if (res.code !== 0) {
      ElMessage.error(res.message || 'Error')
      return Promise.reject(new Error(res.message || 'Error'))
    }
    return res
  },
  (error) => {
    console.error('Response error:', error.response || error)
    const message = error.response?.data?.message || error.message
    ElMessage.error(message)
    return Promise.reject(error)
  }
)

export default service
