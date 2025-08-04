<template>
  <div class="dashboard-container">
    <!-- 统计卡片 -->
    <el-row :gutter="20">
      <el-col :span="6">
        <el-card class="stats-card" shadow="hover">
          <div class="stats-card-inner">
            <div class="stats-icon student">
              <el-icon><User /></el-icon>
            </div>
            <div class="stats-info">
              <div class="stats-title">学生总数</div>
              <div class="stats-value">{{ stats.studentCount }}</div>
              <div class="stats-trend" :class="{ 'up': stats.studentGrowth > 0 }">
                <el-icon><CaretTop /></el-icon>
                {{ stats.studentGrowth }}%
              </div>
            </div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card class="stats-card" shadow="hover">
          <div class="stats-card-inner">
            <div class="stats-icon teacher">
              <el-icon><UserFilled /></el-icon>
            </div>
            <div class="stats-info">
              <div class="stats-title">教师总数</div>
              <div class="stats-value">{{ stats.instructorCount }}</div>
              <div class="stats-trend" :class="{ 'up': stats.instructorGrowth > 0 }">
                <el-icon><CaretTop /></el-icon>
                {{ stats.instructorGrowth }}%
              </div>
            </div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card class="stats-card" shadow="hover">
          <div class="stats-card-inner">
            <div class="stats-icon course">
              <el-icon><Reading /></el-icon>
            </div>
            <div class="stats-info">
              <div class="stats-title">课程总数</div>
              <div class="stats-value">{{ stats.courseCount }}</div>
              <div class="stats-trend" :class="{ 'up': stats.courseGrowth > 0 }">
                <el-icon><CaretTop /></el-icon>
                {{ stats.courseGrowth }}%
              </div>
            </div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card class="stats-card" shadow="hover">
          <div class="stats-card-inner">
            <div class="stats-icon enrollment">
              <el-icon><List /></el-icon>
            </div>
            <div class="stats-info">
              <div class="stats-title">选课总数</div>
              <div class="stats-value">{{ stats.enrollmentCount }}</div>
              <div class="stats-trend" :class="{ 'up': stats.enrollmentGrowth > 0 }">
                <el-icon><CaretTop /></el-icon>
                {{ stats.enrollmentGrowth }}%
              </div>
            </div>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <!-- 图表区域 -->
    <el-row :gutter="20" class="chart-row">
      <el-col :span="12">
        <el-card class="chart-card">
          <template #header>
            <div class="card-header">
              <span>选课趋势</span>
              <el-radio-group v-model="chartTimeRange" size="small">
                <el-radio-button label="week">本周</el-radio-button>
                <el-radio-button label="month">本月</el-radio-button>
                <el-radio-button label="year">本年</el-radio-button>
              </el-radio-group>
            </div>
          </template>
          <div class="chart-container">
            <el-empty v-if="!enrollmentChart.length" description="暂无数据" />
            <div v-else ref="enrollmentChartRef" class="chart"></div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="12">
        <el-card class="chart-card">
          <template #header>
            <div class="card-header">
              <span>院系分布</span>
            </div>
          </template>
          <div class="chart-container">
            <el-empty v-if="!departmentChart.length" description="暂无数据" />
            <div v-else ref="departmentChartRef" class="chart"></div>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <!-- 最近活动 -->
    <el-row>
      <el-col :span="24">
        <el-card class="activity-card">
          <template #header>
            <div class="card-header">
              <span>最近活动</span>
              <el-button type="primary" link>查看全部</el-button>
            </div>
          </template>
          <el-table :data="recentActivities" style="width: 100%" v-loading="loading">
            <el-table-column prop="time" label="时间" width="180" />
            <el-table-column prop="type" label="类型" width="120">
              <template #default="scope">
                <el-tag :type="getActivityTypeTag(scope.row.type)">
                  {{ scope.row.type }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="content" label="内容" />
            <el-table-column prop="operator" label="操作人" width="120" />
          </el-table>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup>
import { ref, onMounted, onBeforeUnmount } from 'vue'
import { getSystemStats } from '@/api/dashboard'
import * as echarts from 'echarts'
import { User, UserFilled, Reading, List, CaretTop } from '@element-plus/icons-vue'

const stats = ref({
  studentCount: 0,
  studentGrowth: 0,
  instructorCount: 0,
  instructorGrowth: 0,
  courseCount: 0,
  courseGrowth: 0,
  enrollmentCount: 0,
  enrollmentGrowth: 0
})

const chartTimeRange = ref('month')
const enrollmentChart = ref([])
const departmentChart = ref([])
const enrollmentChartRef = ref(null)
const departmentChartRef = ref(null)
const loading = ref(true)
const recentActivities = ref([])

let enrollmentChartInstance = null
let departmentChartInstance = null

// 初始化数据
const initData = async () => {
  try {
    loading.value = true
    const response = await getSystemStats()
    if (response.code === 200) {
      stats.value = response.data.stats
      enrollmentChart.value = response.data.enrollmentTrend
      departmentChart.value = response.data.departmentDistribution
      recentActivities.value = response.data.recentActivities
      initCharts()
    }
  } catch (error) {
    console.error('获取数据失败:', error)
  } finally {
    loading.value = false
  }
}

// 初始化图表
const initCharts = () => {
  // 选课趋势图表
  if (enrollmentChartRef.value && enrollmentChart.value.length) {
    enrollmentChartInstance = echarts.init(enrollmentChartRef.value)
    const option = {
      tooltip: {
        trigger: 'axis'
      },
      grid: {
        left: '3%',
        right: '4%',
        bottom: '3%',
        containLabel: true
      },
      xAxis: {
        type: 'category',
        data: enrollmentChart.value.map(item => item.date)
      },
      yAxis: {
        type: 'value'
      },
      series: [{
        data: enrollmentChart.value.map(item => item.value),
        type: 'line',
        smooth: true,
        areaStyle: {
          opacity: 0.3
        },
        lineStyle: {
          width: 2
        },
        itemStyle: {
          color: '#a31f34'
        }
      }]
    }
    enrollmentChartInstance.setOption(option)
  }

  // 院系分布饼图
  if (departmentChartRef.value && departmentChart.value.length) {
    departmentChartInstance = echarts.init(departmentChartRef.value)
    const option = {
      tooltip: {
        trigger: 'item',
        formatter: '{a} <br/>{b}: {c} ({d}%)'
      },
      series: [{
        name: '院系分布',
        type: 'pie',
        radius: ['50%', '70%'],
        avoidLabelOverlap: false,
        itemStyle: {
          borderRadius: 10,
          borderColor: '#fff',
          borderWidth: 2
        },
        label: {
          show: false,
          position: 'center'
        },
        emphasis: {
          label: {
            show: true,
            fontSize: 20,
            fontWeight: 'bold'
          }
        },
        labelLine: {
          show: false
        },
        data: departmentChart.value
      }]
    }
    departmentChartInstance.setOption(option)
  }
}

// 获取活动类型对应的标签类型
const getActivityTypeTag = (type) => {
  const map = {
    '选课': '',
    '退课': 'info',
    '成绩录入': 'success',
    '课程变更': 'warning',
    '系统维护': 'danger'
  }
  return map[type] || ''
}

// 监听窗口大小变化
const handleResize = () => {
  enrollmentChartInstance?.resize()
  departmentChartInstance?.resize()
}

onMounted(() => {
  initData()
  window.addEventListener('resize', handleResize)
})

onBeforeUnmount(() => {
  window.removeEventListener('resize', handleResize)
  enrollmentChartInstance?.dispose()
  departmentChartInstance?.dispose()
})
</script>

<style scoped>
.dashboard-container {
  padding: 20px;
}

.chart-row {
  margin-top: 20px;
}

.stats-card {
  height: 120px;
  transition: all 0.3s;
}

.stats-card-inner {
  display: flex;
  align-items: center;
  height: 100%;
}

.stats-icon {
  width: 64px;
  height: 64px;
  border-radius: 16px;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-right: 16px;
}

.stats-icon :deep(.el-icon) {
  font-size: 32px;
  color: #fff;
}

.stats-icon.student {
  background-color: #a31f34;
}

.stats-icon.teacher {
  background-color: #1f77b4;
}

.stats-icon.course {
  background-color: #2ca02c;
}

.stats-icon.enrollment {
  background-color: #9467bd;
}

.stats-info {
  flex: 1;
}

.stats-title {
  font-size: 14px;
  color: #909399;
  margin-bottom: 8px;
}

.stats-value {
  font-size: 24px;
  font-weight: bold;
  color: #303133;
  margin-bottom: 8px;
}

.stats-trend {
  font-size: 12px;
  color: #f56c6c;
  display: flex;
  align-items: center;
}

.stats-trend.up {
  color: #67c23a;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.chart-card {
  height: 400px;
}

.chart-container {
  height: 320px;
}

.chart {
  height: 100%;
}

.activity-card {
  margin-top: 20px;
}
</style>

