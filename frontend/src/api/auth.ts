import client from './client'
import type {
  ChangePasswordRequest,
  LoginRequest,
  LoginResponse,
  RegisterRequest,
  UpdateProfileRequest,
  User,
} from '@/types'

export const authApi = {
  register(data: RegisterRequest) {
    return client.post('/register', data)
  },

  login(data: LoginRequest) {
    return client.post<LoginResponse>('/token', data)
  },

  getProfile() {
    return client.get<User>('/profile')
  },

  updateProfile(data: UpdateProfileRequest) {
    return client.put<User>('/profile', data)
  },

  changePassword(data: ChangePasswordRequest) {
    return client.post('/profile/password', data)
  },

  uploadAvatar(file: File) {
    const formData = new FormData()
    formData.append('avatar', file)
    return client.post<User>('/profile/avatar', formData, {
      headers: { 'Content-Type': 'multipart/form-data' },
    })
  },

  refresh(refreshToken: string) {
    return client.post<LoginResponse>('/refresh', { refresh_token: refreshToken })
  },

  logout() {
    return client.post('/logout')
  },
}
