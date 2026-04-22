<template>
  <header class="app-header">
    <div class="header-left">
      <div class="logo">
        <span class="logo-icon">📒</span>
        <span class="logo-text">财务助手</span>
      </div>
    </div>

    <div class="header-center">
      <SpaceSelector />
    </div>

    <div class="header-right">
      <!-- 主题切换 -->
      <button class="theme-toggle" @click="toggleTheme" :title="isDark ? '切换浅色' : '切换深色'">
        <span v-if="isDark" class="icon">☀️</span>
        <span v-else class="icon">🌙</span>
      </button>

      <!-- 用户菜单 -->
      <el-dropdown @command="handleCommand" trigger="click">
        <div class="user-info">
          <el-avatar :size="32" class="user-avatar">
            {{ username?.charAt(0)?.toUpperCase() || 'U' }}
          </el-avatar>
          <span class="user-name">{{ username }}</span>
          <el-icon><ArrowDown /></el-icon>
        </div>
        <template #dropdown>
          <el-dropdown-menu>
            <el-dropdown-item command="profile">
              <el-icon><User /></el-icon>
              个人中心
            </el-dropdown-item>
            <el-dropdown-item command="settings">
              <el-icon><Setting /></el-icon>
              设置
            </el-dropdown-item>
            <el-dropdown-item divided command="logout">
              <el-icon><SwitchButton /></el-icon>
              退出登录
            </el-dropdown-item>
          </el-dropdown-menu>
        </template>
      </el-dropdown>
    </div>
  </header>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { ArrowDown, User, Setting, SwitchButton } from '@element-plus/icons-vue'
import SpaceSelector from './SpaceSelector.vue'

const router = useRouter()
const authStore = useAuthStore()

const isDark = ref(false)

const username = computed(() => authStore.user?.username || '用户')

const toggleTheme = () => {
  isDark.value = !isDark.value
  if (window.toggleTheme) {
    window.toggleTheme()
  }
}

const handleCommand = async (command: string) => {
  if (command === 'logout') {
    await authStore.logout()
    router.push('/login')
  }
}

onMounted(() => {
  const theme = localStorage.getItem('theme')
  isDark.value = theme === 'dark'
})
</script>

<style scoped>
.app-header {
  height: 64px;
  background: var(--bg-card);
  border-bottom: 1px solid var(--border-light);
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 24px;
  position: sticky;
  top: 0;
  z-index: 100;
  backdrop-filter: blur(10px);
  background: rgba(var(--bg-card), 0.9);
}

.light .app-header {
  background: rgba(255, 255, 255, 0.9);
}

.dark .app-header {
  background: rgba(45, 43, 41, 0.9);
}

.header-left {
  flex: 0 0 auto;
}

.logo {
  display: flex;
  align-items: center;
  gap: 10px;
}

.logo-icon {
  font-size: 24px;
}

.logo-text {
  font-size: 20px;
  font-weight: 700;
  color: var(--text-primary);
  letter-spacing: -0.5px;
}

.header-center {
  flex: 1;
  display: flex;
  justify-content: center;
  max-width: 400px;
  margin: 0 auto;
}

.header-right {
  flex: 0 0 auto;
  display: flex;
  align-items: center;
  gap: 16px;
}

.theme-toggle {
  width: 36px;
  height: 36px;
  border-radius: var(--radius-md);
  border: 1px solid var(--border-color);
  background: var(--bg-secondary);
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all var(--transition-fast);
  font-size: 16px;
}

.theme-toggle:hover {
  background: var(--color-primary-light);
  border-color: var(--color-primary);
}

.user-info {
  display: flex;
  align-items: center;
  gap: 8px;
  cursor: pointer;
  padding: 6px 12px;
  border-radius: var(--radius-lg);
  transition: background var(--transition-fast);
}

.user-info:hover {
  background: var(--bg-secondary);
}

.user-avatar {
  background: var(--color-primary);
  color: white;
  font-weight: 600;
}

.user-name {
  color: var(--text-primary);
  font-weight: 500;
}
</style>
