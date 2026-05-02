import { beforeEach, describe, expect, it, vi } from 'vitest'
import { createPinia, setActivePinia } from 'pinia'

const authApiMock = vi.hoisted(() => ({
  login: vi.fn(),
  getProfile: vi.fn(),
  logout: vi.fn(),
}))

vi.mock('@/api/auth', () => ({
  authApi: authApiMock,
}))

import { useAuthStore } from '@/stores/auth'
import { authApi } from '@/api/auth'

describe('auth store profile sync', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    localStorage.clear()
    authApiMock.login.mockReset()
    authApiMock.getProfile.mockReset()
    authApiMock.logout.mockReset()
  })

  it('fetches profile after login and stores the user', async () => {
    const store = useAuthStore()

    vi.mocked(authApi.login).mockResolvedValue({
      data: { access_token: 'access', refresh_token: 'refresh', token_type: 'bearer' },
    } as any)
    vi.mocked(authApi.getProfile).mockResolvedValue({
      data: { id: 1, username: 'kitty', email: 'kitty@example.com', avatar_url: '/uploads/avatars/1.png' },
    } as any)

    await store.login({ username: 'kitty', password: 'password123' })

    expect(store.user?.username).toBe('kitty')
    expect(store.user?.avatar_url).toBe('/uploads/avatars/1.png')
    expect(localStorage.getItem('access_token')).toBe('access')
    expect(localStorage.getItem('refresh_token')).toBe('refresh')
  })

  it('restores profile from stored token on initialize', async () => {
    localStorage.setItem('access_token', 'stored-access')
    localStorage.setItem('refresh_token', 'stored-refresh')
    const store = useAuthStore()

    vi.mocked(authApi.getProfile).mockResolvedValue({
      data: { id: 2, username: 'restored', email: 'restored@example.com', avatar_url: '/uploads/avatars/2.png' },
    } as any)

    await store.initializeAuth()

    expect(authApi.getProfile).toHaveBeenCalled()
    expect(store.user?.username).toBe('restored')
  })
})
