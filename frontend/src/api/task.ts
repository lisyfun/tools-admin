import request from '@/utils/request'

export interface TaskQuery {
  page?: number
  pageSize?: number
  name?: string
  type?: number | string
}

export type TaskType = 'regular' | 'urgent' | 'longterm' | 'recurring'
export type TaskPriority = 'low' | 'medium' | 'high'

// 任务状态: 1-启动 2-停止
export type TaskStatus = 1 | 2

// 任务执行状态: 1-待执行 2-执行中 3-执行成功 4-执行失败
export type TaskExecStatus = 1 | 2 | 3 | 4

export interface Task {
  id: number
  name: string
  type: TaskType
  description: string
  status: TaskStatus
  execStatus: TaskExecStatus
  priority: TaskPriority
  createTime: string
  cronExpr: string
  taskContent: string
  taskParams?: string
}

export function getTaskList(params: TaskQuery) {
  return request({
    url: '/task',
    method: 'get',
    params
  })
}

export function getTaskById(id: number) {
  return request({
    url: `/task/${id}`,
    method: 'get'
  })
}

export function createTask(data: Partial<Task>) {
  return request({
    url: '/task',
    method: 'post',
    data
  })
}

export function updateTask(id: number, data: Partial<Task>) {
  return request({
    url: `/task/${id}`,
    method: 'put',
    data
  })
}

export function deleteTask(id: number) {
  return request({
    url: `/task/${id}`,
    method: 'delete'
  })
}

export function batchDeleteTasks(ids: number[]) {
  return request({
    url: '/task/batch',
    method: 'delete',
    data: { ids }
  })
}

export function updateTaskStatus(id: number, status: TaskStatus) {
  return request({
    url: `/task/${id}/status`,
    method: 'patch',
    data: { status }
  })
}

// 获取下次执行时间
export function getNextRunTimes(cronExpr: string) {
  return request({
    url: '/task/next-run-times',
    method: 'get',
    params: { cronExpr }
  })
}
