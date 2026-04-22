import client from './client'
import type { Space, SpaceMember } from '@/types'

export const spaceApi = {
  getList() {
    return client.get<Space[]>('/spaces')
  },

  create(name: string, description?: string) {
    return client.post<Space>('/spaces', { name, description })
  },

  getDetail(id: number) {
    return client.get<Space>(`/spaces/${id}`)
  },

  addMember(spaceId: number, email: string, role: string) {
    return client.post(`/spaces/${spaceId}/members`, { email, role })
  },

  removeMember(spaceId: number, userId: number) {
    return client.delete(`/spaces/${spaceId}/members/${userId}`)
  },

  leave(spaceId: number) {
    return client.delete(`/spaces/${spaceId}/leave`)
  },

  getMembers(spaceId: number) {
    return client.get<{ members: SpaceMember[] }>(`/spaces/${spaceId}/members`)
  },

  updateMemberRole(spaceId: number, userId: number, role: string) {
    return client.put(`/spaces/${spaceId}/members/${userId}`, { role })
  },
}
