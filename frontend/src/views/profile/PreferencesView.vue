<script setup lang="ts">
import { computed, onMounted, reactive, ref } from 'vue'
import { RouterLink } from 'vue-router'
import type { FormInstance, FormRules } from 'element-plus'
import { ElMessage } from 'element-plus'
import { useAuthStore, useUserStore } from '@/stores'
import { getPositionsByCollege } from '@/constants/profile'
import {
  CITY_SUGGESTIONS,
  COMPANY_SUGGESTIONS,
  REMIND_BEFORE_OPTIONS,
} from '@/constants/preferences'
import type { UserPreference } from '@/types'

const authStore = useAuthStore()
const userStore = useUserStore()
const formRef = ref<FormInstance>()
const pageLoading = ref(true)
const saving = ref(false)

const form = reactive<UserPreference>({
  targetPositions: [],
  preferredCities: [],
  preferredCompanies: [],
  focusCompanies: [],
  remindBefore: ['1d'],
})

const focusCompanyOptions = computed(() => {
  const merged = new Set([...COMPANY_SUGGESTIONS, ...form.preferredCompanies])
  return Array.from(merged)
})

const positionOptions = computed(() =>
  getPositionsByCollege(authStore.userInfo?.college),
)

const rules: FormRules<UserPreference> = {
  remindBefore: [
    {
      type: 'array',
      required: true,
      min: 1,
      message: '请至少选择一种提醒提前量',
      trigger: 'change',
    },
  ],
}

function fillForm(data: UserPreference) {
  form.targetPositions = [...(data.targetPositions || [])]
  form.preferredCities = [...(data.preferredCities || [])]
  form.preferredCompanies = [...(data.preferredCompanies || [])]
  form.focusCompanies = [...(data.focusCompanies || [])]
  form.remindBefore = data.remindBefore?.length ? [...data.remindBefore] : ['1d']
}

function buildPayload(): UserPreference {
  return {
    targetPositions: form.targetPositions,
    preferredCities: form.preferredCities,
    preferredCompanies: form.preferredCompanies,
    focusCompanies: form.focusCompanies,
    remindBefore: form.remindBefore,
  }
}

async function loadPreferences() {
  pageLoading.value = true
  try {
    const data = await userStore.fetchPreferences()
    if (data) {
      fillForm(data)
    }
  } finally {
    pageLoading.value = false
  }
}

async function handleSubmit() {
  if (!formRef.value) return
  const valid = await formRef.value.validate().catch(() => false)
  if (!valid) return

  saving.value = true
  try {
    const data = await userStore.updatePreferences(buildPayload())
    if (data) {
      fillForm(data)
    }
    await authStore.fetchUserInfo()
    ElMessage.success('偏好设置已保存')
  } catch {
    // 错误由 axios 拦截器提示
  } finally {
    saving.value = false
  }
}

async function handleReset() {
  await loadPreferences()
  formRef.value?.clearValidate()
  ElMessage.info('已恢复为服务器上的最新偏好')
}

function handleClearOptional() {
  form.preferredCities = []
  form.preferredCompanies = []
  form.focusCompanies = []
  ElMessage.info('已清空城市/公司偏好，请点击保存生效')
}

onMounted(() => {
  loadPreferences()
})
</script>

<template>
  <div v-loading="pageLoading" class="space-y-6">
    <div class="rounded-2xl border border-brand-200 bg-brand-50 px-5 py-4">
      <p class="text-sm font-medium text-brand-800">偏好设置将影响校招推荐与日历提醒</p>
      <p class="mt-1 text-xs text-brand-600">
        意向岗位与
        <RouterLink to="/profile/info" class="font-medium underline hover:text-brand-700">
          基本资料
        </RouterLink>
        保持同步；城市与公司偏好用于个性化推荐排序。
      </p>
    </div>

    <div class="sidebar-card rounded-2xl bg-white p-6">
      <div class="mb-6">
        <h2 class="text-base font-bold text-ink-900">推荐偏好</h2>
        <p class="mt-1 text-xs text-ink-500">设置您关注的岗位、城市与公司，获得更精准的校招推荐</p>
      </div>

      <el-form ref="formRef" :model="form" :rules="rules" label-position="top">
        <el-form-item label="意向岗位">
          <el-select
            v-model="form.targetPositions"
            multiple
            filterable
            allow-create
            default-first-option
            :placeholder="authStore.userInfo?.college ? '选择或输入意向岗位' : '请先在基本资料中选择学院'"
            class="w-full"
            :disabled="!authStore.userInfo?.college"
          >
            <el-option
              v-for="pos in positionOptions"
              :key="pos"
              :label="pos"
              :value="pos"
            />
          </el-select>
          <p class="mt-1.5 text-xs text-ink-400">
            <template v-if="authStore.userInfo?.college">
              根据您的学院（{{ authStore.userInfo.college }}）推荐岗位；保存后将同步至基本资料
            </template>
            <template v-else>
              请先在
              <RouterLink to="/profile/info" class="text-brand-600 underline">基本资料</RouterLink>
              中选择学院
            </template>
          </p>
        </el-form-item>

        <el-form-item label="偏好城市">
          <el-select
            v-model="form.preferredCities"
            multiple
            filterable
            allow-create
            default-first-option
            placeholder="选择或输入偏好城市"
            class="w-full"
          >
            <el-option v-for="city in CITY_SUGGESTIONS" :key="city" :label="city" :value="city" />
          </el-select>
        </el-form-item>

        <el-form-item label="偏好公司">
          <el-select
            v-model="form.preferredCompanies"
            multiple
            filterable
            allow-create
            default-first-option
            placeholder="选择或输入偏好公司"
            class="w-full"
          >
            <el-option
              v-for="company in COMPANY_SUGGESTIONS"
              :key="company"
              :label="company"
              :value="company"
            />
          </el-select>
        </el-form-item>

        <el-form-item label="特别关注公司">
          <el-select
            v-model="form.focusCompanies"
            multiple
            filterable
            allow-create
            default-first-option
            placeholder="从偏好公司中选择或输入，推荐权重更高"
            class="w-full"
          >
            <el-option
              v-for="company in focusCompanyOptions"
              :key="company"
              :label="company"
              :value="company"
            />
          </el-select>
        </el-form-item>

        <el-divider />

        <el-form-item label="日历提醒提前量" prop="remindBefore">
          <el-checkbox-group v-model="form.remindBefore">
            <el-checkbox
              v-for="opt in REMIND_BEFORE_OPTIONS"
              :key="opt.value"
              :label="opt.value"
            >
              {{ opt.label }}
            </el-checkbox>
          </el-checkbox-group>
          <p class="mt-1.5 text-xs text-ink-400">添加到校招日历的活动将按此设置发送邮件提醒</p>
        </el-form-item>

        <div class="flex flex-wrap gap-3 border-t border-ink-100 pt-6">
          <button
            type="button"
            class="btn-primary rounded-xl px-6 py-2.5 text-sm font-medium text-white disabled:opacity-60"
            :disabled="saving || userStore.loading"
            @click="handleSubmit"
          >
            {{ saving ? '保存中...' : '保存偏好' }}
          </button>
          <button
            type="button"
            class="rounded-xl border border-ink-200 px-6 py-2.5 text-sm font-medium text-ink-600 hover:bg-ink-50"
            :disabled="saving"
            @click="handleReset"
          >
            重置
          </button>
          <button
            type="button"
            class="rounded-xl border border-ink-200 px-6 py-2.5 text-sm font-medium text-ink-500 hover:bg-ink-50"
            :disabled="saving"
            @click="handleClearOptional"
          >
            清空城市/公司偏好
          </button>
        </div>
      </el-form>
    </div>
  </div>
</template>

<style scoped>
:deep(.el-form-item__label) {
  font-weight: 500;
  color: #343a40;
}
:deep(.el-input__wrapper),
:deep(.el-select__wrapper) {
  border-radius: 12px;
}
:deep(.el-checkbox) {
  margin-right: 24px;
}
</style>
