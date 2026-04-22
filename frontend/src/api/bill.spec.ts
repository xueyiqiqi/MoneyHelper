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

  it('posts personal analysis requests with explicit is_personal=true', async () => {
    clientMock.post.mockResolvedValue({ data: {} })

    await billApi.analyze({
      period: 'monthly',
      isPersonal: true,
      spaceId: undefined,
    })

    expect(clientMock.post).toHaveBeenCalledWith('/analyze', {
      space_id: undefined,
      period: 'monthly',
      is_personal: true,
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
