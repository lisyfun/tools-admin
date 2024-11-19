export interface LoginParams {
  username: string
  password: string
}

export interface UserInfo {
  id: number
  username: string
  nickname: string
  avatar: string
  email: string
  roles: string[]
  permissions: string[]
}

export interface LoginResponse {
  token: string
  user: UserInfo
}
