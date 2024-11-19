<template>
  <div class="chart-container">
    <div class="chart-header">
      <h4 class="chart-title">任务执行统计</h4>
      <el-radio-group v-model="period" @change="handlePeriodChange">
        <el-radio-button value="week">最近一周</el-radio-button>
        <el-radio-button value="month">最近一月</el-radio-button>
      </el-radio-group>
    </div>
    <div ref="chartRef" class="chart"></div>
  </div>
</template>

<script lang="ts">
import { defineComponent, ref, onMounted, watch } from 'vue'
import * as echarts from 'echarts'
import { getTaskChart } from '@/api/dashboard'
import type { ChartData } from '@/api/dashboard'

export default defineComponent({
  name: 'TaskChart',
  setup() {
    const chartRef = ref<HTMLElement>()
    let chart: echarts.ECharts | null = null
    const period = ref('week')
    const loading = ref(false)

    const initChart = () => {
      if (!chartRef.value) return
      chart = echarts.init(chartRef.value)
    }

    const fetchData = async () => {
      try {
        loading.value = true
        const res = await getTaskChart(period.value)
        if (res.code === 0 && res.data) {
          updateChart(res.data)
        }
      } catch (error) {
        console.error('Failed to fetch task chart data:', error)
      } finally {
        loading.value = false
      }
    }

    const updateChart = (data: ChartData[]) => {
      if (!chart) return

      const dates = data.map(item => item.date)
      const successData = data.map(item => item.success || 0)
      const totalData = data.map(item => {
        const success = item.success || 0
        const fail = item.fail || 0
        return success + fail
      })
      const successRates = data.map(item => {
        const success = item.success || 0
        const fail = item.fail || 0
        const total = success + fail
        return total > 0 ? ((success / total) * 100).toFixed(1) : '100.0'
      })

      // 计算最大任务量并向上取整到最近的50的倍数
      const maxTotal = Math.max(...totalData)
      const maxY = Math.ceil(maxTotal / 50) * 50

      const option = {
        title: {
          show: false
        },
        tooltip: {
          trigger: 'axis',
          axisPointer: {
            type: 'cross',
            crossStyle: {
              color: '#999'
            }
          }
        },
        grid: {
          top: '15%',
          left: '3%',
          right: '4%',
          bottom: '3%',
          containLabel: true
        },
        legend: {
          data: ['任务量', '成功率'],
          top: 5
        },
        xAxis: {
          type: 'category',
          data: dates,
          axisPointer: {
            type: 'shadow'
          }
        },
        yAxis: [
          {
            type: 'value',
            name: '任务量',
            nameLocation: 'end',
            nameGap: 40,
            min: 0,
            max: maxY,
            interval: 50,
            axisLabel: {
              formatter: '{value}'
            }
          },
          {
            type: 'value',
            name: '成功率',
            nameLocation: 'end',
            nameGap: 40,
            min: 0,
            max: 100,
            interval: 20,
            axisLabel: {
              formatter: '{value}%'
            }
          }
        ],
        series: [
          {
            name: '任务量',
            type: 'bar',
            yAxisIndex: 0,
            data: totalData,
            label: {
              show: true,
              position: 'top',
              formatter: '{c}',
              fontSize: 12,
              offset: [0, 0]
            },
            barWidth: '40%',
            itemStyle: {
              color: '#409EFF'
            }
          },
          {
            name: '成功率',
            type: 'line',
            yAxisIndex: 1,
            data: successRates,
            label: {
              show: false
            },
            emphasis: {
              label: {
                show: true,
                position: 'top',
                formatter: '{c}%',
                fontSize: 12
              }
            },
            symbol: 'circle',
            symbolSize: 8,
            itemStyle: {
              color: '#67C23A'
            },
            lineStyle: {
              width: 2
            }
          }
        ]
      }

      chart.setOption(option)
    }

    onMounted(() => {
      initChart()
      fetchData()
      window.addEventListener('resize', () => chart?.resize())
    })

    watch(period, () => {
      fetchData()
    })

    const handlePeriodChange = (newPeriod: string) => {
      period.value = newPeriod
    }

    return {
      chartRef,
      period,
      loading,
      handlePeriodChange
    }
  }
})
</script>

<style scoped>
.chart-container {
  width: 100%;
  padding: 20px;
}

.chart-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.chart-title {
  font-size: 16px;
  font-weight: 500;
  color: #303133;
  margin: 0;
}

.chart {
  width: 100%;
  height: 400px;
}
</style>
