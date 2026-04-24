import { beforeEach, describe, expect, it, vi } from 'vitest'

const clientMock = vi.hoisted(() => ({
  get: vi.fn(),
  post: vi.fn(),
}))

vi.mock('./client', () => ({
  default: clientMock,
}))

import { billApi } from './bill'

describe('billApi analysis report requests', () => {
  beforeEach(() => {
    clientMock.get.mockReset()
    clientMock.post.mockReset()
  })

  it('posts personal analysis requests with explicit date range params', async () => {
    clientMock.post.mockResolvedValue({ data: {} })

    await billApi.analyze({
      isPersonal: true,
      spaceId: undefined,
      startDate: '2026-04-01',
      endDate: '2026-04-30',
    })

    expect(clientMock.post).toHaveBeenCalledWith('/analyze', {
      space_id: undefined,
      is_personal: true,
      start_date: '2026-04-01',
      end_date: '2026-04-30',
    })
  })

  it('queries space reports with explicit is_personal=false', async () => {
    clientMock.get.mockResolvedValue({ data: { reports: [] } })

    await billApi.getReports({
      isPersonal: false,
      spaceId: 7,
    })

    expect(clientMock.get).toHaveBeenCalledWith('/reports', {
      params: {
        space_id: 7,
        is_personal: false,
      },
    })
  })
})
