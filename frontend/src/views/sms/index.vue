<template>
  <div class="sms-container">
    <!-- 搜索区域 -->
    <el-card class="search-card">
      <div class="search-header">
        <div class="title">
          <el-icon><Search /></el-icon>
          <span>筛选查询</span>
        </div>
        <div class="buttons">
          <el-button type="primary" @click="handleSearch">
            <el-icon><Search /></el-icon>
            <span>查询</span>
          </el-button>
          <el-button @click="resetSearch">
            <el-icon><Refresh /></el-icon>
            <span>重置</span>
          </el-button>
        </div>
      </div>
      <el-divider />
      <div class="search-content">
        <el-form :inline="true" :model="searchForm" class="search-form" label-position="top">
          <el-form-item label="手机号码">
            <el-input 
              v-model="searchForm.phone" 
              placeholder="请输入手机号码" 
              clearable
              @keyup.enter="handleSearch"
            >
              <template #prefix>
                <el-icon><Phone /></el-icon>
              </template>
            </el-input>
          </el-form-item>
          <el-form-item label="发送状态">
            <el-select 
              v-model="searchForm.status" 
              placeholder="请选择发送状态" 
              clearable
              class="w-200"
            >
              <el-option
                v-for="option in statusOptions"
                :key="option.value"
                :label="option.label"
                :value="option.value"
              >
                <template #default>
                  <div class="status-option">
                    <el-tag :type="getStatusType(option.value)" size="small" effect="light">
                      {{ option.label }}
                    </el-tag>
                  </div>
                </template>
              </el-option>
            </el-select>
          </el-form-item>
          <el-form-item label="发送时间">
            <el-date-picker
              v-model="searchForm.dateRange"
              type="daterange"
              range-separator="至"
              start-placeholder="开始日期"
              end-placeholder="结束日期"
              value-format="YYYY-MM-DD"
              class="w-300"
              :shortcuts="dateShortcuts"
            />
          </el-form-item>
        </el-form>
      </div>
      <div class="operation-bar">
        <el-divider />
        <div class="operation-buttons">
          <el-button type="primary" @click="handleAdd">
            <el-icon><Plus /></el-icon>
            <span>发送短信</span>
          </el-button>
          <el-button type="success" @click="handleExport">
            <el-icon><Download /></el-icon>
            <span>导出记录</span>
          </el-button>
        </div>
      </div>
    </el-card>

    <!-- 数据列表 -->
    <el-card class="list-card">
      <el-table
        v-loading="loading"
        :data="tableData"
        style="width: 100%"
        border
      >
        <el-table-column type="selection" width="55" />
        <el-table-column prop="id" label="ID" width="80" align="center" />
        <el-table-column prop="phone" label="手机号码" width="120" align="center" />
        <el-table-column prop="content" label="短信内容" min-width="200" show-overflow-tooltip>
          <template #default="{ row }">
            <div class="sms-content">
              <el-icon><Message /></el-icon>
              <span>{{ row.content }}</span>
            </div>
          </template>
        </el-table-column>
        <el-table-column prop="status" label="发送状态" width="100" align="center">
          <template #default="{ row }">
            <el-tag :type="getStatusType(row.status)">
              {{ getStatusText(row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="sendTime" label="发送时间" width="180" align="center">
          <template #default="{ row }">
            <div class="time-column">
              <el-icon><Timer /></el-icon>
              <span>{{ row.sendTime }}</span>
            </div>
          </template>
        </el-table-column>
        <el-table-column prop="operator" label="操作人" width="120" align="center">
          <template #default="{ row }">
            <div class="operator-column">
              <el-icon><User /></el-icon>
              <span>{{ row.operator }}</span>
            </div>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="180" align="center" fixed="right">
          <template #default="{ row }">
            <el-button
              link
              type="primary"
              @click="handleView(row)"
            >
              <el-icon><View /></el-icon>
              查看
            </el-button>
            <el-button
              v-if="row.status === 'failed'"
              link
              type="warning"
              @click="handleResend(row)"
            >
              <el-icon><RefreshRight /></el-icon>
              重发
            </el-button>
            <el-popconfirm
              title="确定要删除该条记录吗？"
              @confirm="handleDelete(row)"
            >
              <template #reference>
                <el-button
                  link
                  type="danger"
                >
                  <el-icon><Delete /></el-icon>
                  删除
                </el-button>
              </template>
            </el-popconfirm>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <!-- 分页 -->
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

    <!-- 发送短信表单 -->
    <el-dialog
      v-model="dialogVisible"
      :title="dialogTitle"
      width="500px"
      destroy-on-close
    >
      <el-form
        ref="formRef"
        :model="form"
        :rules="rules"
        label-width="100px"
      >
        <el-form-item label="手机号码" prop="phone">
          <el-input v-model="form.phone" placeholder="请输入手机号码">
            <template #prefix>
              <el-icon><Phone /></el-icon>
            </template>
          </el-input>
        </el-form-item>
        <el-form-item label="短信内容" prop="content">
          <el-input
            v-model="form.content"
            type="textarea"
            :rows="4"
            placeholder="请输入短信内容"
            show-word-limit
            maxlength="500"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <div class="dialog-footer">
          <el-button @click="dialogVisible = false">取消</el-button>
          <el-button type="primary" @click="handleSubmit">确定</el-button>
        </div>
      </template>
    </el-dialog>

    <!-- 查看短信详情 -->
    <el-dialog
      v-model="detailVisible"
      title="短信详情"
      width="600px"
    >
      <el-descriptions :column="1" border>
        <el-descriptions-item label="手机号码">{{ detail.phone }}</el-descriptions-item>
        <el-descriptions-item label="短信内容">{{ detail.content }}</el-descriptions-item>
        <el-descriptions-item label="发送状态">
          <el-tag :type="getStatusType(detail.status)">
            {{ getStatusText(detail.status) }}
          </el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="发送时间">{{ detail.sendTime }}</el-descriptions-item>
        <el-descriptions-item label="操作人">{{ detail.operator }}</el-descriptions-item>
        <el-descriptions-item v-if="detail.failReason" label="失败原因">
          {{ detail.failReason }}
        </el-descriptions-item>
      </el-descriptions>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue'
import type { FormInstance, FormRules } from 'element-plus'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  Plus,
  Search,
  Refresh,
  Download,
  Delete,
  View,
  RefreshRight,
  Message,
  Phone,
  Timer,
  User
} from '@element-plus/icons-vue'

// 状态选项
const statusOptions = [
  { label: '发送成功', value: 'success' },
  { label: '发送失败', value: 'failed' },
  { label: '待发送', value: 'pending' }
]

// 搜索表单
const searchForm = reactive({
  phone: '',
  status: '',
  dateRange: []
})

// 表格数据
const loading = ref(false)
const tableData = ref([])

// 分页
const currentPage = ref(1)
const pageSize = ref(10)
const total = ref(0)

// 弹窗表单
const dialogVisible = ref(false)
const dialogTitle = ref('发送短信')
const isEdit = ref(false)
const submitLoading = ref(false)
const formRef = ref<FormInstance>()
const form = reactive({
  phone: '',
  content: ''
})

// 表单校验规则
const rules: FormRules = {
  phone: [
    { required: true, message: '请输入手机号码', trigger: 'blur' },
    {
      pattern: /^(\d{11})(,\d{11})*$/,
      message: '手机号码格式不正确',
      trigger: 'blur'
    }
  ],
  content: [
    { required: true, message: '请输入短信内容', trigger: 'blur' },
    { min: 1, max: 500, message: '短信内容长度在1-500个字符', trigger: 'blur' }
  ]
}

// 详情弹窗
const detailVisible = ref(false)
const detail = ref({})

// 日期快捷选项
const dateShortcuts = [
  {
    text: '最近一周',
    value: () => {
      const end = new Date()
      const start = new Date()
      start.setTime(start.getTime() - 3600 * 1000 * 24 * 7)
      return [start, end]
    },
  },
  {
    text: '最近一个月',
    value: () => {
      const end = new Date()
      const start = new Date()
      start.setTime(start.getTime() - 3600 * 1000 * 24 * 30)
      return [start, end]
    },
  },
  {
    text: '最近三个月',
    value: () => {
      const end = new Date()
      const start = new Date()
      start.setTime(start.getTime() - 3600 * 1000 * 24 * 90)
      return [start, end]
    },
  }
]

// 获取状态类型
const getStatusType = (status: string) => {
  const map: Record<string, string> = {
    success: 'success',
    failed: 'danger',
    pending: 'warning'
  }
  return map[status] || 'info'
}

// 获取状态文本
const getStatusText = (status: string) => {
  const map: Record<string, string> = {
    success: '发送成功',
    failed: '发送失败',
    pending: '待发送'
  }
  return map[status] || '未知'
}

// 搜索
const handleSearch = () => {
  currentPage.value = 1
  fetchData()
}

// 重置搜索
const resetSearch = () => {
  searchForm.phone = ''
  searchForm.status = ''
  searchForm.dateRange = []
  handleSearch()
}

// 获取数据
const fetchData = async () => {
  try {
    loading.value = true
    const { startDate, endDate } = searchForm.dateRange || []
    const params = {
      page: currentPage.value,
      pageSize: pageSize.value,
      phone: searchForm.phone,
      status: searchForm.status,
      startDate,
      endDate
    }
    const res = await getSmsPage(params)
    tableData.value = res.data.list
    total.value = res.data.total
  } catch (error) {
    console.error('获取短信列表失败:', error)
    ElMessage.error('获取短信列表失败')
  } finally {
    loading.value = false
  }
}

// 分页大小变化
const handleSizeChange = (val: number) => {
  pageSize.value = val
  fetchData()
}

// 当前页变化
const handleCurrentChange = (val: number) => {
  currentPage.value = val
  fetchData()
}

// 新增短信
const handleAdd = () => {
  isEdit.value = false
  form.phone = ''
  form.content = ''
  dialogVisible.value = true
}

// 查看详情
const handleView = (row: any) => {
  detail.value = row
  detailVisible.value = true
}

// 重发短信
const handleResend = async (row: any) => {
  try {
    await ElMessageBox.confirm('确认重新发送该短信吗？')
    await resendSms(row.id)
    ElMessage.success('重发短信成功')
    fetchData()
  } catch (error) {
    console.error('重发短信失败:', error)
    if (error !== 'cancel') {
      ElMessage.error('重发短信失败')
    }
  }
}

// 删除短信
const handleDelete = async (row: any) => {
  try {
    await ElMessageBox.confirm('确定要删除该条记录吗？')
    await deleteSms(row.id)
    ElMessage.success('删除成功')
    if (tableData.value.length === 1 && currentPage.value > 1) {
      currentPage.value--
    }
    fetchData()
  } catch (error) {
    console.error('删除短信失败:', error)
    ElMessage.error('删除短信失败')
  }
}

// 导出记录
const handleExport = () => {
  ElMessage.warning('导出功能开发中')
}

// 提交表单
const handleSubmit = async () => {
  if (!formRef.value) return
  try {
    await formRef.value.validate()
    submitLoading.value = true
    await sendSms(form)
    ElMessage.success('发送成功')
    dialogVisible.value = false
    fetchData()
  } catch (error) {
    console.error('发送短信失败:', error)
    if (error !== 'cancel') {
      ElMessage.error('发送短信失败')
    }
  } finally {
    submitLoading.value = false
  }
}

// 初始化
fetchData()
</script>

<style scoped lang="scss">
.sms-container {
  padding: 20px;
  min-height: calc(100vh - 115px);
  display: flex;
  flex-direction: column;
  background: transparent;

  .search-card {
    margin-bottom: 20px;
    
    .search-header {
      display: flex;
      justify-content: space-between;
      align-items: center;
      padding: 0 10px;

      .title {
        display: flex;
        align-items: center;
        gap: 8px;
        font-size: 16px;
        font-weight: 500;
        color: var(--el-text-color-primary);

        .el-icon {
          font-size: 18px;
          color: var(--el-color-primary);
        }
      }

      .buttons {
        display: flex;
        gap: 12px;

        .el-button {
          display: flex;
          align-items: center;
          gap: 4px;
          padding: 8px 16px;
        }
      }
    }

    .search-content {
      padding: 16px 10px 0;

      .search-form {
        display: flex;
        flex-wrap: wrap;
        gap: 20px;

        :deep(.el-form-item) {
          margin-bottom: 0;
          margin-right: 0;

          .el-form-item__label {
            padding-bottom: 8px;
            font-weight: 500;
            color: var(--el-text-color-primary);
          }

          .el-input {
            width: 220px;
            
            .el-input__wrapper {
              box-shadow: 0 0 0 1px var(--el-border-color) inset;
              
              &:hover {
                box-shadow: 0 0 0 1px var(--el-color-primary) inset;
              }
            }
          }

          .el-select {
            width: 200px;
          }

          .el-date-editor {
            width: 300px;
          }
        }
      }
    }

    .operation-bar {
      .operation-buttons {
        padding: 16px 10px 0;
        display: flex;
        gap: 12px;

        .el-button {
          display: flex;
          align-items: center;
          gap: 4px;
          padding: 8px 16px;
        }
      }
    }
  }

  .list-card {
    flex: 1;
    margin-bottom: 20px;
  }

  .pagination-card {
    .pagination-container {
      padding: 10px;
      display: flex;
      justify-content: flex-end;
      align-items: center;
      
      :deep(.el-pagination) {
        .el-pagination__total {
          margin-right: 16px;
        }
        
        .el-pagination__sizes {
          margin-right: 16px;
        }
        
        .el-pager {
          li {
            &.is-active {
              background-color: var(--el-color-primary);
              color: #fff;
            }
          }
        }
        
        .btn-prev,
        .btn-next {
          padding: 0 6px;
          background: #fff;
          border: 1px solid #dcdfe6;
          border-radius: 4px;
          
          &:hover {
            color: var(--el-color-primary);
          }
          
          &.is-disabled {
            color: #c0c4cc;
            cursor: not-allowed;
          }
        }
      }
    }
  }

  :deep(.el-table) {
    flex: 1;

    .sms-content {
      display: flex;
      align-items: center;
      gap: 8px;

      .el-icon {
        color: var(--el-color-primary);
      }
    }

    .time-column,
    .operator-column {
      display: flex;
      align-items: center;
      justify-content: center;
      gap: 8px;

      .el-icon {
        font-size: 14px;
      }
    }

    .status-option {
      display: flex;
      align-items: center;
      gap: 8px;
    }
  }

  .dialog-footer {
    display: flex;
    justify-content: flex-end;
    gap: 12px;
    padding-top: 20px;
  }
}
</style>
