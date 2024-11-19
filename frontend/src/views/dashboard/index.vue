<template>
  <div class="dashboard-container">
    <el-row :gutter="20">
      <el-col :span="8">
        <el-card class="box-card">
          <template #header>
            <div class="card-header">
              <span>今日短信发送量</span>
              <el-tag :type="overview.sms_trend >= 0 ? 'success' : 'danger'" size="small">
                {{ formatTrend(overview.sms_trend) }}
              </el-tag>
            </div>
          </template>
          <div class="card-value">{{ overview.sms_count }}</div>
        </el-card>
      </el-col>
      <el-col :span="8">
        <el-card class="box-card">
          <template #header>
            <div class="card-header">
              <span>今日任务数</span>
              <el-tag :type="overview.task_trend >= 0 ? 'success' : 'danger'" size="small">
                {{ formatTrend(overview.task_trend) }}
              </el-tag>
            </div>
          </template>
          <div class="card-value">{{ overview.task_count }}</div>
        </el-card>
      </el-col>
      <el-col :span="8">
        <el-card class="box-card">
          <template #header>
            <div class="card-header">
              <span>任务成功率</span>
              <el-tag :type="overview.success_rate_trend >= 0 ? 'success' : 'danger'" size="small">
                {{ formatTrend(overview.success_rate_trend) }}
              </el-tag>
            </div>
          </template>
          <div class="card-value">{{ overview.success_rate?.toFixed(1) || '0.0' }}%</div>
        </el-card>
      </el-col>
    </el-row>

    <el-row :gutter="20" class="chart-row">
      <el-col :span="12">
        <el-card class="box-card">
          <TaskChart />
        </el-card>
      </el-col>
      <el-col :span="12">
        <el-card class="box-card">
          <SmsChart />
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script lang="ts">
import { defineComponent, ref, onMounted } from 'vue'
import TaskChart from './components/TaskChart.vue'
import SmsChart from './components/SmsChart.vue'
import { getOverview } from '@/api/dashboard'
import type { DashboardOverview } from '@/api/dashboard'

export default defineComponent({
  name: 'Dashboard',
  components: {
    TaskChart,
    SmsChart
  },
  setup() {
    const loading = ref(true)
    const overview = ref<DashboardOverview>({
      sms_count: 0,
      sms_trend: 0,
      task_count: 0,
      task_trend: 0,
      success_rate: 0,
      success_rate_trend: 0
    })

    const fetchData = async () => {
      try {
        loading.value = true
        const { data } = await getOverview()
        overview.value = data
      } catch (error) {
        console.error('Failed to fetch dashboard data:', error)
      } finally {
        loading.value = false
      }
    }

    onMounted(() => {
      fetchData()
    })

    const formatTrend = (trend: number | undefined) => {
      if (trend === undefined || trend === null) {
        return '0.0%'
      }
      if (trend > 0) {
        return `+${trend.toFixed(1)}%`
      }
      return `${trend.toFixed(1)}%`
    }

    return {
      loading,
      overview,
      formatTrend
    }
  }
})
</script>

<style scoped>
.dashboard-container {
  padding: 20px;
}

.box-card {
  margin-bottom: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.card-value {
  font-size: 24px;
  font-weight: bold;
  text-align: center;
  margin-top: 10px;
}

.chart-row {
  margin-top: 20px;
}
</style>
