<template>
  <div class="profile-page">
    <div class="profile-container">
      <section class="profile-section">
        <h2>个人资料</h2>
        <p>用户 ID：{{ user?.id }}</p>
        <div class="profile-fields">
          <label>
            用户名
            <input v-model="profileForm.username" name="username" />
          </label>
          <label>
            邮箱
            <input v-model="profileForm.email" name="email" />
          </label>
        </div>
        <el-button data-test="save-profile" :loading="loading" @click="handleProfileSave">保存资料</el-button>
      </section>

      <section class="profile-section">
        <h2>头像设置</h2>
        <div class="avatar-preview">{{ user?.avatar_url || '默认头像' }}</div>
        <el-upload :show-file-list="false" :before-upload="beforeAvatarUpload" :http-request="handleAvatarRequest">
          <el-button :loading="avatarLoading">上传头像</el-button>
        </el-upload>
      </section>

      <section class="profile-section">
        <h2>安全设置</h2>
        <div class="profile-fields">
          <label>
            原密码
            <input v-model="passwordForm.current_password" name="current_password" type="password" />
          </label>
          <label>
            新密码
            <input v-model="passwordForm.new_password" name="new_password" type="password" />
          </label>
          <label>
            确认新密码
            <input v-model="passwordForm.confirm_password" name="confirm_password" type="password" />
          </label>
        </div>
        <el-button data-test="save-password" :loading="passwordLoading" @click="handlePasswordSave">修改密码</el-button>
      </section>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, reactive, ref } from 'vue'
import { ElMessage } from 'element-plus'
import { authApi } from '@/api/auth'
import { useAuthStore } from '@/stores/auth'

const authStore = useAuthStore()
const loading = ref(false)
const passwordLoading = ref(false)
const avatarLoading = ref(false)

const user = computed(() => authStore.user)
const profileForm = reactive({
  username: '',
  email: '',
})
const passwordForm = reactive({
  current_password: '',
  new_password: '',
  confirm_password: '',
})

const fillProfileForm = () => {
  profileForm.username = user.value?.username || ''
  profileForm.email = user.value?.email || ''
}

const loadProfile = async () => {
  await authStore.fetchProfile()
  fillProfileForm()
}

const handleProfileSave = async () => {
  loading.value = true
  try {
    const { data } = await authApi.updateProfile({
      username: profileForm.username,
      email: profileForm.email,
    })
    authStore.updateProfileState(data)
    fillProfileForm()
    ElMessage.success('资料保存成功')
  } finally {
    loading.value = false
  }
}

const handlePasswordSave = async () => {
  if (passwordForm.new_password !== passwordForm.confirm_password) {
    ElMessage.error('两次输入的新密码不一致')
    return
  }

  passwordLoading.value = true
  try {
    await authApi.changePassword({
      current_password: passwordForm.current_password,
      new_password: passwordForm.new_password,
    })
    ElMessage.success('密码修改成功，请重新登录')
    await authStore.logout()
    window.location.href = '/login'
  } finally {
    passwordLoading.value = false
  }
}

const beforeAvatarUpload = (file: File) => {
  const isAllowed = ['image/jpeg', 'image/png', 'image/webp'].includes(file.type)
  const isLt2M = file.size / 1024 / 1024 < 2

  if (!isAllowed) {
    ElMessage.error('头像仅支持 JPG、PNG、WEBP 格式')
    return false
  }

  if (!isLt2M) {
    ElMessage.error('头像大小不能超过 2MB')
    return false
  }

  return true
}

const handleAvatarRequest = async ({ file }: { file: File }) => {
  avatarLoading.value = true
  try {
    const { data } = await authApi.uploadAvatar(file)
    authStore.updateProfileState(data)
    ElMessage.success('头像上传成功')
  } finally {
    avatarLoading.value = false
  }
}

onMounted(loadProfile)
</script>

<style scoped>
.profile-page {
  padding: 24px;
}

.profile-container {
  display: grid;
  gap: 20px;
}

.profile-section {
  background: var(--bg-card);
  border: 1px solid var(--border-light);
  border-radius: 16px;
  padding: 24px;
}

.profile-section h2 {
  margin: 0 0 12px;
  color: var(--text-primary);
}

.profile-section p {
  margin: 8px 0;
  color: var(--text-secondary);
}

.profile-fields {
  display: grid;
  gap: 12px;
  margin-bottom: 16px;
}

.profile-fields label {
  display: grid;
  gap: 6px;
  color: var(--text-secondary);
}

.profile-fields input {
  border: 1px solid var(--border-color);
  border-radius: 10px;
  padding: 10px 12px;
  background: var(--bg-secondary);
  color: var(--text-primary);
}

.avatar-preview {
  color: var(--text-secondary);
  margin-bottom: 12px;
}
</style>

