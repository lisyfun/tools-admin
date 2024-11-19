<template>
  <div class="task-container">
    <el-card class="search-card" :body-style="{ padding: '12px 20px' }">
      <div class="operation-bar">
        <div class="left-buttons">
          <el-button type="primary" @click="handleAdd" size="default">
            <el-icon><Plus /></el-icon>
            <span>新建任务</span>
          </el-button>
          <el-button
            type="danger"
            :disabled="!selectedTasks.length"
            @click="handleBatchDelete"
            size="default"
          >
            <el-icon><Delete /></el-icon>
            <span>批量删除</span>
          </el-button>
        </div>
        <div class="right-buttons">
          <el-button type="primary" @click="handleQuery" size="default">
            <el-icon><Search /></el-icon>
            查询
          </el-button>
          <el-button @click="resetQuery" size="default">
            <el-icon><RefreshRight /></el-icon>
            重置
          </el-button>
          <el-button
            type="primary"
            plain
            class="collapse-btn"
            @click="toggleCollapse"
            size="default"
          >
            <el-icon class="collapse-icon" :class="{ 'is-active': !isCollapse }">
              <ArrowDown />
            </el-icon>
            {{ isCollapse ? '展开' : '收起' }}
          </el-button>
        </div>
      </div>
      <el-divider v-show="!isCollapse" style="margin: 15px 0" />
      <div v-show="!isCollapse" class="search-content">
        <el-form :inline="true" :model="queryParams" class="search-form" label-position="right">
          <div class="form-left"></div>
          <div class="form-right">
            <el-form-item label="任务名称：" class="form-item">
              <el-input
                v-model="queryParams.name"
                placeholder="请输入任务名称"
                clearable
                @keyup.enter="handleQuery"
                class="custom-input"
                style="width: 180px"
              >
                <template #prefix>
                  <el-icon class="input-icon"><Document /></el-icon>
                </template>
              </el-input>
            </el-form-item>
            <el-form-item label="任务类型：" class="form-item">
              <el-select
                v-model="queryParams.type"
                placeholder="请选择任务类型"
                clearable
                class="custom-select"
                style="width: 180px"
                popper-class="custom-select-dropdown"
              >
                <el-option
                  v-for="option in taskTypeOptions"
                  :key="option.value"
                  :label="option.label"
                  :value="option.value"
                >
                  <template #default>
                    <div class="task-type-option">
                      <el-icon class="option-icon">
                        <component :is="option.icon" />
                      </el-icon>
                      <span>{{ option.label }}</span>
                    </div>
                  </template>
                </el-option>
              </el-select>
            </el-form-item>
          </div>
        </el-form>
      </div>
    </el-card>

    <el-card class="list-card">
      <el-table
        :data="taskList"
        style="width: 100%"
        v-loading="loading"
        @selection-change="handleSelectionChange"
      >
        <el-table-column type="selection" width="55" />
        <el-table-column v-if="false" prop="id" label="ID" width="80" align="center" header-align="center" />
        <el-table-column prop="name" label="任务名称" align="center" header-align="center" />
        <el-table-column prop="type" label="类型" align="center" header-align="center" width="120">
          <template #default="scope">
            <el-tag :type="taskTypeConfig[scope.row.type]?.tagType" effect="plain">
              <el-icon>
                <component :is="taskTypeConfig[scope.row.type]?.icon" />
              </el-icon>
              <span>{{ taskTypeConfig[scope.row.type]?.label }}</span>
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="description" label="任务描述" align="center" header-align="center" show-overflow-tooltip />
        <el-table-column prop="status" label="状态" align="center" header-align="center" width="100">
          <template #default="{ row }">
            <div class="status-switch-container">
              <el-switch
                v-model="row.status"
                :active-value="1"
                :inactive-value="2"
                @change="(val) => handleStatusChange(row, val)"
                active-text="启动"
                inactive-text="停止"
                :active-color="'#13ce66'"
                :inactive-color="'#ff4949'"
                inline-prompt
                class="status-switch"
              />
            </div>
          </template>
        </el-table-column>
        <el-table-column prop="execStatus" label="执行状态" align="center" header-align="center" width="100">
          <template #default="{ row }">
            <el-tag :type="getExecStatusType(row.execStatus)">
              {{ getExecStatusText(row.execStatus) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="priority" label="优先级" align="center" header-align="center" width="100">
          <template #default="{ row }">
            <el-tag
              :type="getPriorityType(row.priority)"
              effect="light"
              size="small"
            >
              <el-icon>
                <component :is="getPriorityIcon(row.priority)" />
              </el-icon>
              {{ getPriorityText(row.priority) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column
          prop="cronExpr"
          label="Cron表达式"
          align="center"
          header-align="center"
          width="180"
          show-overflow-tooltip
        >
          <template #default="{ row }">
            <el-tag size="small">{{ row.cronExpr }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="nextRunTime" label="下次执行" align="center" header-align="center" width="160">
          <template #default="{ row }">
            <el-button
              size="small"
              type="primary"
              @click="showNextRunTimes(row)"
            >
              查看执行时间
            </el-button>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="200" fixed="right">
          <template #default="{ row }">
            <el-button-group>
              <el-button
                size="small"
                type="primary"
                @click="handleEdit(row)"
              >
                <el-icon><Edit /></el-icon>
                编辑
              </el-button>
              <el-button
                size="small"
                type="danger"
                @click="handleDelete(row)"
              >
                <el-icon><Delete /></el-icon>
                删除
              </el-button>
            </el-button-group>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <el-card class="pagination-card">
      <div class="pagination-container">
        <el-pagination
          v-model:currentPage="currentPage"
          v-model:pageSize="pageSize"
          :page-sizes="[10, 20, 30, 50]"
          :small="false"
          :background="true"
          layout="total, sizes, prev, pager, next, jumper"
          :total="total"
          @size-change="handleSizeChange"
          @current-change="handleCurrentChange"
        />
      </div>
    </el-card>

    <TaskForm
      v-model:visible="showTaskForm"
      :task="currentTask"
      @refresh="fetchData"
    />

    <!-- 下次执行时间弹窗 -->
    <el-dialog
      v-model="nextRunTimesVisible"
      title="未来执行时间"
      width="400px"
    >
      <el-timeline>
        <el-timeline-item
          v-for="(time, index) in nextRunTimes"
          :key="index"
          :timestamp="formatDateTime(time)"
          placement="top"
        >
          <span>第 {{ index + 1 }} 次执行</span>
        </el-timeline-item>
      </el-timeline>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import {
  Plus,
  Search,
  Refresh,
  Delete,
  Edit,
  ArrowDown,
  Filter,
  Document,
  RefreshRight,
  Monitor,
  DataAnalysis,
  ArrowUp,
  Remove
} from '@element-plus/icons-vue'
import {
  ref,
  onMounted,
  computed
} from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { getTaskList, deleteTask, createTask, updateTask, batchDeleteTasks, updateTaskStatus, getTaskById, getNextRunTimes } from '@/api/task'
import type { Task, TaskStatus } from '@/api/task'
import TaskForm from './components/TaskForm.vue'

interface Task {
  id: number
  name: string
  type: number
  description: string
  status: TaskStatus
  priority: string
  cronExpr: string
  nextRunTime: string
  execStatus: number
}

const loading = ref(false)
const searchQuery = ref('')
const currentPage = ref(1)
const pageSize = ref(10)
const total = ref(0)
const taskList = ref<Task[]>([])
const taskFormRef = ref()
const showTaskForm = ref(false)
const currentTask = ref<any>(null)
const selectedTasks = ref<Task[]>([])

// 查询参数
const queryParams = ref({
  name: '',
  type: undefined
})

// 展开/收起状态
const isCollapse = ref(true)

// 切换展开/收起
const toggleCollapse = () => {
  isCollapse.value = !isCollapse.value
}

// 计算已选择的筛选条件数量
const filterCount = computed(() => {
  let count = 0
  if (queryParams.value.name) count++
  if (queryParams.value.type !== undefined) count++
  return count
})

// 任务类型配置
const taskTypeConfig = {
  1: { label: 'Shell脚本', icon: Monitor, tagType: 'primary' },
  3: { label: 'Datax任务', icon: DataAnalysis, tagType: 'warning' }
}

const taskTypeOptions = [
  { value: 1, label: 'Shell脚本', icon: Monitor },
  { value: 3, label: 'Datax任务', icon: DataAnalysis }
]

const getTypeIcon = (type: number) => {
  return taskTypeConfig[type]?.icon || Monitor
}

const getStatusType = (status: TaskStatus) => {
  const statusMap: Record<TaskStatus, string> = {
    1: 'running',    // running
    2: 'stopped', // stopped
  }
  return statusMap[status] || 'info'
}

const getStatusText = (status: TaskStatus) => {
  const statusMap: Record<TaskStatus, string> = {
    1: '已启动',
    2: '已停止',
  }
  return statusMap[status] || '未知'
}

// 优先级配置
const priorityConfig = {
  high: { label: '高', icon: ArrowUp, type: 'danger' },
  medium: { label: '中', icon: Remove, type: 'warning' },
  low: { label: '低', icon: ArrowDown, type: 'success' }
}

const getPriorityIcon = (priority: string) => {
  return priorityConfig[priority]?.icon || Remove
}

const getPriorityText = (priority: string) => {
  return priorityConfig[priority]?.label || '未知'
}

const getPriorityType = (priority: string) => {
  return priorityConfig[priority]?.type || 'info'
}

const getExecStatusType = (execStatus: number) => {
  const types = {
    1: 'info',    // 待执行
    2: 'warning', // 执行中
    3: 'success', // 执行成功
    4: 'danger'   // 执行失败
  }
  return types[execStatus] || 'info'
}

const getExecStatusText = (execStatus: number) => {
  const texts = {
    1: '待执行',
    2: '执行中',
    3: '执行成功',
    4: '执行失败'
  }
  return texts[execStatus] || '未知'
}

const formatDateTime = (time: string) => {
  if (!time) return ''
  return new Date(time).toLocaleString('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit',
    second: '2-digit',
    hour12: false
  })
}

const fetchData = async () => {
  loading.value = true
  try {
    const params = {
      page: currentPage.value,
      pageSize: pageSize.value,
      name: queryParams.value.name || '',
      type: queryParams.value.type !== undefined ? queryParams.value.type : ''
    }
    const res = await getTaskList(params)
    if (res.data) {
      taskList.value = res.data.list
      total.value = res.data.total
    }
  } catch (error: any) {
    ElMessage.error(error.response?.data?.message || '获取任务列表失败')
  } finally {
    loading.value = false
  }
}

const handleAdd = () => {
  currentTask.value = null
  showTaskForm.value = true
}

const handleEdit = (row: any) => {
  currentTask.value = row
  showTaskForm.value = true
}

const handleDelete = async (row: Task) => {
  try {
    await ElMessageBox.confirm('确定要删除该任务吗？', '提示', {
      type: 'warning'
    })
    await deleteTask(row.id)
    ElMessage.success('删除成功')
    fetchData()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('删除失败')
    }
  }
}

const handleFormSuccess = async (formData: any) => {
  showTaskForm.value = false
  fetchData()
}

const handleQuery = () => {
  currentPage.value = 1
  fetchData()
}

const resetQuery = () => {
  queryParams.value = {
    name: '',
    type: undefined
  }
  currentPage.value = 1
  fetchData()
}

const handleSizeChange = (val: number) => {
  pageSize.value = val
  fetchData()
}

const handleCurrentChange = (val: number) => {
  currentPage.value = val
  fetchData()
}

const handleBatchDelete = async () => {
  if (selectedTasks.value.length === 0) {
    ElMessage.warning('请选择要删除的任务')
    return
  }

  try {
    await ElMessageBox.confirm('确定要删除选中的任务吗？', '提示', {
      type: 'warning'
    })
    const ids = selectedTasks.value.map(task => task.id)
    await batchDeleteTasks(ids)
    ElMessage.success('删除成功')
    selectedTasks.value = []
    await fetchData()
  } catch (error: any) {
    if (error !== 'cancel') {
      ElMessage.error(error.response?.data?.error || '删除失败')
    }
  }
}

const handleSelectionChange = (selection: Task[]) => {
  selectedTasks.value = selection
}

const handleStatusChange = async (row: Task, status: TaskStatus) => {
  const originalStatus = row.status
  try {
    await updateTaskStatus(row.id, status)
    ElMessage.success(`任务${status === 1 ? '启动' : '停止'}成功`)
  } catch (error: any) {
    row.status = originalStatus
    ElMessage.error(error.response?.data?.message || '状态更新失败')
  }
}

const nextRunTimesVisible = ref(false)
const nextRunTimes = ref<string[]>([])

// 显示下次执行时间
const showNextRunTimes = async (row: Task) => {
  try {
    const { data } = await getNextRunTimes(row.cronExpr)
    nextRunTimes.value = data
    nextRunTimesVisible.value = true
  } catch (error) {
    ElMessage.error('获取执行时间失败')
  }
}

onMounted(() => {
  fetchData()
})
</script>

<style lang="scss" scoped>
@import '@/styles/_variables.scss';

.task-container {
  padding: $content-padding;
  min-height: calc(100vh - #{$navbar-height} - #{$main-padding} * 2);
  display: flex;
  flex-direction: column;
  gap: 12px;

  .search-card {
    margin-bottom: 0;
    transition: all 0.3s ease;
    border-radius: 0;
    box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.05);
    
    &:hover {
      box-shadow: 0 4px 16px 0 rgba(0, 0, 0, 0.1);
    }
    
    .operation-bar {
      display: flex;
      justify-content: space-between;
      align-items: center;
      min-height: 32px;
      padding: 4px 0;
      
      .left-buttons,
      .right-buttons {
        display: flex;
        gap: 8px;
        align-items: center;
        
        .el-button {
          padding: 6px 12px;
          height: 32px;
          font-size: 13px;
          transition: all 0.3s ease;
          
          &:hover {
            transform: translateY(-1px);
          }
          
          .el-icon {
            margin-right: 4px;
            font-size: 14px;
          }
          
          &[type="primary"] {
            background: linear-gradient(45deg, $color-primary, lighten($color-primary, 10%));
            border: none;
            
            &:hover {
              background: linear-gradient(45deg, darken($color-primary, 5%), lighten($color-primary, 5%));
            }
            
            &.collapse-btn {
              background: none;
              border: 1px solid $color-primary;
              color: $color-primary;
              padding: 6px 10px;
              
              &:hover {
                color: white;
                background: linear-gradient(45deg, $color-primary, lighten($color-primary, 10%));
              }
              
              .collapse-icon {
                margin-right: 2px;
                font-size: 12px;
                transition: transform 0.3s ease;
                
                &.is-active {
                  transform: rotate(180deg);
                }
              }
            }
          }
          
          &[type="danger"] {
            background: linear-gradient(45deg, #f56c6c, lighten(#f56c6c, 10%));
            border: none;
            
            &:hover {
              background: linear-gradient(45deg, darken(#f56c6c, 5%), lighten(#f56c6c, 5%));
            }
            
            &:disabled {
              background: #f5f7fa;
              border: 1px solid #dcdfe6;
              color: #c0c4cc;
            }
          }
        }
      }
    }
    
    .search-content {
      margin-top: 8px;
      
      .search-form {
        display: flex;
        justify-content: space-between;
        
        .form-left {
          flex: 1;
        }
        
        .form-right {
          display: flex;
          flex-wrap: wrap;
          gap: 12px 20px;
          
          .form-item {
            margin: 0;
            
            :deep(.el-form-item__label) {
              font-size: 13px;
              padding-right: 6px;
              color: #606266;
            }
            
            .custom-input,
            .custom-select {
              .el-input__wrapper {
                background-color: #f5f7fa;
                border-radius: 0;
                padding: 0 12px;
                height: 32px;
                box-shadow: none;
                border: 1px solid transparent;
                transition: all 0.3s ease;
                
                &:hover {
                  background-color: #fff;
                  border-color: $color-primary;
                }
                
                &.is-focus {
                  background-color: #fff;
                  border-color: $color-primary;
                  box-shadow: 0 0 0 1px rgba(64, 158, 255, 0.1);
                }
                
                .el-input__inner {
                  height: 32px;
                  line-height: 32px;
                  font-size: 13px;
                }
                
                .input-icon {
                  font-size: 14px;
                  color: #909399;
                }
              }
            }
          }
        }
      }
    }
  }

  .list-card {
    flex: 1;
    margin-bottom: 0;
    margin-top: 0;
    border-radius: 0;
    transition: all 0.3s ease;
    box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.05);
    
    &:hover {
      box-shadow: 0 4px 16px 0 rgba(0, 0, 0, 0.1);
    }
    
    :deep(.el-card__body) {
      padding: 12px;
    }
    
    .el-table {
      border-radius: 0;
      
      &::before,
      &::after {
        display: none;
      }
      
      .el-table__inner-wrapper::before {
        display: none;
      }
    }
  }

  .pagination-card {
    display: flex;
    justify-content: flex-end;
    margin-top: 16px;
    padding-top: 16px;
    border-top: 1px solid #ebeef5;
  }
}

// 任务类型选项样式
.task-type-option {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 6px 0;
  transition: all 0.3s ease;
  
  .option-icon {
    font-size: 16px;
    color: #909399;
    transition: all 0.3s ease;
  }
  
  &:hover {
    .option-icon {
      color: $color-primary;
      transform: scale(1.1);
    }
  }
}

:deep(.custom-select-dropdown) {
  border-radius: 6px;
  padding: 4px;
  
  .el-select-dropdown__item {
    border-radius: 4px;
    margin: 2px 0;
    
    &.selected {
      background-color: #ecf5ff;
      
      .option-icon {
        color: $color-primary;
      }
    }
    
    &:hover {
      background-color: #f5f7fa;
    }
  }
}
</style>
