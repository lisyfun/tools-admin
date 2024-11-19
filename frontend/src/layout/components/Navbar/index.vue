<template>
  <div class="navbar">
    <div class="left">
      <el-icon class="hamburger-container" @click="toggleSidebar">
        <Fold v-if="isCollapse" />
        <Expand v-else />
      </el-icon>
      <Breadcrumb class="breadcrumb-container" />
    </div>
    <div class="right-menu">
      <div class="right-menu-item">
        <el-tooltip content="刷新菜单" placement="bottom">
          <el-icon class="right-menu-icon" @click="refreshMenu">
            <Refresh />
          </el-icon>
        </el-tooltip>
      </div>
      <div class="right-menu-item">
        <el-tooltip content="全屏" placement="bottom">
          <el-icon class="right-menu-icon" @click="toggleFullScreen">
            <FullScreen v-if="!isFullscreen" />
            <Aim v-else />
          </el-icon>
        </el-tooltip>
      </div>

      <el-dropdown class="avatar-container" trigger="click" @command="handleCommand">
        <div class="avatar-wrapper">
          <el-avatar :size="32" :src="userInfo?.avatar || defaultAvatar" />
          <span class="user-name">{{ userInfo?.username }}</span>
          <el-icon class="el-icon-caret-bottom">
            <CaretBottom />
          </el-icon>
        </div>
        <template #dropdown>
          <el-dropdown-menu class="user-dropdown">
            <el-dropdown-item command="profile">
              <el-icon><User /></el-icon>
              个人中心
            </el-dropdown-item>
            <el-dropdown-item divided command="logout">
              <el-icon><SwitchButton /></el-icon>
              退出登录
            </el-dropdown-item>
          </el-dropdown-menu>
        </template>
      </el-dropdown>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { useRouter } from 'vue-router'
import { storeToRefs } from 'pinia'
import { ElMessage } from 'element-plus'
import { useAppStore } from '@/stores/app'
import { useUserStore } from '@/stores/user'
import { useMenuStore } from '@/stores/menu'
import defaultAvatar from '@/assets/avatar.svg'
import { 
  Fold, 
  Expand, 
  CaretBottom, 
  FullScreen,
  Aim,
  User,
  SwitchButton,
  Refresh
} from '@element-plus/icons-vue'
import Breadcrumb from '../Breadcrumb/index.vue'

const router = useRouter()
const appStore = useAppStore()
const userStore = useUserStore()
const menuStore = useMenuStore()

const { userInfo } = storeToRefs(userStore)
const isCollapse = computed(() => appStore.sidebar.isCollapse)
const isFullscreen = ref(false)

const toggleSidebar = () => {
  appStore.toggleSidebar()
}

const toggleFullScreen = () => {
  if (!document.fullscreenElement) {
    document.documentElement.requestFullscreen()
    isFullscreen.value = true
  } else {
    document.exitFullscreen()
    isFullscreen.value = false
  }
}

const refreshMenu = async () => {
  try {
    await menuStore.fetchMenuList()
    ElMessage.success('菜单刷新成功')
  } catch (error) {
    console.error('刷新菜单失败:', error)
    ElMessage.error('菜单刷新失败')
  }
}

const handleCommand = async (command: string) => {
  if (command === 'logout') {
    await userStore.logout()
    router.push('/login')
  } else if (command === 'profile') {
    router.push('/profile')
  }
}
</script>

<style lang="scss" scoped>
.navbar {
  height: $navbar-height;
  overflow: hidden;
  position: relative;
  background: $navbar-bg;
  display: flex;
  justify-content: space-between;
  align-items: center;

  .left {
    display: flex;
    align-items: center;
    padding-left: 16px;
  }

  .hamburger-container {
    padding: 0 15px;
    height: 100%;
    cursor: pointer;
    display: flex;
    align-items: center;
    font-size: 20px;
    transition: background 0.3s;

    &:hover {
      background: rgba(0, 0, 0, 0.025);
    }
  }

  .breadcrumb-container {
    padding-left: 8px;
  }

  .right-menu {
    display: flex;
    align-items: center;
    padding-right: 16px;
    height: 100%;

    .right-menu-item {
      padding: 0 8px;
      height: 100%;
      display: flex;
      align-items: center;
      
      .right-menu-icon {
        font-size: 20px;
        cursor: pointer;
        padding: 4px;
        border-radius: 4px;
        
        &:hover {
          background: rgba(0, 0, 0, 0.025);
        }
      }
    }

    .avatar-container {
      .avatar-wrapper {
        display: flex;
        align-items: center;
        padding: 0 8px;
        cursor: pointer;

        .user-name {
          margin: 0 8px;
          font-size: 14px;
          color: var(--el-text-color-primary);
        }

        .el-icon-caret-bottom {
          font-size: 12px;
          color: var(--el-text-color-secondary);
        }

        &:hover {
          background: rgba(0, 0, 0, 0.025);
        }
      }
    }
  }
}

:deep(.user-dropdown) {
  .el-dropdown-menu__item {
    padding: 8px 16px;
    
    .el-icon {
      margin-right: 8px;
      font-size: 16px;
    }
  }
}
</style>
