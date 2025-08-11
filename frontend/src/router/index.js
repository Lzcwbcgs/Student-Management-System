import { createRouter, createWebHistory } from 'vue-router'

const routes = [
  {
    path: '/',
    redirect: '/login'
  },
  {
    path: '/login',
    name: 'Login',
    component: () => import('../views/login/index.vue'),
    meta: { title: '登录', noAuth: true }
  },
  {
    path: '/dashboard',
    component: () => import('../layout/index.vue'),
    redirect: '/dashboard/index',
    children: [
      {
        path: 'index',
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
  const token = localStorage.getItem('token');
  const userType = localStorage.getItem('userType'); // 确保与登录时存储的key一致
  
  // 访问登录页的特殊处理
  if (to.path === '/login') {
    // 如果是从退出登录过来的（from.path是/dashboard），允许访问
    if (from.path === '/dashboard') {
      next();
    } 
    // 如果已登录但非退出操作，重定向到首页
    else if (token) {
      next('/dashboard');
    } else {
      next();
    }
    return;
  }

  // 非登录页需要认证
  if (!token) {
    next('/login');
    return;
  }

  next();
});

export default router
