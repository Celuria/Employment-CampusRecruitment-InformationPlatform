<script setup lang="ts">
import { reactive, ref } from 'vue'
import { useRoute, useRouter, RouterLink } from 'vue-router'
import type { FormInstance, FormRules } from 'element-plus'
import { ElMessage } from 'element-plus'
import { useAuthStore } from '@/stores'
import { APP_NAME } from '@/constants'
import type { LoginForm } from '@/types'

const route = useRoute()
const router = useRouter()
const authStore = useAuthStore()

const formRef = ref<FormInstance>()
const form = reactive<LoginForm>({
  username: '',
  password: '',
  remember: false,
})

const rules: FormRules<LoginForm> = {
  username: [{ required: true, message: '请输入账号', trigger: 'blur' }],
  password: [{ required: true, message: '请输入密码', trigger: 'blur' }],
}

async function handleSubmit() {
  if (!formRef.value) return
  const valid = await formRef.value.validate().catch(() => false)
  if (!valid) return

  try {
    await authStore.login({ ...form })
    ElMessage.success('登录成功')
    let redirect = (route.query.redirect as string) || '/'
    if (!route.query.redirect && authStore.userInfo && !authStore.userInfo.profileCompleted) {
      redirect = '/profile/info'
    }
    router.push(redirect)
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
            d="M19 21V5a2 2 0 00-2-2H7a2 2 0 00-2 2v16m14 0h2m-2 0h-5m-9 0H3m2 0h5M9 7h1m-1 4h1m4-4h1m-1 4h1m-5 10v-5a1 1 0 011-1h2a1 1 0 011 1v5m-4 0h4"
          />
        </svg>
      </div>
      <h1 class="text-2xl font-bold text-ink-900">欢迎回来</h1>
      <p class="mt-1 text-sm text-ink-500">登录 {{ APP_NAME }}</p>
    </div>

    <el-form ref="formRef" :model="form" :rules="rules" label-position="top" @submit.prevent="handleSubmit">
      <el-form-item label="账号" prop="username">
        <el-input
          v-model="form.username"
          size="large"
          placeholder="学号 / 邮箱账号"
          clearable
          @keyup.enter="handleSubmit"
        />
      </el-form-item>

      <el-form-item label="密码" prop="password">
        <el-input
          v-model="form.password"
          type="password"
          size="large"
          placeholder="请输入密码"
          show-password
          @keyup.enter="handleSubmit"
        />
      </el-form-item>

      <div class="mb-6 flex items-center justify-between">
        <el-checkbox v-model="form.remember">记住登录</el-checkbox>
      </div>

      <button
        type="submit"
        class="btn-primary w-full rounded-xl py-3 text-sm font-medium text-white transition-all hover:shadow-lg disabled:cursor-not-allowed disabled:opacity-60"
        :disabled="authStore.loading"
      >
        {{ authStore.loading ? '登录中...' : '登录' }}
      </button>
    </el-form>

    <p class="mt-6 text-center text-sm text-ink-500">
      还没有账号？
      <RouterLink to="/register" class="font-medium text-brand-600 hover:text-brand-700">
        立即注册
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
