import client from './client'
import type { LoginRequest, LoginResponse, RegisterRequest } from '@/types'

export const authApi = {
  register(data: RegisterRequest) {
    return client.post('/register', data)
  },

  login(data: LoginRequest) {
    return client.post<LoginResponse>('/token', data)
  },

  refresh(refreshToken: string) {
    return client.post<LoginResponse>('/refresh', { refresh_token: refreshToken })
  },

  logout() {
    return client.post('/logout')
  },
}
