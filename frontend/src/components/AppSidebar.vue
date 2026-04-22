<template>
  <aside class="app-sidebar">
    <el-menu
      :default-active="activeMenu"
      :collapse="collapsed"
      router
      class="sidebar-menu"
    >
      <el-menu-item index="/bills">
        <el-icon><Wallet /></el-icon>
        <template #title>账单</template>
      </el-menu-item>

      <el-menu-item index="/analysis">
        <el-icon><DataAnalysis /></el-icon>
        <template #title>AI 分析</template>
      </el-menu-item>

      <el-menu-item index="/spaces">
        <el-icon><House /></el-icon>
        <template #title>空间管理</template>
      </el-menu-item>

      <div class="sidebar-footer">
        <el-button
          :icon="collapsed ? Expand : Fold"
          text
          @click="toggleCollapse"
          class="collapse-btn"
        />
      </div>
    </el-menu>
  </aside>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useRoute } from 'vue-router'
import { useAppStore } from '@/stores/app'
import { Wallet, DataAnalysis, House, Expand, Fold } from '@element-plus/icons-vue'

const route = useRoute()
const appStore = useAppStore()

const activeMenu = computed(() => route.path)
const collapsed = computed(() => appStore.sidebarCollapsed)

const toggleCollapse = () => {
  appStore.toggleSidebar()
}
</script>

<style scoped>
.app-sidebar {
  width: auto;
  min-width: 64px;
  background: var(--bg-card);
  border-right: 1px solid var(--border-light);
  display: flex;
  flex-direction: column;
  transition: width var(--transition-normal);
}

.sidebar-menu {
  border-right: none;
  background: transparent;
  flex: 1;
}

.sidebar-menu:not(.el-menu--collapse) {
  width: 180px;
}

.el-menu-item {
  height: 48px;
  margin: 4px 8px;
  border-radius: var(--radius-md);
  color: var(--text-secondary);
  font-weight: 500;
  transition: all var(--transition-fast);
}

.el-menu-item:hover {
  background: var(--bg-secondary);
  color: var(--text-primary);
}

.el-menu-item.is-active {
  background: var(--color-primary-light);
  color: var(--color-primary);
}

.el-menu-item .el-icon {
  font-size: 20px;
}

.sidebar-footer {
  padding: 12px 8px;
  border-top: 1px solid var(--border-light);
}

.collapse-btn {
  width: 100%;
  height: 36px;
  border-radius: var(--radius-md);
  color: var(--text-tertiary);
}

.collapse-btn:hover {
  background: var(--bg-secondary);
  color: var(--text-primary);
}
</style>
