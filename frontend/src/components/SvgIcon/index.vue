<template>
  <el-icon :size="size" :color="color" :class="className">
    <component :is="icon" />
  </el-icon>
</template>

<script lang="ts">
import { defineComponent } from 'vue'

export default defineComponent({
  name: 'SvgIcon',
  props: {
    icon: {
      type: String,
      required: true
    },
    className: {
      type: String,
      default: ''
    },
    color: {
      type: String,
      default: ''
    },
    size: {
      type: [String, Number],
      default: ''
    }
  }
})
</script>

<style scoped>
.svg-icon {
  vertical-align: -0.15em;
}
</style>
