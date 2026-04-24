<template>
  <div class="analysis">
    <div class="header">
      <h2>AI 分析</h2>
    </div>

    <el-card>
      <el-form inline>
        <el-form-item label="报告对象">
          <el-select
            :model-value="scope"
            style="width: 200px"
            @update:model-value="handleScopeChange"
          >
            <el-option label="个人账本" :value="PERSONAL_SCOPE" />
            <el-option
              v-for="space in spaces"
              :key="space.id"
              :label="space.name"
              :value="space.id"
            />
          </el-select>
        </el-form-item>

        <el-form-item label="开始日期">
          <el-date-picker
            v-model="startDate"
            type="date"
            value-format="YYYY-MM-DD"
            placeholder="选择开始日期"
            style="width: 160px"
          />
        </el-form-item>

        <el-form-item label="结束日期">
          <el-date-picker
            v-model="endDate"
            type="date"
            value-format="YYYY-MM-DD"
            placeholder="选择结束日期"
            style="width: 160px"
          />
        </el-form-item>

        <el-form-item>
          <el-button type="primary" :loading="loading" @click="handleAnalyze">
            生成分析
          </el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <el-card class="report-history-card">
      <template #header>
        <span>历史报告</span>
      </template>

      <template v-if="reportHistory.length > 0 && selectedReport">
        <div class="report-history-layout">
          <div class="report-history-list">
            <button
              v-for="historyItem in reportHistory"
              :key="historyItem.id"
              class="history-item"
              :class="{ active: selectedReport?.id === historyItem.id }"
              @click="selectedReport = historyItem"
            >
              <div class="history-item-range">{{ getReportRangeLabel(historyItem) }}</div>
              <div class="history-item-time">{{ historyItem.created_at }}</div>
              <div class="history-item-preview">{{ getPreview(historyItem.content) }}</div>
            </button>
          </div>

          <div class="report-detail">
            <div class="report-detail-header">
              分析报告 · {{ targetLabel }} · {{ getReportRangeLabel(selectedReport) }}
            </div>
            <div class="report-content">{{ selectedReport.content }}</div>
          </div>
        </div>
      </template>

      <div v-else class="empty-state">
        {{ emptyStateText }}
      </div>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { useSpaceStore } from '@/stores/space'
import { billApi } from '@/api/bill'
import { ElMessage } from 'element-plus'
import type { AnalysisReport, Space } from '@/types'
import {
  createReportScopeController,
  PERSONAL_SCOPE,
  type ReportScopeValue,
} from '@/views/analysis/useReportScope'

const spaceStore = useSpaceStore()

const spaces = computed<Space[]>(() => spaceStore.spaces)
const startDate = ref('')
const endDate = ref('')
const loading = ref(false)
const reportHistory = ref<AnalysisReport[]>([])
const selectedReport = ref<AnalysisReport | null>(null)

const { scope, targetLabel, requestParams, setScope } = createReportScopeController({
  spaces: computed<Space[]>(() => spaceStore.spaces),
  currentSpace: computed<Space | null>(() => spaceStore.currentSpace),
  setCurrentSpace: spaceStore.setCurrentSpace,
})

const emptyStateText = computed(() => {
  if (scope.value === PERSONAL_SCOPE) {
    return '暂无个人历史报告'
  }
  return '该空间暂无历史报告'
})

const getReportRangeLabel = (report: AnalysisReport) => {
  if (!report.period_start || !report.period_end) {
    return '未记录区间'
  }
  return `${report.period_start} ~ ${report.period_end}`
}

const getPreview = (content: string) => {
  return content.length > 60 ? `${content.slice(0, 60)}...` : content
}

const normalizeReports = (data: unknown): AnalysisReport[] => {
  if (Array.isArray(data)) {
    return data as AnalysisReport[]
  }

  if (data && typeof data === 'object' && Array.isArray((data as { reports?: unknown[] }).reports)) {
    return (data as { reports: AnalysisReport[] }).reports
  }

  return []
}

const loadHistory = async () => {
  try {
    const { data } = await billApi.getReports(requestParams.value)
    reportHistory.value = normalizeReports(data)
    selectedReport.value = reportHistory.value[0] ?? null
  } catch {
    reportHistory.value = []
    selectedReport.value = null
  }
}

const handleScopeChange = async (value: ReportScopeValue) => {
  setScope(value)
  await loadHistory()
}

const handleAnalyze = async () => {
  if (!startDate.value || !endDate.value) {
    ElMessage.error('请选择开始日期和结束日期')
    return
  }

  if (startDate.value > endDate.value) {
    ElMessage.error('开始日期不能晚于结束日期')
    return
  }

  loading.value = true
  try {
    const { data } = await billApi.analyze({
      ...requestParams.value,
      startDate: startDate.value,
      endDate: endDate.value,
    })
    reportHistory.value = [data, ...reportHistory.value.filter((item) => item.id !== data.id)]
    selectedReport.value = data
    ElMessage.success('分析完成')
  } catch (error: any) {
    ElMessage.error(error.response?.data?.error || `${targetLabel.value}分析失败`)
  } finally {
    loading.value = false
  }
}

onMounted(async () => {
  await loadHistory()
})
</script>

<style scoped>
.analysis .header {
  margin-bottom: 20px;
}

.report-history-card {
  margin-top: 20px;
}

.report-history-layout {
  display: grid;
  grid-template-columns: 280px 1fr;
  gap: 16px;
}

.report-history-list {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.history-item {
  border: 1px solid var(--border-color);
  background: var(--bg-secondary);
  border-radius: 10px;
  padding: 12px;
  text-align: left;
}

.history-item.active {
  border-color: var(--color-primary);
}

.history-item-range,
.history-item-time {
  font-size: 13px;
  margin-bottom: 6px;
}

.history-item-preview {
  font-size: 13px;
  color: var(--text-secondary);
}

.report-detail-header {
  font-weight: 600;
  margin-bottom: 12px;
}

.report-content {
  white-space: pre-wrap;
  line-height: 1.8;
}

.empty-state {
  color: var(--text-secondary);
  padding: 12px 0;
}
</style>
