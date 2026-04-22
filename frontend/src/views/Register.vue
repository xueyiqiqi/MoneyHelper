<template>
  <div class="auth-page">
    <div class="auth-container">
      <div class="auth-brand">
        <div class="brand-icon">📒</div>
        <h1 class="brand-name">财务助手</h1>
        <p class="brand-tagline">开启记账之旅</p>
      </div>

      <div class="auth-card">
        <div class="auth-tabs">
          <div :class="['auth-tab', { active: isLogin }]" @click="router.push('/login')">
            登录
          </div>
          <div :class="['auth-tab', { active: !isLogin }]" @click="router.push('/register')">
            注册
          </div>
        </div>

        <el-form
          ref="formRef"
          :model="form"
          :rules="rules"
          label-position="top"
          @submit.prevent="handleSubmit"
        >
          <el-form-item label="用户名" prop="username">
            <el-input v-model="form.username" placeholder="请输入用户名" size="large" />
          </el-form-item>

          <el-form-item label="邮箱" prop="email">
            <el-input v-model="form.email" placeholder="请输入邮箱" size="large" />
          </el-form-item>

          <el-form-item label="密码" prop="password">
            <el-input v-model="form.password" type="password" placeholder="请输入密码 (至少6位)" size="large" show-password />
          </el-form-item>

          <el-button type="primary" size="large" :loading="loading" native-type="submit" class="submit-btn">
            注 册
          </el-button>
        </el-form>

        <div class="auth-footer">
          <span>已有账号？</span>
          <a @click="router.push('/login')">立即登录</a>
        </div>
      </div>

      <div class="theme-toggle" @click="toggleTheme">
        <span v-if="isDark">☀️</span>
        <span v-else>🌙</span>
      </div>
    </div>

    <div class="bg-shapes">
      <div class="shape shape-1"></div>
      <div class="shape shape-2"></div>
      <div class="shape shape-3"></div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { authApi } from '@/api/auth'
import { ElMessage, type FormInstance, type FormRules } from 'element-plus'

const router = useRouter()

const formRef = ref<FormInstance>()
const loading = ref(false)
const isDark = ref(false)

const isLogin = computed(() => false)

const form = reactive({
  username: '',
  email: '',
  password: '',
})

const rules: FormRules = {
  username: [{ required: true, message: '请输入用户名', trigger: 'blur' }],
  email: [
    { required: true, message: '请输入邮箱', trigger: 'blur' },
    { type: 'email', message: '请输入正确的邮箱格式', trigger: 'blur' },
  ],
  password: [
    { required: true, message: '请输入密码', trigger: 'blur' },
    { min: 6, message: '密码至少6位', trigger: 'blur' },
  ],
}

const handleSubmit = async () => {
  const valid = await formRef.value?.validate().catch(() => false)
  if (!valid) return

  loading.value = true
  try {
    await authApi.register(form)
    ElMessage.success({ message: '注册成功，请登录', duration: 2000 })
    router.push('/login')
  } catch (error: any) {
    ElMessage.error(error.response?.data?.error || '注册失败')
  } finally {
    loading.value = false
  }
}

const toggleTheme = () => {
  isDark.value = !isDark.value
  if (window.toggleTheme) window.toggleTheme()
}

onMounted(() => {
  isDark.value = localStorage.getItem('theme') === 'dark'
})
</script>

<style scoped>
.auth-page {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  position: relative;
  overflow: hidden;
  background: var(--bg-primary);
}

.auth-container {
  position: relative;
  z-index: 10;
  width: 100%;
  max-width: 380px;
  padding: 20px;
}

.auth-brand {
  text-align: center;
  margin-bottom: 32px;
}

.brand-icon {
  font-size: 56px;
  margin-bottom: 12px;
}

.brand-name {
  font-size: 28px;
  font-weight: 800;
  color: var(--text-primary);
  margin: 0 0 8px;
  letter-spacing: -0.5px;
}

.brand-tagline {
  color: var(--text-tertiary);
  font-size: 14px;
  margin: 0;
}

.auth-card {
  background: var(--bg-card);
  border-radius: 24px;
  padding: 28px 24px;
  box-shadow: var(--shadow-lg);
  border: 1px solid var(--border-light);
}

.auth-tabs {
  display: flex;
  margin-bottom: 24px;
  background: var(--bg-secondary);
  border-radius: 12px;
  padding: 4px;
}

.auth-tab {
  flex: 1;
  text-align: center;
  padding: 10px;
  border-radius: 10px;
  font-weight: 600;
  cursor: pointer;
  color: var(--text-secondary);
  transition: all 0.2s ease;
}

.auth-tab.active {
  background: var(--bg-card);
  color: var(--color-primary);
  box-shadow: var(--shadow-sm);
}

.auth-card :deep(.el-form-item__label) {
  font-weight: 500;
  color: var(--text-secondary);
}

.auth-card :deep(.el-input__wrapper) {
  border-radius: 12px;
}

.submit-btn {
  width: 100%;
  height: 48px;
  font-size: 16px;
  font-weight: 700;
  border-radius: 12px;
  margin-top: 8px;
}

.auth-footer {
  text-align: center;
  margin-top: 20px;
  color: var(--text-secondary);
  font-size: 14px;
}

.auth-footer a {
  color: var(--color-primary);
  font-weight: 600;
  cursor: pointer;
  margin-left: 4px;
}

.theme-toggle {
  position: fixed;
  bottom: 24px;
  right: 24px;
  width: 44px;
  height: 44px;
  border-radius: 50%;
  background: var(--bg-card);
  border: 1px solid var(--border-color);
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 20px;
  cursor: pointer;
  box-shadow: var(--shadow-md);
}

.theme-toggle:hover {
  transform: scale(1.1);
}

.bg-shapes {
  position: absolute;
  inset: 0;
  pointer-events: none;
  overflow: hidden;
}

.shape {
  position: absolute;
  border-radius: 50%;
  opacity: 0.4;
}

.shape-1 {
  width: 300px;
  height: 300px;
  background: radial-gradient(circle, var(--color-primary-light) 0%, transparent 70%);
  top: -100px;
  right: -50px;
  animation: morph 20s ease-in-out infinite;
}

.shape-2 {
  width: 250px;
  height: 250px;
  background: radial-gradient(circle, rgba(232, 168, 124, 0.3) 0%, transparent 70%);
  bottom: -80px;
  left: -80px;
  animation: morph 15s ease-in-out infinite reverse;
}

.shape-3 {
  width: 150px;
  height: 150px;
  background: radial-gradient(circle, var(--color-primary) 0%, transparent 70%);
  top: 40%;
  left: 10%;
  animation: pulse 10s ease-in-out infinite;
}

@keyframes morph {
  0%, 100% { transform: scale(1) rotate(0deg); }
  50% { transform: scale(1.1) rotate(10deg); }
}

@keyframes pulse {
  0%, 100% { opacity: 0.2; transform: scale(1); }
  50% { opacity: 0.4; transform: scale(1.2); }
}
</style>
