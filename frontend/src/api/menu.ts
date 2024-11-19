import request from '@/utils/request'
import type { MenuItem } from '@/types/menu'

interface MenuResponse {
  code: number
  data?: MenuItem[]
  message?: string
}

export function getMenuList(): Promise<MenuResponse> {
  return request({
    url: '/menus',
    method: 'get'
  })
}

export function createMenu(data: MenuItem): Promise<MenuResponse> {
  return request({
    url: '/menu',
    method: 'post',
    data
  })
}

export function updateMenu(id: number, data: MenuItem): Promise<MenuResponse> {
  return request({
    url: `/menu/${id}`,
    method: 'put',
    data
  })
}

export function deleteMenu(id: number): Promise<MenuResponse> {
  return request({
    url: `/menu/${id}`,
    method: 'delete'
  })
}
