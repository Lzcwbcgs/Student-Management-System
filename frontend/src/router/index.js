import { createRouter, createWebHistory } from 'vue-router'

const routes = [
  {
    path: '/',
    redirect: '/login'
  },
  {
    path: '/login',
    name: 'Login',
    component: () => import('@/views/login/index.vue'),
    meta: { title: '登录', noAuth: true }
  },
  {
    path: '/dashboard',
    component: () => import('../layout/index.vue'),
    redirect: '/dashboard/index',
    children: [
      {
        path: 'dashboard',
        name: 'Dashboard',
        component: () => import('../views/dashboard/index.vue'),
        meta: { title: '首页', icon: 'HomeFilled' }
      },
      {
        path: 'course',
        name: 'Course',
        component: () => import('../views/course/index.vue'),
        meta: { title: '课程管理', icon: 'Reading' }
      },
      {
        path: 'student',
        name: 'Student',
        component: () => import('../views/student/index.vue'),
        meta: { title: '学生管理', icon: 'User' }
      },
      {
        path: 'instructor',
        name: 'Instructor',
        component: () => import('../views/instructor/index.vue'),
        meta: { title: '教师管理', icon: 'UserFilled' }
      },
      {
        path: 'classroom',
        name: 'Classroom',
        component: () => import('../views/classroom/index.vue'),
        meta: { title: '教室管理', icon: 'School' }
      },
      {
        path: 'registration',
        name: 'Registration',
        component: () => import('../views/registration/index.vue'),
        meta: { title: '选课系统', icon: 'List' }
      }
    ]
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

// 路由守卫
router.beforeEach((to, from, next) => {
  const token = localStorage.getItem('token')
  const userType = localStorage.getItem('userType')
  
  // 如果前往登录页
  if (to.path === '/login') {
    if (token) {
      // 已登录则跳转到仪表盘
      next('/dashboard')
    } else {
      // 未登录则允许访问登录页
      next()
    }
    return
  }

  // 检查是否登录
  if (!token && to.path !== '/login') {
    next('/login')
    return
  }

  // 检查路由权限
  if (to.meta.roles && !to.meta.roles.includes(userType)) {
    next('/403')
    return
  }

  next()
})

export default router
