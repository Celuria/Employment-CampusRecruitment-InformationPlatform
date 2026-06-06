<script setup lang="ts">
import { RouterLink } from 'vue-router'
</script>

<template>
  <div class="flex min-h-screen flex-col items-center justify-center px-4">
    <h1 class="mb-2 text-6xl font-bold text-brand-500">404</h1>
    <p class="mb-6 text-ink-500">抱歉，您访问的页面不存在</p>
    <RouterLink
      to="/"
      class="btn-primary rounded-xl px-6 py-2.5 text-sm font-medium text-white"
    >
      返回首页
    </RouterLink>
  </div>
</template>
