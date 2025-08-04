<template>
  <div class="student-container">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>学生管理</span>
        </div>
      </template>
      <div class="student-content">
        <el-empty description="学生管理功能正在开发中..." />
      </div>
    </el-card>
  </div>
</template>

<script setup>
// 学生管理组件逻辑将在后续实现
</script>

<style scoped>
.student-container {
  padding: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.student-content {
  padding: 20px;
  text-align: center;
}
</style>
