import request from '@/utils/request'
import type { LoginParams, UserInfo } from '@/types/user'

/**
 * 用户登录
 * @param data
 */
export function login(data: LoginParams) {
  return request({
    url: '/login',
    method: 'post',
    data
  })
}

/**
 * 获取用户信息
 */
export function getUserInfo() {
  return request({
    url: '/user/info',
    method: 'get'
  })
}

/**
 * 用户登出
 */
export function logout() {
  return request({
    url: '/logout',
    method: 'post'
  })
}

/**
 * 修改密码
 */
export function changePassword(data: { oldPassword: string; newPassword: string }) {
  return request({
    url: '/user/password',
    method: 'put',
    data
  })
}

/**
 * 重置密码
 */
export function resetPassword(data: { username: string; newPassword: string }) {
  return request({
    url: '/user/reset-password',
    method: 'post',
    data
  })
}

/**
 * 修改用户信息
 */
export function updateUserInfo(data: Partial<UserInfo>) {
  return request({
    url: '/users/profile',
    method: 'put',
    data
  })
}
