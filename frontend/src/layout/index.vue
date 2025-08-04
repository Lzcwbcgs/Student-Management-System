<template>
  <div class="app-wrapper">
    <!-- 侧边栏 -->
    <el-aside :width="isCollapse ? '64px' : '200px'" class="sidebar-container">
      <div class="logo-container">
        <img src="../assets/logo.png" alt="Logo" class="logo" />
        <span v-show="!isCollapse" class="title">教务管理系统</span>
      </div>
      <el-menu
        :default-active="$route.path"
        :collapse="isCollapse"
        :collapse-transition="false"
        router
        class="el-menu-vertical"
      >
        <el-menu-item index="/dashboard">
          <el-icon><HomeFilled /></el-icon>
          <template #title>首页</template>
        </el-menu-item>
        <el-menu-item index="/course">
          <el-icon><Reading /></el-icon>
          <template #title>课程管理</template>
        </el-menu-item>
        <el-menu-item index="/student">
          <el-icon><User /></el-icon>
          <template #title>学生管理</template>
        </el-menu-item>
        <el-menu-item index="/instructor">
          <el-icon><UserFilled /></el-icon>
          <template #title>教师管理</template>
        </el-menu-item>
        <el-menu-item index="/classroom">
          <el-icon><School /></el-icon>
          <template #title>教室管理</template>
        </el-menu-item>
        <el-menu-item index="/registration">
          <el-icon><List /></el-icon>
          <template #title>选课系统</template>
        </el-menu-item>
      </el-menu>
    </el-aside>

    <el-container class="main-container">
      <!-- 顶部导航栏 -->
      <el-header height="60px" class="navbar">
        <div class="left-area">
          <el-icon
            class="collapse-btn"
            @click="isCollapse = !isCollapse"
          >
            <component :is="isCollapse ? 'Expand' : 'Fold'" />
          </el-icon>
        </div>
        <div class="right-menu">
          <el-dropdown @command="handleCommand">
            <span class="user-avatar">
              <el-avatar :size="30" src="https://cube.elemecdn.com/3/7c/3ea6beec64369c2642b92c6726f1epng.png" />
              <span class="user-name">Admin</span>
            </span>
            <template #dropdown>
              <el-dropdown-menu>
                <el-dropdown-item command="profile">个人信息</el-dropdown-item>
                <el-dropdown-item command="logout">退出登录</el-dropdown-item>
              </el-dropdown-menu>
            </template>
          </el-dropdown>
        </div>
      </el-header>

      <!-- 主内容区 -->
      <el-main>
        <router-view></router-view>
      </el-main>
    </el-container>
  </div>
</template>

<style scoped>
.app-wrapper {
  display: flex;
  height: 100vh;
  width: 100%;
  overflow: hidden;
}

.sidebar-container {
  background-color: #1f2937;
  transition: width 0.3s;
  overflow-x: hidden;
}

.logo-container {
  height: 60px;
  padding: 10px;
  display: flex;
  align-items: center;
  background: #111827;
}

.logo {
  width: 40px;
  height: 40px;
}

.title {
  color: #fff;
  margin-left: 10px;
  font-size: 18px;
  font-weight: 600;
}

.el-menu-vertical {
  border-right: none;
  background-color: #1f2937;
}

.navbar {
  background-color: #fff;
  box-shadow: 0 1px 4px rgba(0,0,0,.1);
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 20px;
}

.collapse-btn {
  font-size: 20px;
  cursor: pointer;
  transition: color 0.3s;
}

.collapse-btn:hover {
  color: var(--el-color-primary);
}

.right-menu {
  display: flex;
  align-items: center;
}

.user-avatar {
  display: flex;
  align-items: center;
  cursor: pointer;
}

.user-name {
  margin-left: 8px;
  font-size: 14px;
}

.main-container {
  flex: 1;
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.el-main {
  padding: 20px;
  background-color: #f3f4f6;
  overflow-y: auto;
}
</style>
