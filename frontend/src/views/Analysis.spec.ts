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

const spaceReport: AnalysisReport = {
  id: 2,
  user_id: 1,
  space_id: 7,
  content: '空间报告正文',
  created_at: '2026-04-23T00:00:00Z',
  period_start: '2026-04-10',
  period_end: '2026-04-20',
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

  it('renders history when reports endpoint returns a raw array', async () => {
    apiMock.getReports.mockResolvedValue({ data: [spaceReport] })

    const wrapper = mount(Analysis, {
      global: {
        stubs,
      },
    })

    await flushPromises()

    expect(wrapper.text()).toContain('空间报告正文')
    expect(wrapper.text()).toContain('2026-04-10 ~ 2026-04-20')
  })
})
