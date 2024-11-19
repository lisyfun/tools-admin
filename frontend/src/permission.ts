import router from './router'
import { useUserStore } from '@/stores/user'
import { useMenuStore } from '@/stores/menu'
import { getToken } from '@/utils/auth'
import { ElMessage } from 'element-plus'
import NProgress from 'nprogress'
import 'nprogress/nprogress.css'

NProgress.configure({ showSpinner: false })

const whiteList = ['/login'] // 不需要重定向的白名单

router.beforeEach(async (to, from, next) => {
  NProgress.start()

  const hasToken = getToken()
  const userStore = useUserStore()
  const menuStore = useMenuStore()
  
  if (hasToken) {
    if (to.path === '/login') {
      // 已登录，跳转到首页
      next({ path: '/' })
      NProgress.done()
    } else {
      // 判断是否已获取用户信息
      const hasUserInfo = userStore.getUserInfo !== null
      if (hasUserInfo) {
        // 判断是否已获取菜单
        const hasMenu = menuStore.menus.length > 0
        if (!hasMenu) {
          try {
            // 获取菜单数据
            await menuStore.fetchMenuList()
            // 获取成功后重新加载当前路由
            next({ ...to, replace: true })
          } catch (error) {
            // 获取菜单失败，清空token并跳转到登录页
            userStore.resetToken()
            ElMessage.error('获取菜单失败，请重新登录')
            next(`/login?redirect=${to.path}`)
            NProgress.done()
          }
        } else {
          next()
        }
      } else {
        try {
          // 获取用户信息
          await userStore.fetchUserInfo()
          // 获取菜单数据
          await menuStore.fetchMenuList()
          next({ ...to, replace: true })
        } catch (error) {
          // 获取用户信息失败，清空token并跳转到登录页
          userStore.resetToken()
          ElMessage.error('获取用户信息失败，请重新登录')
          next(`/login?redirect=${to.path}`)
          NProgress.done()
        }
      }
    }
  } else {
    if (whiteList.includes(to.path)) {
      // 在免登录白名单中，直接进入
      next()
    } else {
      // 其他没有访问权限的页面将被重定向到登录页面
      next(`/login?redirect=${to.path}`)
      NProgress.done()
    }
  }
})

router.afterEach(() => {
  NProgress.done()
})
