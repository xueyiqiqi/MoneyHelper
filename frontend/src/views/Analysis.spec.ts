import { ref } from 'vue'
import { flushPromises, mount } from '@vue/test-utils'
import { beforeEach, describe, expect, it, vi } from 'vitest'
import type { AnalysisReport, Space } from '@/types'
import Analysis from './Analysis.vue'

const apiMock = vi.hoisted(() => ({
  getReports: vi.fn(),
  analyze: vi.fn(),
}))

const messageMock = vi.hoisted(() => ({
  success: vi.fn(),
  error: vi.fn(),
}))

const familySpace: Space = {
  id: 7,
  name: '家庭空间',
  description: '家庭共享账本',
  role: 'member',
  created_at: '2026-04-22T00:00:00Z',
}

const spaces = ref<Space[]>([familySpace])
const currentSpace = ref<Space | null>(familySpace)
const setCurrentSpace = vi.fn((space: Space | null) => {
  currentSpace.value = space
})

vi.mock('@/api/bill', () => ({
  billApi: apiMock,
}))

vi.mock('@/stores/space', () => ({
  useSpaceStore: () => ({
    spaces: spaces.value,
    currentSpace: currentSpace.value,
    setCurrentSpace,
  }),
}))

vi.mock('element-plus', () => ({
  ElMessage: messageMock,
}))

const stubs = {
  'el-card': { template: '<div><slot name="header" /><slot /></div>' },
  'el-form': { template: '<form><slot /></form>' },
  'el-form-item': { template: '<div><slot /></div>' },
  'el-select': {
    props: ['modelValue'],
    emits: ['update:modelValue'],
    template: `<select :value="modelValue" @change="$emit('update:modelValue', $event.target.value === 'personal' ? 'personal' : Number($event.target.value))"><slot /></select>`,
  },
  'el-option': {
    props: ['label', 'value'],
    template: '<option :value="value">{{ label }}</option>',
  },
  'el-button': {
    props: ['loading'],
    emits: ['click'],
    template: '<button :disabled="loading" @click="$emit(\'click\')"><slot /></button>',
  },
}

describe('Analysis.vue', () => {
  beforeEach(() => {
    spaces.value = [familySpace]
    currentSpace.value = familySpace
    setCurrentSpace.mockClear()
    apiMock.getReports.mockReset()
    apiMock.analyze.mockReset()
    messageMock.success.mockReset()
    messageMock.error.mockReset()
  })

  it('loads history for the current space on mount', async () => {
    apiMock.getReports.mockResolvedValue({ data: { reports: [] } })

    mount(Analysis, {
      global: {
        stubs,
      },
    })

    await flushPromises()

    expect(apiMock.getReports).toHaveBeenCalledWith({
      isPersonal: false,
      spaceId: 7,
    })
  })

  it('switches to personal scope and sends explicit personal params', async () => {
    const generatedReport: AnalysisReport = {
      id: 1,
      user_id: 1,
      period: 'monthly',
      content: '个人报告',
      created_at: '2026-04-22T00:00:00Z',
    }

    apiMock.getReports.mockResolvedValue({ data: { reports: [] } })
    apiMock.analyze.mockResolvedValue({ data: generatedReport })

    const wrapper = mount(Analysis, {
      global: {
        stubs,
      },
    })

    await flushPromises()

    const [scopeSelect] = wrapper.findAll('select')
    await scopeSelect.setValue('personal')
    await wrapper.find('button').trigger('click')
    await flushPromises()

    expect(setCurrentSpace).toHaveBeenCalledWith(null)
    expect(apiMock.analyze).toHaveBeenCalledWith({
      period: 'monthly',
      isPersonal: true,
      spaceId: undefined,
    })
    expect(messageMock.success).toHaveBeenCalledWith('分析完成')
    expect(wrapper.text()).toContain('分析报告 · 个人账本 · 月报')
  })
})
