<template>
  <div class="app-wrapper" :class="{ 'hide-sidebar': isCollapse }">
    <Sidebar class="sidebar-container" />
    <div class="main-container">
      <div class="fixed-header">
        <Navbar />
      </div>
      <AppMain />
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useAppStore } from '@/stores/app'
import { Sidebar, Navbar, AppMain } from './components'

const appStore = useAppStore()
const isCollapse = computed(() => appStore.sidebar.opened)
</script>

<style lang="scss" scoped>
.app-wrapper {
  position: relative;
  height: 100%;
  width: 100%;

  .sidebar-container {
    position: fixed;
    top: 0;
    left: 0;
    height: 100%;
    width: $sidebar-width;
    background-color: $menu-bg;
    transition: width $sidebar-transition-duration;
    z-index: 1001;
  }

  .main-container {
    min-height: 100%;
    margin-left: $sidebar-width;
    position: relative;
    background-color: $main-bg;
    transition: margin-left $sidebar-transition-duration;

    &.no-sidebar {
      margin-left: 0;
    }
  }

  &.hide-sidebar {
    .sidebar-container {
      width: $sidebar-hide-width;
    }

    .main-container {
      margin-left: $sidebar-hide-width;
    }
  }

  .fixed-header {
    position: fixed;
    top: 0;
    right: 0;
    z-index: 1000;
    width: calc(100% - #{$sidebar-width});
    transition: width $sidebar-transition-duration;

    .hide-sidebar & {
      width: calc(100% - #{$sidebar-hide-width});
    }
  }
}

// 为固定头部下的主内容添加上边距
:deep(.app-main) {
  padding-top: $navbar-height;
}
</style>
