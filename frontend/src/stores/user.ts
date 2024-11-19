import { defineStore } from 'pinia'
import { ElMessage } from 'element-plus'
import { login, resetPassword, getUserInfo as fetchUserInfo } from '@/api/user'
import { useMenuStore } from './menu'
import { setToken, removeToken } from '@/utils/auth'

interface UserInfo {
  username: string
  avatar?: string
  roles?: string[]
}

interface UserState {
  token: string
  userInfo: UserInfo | null
}

interface LoginForm {
  username: string
  password: string
}

interface ResetPasswordForm {
  username: string
  newPassword: string
}

export const useUserStore = defineStore('user', {
  state: (): UserState => ({
    token: '',
    userInfo: null
  }),

  getters: {
    getToken: (state): string => state.token,
    getUserInfo: (state): UserInfo | null => state.userInfo
  },

  actions: {
    async login(loginForm: LoginForm) {
      try {
        const { data, code } = await login(loginForm)
        if (code === 0 && data) {
          this.token = data.token
          this.userInfo = {
            username: loginForm.username,
            ...data.user
          }
          setToken(data.token)
          
          // 登录成功后获取菜单列表
          const menuStore = useMenuStore()
          await menuStore.fetchMenuList()
          return true
        }
        return false
      } catch (error) {
        console.error('Login error:', error)
        return false
      }
    },

    async fetchUserInfo() {
      try {
        const { data, code } = await fetchUserInfo()
        if (code === 0 && data) {
          this.userInfo = data
          return data
        }
        return null
      } catch (error) {
        console.error('Get user info error:', error)
        return null
      }
    },

    async resetPassword(form: ResetPasswordForm) {
      try {
        const { code, message } = await resetPassword(form)
        if (code === 0) {
          ElMessage.success('密码重置成功')
          return true
        } else {
          ElMessage.error(message || '密码重置失败')
          return false
        }
      } catch (error) {
        console.error('Reset password error:', error)
        ElMessage.error('密码重置失败')
        return false
      }
    },

    resetToken() {
      this.token = ''
      this.userInfo = null
      removeToken()
    },

    async logout() {
      this.resetToken()
      // 重置路由
      location.reload()
    }
  },

  persist: {
    enabled: true,
    strategies: [
      {
        key: 'user',
        storage: localStorage,
        paths: ['token', 'userInfo']
      }
    ]
  }
})
