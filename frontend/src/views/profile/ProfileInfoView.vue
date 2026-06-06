<script setup lang="ts">
import { computed, onMounted, reactive, ref, watch } from 'vue'
import type { FormInstance, FormRules } from 'element-plus'
import { ElMessage } from 'element-plus'
import { useAuthStore } from '@/stores'
import {
  COLLEGE_MAJOR_MAP,
  COLLEGE_OPTIONS,
  GRADE_OPTIONS,
  POSITION_SUGGESTIONS,
} from '@/constants/profile'
import type { UpdateProfileForm, UserInfo } from '@/types'

const authStore = useAuthStore()
const formRef = ref<FormInstance>()
const pageLoading = ref(true)
const saving = ref(false)

const form = reactive<UpdateProfileForm>({
  name: '',
  college: '',
  major: '',
  grade: '',
  targetPositions: [],
  phone: '',
  email: '',
})

const accountInfo = ref<Pick<UserInfo, 'username' | 'role' | 'profileCompleted'>>({
  username: '',
  role: 'student',
  profileCompleted: false,
})

const majorOptions = computed(() => {
  if (!form.college) return []
  return COLLEGE_MAJOR_MAP[form.college] || []
})

const rules: FormRules<UpdateProfileForm> = {
  name: [{ required: true, message: '请输入姓名', trigger: 'blur' }],
  college: [{ required: true, message: '请选择学院', trigger: 'change' }],
  major: [{ required: true, message: '请选择专业', trigger: 'change' }],
  targetPositions: [
    {
      type: 'array',
      required: true,
      min: 1,
      message: '请至少选择一个意向岗位',
      trigger: 'change',
    },
  ],
  email: [
    { required: true, message: '请输入邮箱', trigger: 'blur' },
    { type: 'email', message: '邮箱格式不正确', trigger: 'blur' },
  ],
  phone: [
    {
      pattern: /^$|^1[3-9]\d{9}$/,
      message: '请输入正确的手机号',
      trigger: 'blur',
    },
  ],
}

watch(
  () => form.college,
  (college, prev) => {
    if (prev && college !== prev) {
      form.major = ''
    }
  },
)

function fillForm(data: UserInfo) {
  form.name = data.name || ''
  form.college = data.college || ''
  form.major = data.major || ''
  form.grade = data.grade || ''
  form.targetPositions = [...(data.targetPositions || [])]
  form.phone = data.phone || ''
  form.email = data.email || ''
  accountInfo.value = {
    username: data.username,
    role: data.role,
    profileCompleted: !!data.profileCompleted,
  }
}

async function loadProfile() {
  pageLoading.value = true
  try {
    await authStore.fetchUserInfo()
    if (authStore.userInfo) {
      fillForm(authStore.userInfo)
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
    await authStore.updateProfile({
      name: form.name.trim(),
      college: form.college,
      major: form.major,
      grade: form.grade || undefined,
      targetPositions: form.targetPositions,
      phone: form.phone?.trim() || undefined,
      email: form.email.trim(),
    })
    if (authStore.userInfo) {
      fillForm(authStore.userInfo)
    }
    ElMessage.success('资料已保存')
  } catch {
    // 错误由 axios 拦截器提示
  } finally {
    saving.value = false
  }
}

async function handleReset() {
  await loadProfile()
  formRef.value?.clearValidate()
  ElMessage.info('已恢复为服务器上的最新资料')
}

async function handleClearOptional() {
  form.grade = ''
  form.phone = ''
  formRef.value?.clearValidate(['phone'])
  ElMessage.info('已清空选填项，请点击保存生效')
}

onMounted(() => {
  loadProfile()
})
</script>

<template>
  <div v-loading="pageLoading" class="space-y-6">
    <!-- 完善引导 -->
    <div
      v-if="!accountInfo.profileCompleted"
      class="rounded-2xl border border-brand-200 bg-brand-50 px-5 py-4"
    >
      <p class="text-sm font-medium text-brand-800">完善资料以获得更精准的校招推荐</p>
      <p class="mt-1 text-xs text-brand-600">
        请填写姓名、学院、专业与至少一个意向岗位，保存后即可使用个性化推荐功能。
      </p>
    </div>

    <!-- 账号信息（只读） -->
    <div class="sidebar-card rounded-2xl bg-white p-6">
      <h2 class="mb-4 text-base font-bold text-ink-900">账号信息</h2>
      <dl class="grid gap-4 sm:grid-cols-2">
        <div>
          <dt class="text-xs text-ink-500">登录账号</dt>
          <dd class="mt-1 text-sm font-medium text-ink-800">{{ accountInfo.username || '—' }}</dd>
        </div>
        <div>
          <dt class="text-xs text-ink-500">账号角色</dt>
          <dd class="mt-1 text-sm font-medium text-ink-800">
            {{ accountInfo.role === 'admin' ? '管理员' : '学生用户' }}
          </dd>
        </div>
      </dl>
    </div>

    <!-- 基本资料表单 -->
    <div class="sidebar-card rounded-2xl bg-white p-6">
      <div class="mb-6 flex items-center justify-between">
        <div>
          <h2 class="text-base font-bold text-ink-900">基本资料</h2>
          <p class="mt-1 text-xs text-ink-500">以下信息将用于个性化推荐与邮件提醒</p>
        </div>
        <span
          v-if="accountInfo.profileCompleted"
          class="rounded-full bg-brand-100 px-3 py-1 text-xs font-medium text-brand-700"
        >
          已完善
        </span>
      </div>

      <el-form ref="formRef" :model="form" :rules="rules" label-position="top">
        <div class="grid gap-x-6 sm:grid-cols-2">
          <el-form-item label="姓名" prop="name">
            <el-input v-model="form.name" placeholder="真实姓名" clearable />
          </el-form-item>

          <el-form-item label="邮箱" prop="email">
            <el-input v-model="form.email" placeholder="用于邮件提醒" clearable />
          </el-form-item>

          <el-form-item label="学院" prop="college">
            <el-select v-model="form.college" placeholder="请选择学院" class="w-full" filterable>
              <el-option v-for="c in COLLEGE_OPTIONS" :key="c" :label="c" :value="c" />
            </el-select>
          </el-form-item>

          <el-form-item label="专业" prop="major">
            <el-select
              v-model="form.major"
              placeholder="请先选择学院"
              class="w-full"
              filterable
              :disabled="!form.college"
            >
              <el-option v-for="m in majorOptions" :key="m" :label="m" :value="m" />
            </el-select>
          </el-form-item>

          <el-form-item label="年级">
            <el-select v-model="form.grade" placeholder="选填" class="w-full" clearable>
              <el-option
                v-for="g in GRADE_OPTIONS"
                :key="g.value"
                :label="g.label"
                :value="g.value"
              />
            </el-select>
          </el-form-item>

          <el-form-item label="联系电话" prop="phone">
            <el-input v-model="form.phone" placeholder="选填" clearable maxlength="11" />
          </el-form-item>
        </div>

        <el-form-item label="就业意向岗位" prop="targetPositions">
          <el-select
            v-model="form.targetPositions"
            multiple
            filterable
            allow-create
            default-first-option
            placeholder="选择或输入意向岗位，至少 1 个"
            class="w-full"
          >
            <el-option
              v-for="pos in POSITION_SUGGESTIONS"
              :key="pos"
              :label="pos"
              :value="pos"
            />
          </el-select>
        </el-form-item>

        <div class="flex flex-wrap gap-3 border-t border-ink-100 pt-6">
          <button
            type="button"
            class="btn-primary rounded-xl px-6 py-2.5 text-sm font-medium text-white disabled:opacity-60"
            :disabled="saving"
            @click="handleSubmit"
          >
            {{ saving ? '保存中...' : '保存资料' }}
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
            清空选填项
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
</style>
