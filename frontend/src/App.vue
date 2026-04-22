<template>
  <div :class="['app-theme', theme]">
    <router-view v-slot="{ Component }">
      <transition name="fade" mode="out-in">
        <component :is="Component" />
      </transition>
    </router-view>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, watch } from 'vue'
import { useAppStore } from '@/stores/app'

const appStore = useAppStore()
const theme = ref<'light' | 'dark'>('light')

// 从 localStorage 读取主题
onMounted(() => {
  const savedTheme = localStorage.getItem('theme') as 'light' | 'dark' | null
  if (savedTheme) {
    theme.value = savedTheme
  } else {
    // 跟随系统
    const prefersDark = window.matchMedia('(prefers-color-scheme: dark)').matches
    theme.value = prefersDark ? 'dark' : 'light'
  }
})

// 监听主题变化
watch(theme, (newTheme) => {
  localStorage.setItem('theme', newTheme)
  appStore.setTheme(newTheme)
})

// 暴露切换主题的方法供全局使用
const toggleTheme = () => {
  theme.value = theme.value === 'light' ? 'dark' : 'light'
}

// 将 toggleTheme 挂载到 window 以便在组件中调用
;(window as any).toggleTheme = toggleTheme
</script>

<style>
/* CSS 变量定义 - 温暖质感 */
:root,
.light {
  /* 主色调 - 暖白 */
  --bg-primary: #faf8f5;
  --bg-secondary: #f0ede8;
  --bg-card: #ffffff;
  --bg-input: #fdfcfa;

  /* 文字颜色 - 暖灰 */
  --text-primary: #3d3d3d;
  --text-secondary: #6b6b6b;
  --text-tertiary: #9a9a9a;
  --text-placeholder: #bfbfbf;

  /* 强调色 - 温暖绿 + 珊瑚 */
  --color-primary: #7eb08a;
  --color-primary-hover: #6a9c76;
  --color-primary-light: #e8f2eb;
  --color-accent: #e8a87c;
  --color-accent-hover: #d9956a;

  /* 功能色 */
  --color-success: #7eb08a;
  --color-warning: #e8a87c;
  --color-danger: #d98585;
  --color-info: #8faecb;

  /* 收入/支出 */
  --color-income: #7eb08a;
  --color-expense: #d98585;

  /* 边框 */
  --border-color: #e8e4df;
  --border-light: #f0ede8;

  /* 阴影 - 柔和纸质感 */
  --shadow-sm: 0 1px 3px rgba(0, 0, 0, 0.04);
  --shadow-md: 0 4px 12px rgba(0, 0, 0, 0.06);
  --shadow-lg: 0 8px 24px rgba(0, 0, 0, 0.08);

  /* 圆角 */
  --radius-sm: 8px;
  --radius-md: 12px;
  --radius-lg: 16px;
  --radius-xl: 24px;

  /* 过渡 */
  --transition-fast: 0.15s ease;
  --transition-normal: 0.25s ease;
  --transition-slow: 0.4s ease;

  /* 背景纹理 */
  --bg-pattern: url("data:image/svg+xml,%3Csvg width='60' height='60' viewBox='0 0 60 60' xmlns='http://www.w3.org/2000/svg'%3E%3Cg fill='none' fill-rule='evenodd'%3E%3Cg fill='%23e8e4df' fill-opacity='0.4'%3E%3Cpath d='M36 34v-4h-2v4h-4v2h4v4h2v-4h4v-2h-4zm0-30V0h-2v4h-4v2h4v4h2V6h4V4h-4zM6 34v-4H4v4H0v2h4v4h2v-4h4v-2H6zM6 4V0H4v4H0v2h4v4h2V6h4V4H6z'/%3E%3C/g%3E%3C/g%3E%3C/svg%3E");
}

/* 深色模式 - 温暖深色 */
.dark {
  /* 主色调 - 暖深色 */
  --bg-primary: #1a1918;
  --bg-secondary: #252422;
  --bg-card: #2d2b29;
  --bg-input: #353230;

  /* 文字颜色 - 暖白 */
  --text-primary: #f0ede8;
  --text-secondary: #a8a49c;
  --text-tertiary: #7a756d;
  --text-placeholder: #5a554d;

  /* 强调色 - 柔和绿 + 珊瑚 */
  --color-primary: #7eb08a;
  --color-primary-hover: #8fc99a;
  --color-primary-light: #3d4a3f;
  --color-accent: #e8a87c;
  --color-accent-hover: #f0bb92;

  /* 功能色 */
  --color-success: #7eb08a;
  --color-warning: #e8a87c;
  --color-danger: #d98585;
  --color-info: #8faecb;

  /* 收入/支出 */
  --color-income: #7eb08a;
  --color-expense: #d98585;

  /* 边框 */
  --border-color: #3d3936;
  --border-light: #353230;

  /* 阴影 */
  --shadow-sm: 0 1px 3px rgba(0, 0, 0, 0.2);
  --shadow-md: 0 4px 12px rgba(0, 0, 0, 0.3);
  --shadow-lg: 0 8px 24px rgba(0, 0, 0, 0.4);

  /* 背景纹理 */
  --bg-pattern: url("data:image/svg+xml,%3Csvg width='60' height='60' viewBox='0 0 60 60' xmlns='http://www.w3.org/2000/svg'%3E%3Cg fill='none' fill-rule='evenodd'%3E%3Cg fill='%233d3936' fill-opacity='0.4'%3E%3Cpath d='M36 34v-4h-2v4h-4v2h4v4h2v-4h4v-2h-4zm0-30V0h-2v4h-4v2h4v4h2V6h4V4h-4zM6 34v-4H4v4H0v2h4v4h2v-4h4v-2H6zM6 4V0H4v4H0v2h4v4h2V6h4V4H6z'/%3E%3C/g%3E%3C/g%3E%3C/svg%3E");
}

/* 全局样式 */
* {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
}

html, body {
  height: 100%;
}

body {
  font-family: -apple-system, BlinkMacSystemFont, 'PingFang SC', 'Microsoft YaHei', 'Segoe UI', Roboto, sans-serif;
  font-size: 14px;
  line-height: 1.6;
  color: var(--text-primary);
  background-color: var(--bg-primary);
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
}

#app {
  height: 100%;
}

/* 背景纹理 */
.app-theme {
  min-height: 100vh;
  background-color: var(--bg-primary);
  background-image: var(--bg-pattern);
}

/* 过渡动画 */
.fade-enter-active,
.fade-leave-active {
  transition: opacity var(--transition-normal), transform var(--transition-normal);
}

.fade-enter-from {
  opacity: 0;
  transform: translateY(8px);
}

.fade-leave-to {
  opacity: 0;
  transform: translateY(-8px);
}

/* 滚动条样式 */
::-webkit-scrollbar {
  width: 6px;
  height: 6px;
}

::-webkit-scrollbar-track {
  background: var(--bg-secondary);
}

::-webkit-scrollbar-thumb {
  background: var(--border-color);
  border-radius: 3px;
}

::-webkit-scrollbar-thumb:hover {
  background: var(--text-tertiary);
}

/* Element Plus 主题覆盖 */
.el-button--primary {
  --el-button-bg-color: var(--color-primary);
  --el-button-border-color: var(--color-primary);
  --el-button-hover-bg-color: var(--color-primary-hover);
  --el-button-hover-border-color: var(--color-primary-hover);
  --el-button-active-bg-color: var(--color-primary-hover);
  --el-button-active-border-color: var(--color-primary-hover);
}

.el-input__wrapper,
.el-textarea__inner {
  background-color: var(--bg-input) !important;
  box-shadow: none !important;
  border: 1px solid var(--border-color) !important;
}

.el-input__wrapper:hover,
.el-textarea__inner:hover {
  border-color: var(--color-primary) !important;
}

.el-input__wrapper.is-focus,
.el-textarea__inner:focus {
  border-color: var(--color-primary) !important;
  box-shadow: 0 0 0 3px var(--color-primary-light) !important;
}

.el-input__inner,
.el-textarea__inner {
  color: var(--text-primary) !important;
}

.el-input__inner::placeholder,
.el-textarea__inner::placeholder {
  color: var(--text-placeholder) !important;
}

.el-form-item__label {
  color: var(--text-secondary) !important;
}

.el-card {
  --el-card-bg-color: var(--bg-card);
  --el-card-border-color: var(--border-color);
  border-radius: var(--radius-lg) !important;
  box-shadow: var(--shadow-sm) !important;
  border: 1px solid var(--border-light) !important;
}

.el-card:hover {
  box-shadow: var(--shadow-md) !important;
  transition: box-shadow var(--transition-normal);
}

.el-table {
  --el-table-bg-color: var(--bg-card);
  --el-table-tr-bg-color: var(--bg-card);
  --el-table-header-bg-color: var(--bg-secondary);
  --el-table-row-hover-bg-color: var(--bg-secondary);
  --el-table-border-color: var(--border-light);
  --el-table-text-color: var(--text-primary);
  --el-table-header-text-color: var(--text-secondary);
}

.el-tag {
  border-radius: var(--radius-sm) !important;
}

.el-dialog {
  --el-dialog-bg-color: var(--bg-card);
  --el-dialog-title-font-size: 18px;
  border-radius: var(--radius-lg) !important;
}

.el-menu {
  --el-menu-bg-color: transparent;
  --el-menu-text-color: var(--text-secondary);
  --el-menu-hover-bg-color: var(--bg-secondary);
  --el-menu-active-color: var(--color-primary);
  border: none !important;
}

.el-menu-item.is-active {
  background-color: var(--color-primary-light) !important;
}

.el-dropdown-menu {
  --el-dropdown-menuItem-hover-fill: var(--bg-secondary);
  --el-dropdown-menuItem-hover-color: var(--text-primary);
}

.el-select-dropdown__item.is-selected {
  color: var(--color-primary) !important;
}

.el-pagination {
  --el-pagination-bg-color: var(--bg-card);
  --el-pagination-button-bg-color: var(--bg-card);
  --el-pagination-hover-color: var(--color-primary);
}

.el-message {
  border-radius: var(--radius-md) !important;
}

/* 日期选择器 */
.el-date-editor {
  --el-date-editor-width: 240px;
}

/* 链接 */
a {
  color: var(--color-primary);
  text-decoration: none;
  transition: color var(--transition-fast);
}

a:hover {
  color: var(--color-primary-hover);
}

/* 标题 */
h1, h2, h3, h4, h5, h6 {
  color: var(--text-primary);
  font-weight: 600;
}

/* 数字强调 */
.amount {
  font-family: 'SF Mono', 'Monaco', 'Menlo', monospace;
  font-weight: 500;
}

.amount.income {
  color: var(--color-income);
}

.amount.expense {
  color: var(--color-expense);
}
</style>
