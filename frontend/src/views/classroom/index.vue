<template>
  <div class="classroom-container">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>教室管理</span>
        </div>
      </template>
      <div class="classroom-content">
        <el-empty description="教室管理功能正在开发中..." />
      </div>
    </el-card>
  </div>
</template>

<script setup>
// 教室管理组件逻辑将在后续实现
</script>

<style scoped>
.classroom-container {
  padding: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.classroom-content {
  padding: 20px;
  text-align: center;
}
</style>
