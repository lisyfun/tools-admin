import request from '@/utils/request'

// 获取短信列表
export function getSmsPage(params: any) {
  return request({
    url: '/api/v1/sms/page',
    method: 'get',
    params
  })
}

// 发送短信
export function sendSms(data: any) {
  return request({
    url: '/api/v1/sms/send',
    method: 'post',
    data
  })
}

// 删除短信
export function deleteSms(id: number) {
  return request({
    url: `/api/v1/sms/${id}`,
    method: 'delete'
  })
}

// 重发短信
export function resendSms(id: number) {
  return request({
    url: `/api/v1/sms/${id}/resend`,
    method: 'post'
  })
}

// 导出短信记录
export function exportSms(params: any) {
  return request({
    url: '/api/v1/sms/export',
    method: 'get',
    params,
    responseType: 'blob'
  })
}
