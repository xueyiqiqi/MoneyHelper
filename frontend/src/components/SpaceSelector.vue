<template>
  <div class="space-selector">
    <el-select
      v-model="currentSpaceId"
      placeholder="选择空间"
      :filterable="true"
      :clearable="true"
      @change="handleSpaceChange"
      class="space-select"
    >
      <el-option label="个人账本" value="personal">
        <div class="space-option">
          <span class="space-icon">📒</span>
          <span class="space-name">个人账本</span>
          <span class="space-badge personal">个人</span>
        </div>
      </el-option>
      <el-option
        v-for="space in spaces"
        :key="space.id"
        :label="space.name"
        :value="space.id"
      >
        <div class="space-option">
          <span class="space-icon">🏠</span>
          <span class="space-name">{{ space.name }}</span>
          <span :class="['space-badge', space.role]">{{ getRoleText(space.role) }}</span>
        </div>
      </el-option>
    </el-select>

    <el-button
      type="primary"
      text
      @click="showCreateDialog = true"
      class="add-space-btn"
    >
      <el-icon><Plus /></el-icon>
    </el-button>

    <el-dialog
      v-model="showCreateDialog"
      title="创建新空间"
      width="420px"
      :close-on-click-modal="false"
      append-to-body
      destroy-on-close
      @closed="handleDialogClosed"
    >
      <el-form
        ref="formRef"
        :model="createForm"
        :rules="rules"
        label-position="top"
      >
        <el-form-item label="空间名称" prop="name">
          <el-input
            v-model="createForm.name"
            placeholder="给空间起个名字"
            maxlength="30"
            show-word-limit
          />
        </el-form-item>
        <el-form-item label="空间描述 (可选)">
          <el-input
            v-model="createForm.description"
            type="textarea"
            :rows="3"
            placeholder="描述一下这个空间的用途"
            maxlength="200"
            show-word-limit
          />
        </el-form-item>
      </el-form>

      <template #footer>
        <el-button @click="showCreateDialog = false">取消</el-button>
        <el-button type="primary" :loading="loading" @click="handleCreate">
          创建空间
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted, watch } from 'vue'
import { useSpaceStore } from '@/stores/space'
import { spaceApi } from '@/api/space'
import { ElMessage, type FormInstance, type FormRules } from 'element-plus'
import { Plus } from '@element-plus/icons-vue'
import type { Space } from '@/types'

const spaceStore = useSpaceStore()

const spaces = computed(() => spaceStore.spaces)
const showCreateDialog = ref(false)
const loading = ref(false)
const formRef = ref<FormInstance>()

const currentSpaceId = ref<number | 'personal' | null>('personal')

const createForm = reactive({
  name: '',
  description: '',
})

const rules: FormRules = {
  name: [
    { required: true, message: '请输入空间名称', trigger: 'blur' },
    { min: 2, max: 30, message: '名称长度 2-30 个字符', trigger: 'blur' },
  ],
}

const resetCreateForm = () => {
  createForm.name = ''
  createForm.description = ''
  formRef.value?.clearValidate()
}

const handleDialogClosed = () => {
  loading.value = false
  resetCreateForm()
}

const getRoleText = (role: string) => {
  const map: Record<string, string> = {
    creator: '创建者',
    admin: '管理员',
    member: '成员',
    observer: '观察者',
  }
  return map[role] || role
}

const handleSpaceChange = async (value: number | 'personal' | null) => {
  if (value === 'personal' || value === null) {
    spaceStore.setCurrentSpace(null)
  } else {
    const space = spaces.value.find(s => s.id === value)
    if (space) {
      spaceStore.setCurrentSpace(space)
    }
  }
}

const handleCreate = async () => {
  const valid = await formRef.value?.validate().catch(() => false)
  if (!valid) return

  loading.value = true
  try {
    await spaceApi.create(createForm.name, createForm.description)
    ElMessage.success('空间创建成功')
    showCreateDialog.value = false
    await spaceStore.fetchSpaces()
    if (spaceStore.spaces.length > 0) {
      currentSpaceId.value = spaceStore.spaces[spaceStore.spaces.length - 1].id
    }
  } catch (error: any) {
    ElMessage.error(error.response?.data?.error || '创建失败')
  } finally {
    loading.value = false
  }
}

onMounted(async () => {
  await spaceStore.fetchSpaces()
  if (spaceStore.currentSpace) {
    currentSpaceId.value = spaceStore.currentSpace.id
  } else {
    currentSpaceId.value = 'personal'
  }
})

watch(() => spaceStore.currentSpace, (newSpace) => {
  if (newSpace) {
    currentSpaceId.value = newSpace.id
  } else {
    currentSpaceId.value = 'personal'
  }
})
</script>

<style scoped>
.space-selector {
  display: flex;
  align-items: center;
  gap: 8px;
  width: 100%;
}

.space-select {
  flex: 1;
  max-width: 280px;
}

.space-select :deep(.el-select__wrapper) {
  background: var(--bg-secondary);
  border: 1px solid var(--border-color);
  border-radius: var(--radius-lg);
  padding: 4px 12px;
  height: 40px;
}

.space-select :deep(.el-select__placeholder) {
  color: var(--text-secondary);
}

.space-option {
  display: flex;
  align-items: center;
  gap: 10px;
}

.space-icon {
  font-size: 16px;
}

.space-name {
  flex: 1;
  color: var(--text-primary);
}

.space-badge {
  font-size: 11px;
  padding: 2px 8px;
  border-radius: 10px;
  font-weight: 500;
}

.space-badge.personal {
  background: var(--color-primary-light);
  color: var(--color-primary);
}

.space-badge.creator {
  background: #fce4ec;
  color: #e91e63;
}

.space-badge.admin {
  background: #fff3e0;
  color: #ff9800;
}

.space-badge.member {
  background: var(--color-primary-light);
  color: var(--color-primary);
}

.space-badge.observer {
  background: var(--bg-secondary);
  color: var(--text-tertiary);
}

.add-space-btn {
  width: 40px;
  height: 40px;
  border-radius: var(--radius-lg);
  background: var(--color-primary);
  color: white;
  border: none;
}

.add-space-btn:hover {
  background: var(--color-primary-hover);
}
</style>

