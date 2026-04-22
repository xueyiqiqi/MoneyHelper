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

        <el-form-item label="分析周期">
          <el-select v-model="period" style="width: 120px">
            <el-option label="日报" value="daily" />
            <el-option label="周报" value="weekly" />
            <el-option label="月报" value="monthly" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" :loading="loading" @click="handleAnalyze">
            生成分析
          </el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <el-card v-if="report" class="report-card">
      <template #header>
        <span>分析报告 · {{ targetLabel }} · {{ getPeriodText(report.period) }}</span>
      </template>
      <div class="report-content">
        {{ report.content }}
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
const period = ref<'daily' | 'weekly' | 'monthly'>('monthly')
const loading = ref(false)
const report = ref<AnalysisReport | null>(null)

const { scope, targetLabel, requestParams, setScope } = createReportScopeController({
  spaces: computed<Space[]>(() => spaceStore.spaces),
  currentSpace: computed<Space | null>(() => spaceStore.currentSpace),
  setCurrentSpace: spaceStore.setCurrentSpace,
})

const loadLatestReport = async () => {
  try {
    const { data } = await billApi.getReports(requestParams.value)
    report.value = data.reports.length > 0 ? data.reports[0] : null
  } catch {
    report.value = null
  }
}

const handleScopeChange = async (value: ReportScopeValue) => {
  setScope(value)
  await loadLatestReport()
}

const handleAnalyze = async () => {
  loading.value = true
  try {
    const { data } = await billApi.analyze({
      ...requestParams.value,
      period: period.value,
    })
    report.value = data
    ElMessage.success('分析完成')
  } catch (error: any) {
    ElMessage.error(error.response?.data?.error || `${targetLabel.value}分析失败`)
  } finally {
    loading.value = false
  }
}

const getPeriodText = (value: string) => {
  const map: Record<string, string> = {
    daily: '日报',
    weekly: '周报',
    monthly: '月报',
  }
  return map[value] || value
}

onMounted(async () => {
  await loadLatestReport()
})
</script>

<style scoped>
.analysis .header {
  margin-bottom: 20px;
}

.report-card {
  margin-top: 20px;
}

.report-content {
  white-space: pre-wrap;
  line-height: 1.8;
}
</style>
