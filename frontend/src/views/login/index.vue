<template>
  <div class="login-page">
    <div class="login-container">
      <div class="login-content">
        <!-- 左侧信息 -->
        <div class="login-banner">
          <h1 class="system-name">教务管理系统</h1>
          <p class="system-desc">新一代智能教学管理平台</p>
          <div class="decoration">
            <div class="circle circle-1"></div>
            <div class="circle circle-2"></div>
            <div class="circle circle-3"></div>
          </div>
        </div>

        <!-- 右侧登录表单 -->
        <div class="login-form-container">
          <h2 class="welcome-text">欢迎登录</h2>
          <p class="login-desc">请选择您的身份并登录系统</p>
          
          <el-form 
            ref="loginFormRef"
            :model="loginForm"
            :rules="rules"
            class="login-form"
          >
            <el-form-item class="login-type">
              <el-radio-group v-model="loginForm.type" size="large">
                <el-radio-button label="student">
                  <el-icon><User /></el-icon>
                  学生
                </el-radio-button>
                <el-radio-button label="instructor">
                  <el-icon><UserFilled /></el-icon>
                  教师
                </el-radio-button>
                <el-radio-button label="admin">
                  <el-icon><Key /></el-icon>
                  管理员
                </el-radio-button>
              </el-radio-group>
            </el-form-item>

            <el-form-item prop="username">
              <el-input
                v-model="loginForm.username"
                :placeholder="usernamePlaceholder"
                prefix-icon="User"
                size="large"
                clearable
              />
            </el-form-item>
            
            <el-form-item prop="password">
              <el-input
                v-model="loginForm.password"
                type="password"
                placeholder="请输入密码"
                prefix-icon="Lock"
                size="large"
                show-password
                clearable
                @keyup.enter="handleLogin"
              />
            </el-form-item>

            <el-form-item class="remember-me">
              <el-checkbox v-model="rememberMe">记住我</el-checkbox>
              <el-button link type="primary" class="forget-password">
                忘记密码？
              </el-button>
            </el-form-item>

            <el-form-item>
              <el-button
                type="primary"
                size="large"
                :loading="loading"
                class="login-button"
                @click="handleLogin"
              >
                登录
              </el-button>
            </el-form-item>
          </el-form>

          <div class="login-footer">
            <el-row justify="center" class="register-hint">
              <el-col :span="24" style="text-align: center">
                <p>还没有账号？
                  <el-button type="primary" link @click="handleRegister">
                    立即注册
                  </el-button>
                </p>
              </el-col>
            </el-row>
            <p class="copyright">© {{ currentYear }} MIT-Style Education System. All rights reserved.</p>
          </div>
        </div>
      </div>
    </div>

    <!-- 注册对话框 -->
    <el-dialog
      v-model="registerDialogVisible"
      title="用户注册"
      width="500px"
      destroy-on-close
    >
      <el-form
        ref="registerFormRef"
        :model="registerForm"
        :rules="registerRules"
        label-width="100px"
      >
        <el-form-item label="用户类型" prop="type">
          <el-radio-group v-model="registerForm.type">
            <el-radio-button label="student">学生</el-radio-button>
            <el-radio-button label="instructor">教师</el-radio-button>
          </el-radio-group>
        </el-form-item>

        <el-form-item :label="registerForm.type === 'student' ? '学号' : '工号'" prop="id">
          <el-input v-model="registerForm.id" />
        </el-form-item>

        <el-form-item label="姓名" prop="name">
          <el-input v-model="registerForm.name" />
        </el-form-item>

        <el-form-item label="密码" prop="password">
          <el-input v-model="registerForm.password" type="password" show-password />
        </el-form-item>

        <el-form-item label="确认密码" prop="confirmPassword">
          <el-input v-model="registerForm.confirmPassword" type="password" show-password />
        </el-form-item>

        <el-form-item label="院系" prop="department">
          <el-select v-model="registerForm.department" placeholder="请选择院系">
            <el-option
              v-for="dept in departmentOptions"
              :key="dept.name"
              :label="dept.name"
              :value="dept.name"
            />
          </el-select>
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="registerDialogVisible = false">取消</el-button>
          <el-button type="primary" @click="handleRegisterSubmit" :loading="registerLoading">
            注册
          </el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { User, UserFilled, Key } from '@element-plus/icons-vue'
import { login } from '@/api/auth'
import { getDepartmentList } from '@/api/department'

const router = useRouter()
const loginFormRef = ref(null)
const loading = ref(false)
const rememberMe = ref(false)

const loginForm = reactive({
  username: '',
  password: '',
  type: 'student'
})

const rules = {
  username: [
    { required: true, message: '请输入用户名', trigger: 'blur' },
    { min: 3, max: 20, message: '长度在 3 到 20 个字符', trigger: 'blur' }
  ],
  password: [
    { required: true, message: '请输入密码', trigger: 'blur' },
    { min: 6, max: 20, message: '长度在 6 到 20 个字符', trigger: 'blur' }
  ]
}

const usernamePlaceholder = computed(() => {
  const map = {
    student: '请输入学号',
    instructor: '请输入工号',
    admin: '请输入管理员账号'
  }
  return map[loginForm.type]
})

const currentYear = computed(() => new Date().getFullYear())

// 注册相关
const registerDialogVisible = ref(false)
const registerFormRef = ref(null)
const registerLoading = ref(false)
const departmentOptions = ref([])

const registerForm = reactive({
  type: 'student',
  id: '',
  name: '',
  password: '',
  confirmPassword: '',
  department: ''
})

const validateConfirmPassword = (rule, value, callback) => {
  if (value !== registerForm.password) {
    callback(new Error('两次输入的密码不一致'))
  } else {
    callback()
  }
}

const registerRules = {
  id: [
    { required: true, message: '请输入学号/工号', trigger: 'blur' },
    { min: 3, max: 20, message: '长度在 3 到 20 个字符', trigger: 'blur' }
  ],
  name: [
    { required: true, message: '请输入姓名', trigger: 'blur' },
    { min: 2, max: 20, message: '长度在 2 到 20 个字符', trigger: 'blur' }
  ],
  password: [
    { required: true, message: '请输入密码', trigger: 'blur' },
    { min: 6, max: 20, message: '长度在 6 到 20 个字符', trigger: 'blur' }
  ],
  confirmPassword: [
    { required: true, message: '请再次输入密码', trigger: 'blur' },
    { validator: validateConfirmPassword, trigger: 'blur' }
  ],
  department: [
    { required: true, message: '请选择院系', trigger: 'change' }
  ]
}

// 获取院系列表
const fetchDepartmentList = async () => {
  try {
    const response = await getDepartmentList()
    if (response.code === 200) {
      departmentOptions.value = response.data
    }
  } catch (error) {
    console.error('获取院系列表失败:', error)
  }
}

// 显示注册对话框
const handleRegister = () => {
  registerDialogVisible.value = true
}

// 提交注册
const handleRegisterSubmit = async () => {
  if (!registerFormRef.value) return

  await registerFormRef.value.validate(async (valid) => {
    if (valid) {
      registerLoading.value = true
      try {
        const response = await register(registerForm)
        if (response.code === 200) {
          ElMessage.success('注册成功，请登录')
          registerDialogVisible.value = false
          // 清空表单
          Object.keys(registerForm).forEach(key => {
            registerForm[key] = ''
          })
          registerForm.type = 'student'
        }
      } catch (error) {
        console.error('注册失败:', error)
        ElMessage.error(error.response?.data?.message || '注册失败')
      } finally {
        registerLoading.value = false
      }
    }
  })
}

// 登录处理
const handleLogin = async () => {
  if (!loginFormRef.value) return
  
  await loginFormRef.value.validate(async (valid) => {
    if (valid) {
      loading.value = true
      try {
        const response = await login({
          ...loginForm,
          remember: rememberMe.value
        })
        
        if (response.code === 200) {
          const { token, userInfo } = response.data
          // 保存登录信息
          localStorage.setItem('token', token)
          localStorage.setItem('userType', loginForm.type)
          localStorage.setItem('userInfo', JSON.stringify(userInfo))
          
          // 记住用户名
          if (rememberMe.value) {
            localStorage.setItem('remembered_username', loginForm.username)
            localStorage.setItem('remembered_type', loginForm.type)
          } else {
            localStorage.removeItem('remembered_username')
            localStorage.removeItem('remembered_type')
          }
          
          ElMessage.success('登录成功')
          router.push('/')
        }
      } catch (error) {
        console.error('登录失败:', error)
        ElMessage.error(error.response?.data?.message || '登录失败')
      } finally {
        loading.value = false
      }
    }
  })
}

// 初始化记住的用户名
const initRememberedUser = () => {
  const rememberedUsername = localStorage.getItem('remembered_username')
  const rememberedType = localStorage.getItem('remembered_type')
  if (rememberedUsername && rememberedType) {
    loginForm.username = rememberedUsername
    loginForm.type = rememberedType
    rememberMe.value = true
  }
}

initRememberedUser()
</script>

<style scoped>
.login-page {
  min-height: 100vh;
  background: linear-gradient(135deg, #f5f7fa 0%, #c3cfe2 100%);
  display: flex;
  align-items: center;
  justify-content: center;
}

.login-container {
  width: 1000px;
  height: 600px;
  background: #fff;
  border-radius: 16px;
  box-shadow: 0 8px 30px rgba(0, 0, 0, 0.1);
  overflow: hidden;
}

.login-content {
  display: flex;
  height: 100%;
}

.login-banner {
  flex: 1;
  background: linear-gradient(135deg, #a31f34 0%, #8a1a2b 100%);
  padding: 40px;
  color: #fff;
  position: relative;
  overflow: hidden;
  display: flex;
  flex-direction: column;
  justify-content: center;
}

.system-name {
  font-size: 36px;
  font-weight: 600;
  margin-bottom: 16px;
  position: relative;
  z-index: 1;
}

.system-desc {
  font-size: 18px;
  opacity: 0.9;
  margin-bottom: 40px;
  position: relative;
  z-index: 1;
}

.decoration {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
}

.circle {
  position: absolute;
  border-radius: 50%;
  background: rgba(255, 255, 255, 0.1);
}

.circle-1 {
  width: 200px;
  height: 200px;
  top: -100px;
  right: -100px;
}

.circle-2 {
  width: 300px;
  height: 300px;
  bottom: -150px;
  left: -150px;
}

.circle-3 {
  width: 150px;
  height: 150px;
  bottom: 50px;
  right: 50px;
}

.login-form-container {
  flex: 1;
  padding: 40px;
  display: flex;
  flex-direction: column;
}

.welcome-text {
  font-size: 24px;
  font-weight: 600;
  margin-bottom: 8px;
  color: #303133;
}

.login-desc {
  font-size: 14px;
  color: #909399;
  margin-bottom: 40px;
}

.login-form {
  width: 100%;
  max-width: 360px;
}

.login-type {
  margin-bottom: 24px;
  display: flex;
  justify-content: center;
}

.remember-me {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 24px;
}

.login-button {
  width: 100%;
  height: 44px;
  font-size: 16px;
  background: linear-gradient(135deg, #a31f34 0%, #8a1a2b 100%);
  border: none;
}

.login-button:hover {
  background: linear-gradient(135deg, #8a1a2b 0%, #731623 100%);
}

.login-footer {
  margin-top: auto;
  text-align: center;
}

.copyright {
  font-size: 12px;
  color: #909399;
}

:deep(.el-radio-button__inner) {
  display: flex;
  align-items: center;
  gap: 4px;
}

:deep(.el-input__wrapper) {
  padding-left: 11px;
}

:deep(.el-input__prefix) {
  font-size: 18px;
}
</style>
