import client from './client'
import type { Bill } from '@/types'

type ReportScopeParams = {
  spaceId?: number
  isPersonal: boolean
}

type AnalyzeParams = ReportScopeParams & {
  period?: string
}

export const billApi = {
  getList(params?: { space_id?: number; is_personal?: boolean; start_date?: string; end_date?: string; type?: string; category?: string; page?: number; page_size?: number }) {
    // 处理 is_personal 参数
    const queryParams: Record<string, unknown> = { ...params }
    if (params?.is_personal !== undefined) {
      queryParams.is_personal = params.is_personal.toString()
    }

    return client.get('/bills', { params: queryParams })
  },

  getById(id: number) {
    return client.get(`/bills/${id}`)
  },

  create(data: { is_personal?: boolean; space_id?: number; amount: number; type: string; category: string; description: string; bill_date: string }) {
    return client.post('/bills', {
      amount: data.amount,
      category: data.category,
      type: data.type,
      description: data.description,
      is_personal: data.is_personal ?? true,
      space_id: data.space_id,
    })
  },

  update(id: number, data: Partial<Bill>) {
    return client.put(`/bills/${id}`, data)
  },

  delete(id: number) {
    return client.delete(`/bills/${id}`)
  },

  analyze({ spaceId, period = 'monthly', isPersonal }: AnalyzeParams) {
    return client.post('/analyze', {
      space_id: spaceId,
      period,
      is_personal: isPersonal,
    })
  },

  getReports({ spaceId, isPersonal }: ReportScopeParams) {
    return client.get('/reports', {
      params: {
        space_id: spaceId,
        is_personal: isPersonal,
      },
    })
  },
}
