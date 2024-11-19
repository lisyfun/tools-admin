<template>
  <component :is="type" v-bind="linkProps(to)">
    <slot />
  </component>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { isExternal } from '@/utils/validate'

interface Props {
  to: string
}

const props = defineProps<Props>()

const type = computed((): string => {
  if (isExternal(props.to)) {
    return 'a'
  }
  return 'router-link'
})

interface LinkProps {
  href?: string
  target?: string
  rel?: string
  to?: string
}

function linkProps(to: string): LinkProps {
  if (isExternal(to)) {
    return {
      href: to,
      target: '_blank',
      rel: 'noopener'
    }
  }
  return {
    to: to
  }
}
</script>
