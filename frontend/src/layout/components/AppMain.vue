<template>
  <section class="app-main">
    <router-view v-slot="{ Component }">
      <transition name="fade-transform" mode="out-in">
        <component :is="Component" />
      </transition>
    </router-view>
  </section>
</template>

<script setup lang="ts">
// AppMain component
</script>

<style lang="scss" scoped>
.app-main {
  min-height: calc(100vh - #{$navbar-height});
  width: 100%;
  position: relative;
  overflow: hidden;
  padding: $main-padding;
  box-sizing: border-box;
  background-color: $main-bg;
}

.fade-transform-enter-active,
.fade-transform-leave-active {
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
}

.fade-transform-enter-from {
  opacity: 0;
  transform: translateX(-20px);
}

.fade-transform-leave-to {
  opacity: 0;
  transform: translateX(20px);
}
</style>
