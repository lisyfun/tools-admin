import request from '@/utils/request'
import type { DashboardOverview, ChartData, ApiResponse } from '@/types/api'

export interface DashboardOverview {
  sms_count: number
  sms_trend: number
  task_count: number
  task_trend: number
  success_rate: number
  success_rate_trend: number
}

export interface ChartData {
  date: string
  success: number
  fail?: number
}

export interface ApiResponse<T> {
  code: number
  data: T
  msg: string
}

export function getOverview() {
  return request<ApiResponse<DashboardOverview>>({
    url: '/dashboard/overview',
    method: 'get'
  })
}

export function getTaskChart(period: string) {
  return request<ApiResponse<ChartData[]>>({
    url: '/dashboard/task-chart',
    method: 'get',
    params: {
      period
    }
  })
}

export function getSmsChart(period: string) {
  return request<ApiResponse<ChartData[]>>({
    url: '/dashboard/sms-chart',
    method: 'get',
    params: {
      period
    }
  })
}
