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

    <el-card v-if="report" class="report-card">
      <template #header>
        <span>分析报告 · {{ targetLabel }} · {{ reportRangeLabel }}</span>
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
const startDate = ref('')
const endDate = ref('')
const loading = ref(false)
const report = ref<AnalysisReport | null>(null)

const { scope, targetLabel, requestParams, setScope } = createReportScopeController({
  spaces: computed<Space[]>(() => spaceStore.spaces),
  currentSpace: computed<Space | null>(() => spaceStore.currentSpace),
  setCurrentSpace: spaceStore.setCurrentSpace,
})

const reportRangeLabel = computed(() => {
  if (!report.value?.period_start || !report.value?.period_end) {
    return '未记录区间'
  }
  return `${report.value.period_start} ~ ${report.value.period_end}`
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
    report.value = data
    ElMessage.success('分析完成')
  } catch (error: any) {
    ElMessage.error(error.response?.data?.error || `${targetLabel.value}分析失败`)
  } finally {
    loading.value = false
  }
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
