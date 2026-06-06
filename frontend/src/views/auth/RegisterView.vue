<script setup lang="ts">
import { reactive, ref } from 'vue'
import { useRouter, RouterLink } from 'vue-router'
import type { FormInstance, FormRules } from 'element-plus'
import { ElMessage } from 'element-plus'
import { useAuthStore } from '@/stores'
import { APP_NAME } from '@/constants'
import type { RegisterForm } from '@/types'
import {
  passwordRuleMessage,
  usernameRuleMessage,
  validatePassword,
  validateUsername,
} from '@/utils/validate'

const router = useRouter()
const authStore = useAuthStore()

const formRef = ref<FormInstance>()
const form = reactive<RegisterForm>({
  username: '',
  password: '',
  confirmPassword: '',
  email: '',
})

const validateConfirmPassword = (_rule: unknown, value: string, callback: (err?: Error) => void) => {
  if (!value) {
    callback(new Error('请再次输入密码'))
  } else if (value !== form.password) {
    callback(new Error('两次输入的密码不一致'))
  } else {
    callback()
  }
}

const rules: FormRules<RegisterForm> = {
  username: [
    { required: true, message: '请输入账号', trigger: 'blur' },
    {
      validator: (_r, v, cb) => (validateUsername(v) ? cb() : cb(new Error(usernameRuleMessage))),
      trigger: 'blur',
    },
  ],
  email: [
    { required: true, message: '请输入邮箱', trigger: 'blur' },
    { type: 'email', message: '邮箱格式不正确', trigger: 'blur' },
  ],
  password: [
    { required: true, message: '请输入密码', trigger: 'blur' },
    {
      validator: (_r, v, cb) => (validatePassword(v) ? cb() : cb(new Error(passwordRuleMessage))),
      trigger: 'blur',
    },
  ],
  confirmPassword: [{ validator: validateConfirmPassword, trigger: 'blur' }],
}

async function handleSubmit() {
  if (!formRef.value) return
  const valid = await formRef.value.validate().catch(() => false)
  if (!valid) return

  try {
    await authStore.register({ ...form })
    ElMessage.success('注册成功，请登录')
    router.push('/login')
  } catch {
    // 错误提示由 axios 拦截器处理
  }
}
</script>

<template>
  <div class="card-shadow rounded-2xl bg-white p-8">
    <div class="mb-6 text-center">
      <div class="btn-primary mx-auto mb-4 flex h-12 w-12 items-center justify-center rounded-xl">
        <svg class="h-6 w-6 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path
            stroke-linecap="round"
            stroke-linejoin="round"
            stroke-width="2"
            d="M18 9v3m0 0v3m0-3h3m-3 0h-3m-2-5a4 4 0 11-8 0 4 4 0 018 0zM3 20a6 6 0 0112 0v1H3v-1z"
          />
        </svg>
      </div>
      <h1 class="text-2xl font-bold text-ink-900">创建账号</h1>
      <p class="mt-1 text-sm text-ink-500">加入 {{ APP_NAME }}</p>
    </div>

    <el-form ref="formRef" :model="form" :rules="rules" label-position="top" @submit.prevent="handleSubmit">
      <el-form-item label="账号" prop="username">
        <el-input v-model="form.username" size="large" placeholder="4-32 位字母、数字或下划线" clearable />
      </el-form-item>

      <el-form-item label="邮箱" prop="email">
        <el-input v-model="form.email" size="large" placeholder="用于接收邮件提醒" clearable />
      </el-form-item>

      <el-form-item label="密码" prop="password">
        <el-input
          v-model="form.password"
          type="password"
          size="large"
          placeholder="至少 8 位，含字母和数字"
          show-password
        />
      </el-form-item>

      <el-form-item label="确认密码" prop="confirmPassword">
        <el-input
          v-model="form.confirmPassword"
          type="password"
          size="large"
          placeholder="请再次输入密码"
          show-password
          @keyup.enter="handleSubmit"
        />
      </el-form-item>

      <button
        type="submit"
        class="btn-primary mt-2 w-full rounded-xl py-3 text-sm font-medium text-white transition-all hover:shadow-lg disabled:cursor-not-allowed disabled:opacity-60"
        :disabled="authStore.loading"
      >
        {{ authStore.loading ? '注册中...' : '注册' }}
      </button>
    </el-form>

    <p class="mt-6 text-center text-sm text-ink-500">
      已有账号？
      <RouterLink to="/login" class="font-medium text-brand-600 hover:text-brand-700">
        去登录
      </RouterLink>
    </p>
  </div>
</template>

<style scoped>
:deep(.el-form-item__label) {
  font-weight: 500;
  color: #343a40;
}
:deep(.el-input__wrapper) {
  border-radius: 12px;
  padding: 4px 12px;
}
</style>
