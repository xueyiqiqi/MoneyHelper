<template>
  <div class="bill-page">
    <!-- 页面头部 -->
    <div class="page-header">
      <div class="header-content">
        <h1 class="page-title">账单</h1>
        <p class="page-subtitle">
          <span v-if="isPersonal">📒 个人账本</span>
          <span v-else>🏠 {{ currentSpace?.name }}</span>
        </p>
      </div>
      <el-button
        type="primary"
        size="large"
        class="add-btn"
        @click="router.push('/bills/new')"
      >
        <el-icon><Plus /></el-icon>
        添加账单
      </el-button>
    </div>

    <!-- 统计概览 -->
    <div class="overview-section">
      <div class="overview-card">
        <div class="overview-header">
          <span class="overview-label">本月收入</span>
          <span class="overview-icon">📈</span>
        </div>
        <div class="overview-value income">¥{{ formatAmount(totalIncome) }}</div>
        <div class="overview-detail">共 {{ incomeCount }} 笔收入</div>
      </div>

      <div class="overview-card">
        <div class="overview-header">
          <span class="overview-label">本月支出</span>
          <span class="overview-icon">📉</span>
        </div>
        <div class="overview-value expense">¥{{ formatAmount(totalExpense) }}</div>
        <div class="overview-detail">共 {{ expenseCount }} 笔支出</div>
      </div>

      <div class="overview-card highlight">
        <div class="overview-header">
          <span class="overview-label">本月结余</span>
          <span class="overview-icon">💰</span>
        </div>
        <div :class="['overview-value', balance >= 0 ? 'positive' : 'negative']">
          ¥{{ formatAmount(Math.abs(balance)) }}
        </div>
        <div class="overview-detail">{{ balance >= 0 ? '储蓄充裕 💪' : '入不敷出 😅' }}</div>
      </div>
    </div>

    <!-- 筛选栏 -->
    <div class="filter-bar">
      <div class="filter-chips">
        <div
          :class="['chip', { active: !filters.type }]"
          @click="filters.type = ''; fetchBills()"
        >
          全部
        </div>
        <div
          :class="['chip income', { active: filters.type === 'income' }]"
          @click="filters.type = 'income'; fetchBills()"
        >
          📈 收入
        </div>
        <div
          :class="['chip expense', { active: filters.type === 'expense' }]"
          @click="filters.type = 'expense'; fetchBills()"
        >
          📉 支出
        </div>
      </div>

      <div class="filter-actions">
        <el-select
          v-model="filters.category"
          placeholder="选择分类"
          clearable
          filterable
          class="search-input"
          @change="fetchBills"
        >
          <el-option
            v-for="cat in allCategories"
            :key="cat"
            :label="cat"
            :value="cat"
          />
        </el-select>
        <el-date-picker
          v-model="dateRange"
          type="daterange"
          range-separator="至"
          start-placeholder="开始"
          end-placeholder="结束"
          value-format="YYYY-MM-DD"
          class="date-picker"
          @change="fetchBills"
        />
      </div>
    </div>

    <!-- 账单列表 -->
    <div class="bills-container" v-if="bills.length > 0">
      <TransitionGroup name="bill" tag="div" class="bill-list">
        <div
          v-for="(bill, index) in bills"
          :key="bill.id"
          class="bill-item"
          :style="{ animationDelay: `${index * 50}ms` }"
        >
          <div class="bill-main" @click="editBill(bill)">
            <div :class="['bill-icon', bill.type]">
              {{ bill.type === 'income' ? '📈' : '📉' }}
            </div>

            <div class="bill-info">
              <div class="bill-category-row">
                <span class="bill-category">{{ bill.category }}</span>
                <span :class="['bill-type-tag', bill.type]">
                  {{ bill.type === 'income' ? '收入' : '支出' }}
                </span>
              </div>
              <div class="bill-desc" v-if="bill.description">
                {{ bill.description }}
              </div>
              <div class="bill-date">{{ formatFullDate(bill.bill_date) }}</div>
            </div>

            <div class="bill-amount-section">
              <span :class="['bill-amount', bill.type]">
                {{ bill.type === 'income' ? '+' : '-' }}¥{{ formatAmount(bill.amount) }}
              </span>
            </div>
          </div>

          <div class="bill-actions">
            <el-button
              type="danger"
              :icon="Delete"
              circle
              size="small"
              plain
              @click.stop="handleDelete(bill.id)"
            />
          </div>
        </div>
      </TransitionGroup>

      <!-- 分页 -->
      <div class="pagination" v-if="pagination.total > pagination.pageSize">
        <el-pagination
          v-model:current-page="pagination.page"
          :page-size="pagination.pageSize"
          :total="pagination.total"
          layout="prev, pager, next"
          @current-change="fetchBills"
          background
        />
      </div>
    </div>

    <!-- 空状态 -->
    <div class="empty-state" v-else-if="!loading">
      <div class="empty-illustration">
        <div class="empty-cloud"></div>
        <div class="empty-bill">📝</div>
      </div>
      <h3>还没有账单记录</h3>
      <p>点击"添加账单"开始记录你的第一笔收支</p>
      <el-button type="primary" size="large" @click="router.push('/bills/new')">
        <el-icon><Plus /></el-icon>
        添加第一笔账单
      </el-button>
    </div>

    <!-- 加载中 -->
    <div class="loading-state" v-if="loading">
      <el-skeleton :rows="5" animated />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted, watch } from 'vue'
import { useRouter } from 'vue-router'
import { useSpaceStore } from '@/stores/space'
import { billApi } from '@/api/bill'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus, Search, Delete } from '@element-plus/icons-vue'
import type { Bill } from '@/types'

const router = useRouter()
const spaceStore = useSpaceStore()

const loading = ref(false)
const bills = ref<Bill[]>([])

const currentSpace = computed(() => spaceStore.currentSpace)
const isPersonal = computed(() => !currentSpace.value)

const filters = reactive({
  type: '' as '' | 'income' | 'expense',
  category: '',
})

const dateRange = ref<[string, string] | null>(null)

const pagination = reactive({
  page: 1,
  pageSize: 15,
  total: 0,
})

const totalIncome = computed(() => bills.value.filter(b => b.type === 'income').reduce((s, b) => s + Math.abs(b.amount), 0))
const totalExpense = computed(() => bills.value.filter(b => b.type === 'expense').reduce((s, b) => s + Math.abs(b.amount), 0))
const balance = computed(() => totalIncome.value - totalExpense.value)
const incomeCount = computed(() => bills.value.filter(b => b.type === 'income').length)
const expenseCount = computed(() => bills.value.filter(b => b.type === 'expense').length)

// 所有分类选项（从预定义分类和实际账单中提取）
const allCategories = computed(() => {
  const predefined = ['工资', '奖金', '投资收益', '兼职', '其他收入', '餐饮', '交通', '购物', '娱乐', '居住', '医疗', '教育', '其他支出']
  const existing = [...new Set(bills.value.map(b => b.category))]
  return [...new Set([...predefined, ...existing])]
})

const formatAmount = (amount: number) => amount.toLocaleString('zh-CN', { minimumFractionDigits: 2, maximumFractionDigits: 2 })
const formatFullDate = (date: string) => new Date(date).toLocaleDateString('zh-CN', { year: 'numeric', month: 'long', day: 'numeric' })

let fetchTimer: ReturnType<typeof setTimeout>
const debouncedFetch = () => {
  clearTimeout(fetchTimer)
  fetchTimer = setTimeout(fetchBills, 300)
}

const fetchBills = async () => {
  loading.value = true
  try {
    const params: any = {
      page: pagination.page,
      page_size: pagination.pageSize,
      is_personal: isPersonal.value ? 'true' : 'false'
    }
    if (filters.type) params.type = filters.type
    if (filters.category) params.category = filters.category
    if (dateRange.value) {
      params.start_date = dateRange.value[0]
      params.end_date = dateRange.value[1]
    }
    if (currentSpace.value) params.space_id = currentSpace.value.id

    const { data } = await billApi.getList(params)
    // 后端直接返回数组
    bills.value = Array.isArray(data) ? data : data.bills || []
    pagination.total = bills.value.length
  } catch (error) {
    ElMessage.error('获取账单失败')
    bills.value = []
  } finally {
    loading.value = false
  }
}

const editBill = (bill: Bill) => {
  // 使用 sessionStorage 传递账单数据
  sessionStorage.setItem('editBill', JSON.stringify(bill))
  router.push(`/bills/${bill.id}`)
}

const handleDelete = async (id: number) => {
  await ElMessageBox.confirm('确定删除这条账单吗？', '提示', { type: 'warning' })
  try {
    await billApi.delete(id)
    ElMessage.success('删除成功')
    fetchBills()
  } catch (error: any) {
    ElMessage.error(error.response?.data?.error || '删除失败')
  }
}

onMounted(fetchBills)

// 监听空间切换，自动刷新账单
watch(() => spaceStore.currentSpace, () => {
  pagination.page = 1
  fetchBills()
})
</script>

<style scoped>
.bill-page {
  max-width: 800px;
  margin: 0 auto;
}

/* 页面头部 */
.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 28px;
}

.page-title {
  font-size: 32px;
  font-weight: 800;
  color: var(--text-primary);
  margin: 0;
  letter-spacing: -1px;
}

.page-subtitle {
  font-size: 14px;
  color: var(--text-tertiary);
  margin-top: 4px;
}

.add-btn {
  height: 44px;
  padding: 0 20px;
  font-size: 15px;
  border-radius: 22px;
  box-shadow: 0 4px 12px rgba(126, 176, 138, 0.3);
}

/* 概览卡片 */
.overview-section {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 16px;
  margin-bottom: 24px;
}

.overview-card {
  background: var(--bg-card);
  border-radius: 20px;
  padding: 20px;
  border: 1px solid var(--border-light);
  transition: all 0.3s ease;
}

.overview-card:hover {
  transform: translateY(-4px);
  box-shadow: var(--shadow-md);
}

.overview-card.highlight {
  background: linear-gradient(135deg, var(--color-primary-light) 0%, var(--bg-card) 100%);
  border-color: var(--color-primary);
}

.overview-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 12px;
}

.overview-label {
  font-size: 13px;
  color: var(--text-tertiary);
  font-weight: 500;
}

.overview-icon {
  font-size: 20px;
}

.overview-value {
  font-size: 26px;
  font-weight: 800;
  font-family: 'SF Mono', monospace;
  margin-bottom: 4px;
}

.overview-value.income { color: var(--color-income); }
.overview-value.expense { color: var(--color-expense); }
.overview-value.positive { color: var(--color-income); }
.overview-value.negative { color: var(--color-expense); }

.overview-detail {
  font-size: 12px;
  color: var(--text-tertiary);
}

/* 筛选栏 */
.filter-bar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
  gap: 16px;
  flex-wrap: wrap;
}

.filter-chips {
  display: flex;
  gap: 8px;
}

.chip {
  padding: 8px 16px;
  border-radius: 20px;
  font-size: 13px;
  font-weight: 500;
  cursor: pointer;
  background: var(--bg-card);
  border: 1px solid var(--border-color);
  color: var(--text-secondary);
  transition: all 0.2s ease;
}

.chip:hover {
  border-color: var(--color-primary);
}

.chip.active {
  background: var(--color-primary);
  border-color: var(--color-primary);
  color: white;
}

.chip.active.income {
  background: var(--color-income);
  border-color: var(--color-income);
}

.chip.active.expense {
  background: var(--color-expense);
  border-color: var(--color-expense);
}

.filter-actions {
  display: flex;
  gap: 12px;
}

.search-input {
  width: 180px;
}

.date-picker {
  width: 240px;
}

/* 账单列表 */
.bills-container {
  background: var(--bg-card);
  border-radius: 20px;
  padding: 8px;
  border: 1px solid var(--border-light);
}

.bill-list {
  display: flex;
  flex-direction: column;
}

.bill-item {
  display: flex;
  align-items: center;
  padding: 16px;
  border-radius: 16px;
  transition: all 0.2s ease;
  animation: slideIn 0.4s ease forwards;
  opacity: 0;
}

@keyframes slideIn {
  from {
    opacity: 0;
    transform: translateX(-20px);
  }
  to {
    opacity: 1;
    transform: translateX(0);
  }
}

.bill-item:hover {
  background: var(--bg-secondary);
}

.bill-main {
  display: flex;
  align-items: center;
  flex: 1;
  gap: 16px;
  cursor: pointer;
}

.bill-icon {
  width: 48px;
  height: 48px;
  border-radius: 14px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 22px;
  flex-shrink: 0;
}

.bill-icon.income {
  background: linear-gradient(135deg, #e8f5e9 0%, #c8e6c9 100%);
}

.bill-icon.expense {
  background: linear-gradient(135deg, #ffebee 0%, #ffcdd2 100%);
}

.bill-info {
  flex: 1;
  min-width: 0;
}

.bill-category-row {
  display: flex;
  align-items: center;
  gap: 8px;
}

.bill-category {
  font-size: 15px;
  font-weight: 600;
  color: var(--text-primary);
}

.bill-type-tag {
  font-size: 11px;
  padding: 2px 8px;
  border-radius: 10px;
  font-weight: 500;
}

.bill-type-tag.income {
  background: rgba(126, 176, 138, 0.15);
  color: var(--color-income);
}

.bill-type-tag.expense {
  background: rgba(217, 133, 133, 0.15);
  color: var(--color-expense);
}

.bill-desc {
  font-size: 13px;
  color: var(--text-secondary);
  margin-top: 2px;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.bill-date {
  font-size: 12px;
  color: var(--text-tertiary);
  margin-top: 4px;
}

.bill-amount-section {
  text-align: right;
}

.bill-amount {
  font-size: 18px;
  font-weight: 700;
  font-family: 'SF Mono', monospace;
}

.bill-amount.income { color: var(--color-income); }
.bill-amount.expense { color: var(--color-expense); }

.bill-actions {
  opacity: 0;
  transition: opacity 0.2s;
  margin-left: 12px;
}

.bill-item:hover .bill-actions {
  opacity: 1;
}

/* 分页 */
.pagination {
  display: flex;
  justify-content: center;
  padding: 20px;
  border-top: 1px solid var(--border-light);
  margin-top: 8px;
}

/* 空状态 */
.empty-state {
  text-align: center;
  padding: 60px 20px;
  background: var(--bg-card);
  border-radius: 24px;
  border: 1px dashed var(--border-color);
}

.empty-illustration {
  position: relative;
  width: 120px;
  height: 100px;
  margin: 0 auto 24px;
}

.empty-cloud {
  position: absolute;
  width: 80px;
  height: 40px;
  background: var(--bg-secondary);
  border-radius: 40px;
  top: 20px;
  left: 20px;
  animation: float 3s ease-in-out infinite;
}

.empty-bill {
  position: absolute;
  font-size: 48px;
  bottom: 0;
  left: 50%;
  transform: translateX(-50%);
  animation: bounce 2s ease-in-out infinite;
}

@keyframes float {
  0%, 100% { transform: translateY(0); }
  50% { transform: translateY(-10px); }
}

@keyframes bounce {
  0%, 100% { transform: translateX(-50%) translateY(0); }
  50% { transform: translateX(-50%) translateY(-10px); }
}

.empty-state h3 {
  font-size: 20px;
  color: var(--text-primary);
  margin-bottom: 8px;
}

.empty-state p {
  color: var(--text-tertiary);
  margin-bottom: 24px;
}

/* 加载中 */
.loading-state {
  background: var(--bg-card);
  border-radius: 20px;
  padding: 20px;
}
</style>
