export interface User {
  id: number
  username: string
  email: string
}

export interface Space {
  id: number
  name: string
  description: string
  role: Role
  member_count?: number
  created_at: string
}

export interface SpaceMember {
  user_id: number
  username: string
  role: Role
}

export type Role = 'creator' | 'admin' | 'member' | 'observer'

export interface Bill {
  id: number
  space_id?: number
  user_id: number
  amount: number
  type: 'income' | 'expense'
  category: string
  description?: string
  remarks?: string
  bill_date: string
  date?: string
  created_at: string
  creator_name?: string
}

export interface AnalysisReport {
  id: number
  space_id?: number
  user_id: number
  content: string
  period_start?: string
  period_end?: string
  created_at: string
}

export interface LoginRequest {
  username: string
  password: string
}

export interface LoginResponse {
  access_token: string
  refresh_token: string
  token_type: string
}

export interface RegisterRequest {
  username: string
  password: string
  email: string
}
