<template>
  <div class="bill-edit-page">
    <div class="page-header">
      <el-button :icon="ArrowLeft" text @click="router.push('/bills')">
        返回
      </el-button>
      <h1 class="page-title">{{ isEdit ? '编辑账单' : '添加账单' }}</h1>
    </div>

    <div class="edit-card">
      <el-form
        ref="formRef"
        :model="form"
        :rules="rules"
        label-position="top"
        class="edit-form"
      >
        <!-- 类型选择 -->
        <el-form-item label="类型" prop="type">
          <div class="type-selector">
            <div
              :class="['type-option', { active: form.type === 'expense' }]"
              @click="form.type = 'expense'"
            >
              <span class="type-icon">📉</span>
              <span class="type-label">支出</span>
            </div>
            <div
              :class="['type-option', { active: form.type === 'income' }]"
              @click="form.type = 'income'"
            >
              <span class="type-icon">📈</span>
              <span class="type-label">收入</span>
            </div>
          </div>
        </el-form-item>

        <!-- 金额 -->
        <el-form-item label="金额" prop="amount">
          <el-input
            v-model="form.amount"
            type="number"
            placeholder="0.00"
            class="amount-input"
          >
            <template #prefix>¥</template>
          </el-input>
        </el-form-item>

        <!-- 分类 -->
        <el-form-item label="分类" prop="category">
          <div class="category-grid">
            <div
              v-for="cat in categories[form.type]"
              :key="cat"
              :class="['category-item', { active: form.category === cat }]"
              @click="form.category = cat"
            >
              {{ cat }}
            </div>
          </div>
        </el-form-item>

        <!-- 日期 -->
        <el-form-item label="日期" prop="bill_date">
          <el-date-picker
            v-model="form.bill_date"
            type="date"
            placeholder="选择日期"
            value-format="YYYY-MM-DD"
            class="date-picker"
          />
        </el-form-item>

        <!-- 备注 -->
        <el-form-item label="备注 (可选)">
          <el-input
            v-model="form.description"
            type="textarea"
            :rows="3"
            placeholder="添加备注..."
            maxlength="200"
            show-word-limit
          />
        </el-form-item>

        <!-- 提交按钮 -->
        <el-form-item>
          <div class="form-actions">
            <el-button size="large" @click="router.push('/bills')">
              取消
            </el-button>
            <el-button
              type="primary"
              size="large"
              :loading="loading"
              @click="handleSubmit"
              class="submit-btn"
            >
              {{ isEdit ? '保存修改' : '添加账单' }}
            </el-button>
          </div>
        </el-form-item>
      </el-form>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useSpaceStore } from '@/stores/space'
import { billApi } from '@/api/bill'
import { ElMessage, type FormInstance, type FormRules } from 'element-plus'
import { ArrowLeft } from '@element-plus/icons-vue'
import type { Bill } from '@/types'

const route = useRoute()
const router = useRouter()
const spaceStore = useSpaceStore()

const formRef = ref<FormInstance>()
const loading = ref(false)
const isEdit = computed(() => !!route.params.id)

const form = reactive({
  type: 'expense' as 'income' | 'expense',
  amount: 0,
  category: '',
  description: '',
  bill_date: '',
})

const rules: FormRules = {
  type: [{ required: true, message: '请选择类型', trigger: 'change' }],
  amount: [{ required: true, message: '请输入金额', trigger: 'blur' }],
  category: [{ required: true, message: '请选择或输入分类', trigger: 'change' }],
  bill_date: [{ required: true, message: '请选择日期', trigger: 'change' }],
}

// 预定义分类
const categories = {
  income: ['工资', '奖金', '投资收益', '兼职', '其他收入'],
  expense: ['餐饮', '交通', '购物', '娱乐', '居住', '医疗', '教育', '其他支出'],
}

const handleSubmit = async () => {
  const valid = await formRef.value?.validate().catch(() => false)
  if (!valid) return

  loading.value = true
  try {
    const data: any = { ...form }
    if (spaceStore.currentSpace) {
      data.space_id = spaceStore.currentSpace.id
    }

    if (isEdit.value) {
      await billApi.update(Number(route.params.id), data)
      ElMessage.success('更新成功')
    } else {
      await billApi.create(data)
      ElMessage.success('添加成功')
    }
    router.push('/bills')
  } catch (error: any) {
    ElMessage.error(error.response?.data?.error || '操作失败')
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  if (isEdit.value) {
    // 从 sessionStorage 获取账单数据
    const savedBill = sessionStorage.getItem('editBill')
    if (savedBill) {
      const billData = JSON.parse(savedBill) as Bill
      sessionStorage.removeItem('editBill')

      form.type = billData.type || (billData.amount >= 0 ? 'income' : 'expense')
      form.amount = Math.abs(billData.amount)
      form.category = billData.category || ''
      form.description = billData.description || billData.remarks || ''
      form.bill_date = billData.bill_date || billData.date || ''
    }
  }
})

// 监听类型变化，清空分类
watch(() => form.type, () => {
  form.category = ''
})
</script>

<style scoped>
.bill-edit-page {
  max-width: 560px;
  margin: 0 auto;
}

.page-header {
  display: flex;
  align-items: center;
  gap: 16px;
  margin-bottom: 24px;
}

.page-title {
  font-size: 24px;
  font-weight: 700;
  color: var(--text-primary);
  margin: 0;
}

.edit-card {
  background: var(--bg-card);
  border-radius: var(--radius-xl);
  padding: 32px;
  box-shadow: var(--shadow-md);
  border: 1px solid var(--border-light);
}

/* 类型选择器 */
.type-selector {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 12px;
}

.type-option {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 10px;
  padding: 16px;
  border-radius: var(--radius-md);
  border: 2px solid var(--border-color);
  cursor: pointer;
  transition: all var(--transition-fast);
}

.type-option:hover {
  border-color: var(--color-primary);
}

.type-option.active {
  border-color: var(--color-primary);
  background: var(--color-primary-light);
}

.type-icon {
  font-size: 24px;
}

.type-label {
  font-size: 16px;
  font-weight: 600;
  color: var(--text-primary);
}

/* 金额输入 */
.amount-input :deep(.el-input__wrapper) {
  padding: 12px 16px;
}

.amount-input :deep(.el-input__inner) {
  font-size: 24px;
  font-weight: 700;
  font-family: 'SF Mono', monospace;
}

/* 分类选择 */
.category-grid {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.category-item {
  padding: 8px 16px;
  border-radius: var(--radius-lg);
  background: var(--bg-secondary);
  color: var(--text-secondary);
  font-size: 14px;
  cursor: pointer;
  transition: all var(--transition-fast);
  border: 1px solid transparent;
}

.category-item:hover {
  background: var(--color-primary-light);
  color: var(--color-primary);
}

.category-item.active {
  background: var(--color-primary);
  color: white;
  border-color: var(--color-primary);
}

/* 日期选择器 */
.date-picker {
  width: 100%;
}

/* 表单操作 */
.form-actions {
  display: flex;
  gap: 12px;
  margin-top: 24px;
}

.form-actions .el-button {
  flex: 1;
  height: 48px;
  font-size: 16px;
  font-weight: 600;
  border-radius: var(--radius-md);
}

.submit-btn {
  background: var(--color-primary);
  border-color: var(--color-primary);
}

.submit-btn:hover {
  background: var(--color-primary-hover);
  border-color: var(--color-primary-hover);
}
</style>
