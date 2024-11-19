<template>
  <div class="cron-container">
    <div class="cron-content">
      <el-tabs v-model="activeTab">
        <el-tab-pane
          v-for="unit in timeUnits"
          :key="unit.value"
          :label="unit.label"
          :name="unit.value"
        >
          <div class="time-unit-config">
            <!-- 通用配置选项 -->
            <el-radio-group v-model="cronState[unit.value].type">
              <template v-if="unit.value === 'week'">
                <el-radio value="question">不指定(?)</el-radio>
                <el-radio value="every">每{{ unit.label }}</el-radio>
                <el-radio value="range">周期</el-radio>
                <el-radio value="interval">循环</el-radio>
                <el-radio value="specific">指定</el-radio>
              </template>
              <template v-else>
                <el-radio value="every">每{{ unit.label }}</el-radio>
                <el-radio value="range">周期</el-radio>
                <el-radio value="interval">循环</el-radio>
                <el-radio value="specific">指定</el-radio>
              </template>
            </el-radio-group>

            <!-- 周期选择 -->
            <div v-if="cronState[unit.value].type === 'range'" class="config-item">
              从
              <el-input-number
                v-model="cronState[unit.value].start"
                :min="unit.min"
                :max="unit.range"
                controls-position="right"
              />
              到
              <el-input-number
                v-model="cronState[unit.value].end"
                :min="unit.min"
                :max="unit.range"
                controls-position="right"
              />
            </div>

            <!-- 循环选择 -->
            <div v-if="cronState[unit.value].type === 'interval'" class="config-item">
              从
              <el-input-number
                v-model="cronState[unit.value].start"
                :min="unit.min"
                :max="unit.range"
                controls-position="right"
              />
              开始，每
              <el-input-number
                v-model="cronState[unit.value].interval"
                :min="1"
                :max="unit.range"
                controls-position="right"
              />
              {{ unit.label }}执行一次
            </div>

            <!-- 指定选择 -->
            <div v-if="cronState[unit.value].type === 'specific'" class="config-item specific-config">
              <el-checkbox-group v-model="cronState[unit.value].specific">
                <el-checkbox
                  v-for="i in unit.range + 1"
                  :key="i - 1"
                  :value="i - 1"
                >
                  {{ String(i - 1).padStart(2, '0') }}
                </el-checkbox>
              </el-checkbox-group>
            </div>
          </div>
        </el-tab-pane>
      </el-tabs>
    </div>

    <div class="cron-footer">
      <div class="expression">
        表达式: <span class="expression-text">{{ cronExpression }}</span>
      </div>
      <div class="buttons">
        <el-button @click="$emit('close')">取消</el-button>
        <el-button type="primary" @click="handleConfirm">确定</el-button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue'

interface TimeUnit {
  label: string
  value: string
  min: number
  range: number
}

interface UnitState {
  type: 'every' | 'specific' | 'range' | 'interval' | 'question' | 'disabled'
  specific: number[]
  start: number
  end: number
  interval: number
}

interface CronState {
  [key: string]: UnitState
}

const props = defineProps<{
  modelValue: string
}>()

const emit = defineEmits<{
  (e: 'update:modelValue', value: string): void
  (e: 'close'): void
}>()

const timeUnits: TimeUnit[] = [
  { label: '秒', value: 'second', min: 0, range: 59 },
  { label: '分', value: 'minute', min: 0, range: 59 },
  { label: '时', value: 'hour', min: 0, range: 23 },
  { label: '日', value: 'day', min: 1, range: 31 },
  { label: '月', value: 'month', min: 1, range: 12 },
  { label: '周', value: 'week', min: 0, range: 6 }
]
const activeTab = ref('second')
const cronState = ref<CronState>({})

// 初始化cronState
const initCronState = () => {
  timeUnits.forEach(unit => {
    cronState.value[unit.value] = {
      type: unit.value === 'week' ? 'question' : 'every',
      specific: [],
      start: unit.min,
      end: unit.range,
      interval: 1
    }
  })
}

initCronState()

const cronExpression = computed(() => {
  return timeUnits
    .map(unit => {
      const state = cronState.value[unit.value]
      if (unit.value === 'week' && state.type === 'question') {
        return '?'
      }
      switch (state.type) {
        case 'every':
          return '*'
        case 'specific':
          return state.specific.length ? state.specific.sort((a, b) => a - b).join(',') : '*'
        case 'range':
          if (state.start === unit.min && state.end === unit.range) {
            return '*'
          }
          return `${state.start}-${state.end}`
        case 'interval':
          const start = state.start === unit.min ? '*' : state.start
          return `${start}/${state.interval}`
        default:
          return '*'
      }
    })
    .join(' ')
})

const handleConfirm = () => {
  emit('update:modelValue', cronExpression.value)
  emit('close')
}

const parseCronExpression = (expression: string) => {
  const parts = expression.split(' ')
  if (parts.length === 6) {
    timeUnits.forEach((unit, index) => {
      const part = parts[index]
      const state = cronState.value[unit.value]
      
      // 重置状态
      state.type = 'every'
      state.specific = []
      state.start = unit.min
      state.end = unit.range
      state.interval = 1

      if (part === '*') {
        state.type = 'every'
      } else if (part.includes(',')) {
        state.type = 'specific'
        state.specific = part.split(',').map(Number).filter(n => !isNaN(n))
      } else if (part.includes('-')) {
        state.type = 'range'
        const [start, end] = part.split('-').map(Number)
        if (!isNaN(start) && !isNaN(end)) {
          state.start = Math.max(unit.min, Math.min(unit.range, start))
          state.end = Math.max(unit.min, Math.min(unit.range, end))
        }
      } else if (part.includes('/')) {
        state.type = 'interval'
        const [start, interval] = part.split('/')
        if (start === '*') {
          state.start = unit.min
        } else {
          const startNum = parseInt(start)
          if (!isNaN(startNum)) {
            state.start = Math.max(unit.min, Math.min(unit.range, startNum))
          }
        }
        const intervalNum = parseInt(interval)
        if (!isNaN(intervalNum)) {
          state.interval = Math.max(1, intervalNum)
        }
      } else if (part !== '*') {
        const num = parseInt(part)
        if (!isNaN(num)) {
          state.type = 'specific'
          state.specific = [Math.max(unit.min, Math.min(unit.range, num))]
        }
      }
    })
  }
}

watch(
  () => props.modelValue,
  (newVal) => {
    if (newVal) {
      parseCronExpression(newVal)
    }
  },
  { immediate: true }
)
</script>

<style scoped>
.cron-container {
  padding: 20px;
}

.cron-content {
  margin-bottom: 20px;
}

.time-unit-config {
  padding: 20px;
  background-color: #f8f9fa;
  border-radius: 4px;
}

.config-type {
  display: flex;
  flex-direction: row;
  gap: 24px;
  margin-bottom: 20px;
}

:deep(.el-radio-group) {
  display: flex;
  flex-direction: row;
  gap: 24px;
}

:deep(.el-radio) {
  margin-right: 0;
  height: auto;
}

.config-item {
  padding: 16px;
  background-color: #fff;
  border-radius: 4px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.05);
}

.specific-config {
  max-height: 200px;
  overflow-y: auto;
}

:deep(.el-checkbox-group) {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(60px, 1fr));
  gap: 8px;
}

:deep(.el-checkbox) {
  margin-right: 0;
  margin-bottom: 8px;
}

:deep(.el-input-number) {
  width: 100px;
  margin: 0 8px;
}

.cron-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-top: 24px;
  padding-top: 16px;
  border-top: 1px solid #eee;
}

.expression {
  font-size: 14px;
}

.expression-text {
  color: #409EFF;
  font-family: monospace;
  margin-left: 8px;
}

.buttons {
  display: flex;
  gap: 12px;
}

:deep(.el-tabs__content) {
  padding: 20px 0;
}
</style>