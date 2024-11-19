<template>
  <div class="sidebar-container">
    <div class="logo-container">
      <router-link to="/" class="logo-link">
        <el-icon class="logo-icon"><Monitor /></el-icon>
        <span class="logo-title">运维管理系统</span>
      </router-link>
    </div>
    <el-scrollbar>
      <el-menu
        :default-active="activeMenu"
        :collapse="isCollapse"
        text-color="#1d2129"
        active-text-color="#1d2129"
        :background-color="variables.menuBg"
        :unique-opened="true"
        :collapse-transition="false"
        mode="vertical"
      >
        <sidebar-item
          v-for="route in menuList"
          :key="route.path"
          :item="route"
          :base-path="route.path"
        />
      </el-menu>
    </el-scrollbar>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useRoute } from 'vue-router'
import { storeToRefs } from 'pinia'
import { useMenuStore } from '@/stores/menu'
import { useAppStore } from '@/stores/app'
import { Monitor } from '@element-plus/icons-vue'
import type { MenuItem } from '@/types/menu'
import SidebarItem from './SidebarItem.vue'
import variables from '@/styles/variables.module.scss'

const route = useRoute()
const menuStore = useMenuStore()
const appStore = useAppStore()

// 从 pinia 中获取菜单列表
const { menuList } = storeToRefs(menuStore)

// 当前激活的菜单
const activeMenu = computed((): string => {
  const { meta, path } = route
  if (meta.activeMenu) {
    return meta.activeMenu
  }
  return path
})

// 是否折叠菜单
const isCollapse = computed((): boolean => appStore.sidebar.isCollapse)
</script>

<style lang="scss" scoped>
.sidebar-container {
  background-color: $menu-bg;
  width: v-bind('variables.sidebarWidth');
  height: 100%;
  position: fixed;
  font-size: 0;
  top: 0;
  bottom: 0;
  left: 0;
  z-index: 1001;
  overflow: hidden;
  transition: width 0.3s;

  .logo-container {
    height: 60px;
    padding: 10px 0;
    margin-bottom: 8px;
    background-color: $menu-bg;
    overflow: hidden;

    .logo-link {
      height: 100%;
      padding: 0 16px;
      display: flex;
      align-items: center;
      justify-content: center;
      text-decoration: none;

      .logo-icon {
        font-size: 24px;
        color: var(--el-color-primary);
        margin-right: 8px;
      }

      .logo-title {
        font-size: 18px;
        font-weight: 600;
        color: #1d2129;
        white-space: nowrap;
        transition: all 0.3s;
      }
    }
  }

  :deep(a) {
    text-decoration: none;
  }

  .el-scrollbar {
    height: calc(100% - 69px);
    background: $menu-bg;

    .el-menu {
      border: none;
      height: 100%;
      width: 100% !important;
      padding: 4px 2px;
      background-color: $menu-bg;
      
      .el-menu-item, .el-sub-menu__title {
        height: 50px !important;
        line-height: 50px !important;
        color: #1d2129;
        font-size: 14px;
        padding: 0 8px !important;
        margin: 4px 0;
        border-radius: 4px;
        display: flex;
        align-items: center;
        justify-content: center;
        
        .el-icon {
          font-size: 16px;
          color: #1d2129;
          margin-right: 6px;
          transition: color 0.3s;
          flex-shrink: 0;
        }

        :deep(.el-icon) {
          margin-right: 6px;
          font-size: 16px;
          color: #1d2129;
        }
        
        span {
          flex: 1;
          text-align: center;
          color: #1d2129;
        }
        
        &:hover {
          background-color: #e9ecf2;
          
          .el-icon {
            color: #1d2129;
          }
        }
        
        &.is-active {
          background-color: var(--el-color-primary-light-9);
          color: #1d2129;
          font-weight: 500;
          
          .el-icon {
            color: #1d2129;
          }
        }
      }

      .el-sub-menu {
        &.is-active {
          > .el-sub-menu__title {
            color: #1d2129;
            
            .el-icon {
              color: #1d2129;
            }
          }
        }

        .el-menu {
          padding: 0;
          background: transparent;
          
          .el-menu-item {
            padding-left: 48px !important;
            height: 40px;
            line-height: 40px;
            
            &:hover {
              color: #1d2129;
            }
          }
        }

        :deep(.el-sub-menu__title) {
          display: flex;
          align-items: center;
          justify-content: center;
          color: #1d2129;

          span {
            flex: 1;
            text-align: center;
            color: #1d2129;
          }
        }
      }
    }
  }
}

.el-menu--collapse {
  .logo-container {
    .logo-title {
      display: none;
    }
    .logo-icon {
      margin: 0;
    }
  }

  .el-menu-item, .el-sub-menu__title {
    padding: 0 !important;
    justify-content: center;

    .el-icon {
      margin: 0;
      color: #1d2129;
    }

    span {
      display: none;
    }
  }
}

:deep(.el-menu--popup) {
  padding: 4px 0;
  min-width: 160px;
  border-radius: 4px;
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
  
  .el-menu-item {
    height: 40px !important;
    line-height: 40px !important;
    padding: 0 16px !important;
    display: flex;
    align-items: center;
    justify-content: center;
    color: #1d2129;
    
    span {
      flex: 1;
      text-align: center;
      color: #1d2129;
    }
    
    &:hover {
      background-color: #e9ecf2 !important;
      color: #1d2129 !important;
    }
    
    &.is-active {
      background-color: var(--el-color-primary-light-9) !important;
      color: #1d2129 !important;
    }
  }
}

:deep(.el-menu) {
  border: none;
  height: 100%;
  width: 100% !important;
  padding: 4px 2px;
  background-color: $menu-bg;
}

:deep(.el-menu-item), :deep(.el-sub-menu__title) {
  display: flex;
  align-items: center;
  justify-content: flex-start;
  border-radius: 4px;
  margin: 4px 20px;
  position: relative;
  padding: 0 8px 0 36px !important;
  height: 50px !important;
  line-height: 50px !important;
  color: #1d2129;
  
  &::before {
    content: '';
    position: absolute;
    left: 20px;
    width: 6px;
    height: 6px;
    border-radius: 50%;
    background-color: #909399;
    opacity: 0.5;
  }
  
  .el-icon {
    margin-right: 8px;
    font-size: 16px;
    color: #1d2129;
  }

  &:hover {
    background-color: #e9ecf2 !important;
    color: #1d2129 !important;

    &::before {
      background-color: #1d2129;
      opacity: 1;
    }
  }
}

:deep(.el-sub-menu .el-menu-item) {
  min-width: unset !important;
  display: flex;
  align-items: center;
  justify-content: flex-start;
  border-radius: 4px;
  margin: 4px 20px;
  padding: 0 8px 0 36px !important;
  height: 50px !important;
  line-height: 50px !important;
  color: #1d2129;

  &.is-active {
    background-color: #e9ecf2 !important;
    color: #1d2129 !important;

    &::before {
      background-color: #1d2129;
      opacity: 1;
    }
  }
}

:deep(.el-sub-menu) {
  &.is-active {
    > .el-sub-menu__title {
      color: #1d2129 !important;

      &::before {
        background-color: #1d2129;
        opacity: 1;
      }
    }
  }
}

:deep(.el-menu-item.is-active) {
  background-color: #e9ecf2 !important;
  color: #1d2129 !important;
  border-radius: 4px;
  
  &::before {
    background-color: #1d2129;
    opacity: 1;
  }
}

:deep(.el-sub-menu .el-menu-item) {
  min-width: unset !important;
  display: flex;
  align-items: center;
  justify-content: flex-start;
  border-radius: 4px;
  margin: 4px 20px;
  color: #1d2129;

  &.is-active {
    background-color: #e9ecf2 !important;
    color: #1d2129 !important;

    &::before {
      background-color: #1d2129;
      opacity: 1;
    }
  }
}

:deep(.el-sub-menu) {
  &.is-active {
    > .el-sub-menu__title {
      color: #1d2129 !important;

      &::before {
        background-color: #1d2129;
        opacity: 1;
      }
    }
  }
}
</style>
