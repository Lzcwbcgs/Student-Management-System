<template>
  <div class="course-container">
    <!-- 搜索和过滤区域 -->
    <el-card class="filter-card" shadow="hover">
      <el-form :inline="true" :model="searchForm" @submit.prevent="handleSearch">
        <el-form-item label="课程代码">
          <el-input v-model="searchForm.code" placeholder="请输入课程代码" clearable />
        </el-form-item>
        <el-form-item label="课程名称">
          <el-input v-model="searchForm.name" placeholder="请输入课程名称" clearable />
        </el-form-item>
        <el-form-item label="学分">
          <el-select v-model="searchForm.credits" placeholder="选择学分" clearable>
            <el-option v-for="i in 5" :key="i" :label="i" :value="i" />
          </el-select>
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="searchForm.status" placeholder="选择状态" clearable>
            <el-option label="开课中" value="active" />
            <el-option label="未开课" value="inactive" />
            <el-option label="已结课" value="completed" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSearch">
            <el-icon><Search /></el-icon>
            搜索
          </el-button>
          <el-button @click="resetSearch">
            <el-icon><Refresh /></el-icon>
            重置
          </el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <!-- 课程列表 -->
    <el-card class="list-card" shadow="hover">
      <template #header>
        <div class="card-header">
          <span class="header-title">课程列表</span>
          <el-button type="primary" @click="handleAdd">
            <el-icon><Plus /></el-icon>
            新增课程
          </el-button>
        </div>
      </template>

      <el-table
        v-loading="loading"
        :data="courseList"
        style="width: 100%"
        :default-sort="{ prop: 'code', order: 'ascending' }"
      >
        <el-table-column prop="code" label="课程代码" sortable />
        <el-table-column prop="name" label="课程名称" min-width="180" />
        <el-table-column prop="credits" label="学分" width="80" align="center" />
        <el-table-column prop="capacity" label="容量" width="100" align="center" />
        <el-table-column prop="enrolled" label="已选人数" width="100" align="center" />
        <el-table-column prop="instructors" label="授课教师" min-width="150">
          <template #default="{ row }">
            <el-tag
              v-for="instructor in row.instructors"
              :key="instructor.id"
              size="small"
              class="instructor-tag"
            >
              {{ instructor.name }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="100" align="center">
          <template #default="{ row }">
            <el-tag :type="getStatusType(row.status)">
              {{ getStatusText(row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="260" fixed="right">
          <template #default="{ row }">
            <el-button-group>
              <el-button type="primary" link @click="handleEdit(row)">
                <el-icon><Edit /></el-icon>
                编辑
              </el-button>
              <el-button type="primary" link @click="handleViewDetail(row)">
                <el-icon><View /></el-icon>
                详情
              </el-button>
              <el-button type="primary" link @click="handleAssignInstructor(row)">
                <el-icon><UserFilled /></el-icon>
                分配教师
              </el-button>
              <el-button type="danger" link @click="handleDelete(row)">
                <el-icon><Delete /></el-icon>
                删除
              </el-button>
            </el-button-group>
          </template>
        </el-table-column>
      </el-table>

      <!-- 分页 -->
      <div class="pagination-container">
        <el-pagination
          v-model:current-page="page"
          v-model:page-size="pageSize"
          :page-sizes="[10, 20, 50, 100]"
          :total="total"
          layout="total, sizes, prev, pager, next, jumper"
          @size-change="handleSizeChange"
          @current-change="handleCurrentChange"
        />
      </div>
    </el-card>

    <!-- 课程表单对话框 -->
    <el-dialog
      v-model="dialogVisible"
      :title="dialogType === 'add' ? '新增课程' : '编辑课程'"
      width="600px"
      destroy-on-close
    >
      <el-form
        ref="courseFormRef"
        :model="courseForm"
        :rules="courseRules"
        label-width="100px"
        label-position="right"
      >
        <el-form-item label="课程代码" prop="code">
          <el-input v-model="courseForm.code" placeholder="请输入课程代码" />
        </el-form-item>
        <el-form-item label="课程名称" prop="name">
          <el-input v-model="courseForm.name" placeholder="请输入课程名称" />
        </el-form-item>
        <el-form-item label="学分" prop="credits">
          <el-input-number v-model="courseForm.credits" :min="1" :max="5" />
        </el-form-item>
        <el-form-item label="课程容量" prop="capacity">
          <el-input-number v-model="courseForm.capacity" :min="1" :max="500" />
        </el-form-item>
        <el-form-item label="先修课程" prop="prerequisites">
          <el-select
            v-model="courseForm.prerequisites"
            multiple
            filterable
            placeholder="请选择先修课程"
            style="width: 100%"
          >
            <el-option
              v-for="course in prerequisiteOptions"
              :key="course.id"
              :label="course.name"
              :value="course.id"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="课程描述" prop="description">
          <el-input
            v-model="courseForm.description"
            type="textarea"
            rows="4"
            placeholder="请输入课程描述"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="dialogVisible = false">取消</el-button>
          <el-button type="primary" @click="handleSubmit" :loading="submitting">
            确定
          </el-button>
        </span>
      </template>
    </el-dialog>

    <!-- 分配教师对话框 -->
    <el-dialog
      v-model="assignDialogVisible"
      title="分配教师"
      width="500px"
      destroy-on-close
    >
      <el-transfer
        v-model="selectedInstructors"
        :data="instructorOptions"
        :titles="['可选教师', '已选教师']"
      />
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="assignDialogVisible = false">取消</el-button>
          <el-button
            type="primary"
            @click="handleAssignSubmit"
            :loading="submitting"
          >
            确定
          </el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  Search,
  Refresh,
  Plus,
  Edit,
  View,
  Delete,
  UserFilled
} from '@element-plus/icons-vue'
import {
  getCourseList,
  createCourse,
  updateCourse,
  deleteCourse,
  getCoursePrerequisites,
  getInstructorList,
  assignInstructor
} from '@/api/course'

// 搜索表单
const searchForm = reactive({
  code: '',
  name: '',
  credits: '',
  status: ''
})

// 分页参数
const page = ref(1)
const pageSize = ref(10)
const total = ref(0)

// 课程列表数据
const loading = ref(false)
const courseList = ref([])

// 课程表单
const dialogVisible = ref(false)
const dialogType = ref('add')
const courseFormRef = ref(null)
const courseForm = reactive({
  code: '',
  name: '',
  credits: 3,
  capacity: 100,
  prerequisites: [],
  description: ''
})

// 表单校验规则
const courseRules = {
  code: [
    { required: true, message: '请输入课程代码', trigger: 'blur' },
    { pattern: /^[A-Z0-9]+$/, message: '课程代码只能包含大写字母和数字', trigger: 'blur' }
  ],
  name: [
    { required: true, message: '请输入课程名称', trigger: 'blur' },
    { min: 2, max: 50, message: '长度在 2 到 50 个字符', trigger: 'blur' }
  ],
  credits: [
    { required: true, message: '请选择学分', trigger: 'change' }
  ],
  capacity: [
    { required: true, message: '请输入课程容量', trigger: 'change' }
  ]
}

// 先修课程选项
const prerequisiteOptions = ref([])

// 教师分配对话框
const assignDialogVisible = ref(false)
const selectedInstructors = ref([])
const instructorOptions = ref([])
const currentCourseId = ref(null)
const submitting = ref(false)

// 获取课程列表
const fetchCourseList = async () => {
  loading.value = true
  try {
    const params = {
      page: page.value,
      pageSize: pageSize.value,
      ...searchForm
    }
    const response = await getCourseList(params)
    if (response.code === 200) {
      courseList.value = response.data.list
      total.value = response.data.total
    }
  } catch (error) {
    console.error('获取课程列表失败:', error)
    ElMessage.error('获取课程列表失败')
  } finally {
    loading.value = false
  }
}

// 获取教师列表
const fetchInstructorList = async () => {
  try {
    const response = await getInstructorList()
    if (response.code === 200) {
      instructorOptions.value = response.data.map(item => ({
        key: item.id,
        label: item.name,
        disabled: false
      }))
    }
  } catch (error) {
    console.error('获取教师列表失败:', error)
  }
}

// 搜索处理
const handleSearch = () => {
  page.value = 1
  fetchCourseList()
}

// 重置搜索
const resetSearch = () => {
  Object.keys(searchForm).forEach(key => {
    searchForm[key] = ''
  })
  handleSearch()
}

// 分页处理
const handleSizeChange = (val) => {
  pageSize.value = val
  fetchCourseList()
}

const handleCurrentChange = (val) => {
  page.value = val
  fetchCourseList()
}

// 新增课程
const handleAdd = () => {
  dialogType.value = 'add'
  Object.keys(courseForm).forEach(key => {
    courseForm[key] = key === 'credits' ? 3 : key === 'capacity' ? 100 : ''
  })
  dialogVisible.value = true
}

// 编辑课程
const handleEdit = (row) => {
  dialogType.value = 'edit'
  Object.keys(courseForm).forEach(key => {
    courseForm[key] = row[key]
  })
  dialogVisible.value = true
}

// 查看详情
const handleViewDetail = (row) => {
  // TODO: 实现查看详情功能
}

// 删除课程
const handleDelete = (row) => {
  ElMessageBox.confirm(
    `确定要删除课程 "${row.name}" 吗？`,
    '警告',
    {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    }
  ).then(async () => {
    try {
      const response = await deleteCourse(row.id)
      if (response.code === 200) {
        ElMessage.success('删除成功')
        fetchCourseList()
      }
    } catch (error) {
      console.error('删除课程失败:', error)
      ElMessage.error('删除课程失败')
    }
  }).catch(() => {
    // 取消删除
  })
}

// 提交表单
const handleSubmit = async () => {
  if (!courseFormRef.value) return
  
  await courseFormRef.value.validate(async (valid) => {
    if (valid) {
      submitting.value = true
      try {
        const handler = dialogType.value === 'add' ? createCourse : updateCourse
        const response = await handler(courseForm)
        if (response.code === 200) {
          ElMessage.success(dialogType.value === 'add' ? '添加成功' : '更新成功')
          dialogVisible.value = false
          fetchCourseList()
        }
      } catch (error) {
        console.error('提交失败:', error)
        ElMessage.error('提交失败')
      } finally {
        submitting.value = false
      }
    }
  })
}

// 分配教师
const handleAssignInstructor = (row) => {
  currentCourseId.value = row.id
  selectedInstructors.value = row.instructors.map(item => item.id)
  assignDialogVisible.value = true
}

// 提交教师分配
const handleAssignSubmit = async () => {
  if (!currentCourseId.value) return
  
  submitting.value = true
  try {
    const response = await assignInstructor(currentCourseId.value, {
      instructorIds: selectedInstructors.value
    })
    if (response.code === 200) {
      ElMessage.success('分配成功')
      assignDialogVisible.value = false
      fetchCourseList()
    }
  } catch (error) {
    console.error('分配教师失败:', error)
    ElMessage.error('分配教师失败')
  } finally {
    submitting.value = false
  }
}

// 获取状态类型
const getStatusType = (status) => {
  const map = {
    active: 'success',
    inactive: 'info',
    completed: ''
  }
  return map[status] || 'info'
}

// 获取状态文本
const getStatusText = (status) => {
  const map = {
    active: '开课中',
    inactive: '未开课',
    completed: '已结课'
  }
  return map[status] || status
}

onMounted(() => {
  fetchCourseList()
  fetchInstructorList()
})
</script>

<style scoped>
.course-container {
  padding: 20px;
}

.filter-card {
  margin-bottom: 20px;
}

.list-card {
  margin-top: 20px;
}

.stats-card {
  height: 120px;
  display: flex;
  align-items: center;
  padding: 20px;
  box-sizing: border-box;
}

.stats-card-inner {
  display: flex;
  align-items: center;
  width: 100%;
}

.stats-icon-wrapper {
  width: 60px;
  height: 60px;
  border-radius: 50%;
  display: flex;
  justify-content: center;
  align-items: center;
  margin-right: 15px;
  font-size: 30px;
  color: #fff;
}

.stats-content {
  display: flex;
  flex-direction: column;
  justify-content: center;
  flex-grow: 1;
}

.stats-title {
  font-size: 16px;
  color: #666;
  margin-bottom: 5px;
}

.stats-value {
  font-size: 24px;
  font-weight: bold;
  color: #333;
}

.stats-trend {
  display: flex;
  align-items: center;
  font-size: 14px;
  margin-top: 5px;
}

.trend-icon {
  margin-right: 5px;
}

.trend-up {
  color: #67c23a;
}

.trend-down {
  color: #f56c6c;
}

.chart-card {
  height: 350px;
  margin-bottom: 20px;
}

.chart-container {
  width: 100%;
  height: 100%;
}

.activity-card {
  margin-bottom: 20px;
}

.activity-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 15px;
}

.activity-title {
  font-size: 18px;
  font-weight: bold;
}

.activity-table .el-tag {
  margin-right: 5px;
}

.empty-chart {
  display: flex;
  justify-content: center;
  align-items: center;
  height: 100%;
  color: #909399;
  font-size: 16px;
}


.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.header-title {
  font-size: 16px;
  font-weight: 500;
}

.instructor-tag {
  margin-right: 4px;
  margin-bottom: 4px;
}

.pagination-container {
  margin-top: 20px;
  display: flex;
  justify-content: flex-end;
}

:deep(.el-transfer) {
  width: 100%;
  display: flex;
  justify-content: space-between;
}

:deep(.el-transfer__buttons) {
  padding: 0 20px;
}

:deep(.el-transfer-panel) {
  width: 40%;
}

:deep(.el-button-group) .el-button--primary.is-link:first-child:not(:last-child) {
  border-right: 1px solid var(--el-button-border-color);
  padding-right: 8px;
  margin-right: 8px;
}
</style>

