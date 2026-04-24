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
  'el-date-picker': {
    props: ['modelValue'],
    emits: ['update:modelValue'],
    template: '<input type="date" :value="modelValue" @input="$emit(\'update:modelValue\', $event.target.value)" />',
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

  it('switches to personal scope and sends explicit date range params', async () => {
    const generatedReport: AnalysisReport = {
      id: 1,
      user_id: 1,
      content: '个人报告',
      created_at: '2026-04-22T00:00:00Z',
      period_start: '2026-04-01',
      period_end: '2026-04-30',
    }

    apiMock.getReports.mockResolvedValue({ data: { reports: [] } })
    apiMock.analyze.mockResolvedValue({ data: generatedReport })

    const wrapper = mount(Analysis, {
      global: {
        stubs,
      },
    })

    await flushPromises()

    const scopeSelect = wrapper.find('select')
    const dateInputs = wrapper.findAll('input[type="date"]')
    const [startDateInput, endDateInput] = dateInputs

    await scopeSelect.setValue('personal')
    await startDateInput.setValue('2026-04-01')
    await endDateInput.setValue('2026-04-30')
    await wrapper.find('button').trigger('click')
    await flushPromises()

    expect(setCurrentSpace).toHaveBeenCalledWith(null)
    expect(apiMock.analyze).toHaveBeenCalledWith({
      isPersonal: true,
      spaceId: undefined,
      startDate: '2026-04-01',
      endDate: '2026-04-30',
    })
    expect(messageMock.success).toHaveBeenCalledWith('分析完成')
    expect(wrapper.text()).toContain('分析报告 · 个人账本 · 2026-04-01 ~ 2026-04-30')
  })

  it('does not submit when dates are missing', async () => {
    apiMock.getReports.mockResolvedValue({ data: { reports: [] } })

    const wrapper = mount(Analysis, {
      global: {
        stubs,
      },
    })

    await flushPromises()
    await wrapper.find('button').trigger('click')
    await flushPromises()

    expect(apiMock.analyze).not.toHaveBeenCalled()
    expect(messageMock.error).toHaveBeenCalledWith('请选择开始日期和结束日期')
  })

  it('does not submit when start date is after end date', async () => {
    apiMock.getReports.mockResolvedValue({ data: { reports: [] } })

    const wrapper = mount(Analysis, {
      global: {
        stubs,
      },
    })

    await flushPromises()

    const dateInputs = wrapper.findAll('input[type="date"]')
    const [startDateInput, endDateInput] = dateInputs

    await startDateInput.setValue('2026-05-01')
    await endDateInput.setValue('2026-04-01')
    await wrapper.find('button').trigger('click')
    await flushPromises()

    expect(apiMock.analyze).not.toHaveBeenCalled()
    expect(messageMock.error).toHaveBeenCalledWith('开始日期不能晚于结束日期')
  })
})
