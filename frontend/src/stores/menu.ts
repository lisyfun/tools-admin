import { defineStore } from 'pinia'
import { getMenuList } from '@/api/menu'
import type { MenuItem } from '@/types/menu'
import router from '@/router'
import type { RouteRecordRaw } from 'vue-router'

interface MenuState {
  menuList: MenuItem[]
  routes: MenuItem[]
}

// 处理菜单数据，确保所有必要的字段都存在
const processMenuItem = (item: MenuItem): MenuItem => {
  console.log('Processing menu item:', item)
  const processed: MenuItem = {
    ...item,
    path: item.path || '',
    name: item.name || '',
    meta: {
      title: item.name || '',
      icon: item.meta?.icon || item.icon || '',  // 优先使用 meta.icon，如果不存在则使用 item.icon
      noCache: item.keepAlive === 0,
      breadcrumb: true,
      hidden: item.visible === 0
    },
    children: (item.children || []).map(child => processMenuItem(child))
  }
  console.log('Processed menu item:', processed)
  return processed
}

const filterHiddenRoutes = (routes: RouteRecordRaw[]): MenuItem[] => {
  return routes
    .filter((route: RouteRecordRaw) => {
      const { meta } = route
      return !meta?.hidden
    })
    .map((route: RouteRecordRaw) => {
      const menuItem: MenuItem = {
        path: route.path,
        name: route.name as string,
        meta: {
          title: route.meta?.title || route.name,
          icon: route.meta?.icon || '',
          noCache: route.meta?.noCache || false,
          breadcrumb: route.meta?.breadcrumb !== false,
          hidden: route.meta?.hidden || false
        },
        children: route.children ? filterHiddenRoutes(route.children) : []
      }
      return menuItem
    })
}

export const useMenuStore = defineStore('menu', {
  state: (): MenuState => ({
    menuList: [],
    routes: []
  }),

  getters: {
    menus: (state): MenuItem[] => {
      if (state.menuList.length === 0) {
        return filterHiddenRoutes(router.options.routes)
      }
      return state.menuList.map(item => processMenuItem(item))
    }
  },

  actions: {
    async fetchMenuList() {
      try {
        const { data, code } = await getMenuList()
        if (code === 0 && data) {
          this.menuList = data.map(item => processMenuItem(item))
          return this.menuList
        }
        this.menuList = filterHiddenRoutes(router.options.routes)
        return this.menuList
      } catch (error) {
        console.error('Get menu list error:', error)
        this.menuList = filterHiddenRoutes(router.options.routes)
        return this.menuList
      }
    },

    setMenuList(menuList: MenuItem[]) {
      this.menuList = menuList.map(item => processMenuItem(item))
    },

    setRoutes(routes: MenuItem[]) {
      this.routes = routes.map(item => processMenuItem(item))
    }
  },

  persist: {
    enabled: true,
    strategies: [
      {
        key: 'menu',
        storage: localStorage,
        paths: ['menuList', 'routes']
      }
    ]
  }
})

// 将后端菜单数据转换为路由配置
function generateRoutesFromMenu(menuList: MenuItem[], parentPath = ''): MenuItem[] {
  return menuList.map(menu => {
    const route: MenuItem = {
      path: parentPath + '/' + menu.path,
      name: menu.name,
      component: loadComponent(menu.component),
      meta: { ...menu.meta }
    }

    if (menu.children && menu.children.length > 0) {
      route.children = generateRoutesFromMenu(menu.children, route.path)
    }

    return route
  })
}

// 动态加载组件
function loadComponent(component: string): string {
  if (!component) return ''
  // 这里假设所有组件都在 views 目录下
  return `@/views/${component}`
}
