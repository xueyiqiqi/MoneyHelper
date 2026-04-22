import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { authApi } from '@/api/auth'
import type { User, LoginRequest } from '@/types'

export const useAuthStore = defineStore('auth', () => {
  const user = ref<User | null>(null)
  const token = ref<string | null>(localStorage.getItem('access_token'))
  const refreshToken = ref<string | null>(localStorage.getItem('refresh_token'))

  const isLoggedIn = computed(() => !!token.value)

  const setTokens = (access: string, refresh: string) => {
    token.value = access
    refreshToken.value = refresh
    localStorage.setItem('access_token', access)
    localStorage.setItem('refresh_token', refresh)
  }

  const login = async (credentials: LoginRequest) => {
    const { data } = await authApi.login(credentials)
    setTokens(data.access_token, data.refresh_token)
  }

  const logout = async () => {
    try {
      if (token.value) {
        await authApi.logout()
      }
    } finally {
      user.value = null
      token.value = null
      refreshToken.value = null
      localStorage.removeItem('access_token')
      localStorage.removeItem('refresh_token')
    }
  }

  return { user, token, isLoggedIn, setTokens, login, logout }
})
