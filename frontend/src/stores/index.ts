import { createPinia } from 'pinia'
import piniaPluginPersistedstate from 'pinia-plugin-persistedstate'

// 创建 pinia 实例
const pinia = createPinia()

// 使用持久化插件
pinia.use(piniaPluginPersistedstate)

export default pinia

// 导出类型
declare module 'pinia' {
  export interface PiniaCustomProperties {
    // 在这里添加你的自定义属性
  }
}
