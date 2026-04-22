import axios from 'axios'
import { useAuthStore } from '@/stores/auth'

const client = axios.create({
  baseURL: '/api/v1',
  timeout: 10000,
})

let refreshPromise: Promise<{ access_token: string; refresh_token: string }> | null = null

client.interceptors.request.use((config) => {
  const token = localStorage.getItem('access_token')
  if (token) {
    config.headers.Authorization = `Bearer ${token}`
  }
  return config
})

client.interceptors.response.use(
  (response) => response,
  async (error) => {
    const originalRequest = error.config
    const requestUrl = String(originalRequest?.url || '')
    const isAuthRefreshRequest = requestUrl.includes('/refresh')
    const isAuthLogoutRequest = requestUrl.includes('/logout')

    if (isAuthRefreshRequest || isAuthLogoutRequest) {
      return Promise.reject(error)
    }

    if (error.response?.status === 401 && originalRequest && !originalRequest._retry) {
      const storedRefreshToken = localStorage.getItem('refresh_token')
      const authStore = useAuthStore()

      if (!storedRefreshToken) {
        await authStore.logout()
        window.location.href = '/login'
        return Promise.reject(error)
      }

      originalRequest._retry = true

      try {
        if (!refreshPromise) {
          refreshPromise = axios
            .post('/api/v1/refresh', { refresh_token: storedRefreshToken })
            .then(({ data }) => ({
              access_token: data.access_token,
              refresh_token: data.refresh_token,
            }))
            .finally(() => {
              refreshPromise = null
            })
        }

        const tokens = await refreshPromise
        authStore.setTokens(tokens.access_token, tokens.refresh_token)
        originalRequest.headers.Authorization = `Bearer ${tokens.access_token}`
        return client(originalRequest)
      } catch (refreshError) {
        await authStore.logout()
        window.location.href = '/login'
        return Promise.reject(refreshError)
      }
    }

    return Promise.reject(error)
  }
)

export default client
