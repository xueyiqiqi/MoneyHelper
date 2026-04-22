import { nextTick, ref } from 'vue'
import { describe, expect, it, vi } from 'vitest'
import type { Space } from '@/types'
import { createReportScopeController, PERSONAL_SCOPE } from './useReportScope'

const familySpace: Space = {
  id: 7,
  name: '家庭空间',
  description: '家庭共享账本',
  role: 'member',
  created_at: '2026-04-22T00:00:00Z',
}

describe('createReportScopeController', () => {
  it('defaults to personal when there is no current space', () => {
    const spaces = ref<Space[]>([familySpace])
    const currentSpace = ref<Space | null>(null)
    const setCurrentSpace = vi.fn((space: Space | null) => {
      currentSpace.value = space
    })

    const controller = createReportScopeController({
      spaces,
      currentSpace,
      setCurrentSpace,
    })

    expect(controller.scope.value).toBe(PERSONAL_SCOPE)
    expect(controller.targetLabel.value).toBe('个人账本')
    expect(controller.requestParams.value).toEqual({
      isPersonal: true,
      spaceId: undefined,
    })
    expect(setCurrentSpace).not.toHaveBeenCalled()
  })

  it('maps a selected space to space-scoped request params', async () => {
    const spaces = ref<Space[]>([familySpace])
    const currentSpace = ref<Space | null>(null)
    const setCurrentSpace = vi.fn((space: Space | null) => {
      currentSpace.value = space
    })

    const controller = createReportScopeController({
      spaces,
      currentSpace,
      setCurrentSpace,
    })

    controller.setScope(7)
    await nextTick()

    expect(setCurrentSpace).toHaveBeenCalledWith(familySpace)
    expect(controller.scope.value).toBe(7)
    expect(controller.targetLabel.value).toBe('家庭空间')
    expect(controller.requestParams.value).toEqual({
      isPersonal: false,
      spaceId: 7,
    })
  })

  it('falls back to personal when the current space disappears from the store', async () => {
    const missingSpace: Space = {
      id: 99,
      name: '失效空间',
      description: '已失效',
      role: 'member',
      created_at: '2026-04-22T00:00:00Z',
    }

    const spaces = ref<Space[]>([familySpace])
    const currentSpace = ref<Space | null>(missingSpace)
    const setCurrentSpace = vi.fn((space: Space | null) => {
      currentSpace.value = space
    })

    const controller = createReportScopeController({
      spaces,
      currentSpace,
      setCurrentSpace,
    })

    await nextTick()

    expect(setCurrentSpace).toHaveBeenCalledWith(null)
    expect(controller.scope.value).toBe(PERSONAL_SCOPE)
    expect(controller.targetLabel.value).toBe('个人账本')
  })
})
