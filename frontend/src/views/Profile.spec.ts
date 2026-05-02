import { flushPromises, mount } from '@vue/test-utils'
import { beforeEach, describe, expect, it, vi } from 'vitest'
import Profile from './Profile.vue'

const fetchProfile = vi.fn()
const updateProfileState = vi.fn()
const logout = vi.fn()

const authStoreMock = {
  user: { id: 1, username: 'kitty', email: 'kitty@example.com', avatar_url: '' },
  fetchProfile,
  updateProfileState,
  logout,
}

const authApiMock = vi.hoisted(() => ({
  updateProfile: vi.fn(),
  changePassword: vi.fn(),
  uploadAvatar: vi.fn(),
}))

const messageMock = vi.hoisted(() => ({
  success: vi.fn(),
  error: vi.fn(),
}))

vi.mock('@/stores/auth', () => ({
  useAuthStore: () => authStoreMock,
}))

vi.mock('@/api/auth', () => ({
  authApi: authApiMock,
}))

vi.mock('element-plus', () => ({
  ElMessage: messageMock,
}))

const stubs = {
  'el-card': { template: '<div><slot name="header" /><slot /></div>' },
  'el-form': { template: '<form><slot /></form>' },
  'el-form-item': { template: '<div><slot /></div>' },
  'el-input': {
    props: ['modelValue', 'name'],
    emits: ['update:modelValue'],
    template: '<input :name="name" :value="modelValue" @input="$emit(\'update:modelValue\', $event.target.value)" />',
  },
  'el-button': {
    props: ['loading'],
    emits: ['click'],
    template: '<button :disabled="loading" @click="$emit(\'click\')"><slot /></button>',
  },
  'el-upload': {
    props: ['beforeUpload', 'httpRequest'],
    template: '<button data-test="upload-avatar" @click="triggerUpload"><slot /></button>',
    methods: {
      async triggerUpload() {
        const file = new File(['avatar'], 'avatar.png', { type: 'image/png' })
        if (!this.beforeUpload || this.beforeUpload(file) !== false) {
          await this.httpRequest?.({ file })
        }
      },
    },
  },
  'el-avatar': { template: '<div><slot /></div>' },
}

describe('Profile.vue', () => {
  beforeEach(() => {
    authStoreMock.user = { id: 1, username: 'kitty', email: 'kitty@example.com', avatar_url: '' }
    fetchProfile.mockResolvedValue(authStoreMock.user)
    fetchProfile.mockClear()
    updateProfileState.mockClear()
    logout.mockClear()
    authApiMock.updateProfile.mockReset()
    authApiMock.changePassword.mockReset()
    authApiMock.uploadAvatar.mockReset()
    messageMock.success.mockReset()
    messageMock.error.mockReset()
  })

  it('loads and displays the current profile data on mount', async () => {
    const wrapper = mount(Profile, {
      global: { stubs },
    })

    await flushPromises()

    expect(fetchProfile).toHaveBeenCalled()
    expect((wrapper.find('input[name="email"]').element as HTMLInputElement).value).toBe('kitty@example.com')
  })

  it('updates profile and syncs auth store after save', async () => {
    authApiMock.updateProfile.mockResolvedValue({
      data: { id: 1, username: 'kitty-2', email: 'kitty-2@example.com', avatar_url: '' },
    })

    const wrapper = mount(Profile, { global: { stubs } })
    await flushPromises()

    await wrapper.find('input[name="username"]').setValue('kitty-2')
    await wrapper.find('input[name="email"]').setValue('kitty-2@example.com')
    await wrapper.find('[data-test="save-profile"]').trigger('click')

    expect(authApiMock.updateProfile).toHaveBeenCalledWith({
      username: 'kitty-2',
      email: 'kitty-2@example.com',
    })
    expect(updateProfileState).toHaveBeenCalled()
  })

  it('logs out after a successful password change', async () => {
    authApiMock.changePassword.mockResolvedValue({ data: { message: 'ok' } })
    logout.mockResolvedValue(undefined)

    const wrapper = mount(Profile, { global: { stubs } })
    await flushPromises()

    await wrapper.find('input[name="current_password"]').setValue('password123')
    await wrapper.find('input[name="new_password"]').setValue('new-password-123')
    await wrapper.find('input[name="confirm_password"]').setValue('new-password-123')
    await wrapper.find('[data-test="save-password"]').trigger('click')

    expect(logout).toHaveBeenCalled()
    expect(authApiMock.changePassword).toHaveBeenCalledWith({
      current_password: 'password123',
      new_password: 'new-password-123',
    })
  })

  it('uploads avatar and syncs auth store', async () => {
    authApiMock.uploadAvatar.mockResolvedValue({
      data: { id: 1, username: 'kitty', email: 'kitty@example.com', avatar_url: '/uploads/avatars/new.png' },
    })

    const wrapper = mount(Profile, { global: { stubs } })
    await flushPromises()

    await wrapper.find('[data-test="upload-avatar"]').trigger('click')
    await flushPromises()

    expect(authApiMock.uploadAvatar).toHaveBeenCalled()
    expect(updateProfileState).toHaveBeenCalledWith({
      id: 1,
      username: 'kitty',
      email: 'kitty@example.com',
      avatar_url: '/uploads/avatars/new.png',
    })
  })
})

