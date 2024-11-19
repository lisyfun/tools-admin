import { defineStore } from 'pinia'

interface AppState {
  sidebar: {
    isCollapse: boolean
  }
}

// 持久化配置的类型定义
interface PersistStrategy {
  key?: string
  storage?: Storage
  paths?: string[]
}

interface PersistOptions {
  enabled: boolean
  strategies: PersistStrategy[]
}

// 扩展 Pinia store 选项，添加持久化配置支持
declare module 'pinia' {
  interface DefineStoreOptionsBase<S, Store> {
    persist?: PersistOptions
  }
}

export const useAppStore = defineStore('app', {
  state: (): AppState => ({
    sidebar: {
      isCollapse: false
    }
  }),

  getters: {
    getSidebarStatus: (state: AppState): boolean => state.sidebar.isCollapse
  },

  actions: {
    toggleSidebar(): void {
      this.sidebar.isCollapse = !this.sidebar.isCollapse
    }
  },

  persist: {
    enabled: true,
    strategies: [
      {
        key: 'app-store',
        storage: localStorage,
        paths: ['sidebar.isCollapse']
      }
    ]
  }
})
