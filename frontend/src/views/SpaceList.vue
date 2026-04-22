<template>
  <div class="space-list">
    <div class="header">
      <h2>我的空间</h2>
      <el-button type="primary" @click="showCreateDialog = true">
        创建空间
      </el-button>
    </div>

    <el-row :gutter="20">
      <el-col
        v-for="space in spaces"
        :key="space.id"
        :span="8"
      >
        <el-card class="space-card" @click="goToSpace(space)">
          <h3>{{ space.name }}</h3>
          <p>{{ space.description || '暂无描述' }}</p>
          <div class="space-meta">
            <el-tag :type="getRoleType(space.role)">
              {{ getRoleText(space.role) }}
            </el-tag>
            <span class="member-count">
              {{ space.member_count || 0 }} 人
            </span>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <!-- 创建空间对话框 -->
    <el-dialog v-model="showCreateDialog" title="创建空间" width="400px">
      <el-form
        ref="formRef"
        :model="createForm"
        :rules="rules"
        label-width="80px"
      >
        <el-form-item label="名称" prop="name">
          <el-input v-model="createForm.name" placeholder="空间名称" />
        </el-form-item>
        <el-form-item label="描述">
          <el-input
            v-model="createForm.description"
            type="textarea"
            placeholder="空间描述"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showCreateDialog = false">取消</el-button>
        <el-button type="primary" :loading="loading" @click="handleCreate">
          创建
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useSpaceStore } from '@/stores/space'
import { spaceApi } from '@/api/space'
import { ElMessage } from 'element-plus'
import type { Space, FormInstance } from 'element-plus'

const router = useRouter()
const spaceStore = useSpaceStore()

const spaces = ref<Space[]>([])
const showCreateDialog = ref(false)
const loading = ref(false)
const formRef = ref<FormInstance>()

const createForm = reactive({
  name: '',
  description: '',
})

const rules = {
  name: [{ required: true, message: '请输入空间名称', trigger: 'blur' }],
}

const fetchSpaces = async () => {
  await spaceStore.fetchSpaces()
  spaces.value = spaceStore.spaces
}

const handleCreate = async () => {
  const valid = await formRef.value?.validate().catch(() => false)
  if (!valid) return

  loading.value = true
  try {
    await spaceApi.create(createForm.name, createForm.description)
    ElMessage.success('创建成功')
    showCreateDialog.value = false
    fetchSpaces()
  } catch (error: any) {
    ElMessage.error(error.response?.data?.error || '创建失败')
  } finally {
    loading.value = false
  }
}

const goToSpace = (space: Space) => {
  spaceStore.setCurrentSpace(space)
  router.push(`/spaces/${space.id}`)
}

const getRoleType = (role: string) => {
  const map: Record<string, string> = {
    creator: 'danger',
    admin: 'warning',
    member: 'success',
    observer: 'info',
  }
  return map[role] || 'info'
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

onMounted(fetchSpaces)
</script>

<style scoped>
.space-list .header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.space-card {
  cursor: pointer;
  margin-bottom: 20px;
  transition: transform 0.2s;
}

.space-card:hover {
  transform: translateY(-4px);
}

.space-card h3 {
  margin: 0 0 10px;
}

.space-card p {
  color: #666;
  font-size: 14px;
  margin-bottom: 15px;
}

.space-meta {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.member-count {
  color: #999;
  font-size: 14px;
}
</style>
