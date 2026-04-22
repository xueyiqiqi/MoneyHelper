<template>
  <div class="space-detail">
    <div class="header">
      <el-button @click="router.push('/spaces')">返回</el-button>
      <h2>{{ space?.name || '加载中...' }}</h2>
    </div>

    <el-card v-loading="loading">
      <template #header>
        <span>成员管理</span>
      </template>
      <el-table :data="members" stripe>
        <el-table-column prop="username" label="用户名" />
        <el-table-column prop="role" label="角色">
          <template #default="{ row }">
            <div v-if="editingMemberId === row.user_id && canEditMemberRole(row.role, row.user_id)" style="display: flex; gap: 8px; align-items: center;">
              <el-select v-model="editingRole" style="width: 120px">
                <el-option
                  v-for="role in assignableRoles"
                  :key="role"
                  :label="getRoleText(role)"
                  :value="role"
                />
              </el-select>
              <el-button size="small" type="primary" @click="handleUpdateRole(row.user_id)">保存</el-button>
              <el-button size="small" @click="cancelEditRole">取消</el-button>
            </div>
            <el-tag v-else :type="getRoleType(row.role)">
              {{ getRoleText(row.role) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="260">
          <template #default="{ row }">
            <el-button
              v-if="canEditMemberRole(row.role, row.user_id) && editingMemberId !== row.user_id"
              size="small"
              @click="startEditRole(row.user_id, row.role)"
            >
              修改角色
            </el-button>
            <el-button
              v-if="canRemoveMember(row.role, row.user_id)"
              type="danger"
              size="small"
              @click="handleRemove(row.user_id)"
            >
              移除
            </el-button>
          </template>
        </el-table-column>
      </el-table>

      <div v-if="canInviteMember" class="add-member">
        <el-input
          v-model="addMemberForm.email"
          placeholder="输入邮箱"
          style="width: 200px"
        />
        <el-select v-model="addMemberForm.role" placeholder="选择角色" style="width: 120px">
          <el-option
            v-for="role in assignableRoles"
            :key="role"
            :label="getRoleText(role)"
            :value="role"
          />
        </el-select>
        <el-button type="primary" @click="handleAddMember">添加</el-button>
      </div>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useSpaceStore } from '@/stores/space'
import { usePermission } from '@/composables/usePermission'
import { spaceApi } from '@/api/space'
import { ElMessage, ElMessageBox } from 'element-plus'
import type { Role, SpaceMember } from '@/types'

const route = useRoute()
const router = useRouter()
const spaceStore = useSpaceStore()
const { canInviteMember, canRemoveMember, canEditMemberRole, assignableRoles } = usePermission()

const space = computed(() => spaceStore.currentSpace)
const loading = ref(false)
const members = ref<SpaceMember[]>([])
const editingMemberId = ref<number | null>(null)
const editingRole = ref<Role>('member')

const addMemberForm = reactive<{ email: string; role: Role }>({
  email: '',
  role: 'member',
})

const ensureCurrentSpace = async () => {
  const id = Number(route.params.id)
  if (!id) return

  if (spaceStore.currentSpace?.id === id) {
    return
  }

  if (spaceStore.spaces.length === 0) {
    await spaceStore.fetchSpaces()
  }

  const matchedSpace = spaceStore.spaces.find((item) => item.id === id) || null
  spaceStore.setCurrentSpace(matchedSpace)
}

const fetchMembers = async () => {
  loading.value = true
  try {
    await ensureCurrentSpace()
    const id = Number(route.params.id)
    const { data } = await spaceApi.getMembers(id)
    members.value = data.members
  } finally {
    loading.value = false
  }
}

const handleAddMember = async () => {
  if (!addMemberForm.email) {
    ElMessage.warning('请输入邮箱')
    return
  }
  const id = Number(route.params.id)
  try {
    await spaceApi.addMember(id, addMemberForm.email, addMemberForm.role)
    ElMessage.success('添加成功')
    addMemberForm.email = ''
    addMemberForm.role = assignableRoles.value[0] || 'member'
    fetchMembers()
  } catch (error: any) {
    ElMessage.error(error.response?.data?.error || '添加失败')
  }
}

const startEditRole = (userId: number, role: Role) => {
  editingMemberId.value = userId
  editingRole.value = role
}

const cancelEditRole = () => {
  editingMemberId.value = null
  editingRole.value = 'member'
}

const handleUpdateRole = async (userId: number) => {
  const id = Number(route.params.id)
  try {
    await spaceApi.updateMemberRole(id, userId, editingRole.value)
    ElMessage.success('角色修改成功')
    cancelEditRole()
    fetchMembers()
  } catch (error: any) {
    ElMessage.error(error.response?.data?.error || '角色修改失败')
  }
}

const handleRemove = async (userId: number) => {
  await ElMessageBox.confirm('确定要移除该成员吗？', '提示')
  const id = Number(route.params.id)
  try {
    await spaceApi.removeMember(id, userId)
    ElMessage.success('移除成功')
    fetchMembers()
  } catch (error: any) {
    ElMessage.error(error.response?.data?.error || '移除失败')
  }
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

onMounted(fetchMembers)
</script>

<style scoped>
.space-detail .header {
  display: flex;
  align-items: center;
  gap: 20px;
  margin-bottom: 20px;
}

.add-member {
  margin-top: 20px;
  display: flex;
  gap: 10px;
}
</style>
