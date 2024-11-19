import { ElMessage } from 'element-plus'

interface ErrorOptions {
  message?: string
  showMessage?: boolean
  level?: 'error' | 'warn' | 'info'
}

export function handleError(error: any, options: ErrorOptions = {}) {
  const {
    message = '操作失败',
    showMessage = true,
    level = 'error'
  } = options

  // 开发环境下打印错误
  if (import.meta.env.DEV) {
    console[level](`[${level.toUpperCase()}]`, error)
  }

  // 显示错误消息
  if (showMessage) {
    ElMessage[level](error?.message || message)
  }

  return error
}

export function createErrorHandler(defaultOptions: ErrorOptions = {}) {
  return (error: any, options: ErrorOptions = {}) => 
    handleError(error, { ...defaultOptions, ...options })
}
