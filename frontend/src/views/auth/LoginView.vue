<script setup lang="ts">
import { RouterLink } from 'vue-router'
import PagePlaceholder from '@/components/common/PagePlaceholder.vue'
import { APP_NAME } from '@/constants'
</script>

<template>
  <div class="card-shadow rounded-2xl bg-white p-8">
    <h1 class="mb-2 text-2xl font-bold text-ink-900">登录</h1>
    <p class="mb-6 text-sm text-ink-500">登录 {{ APP_NAME }} 获取个性化推荐与日程管理</p>
    <PagePlaceholder title="登录表单">
      待接入 Element Plus 表单、loginApi、记住登录与登录失败锁定提示
    </PagePlaceholder>
    <p class="mt-6 text-center text-sm text-ink-500">
      还没有账号？
      <RouterLink to="/register" class="font-medium text-brand-600 hover:text-brand-700">
        去注册
      </RouterLink>
    </p>
  </div>
</template>
