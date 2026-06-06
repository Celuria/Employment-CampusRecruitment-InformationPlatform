<script setup lang="ts">
import { RouterLink } from 'vue-router'
import PagePlaceholder from '@/components/common/PagePlaceholder.vue'
import { APP_NAME } from '@/constants'
</script>

<template>
  <div class="card-shadow rounded-2xl bg-white p-8">
    <h1 class="mb-2 text-2xl font-bold text-ink-900">注册</h1>
    <p class="mb-6 text-sm text-ink-500">创建 {{ APP_NAME }} 账号</p>
    <PagePlaceholder title="注册表单">
      待接入 registerApi、密码强度校验与确认密码校验
    </PagePlaceholder>
    <p class="mt-6 text-center text-sm text-ink-500">
      已有账号？
      <RouterLink to="/login" class="font-medium text-brand-600 hover:text-brand-700">
        去登录
      </RouterLink>
    </p>
  </div>
</template>
