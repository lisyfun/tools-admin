<template>
  <el-dialog
    :title="isEdit ? '编辑任务' : '新建任务'"
    v-model="visible"
    :close-on-click-modal="false"
    width="800px"
    class="task-form"
    destroy-on-close
  >
    <el-form
      ref="formRef"
      :model="form"
      :rules="rules"
      label-width="100px"
      class="form-container"
    >
      <div class="form-section">
        <div class="section-title">
          <el-icon><Document /></el-icon>
          基本信息
        </div>
        
        <el-form-item label="任务名称" prop="name">
          <el-input v-model="form.name" placeholder="请输入任务名称" />
        </el-form-item>

        <el-form-item label="任务类型" prop="type">
          <el-select v-model="form.type" placeholder="请选择任务类型">
            <el-option
              v-for="(label, value) in taskTypes"
              :key="value"
              :label="label"
              :value="Number(value)"
            >
              <div class="task-type-option">
                <el-icon>
                  <component :is="getTypeIcon(Number(value))" />
                </el-icon>
                <span class="label">{{ label }}</span>
              </div>
            </el-option>
          </el-select>
        </el-form-item>

        <el-form-item label="优先级" prop="priority">
          <el-select v-model="form.priority" placeholder="请选择优先级">
            <el-option
              v-for="item in priorityOptions"
              :key="item.value"
              :label="item.label"
              :value="item.value"
            >
              <div class="priority-option">
                <el-icon :class="['priority-icon', item.value]">
                  <component :is="item.icon" />
                </el-icon>
                <span>{{ item.label }}</span>
              </div>
            </el-option>
          </el-select>
        </el-form-item>

        <el-form-item label="Cron表达式" prop="cronExpr">
          <el-input v-model="form.cronExpr" placeholder="请输入Cron表达式">
            <template #append>
              <el-button @click="showCronDialog = true">
                <el-icon><AlarmClock /></el-icon>
                生成
              </el-button>
            </template>
          </el-input>
        </el-form-item>

        <el-form-item label="任务描述" prop="description">
          <el-input
            v-model="form.description"
            type="textarea"
            :rows="2"
            placeholder="请输入任务描述"
          />
        </el-form-item>
      </div>

      <div class="form-section">
        <div class="section-title">
          <el-icon><Setting /></el-icon>
          任务配置
        </div>

        <el-form-item label="任务参数" prop="taskParams">
          <el-input
            v-model="form.taskParams"
            type="textarea"
            :rows="4"
            placeholder="请输入任务参数，使用JSON格式，例如：&#x0A;{&#x0A;  &quot;param1&quot;: &quot;value1&quot;,&#x0A;  &quot;param2&quot;: &quot;value2&quot;&#x0A;}"
          >
            <template #append>
              <el-button @click="formatJson">
                <el-icon><Tools /></el-icon>
                格式化
              </el-button>
            </template>
          </el-input>
        </el-form-item>

        <el-form-item label="任务内容" prop="taskContent">
          <el-input
            v-model="form.taskContent"
            type="textarea"
            :rows="4"
            placeholder="请输入任务内容"
          />
        </el-form-item>
      </div>
    </el-form>

    <template #footer>
      <div class="dialog-footer">
        <el-button class="cancel-button" @click="visible = false">
          <el-icon><CircleClose /></el-icon>
          取消
        </el-button>
        <el-button
          type="primary"
          class="submit-button"
          :loading="loading"
          @click="submitForm"
        >
          <el-icon><Check /></el-icon>
          确定
        </el-button>
      </div>
    </template>

    <!-- Cron表达式生成器 -->
    <el-dialog
      v-model="showCronDialog"
      title="Cron表达式生成器"
      width="680px"
      append-to-body
      destroy-on-close
    >
      <Cron
        v-model="form.cronExpr"
        @confirm="handleCronConfirm"
        @close="handleCronClose"
      />
    </el-dialog>
  </el-dialog>
</template>

<script setup lang="ts">
import { 
  Calendar, Timer, Document, Monitor, Link, DataAnalysis,
  Warning, RefreshRight, ArrowDown, ArrowUp, Remove, InfoFilled, Edit, AlarmClock,
  Check, CircleClose, Setting, Tools
} from '@element-plus/icons-vue'
import { ref, computed, watch, nextTick } from 'vue'
import { ElMessage } from 'element-plus'
import Cron from '@/components/Cron/index.vue'
import type { FormInstance, FormRules } from 'element-plus'
import { createTask, updateTask } from '@/api/task'

// 获取任务类型对应的图标
const getTypeIcon = (type: number) => {
  const iconMap = {
    1: Monitor,    // Shell脚本
    3: DataAnalysis  // Datax任务
  }
  return iconMap[type] || Document
}

const taskTypes = {
  1: 'Shell脚本',
  3: 'Datax任务'
}

const priorityOptions = [
  { label: '高优先级', value: 'high', icon: ArrowUp },
  { label: '中优先级', value: 'medium', icon: Remove },
  { label: '低优先级', value: 'low', icon: ArrowDown }
]

const props = defineProps<{
  visible: boolean
  task?: any
}>()

const emit = defineEmits(['update:visible', 'refresh'])

const visible = computed({
  get: () => props.visible,
  set: (value) => emit('update:visible', value)
})

const isEdit = computed(() => Boolean(props.task?.id))

const form = ref({
  id: undefined,
  name: '',
  type: undefined,
  priority: 'medium',
  description: '',
  status: 1,
  cronExpr: '',
  taskContent: '',
  taskParams: ''
})

const loading = ref(false)
const showCronDialog = ref(false)

const handleCronConfirm = (value: string) => {
  form.value.cronExpr = value
  showCronDialog.value = false
}

const handleCronClose = () => {
  showCronDialog.value = false
}

const formRef = ref<FormInstance>()

// 初始化表单数据
const initForm = () => {
  if (props.task) {
    const { id, name, type, priority, description, status, cronExpr, taskContent, taskParams } = props.task
    form.value = {
      id,
      name,
      type: Number(type),
      priority,
      description,
      status: Number(status),
      cronExpr,
      taskContent,
      taskParams
    }
  } else {
    form.value = {
      id: undefined,
      name: '',
      type: undefined,
      priority: 'medium',
      description: '',
      status: 1,
      cronExpr: '',
      taskContent: '',
      taskParams: ''
    }
  }
}

// 监听 task 变化
watch(() => props.task, () => {
  initForm()
}, { immediate: true })

// 监听弹窗显示
watch(() => visible.value, (val) => {
  if (val) {
    initForm()
  }
})

const rules: FormRules = {
  name: [
    { required: true, message: '请输入任务名称', trigger: 'blur' },
    { min: 2, max: 100, message: '长度在 2 到 100 个字符', trigger: 'blur' }
  ],
  type: [
    { required: true, message: '请选择任务类型', trigger: 'change' }
  ],
  priority: [
    { required: true, message: '请选择优先级', trigger: 'change' }
  ],
  cronExpr: [
    { required: true, message: '请输入Cron表达式', trigger: 'blur' }
  ],
  taskContent: [
    { required: true, message: '请输入任务内容', trigger: 'blur' }
  ],
  taskParams: [{ 
    validator: (rule: any, value: string, callback: any) => {
      if (!value) {
        callback()
        return
      }
      try {
        JSON.parse(value)
        callback()
      } catch (error) {
        callback(new Error('请输入有效的JSON格式'))
      }
    },
    trigger: 'blur'
  }]
}

const submitForm = async () => {
  if (!formRef.value) return
  
  await formRef.value.validate(async (valid, fields) => {
    if (valid) {
      loading.value = true
      try {
        const { id, ...data } = form.value
        if (id) {
          await updateTask(id, data)
          ElMessage.success('更新成功')
        } else {
          await createTask(data)
          ElMessage.success('创建成功')
        }
        visible.value = false
        emit('refresh')
      } catch (error) {
        console.error('提交失败:', error)
        ElMessage.error('操作失败，请重试')
      } finally {
        loading.value = false
      }
    } else {
      console.log('验证失败:', fields)
    }
  })
}

// 处理输入框失焦
const handleCronInputBlur = (e: FocusEvent) => {
  // 检查点击的元素是否在 cron 弹框内
  const cronDialog = document.querySelector('.cron-container')
  if (cronDialog && !cronDialog.contains(e.relatedTarget as Node)) {
    showCronDialog.value = false
  }
}

// 格式化JSON
const formatJson = () => {
  if (!form.value.taskParams) return
  try {
    const parsed = JSON.parse(form.value.taskParams)
    form.value.taskParams = JSON.stringify(parsed, null, 2)
  } catch (error) {
    ElMessage.warning('当前输入不是有效的JSON格式')
  }
}
</script>

<style lang="scss" scoped>
.task-form {
  :deep(.el-dialog) {
    border-radius: 16px;
    box-shadow: 0 12px 32px rgba(0, 0, 0, 0.1);
    overflow: hidden;
    max-width: 90vw;

    .el-dialog__header {
      margin: 0;
      padding: 24px 32px;
      border-bottom: 1px solid var(--el-border-color-lighter);
      background: var(--el-bg-color);
      
      .el-dialog__title {
        font-size: 20px;
        font-weight: 600;
        color: var(--el-text-color-primary);
        letter-spacing: -0.01em;
      }

      .el-dialog__headerbtn {
        top: 24px;
        right: 28px;
        font-size: 18px;
      }
    }

    .el-dialog__body {
      padding: 32px;
      background: var(--el-fill-color-blank);
    }

    .el-dialog__footer {
      padding: 20px 32px;
      border-top: 1px solid var(--el-border-color-lighter);
      background: var(--el-bg-color);
    }
  }

  .form-container {
    .form-section {
      background-color: var(--el-bg-color);
      border-radius: 12px;
      padding: 28px 32px;
      margin-bottom: 24px;
      transition: all 0.3s ease;
      border: 1px solid var(--el-border-color-lighter);

      &:last-child {
        margin-bottom: 0;
      }

      .section-title {
        font-size: 16px;
        font-weight: 600;
        color: var(--el-text-color-primary);
        margin-bottom: 24px;
        display: flex;
        align-items: center;
        gap: 8px;

        .el-icon {
          font-size: 18px;
          color: var(--el-color-primary);
          margin-right: 4px;
        }
      }
    }
  }

  :deep(.el-form-item) {
    margin-bottom: 24px;

    &:last-child {
      margin-bottom: 0;
    }

    .el-form-item__label {
      font-weight: 500;
      color: var(--el-text-color-regular);
      font-size: 14px;
      padding-right: 16px;
    }

    .el-form-item__content {
      .el-input__wrapper,
      .el-textarea__wrapper {
        box-shadow: 0 0 0 1px var(--el-border-color-lighter) inset;
        transition: all 0.2s ease;
        border-radius: 8px;
        padding: 4px 12px;

        &:hover {
          box-shadow: 0 0 0 1px var(--el-border-color-dark) inset;
        }

        &.is-focus {
          box-shadow: 0 0 0 2px var(--el-color-primary-light-8) inset !important;
        }
      }

      .el-textarea__inner {
        padding: 8px;
        font-family: "Menlo", "Monaco", "Consolas", monospace;
      }

      .el-select {
        width: 100%;
      }
    }
  }

  // 任务类型选项样式
  .task-type-option {
    display: flex;
    align-items: center;
    gap: 12px;
    padding: 10px 8px;
    border-radius: 6px;
    transition: all 0.2s ease;

    .el-icon {
      flex-shrink: 0;
      font-size: 18px;
      color: var(--el-text-color-secondary);
      transition: color 0.2s;
    }

    .label {
      flex: 1;
      font-size: 14px;
      color: var(--el-text-color-primary);
      transition: color 0.2s;
    }

    &:hover {
      background-color: var(--el-fill-color-light);
      .el-icon {
        color: var(--el-color-primary);
      }
    }
  }

  // 下拉菜单样式
  :deep(.el-select-dropdown) {
    border-radius: 8px;
    box-shadow: 0 4px 16px rgba(0, 0, 0, 0.08);
    padding: 6px;

    .el-select-dropdown__item {
      border-radius: 4px;
      margin: 2px 0;
      padding: 0 12px;
      height: auto;

      &.selected {
        background-color: var(--el-color-primary-light-9);
        
        .task-type-option {
          .el-icon {
            color: var(--el-color-primary);
          }

          .label {
            color: var(--el-color-primary);
            font-weight: 500;
          }
        }
      }

      &:hover {
        background-color: var(--el-fill-color-light);
      }
    }
  }

  // 优先级图标样式
  .priority-option {
    display: flex;
    align-items: center;
    gap: 8px;
    padding: 6px 0;

    .priority-icon {
      font-size: 16px;

      &.high {
        color: var(--el-color-danger);
      }

      &.medium {
        color: var(--el-color-warning);
      }

      &.low {
        color: var(--el-color-success);
      }
    }
  }

  // 底部按钮样式
  .dialog-footer {
    display: flex;
    justify-content: flex-end;
    gap: 12px;

    .el-button {
      padding: 10px 24px;
      border-radius: 8px;
      font-weight: 500;
      transition: all 0.2s ease;
      min-width: 100px;

      .el-icon {
        margin-right: 6px;
        font-size: 16px;
      }

      &.el-button--primary {
        background: linear-gradient(45deg, var(--el-color-primary), var(--el-color-primary-light-3));
        border: none;
        
        &:hover {
          background: linear-gradient(45deg, var(--el-color-primary-dark-2), var(--el-color-primary));
          transform: translateY(-1px);
        }
      }

      &.cancel-button {
        background-color: var(--el-fill-color-blank);
        border: 1px solid var(--el-border-color);
        color: var(--el-text-color-regular);

        &:hover {
          border-color: var(--el-border-color-darker);
          color: var(--el-text-color-primary);
          background-color: var(--el-fill-color);
        }
      }
    }
  }

  // Cron表达式输入框样式
  :deep(.el-input-group__append) {
    padding: 0;
    .el-button {
      border: none;
      height: 100%;
      border-radius: 0;
      padding: 8px 16px;
      
      .el-icon {
        margin-right: 4px;
      }

      &:hover {
        background-color: var(--el-fill-color-light);
        color: var(--el-color-primary);
      }
    }
  }
}
</style>
