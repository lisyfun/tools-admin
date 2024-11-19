<template>
  <div v-if="!item.hidden">
    <template v-if="hasOneShowingChild(item.children, item) && (!onlyOneChild.children || onlyOneChild.noShowingChildren) && !item.alwaysShow">
      <app-link v-if="onlyOneChild.meta" :to="resolvePath(onlyOneChild.path, onlyOneChild.query)">
        <el-menu-item :index="resolvePath(onlyOneChild.path)" :class="{ 'submenu-title-noDropdown': !isNest }">
          <el-icon v-if="onlyOneChild.meta.icon || (item.meta && item.meta.icon)">
            <component :is="onlyOneChild.meta.icon || (item.meta && item.meta.icon)" />
          </el-icon>
          <template #title>
            <span>{{ onlyOneChild.meta.title }}</span>
          </template>
        </el-menu-item>
      </app-link>
    </template>

    <el-sub-menu v-else ref="subMenu" :index="resolvePath(item.path)" popper-append-to-body>
      <template #title>
        <el-icon v-if="item.meta && item.meta.icon">
          <component :is="item.meta.icon" />
        </el-icon>
        <span>{{ item.meta.title }}</span>
        <span class="custom-arrow">{{ isCollapsed ? '+' : '-' }}</span>
      </template>

      <template v-for="child in item.children" :key="child.path">
        <template v-if="!child.hidden">
          <sidebar-item
            v-if="child.children && child.children.length > 0"
            :key="child.path"
            :item="child"
            :is-nest="true"
            :base-path="resolvePath(child.path)"
            class="nest-menu"
          />
          <app-link
            v-else
            :key="child.path"
            :to="resolvePath(child.path, child.query)"
          >
            <el-menu-item :index="resolvePath(child.path)">
              <el-icon v-if="child.meta && child.meta.icon">
                <component :is="child.meta.icon" />
              </el-icon>
              <template #title>
                <span>{{ child.meta.title }}</span>
              </template>
            </el-menu-item>
          </app-link>
        </template>
      </template>
    </el-sub-menu>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { isExternal } from '@/utils/validate'
import AppLink from './Link.vue'
import path from 'path-browserify'

const props = defineProps({
  item: {
    type: Object,
    required: true
  },
  isNest: {
    type: Boolean,
    default: false
  },
  basePath: {
    type: String,
    default: ''
  }
})

const subMenu = ref()
const onlyOneChild = ref<any>(null)

const isCollapsed = computed(() => {
  return !subMenu.value?.opened
})

const hasOneShowingChild = (children = [], parent) => {
  if (!children) {
    children = []
  }
  const showingChildren = children.filter(item => {
    if (item.hidden) {
      return false
    } else {
      onlyOneChild.value = item
      return true
    }
  })

  if (showingChildren.length === 1) {
    return true
  }

  if (showingChildren.length === 0) {
    onlyOneChild.value = { ...parent, path: '', noShowingChildren: true }
    return true
  }

  return false
}

const resolvePath = (routePath: string, routeQuery?: Record<string, any>): string => {
  if (isExternal(routePath)) {
    return routePath
  }
  if (isExternal(props.basePath)) {
    return props.basePath
  }

  const resolvedPath = path.resolve(props.basePath, routePath)
  
  if (routeQuery) {
    const query = Object.entries(routeQuery)
      .map(([key, value]) => `${key}=${value}`)
      .join('&')
    return `${resolvedPath}?${query}`
  }
  
  return resolvedPath
}
</script>

<style lang="scss" scoped>
:deep(.el-menu-item), :deep(.el-sub-menu) {
  &.is-active > .el-sub-menu__title {
    color: var(--el-menu-active-color);
  }
}

:deep(.el-sub-menu) {
  .el-sub-menu__title {
    position: relative;
    user-select: none;
    -webkit-user-select: none;
    -moz-user-select: none;
    -ms-user-select: none;
    
    span {
      user-select: none;
      -webkit-user-select: none;
      -moz-user-select: none;
      -ms-user-select: none;
    }
    
    .el-sub-menu__icon-arrow {
      display: none !important;
      opacity: 0 !important;
      visibility: hidden !important;
      width: 0 !important;
      height: 0 !important;
    }
    
    .custom-arrow {
      position: absolute;
      right: 20px;
      top: 50%;
      transform: translateY(-50%);
      font-size: 16px;
      font-weight: bold;
      font-family: none;
      line-height: 1;
      width: 16px;
      height: 16px;
      text-align: center;
      transition: all 0.3s;
      user-select: none;
      -webkit-user-select: none;
      -moz-user-select: none;
      -ms-user-select: none;
      margin-top: -1px;
    }
  }
}

:deep(.el-menu) {
  .el-sub-menu__icon-arrow {
    display: none !important;
  }
}

.nest-menu .el-menu-item {
  min-width: 160px;
  padding-left: 44px !important;
  
  &.is-active {
    background-color: var(--el-menu-hover-bg-color);
  }
}

// 额外的样式来确保箭头隐藏
:deep([class*='el-icon-arrow']) {
  display: none !important;
}

:deep(.el-icon) {
  &.el-sub-menu__icon-arrow {
    display: none !important;
  }
}
</style>
