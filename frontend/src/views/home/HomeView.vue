<script setup lang="ts">
import PageHeader from '@/components/business/PageHeader.vue'
</script>

<template>
  <div>
    <PageHeader
      title="欢迎回来"
      description="一站式校招信息浏览、个性化推荐与日程管理"
    />
    <div class="grid gap-6 md:grid-cols-2 lg:grid-cols-4">
      <RouterLink
        to="/recommendations"
        class="card-shadow rounded-2xl bg-white p-6 transition-transform hover:-translate-y-0.5"
      >
        <h3 class="mb-2 font-bold text-ink-900">个性化推荐</h3>
        <p class="text-sm text-ink-500">根据画像匹配宣讲会与双选会</p>
      </RouterLink>
      <RouterLink
        to="/career-talks"
        class="card-shadow rounded-2xl bg-white p-6 transition-transform hover:-translate-y-0.5"
      >
        <h3 class="mb-2 font-bold text-ink-900">宣讲会</h3>
        <p class="text-sm text-ink-500">浏览与检索企业宣讲活动</p>
      </RouterLink>
      <RouterLink
        to="/job-fairs"
        class="card-shadow rounded-2xl bg-white p-6 transition-transform hover:-translate-y-0.5"
      >
        <h3 class="mb-2 font-bold text-ink-900">双选会</h3>
        <p class="text-sm text-ink-500">大型招聘会信息与报名</p>
      </RouterLink>
      <RouterLink
        to="/calendar"
        class="card-shadow rounded-2xl bg-white p-6 transition-transform hover:-translate-y-0.5"
      >
        <h3 class="mb-2 font-bold text-ink-900">我的日历</h3>
        <p class="text-sm text-ink-500">管理日程与邮件提醒</p>
      </RouterLink>
    </div>
  </div>
</template>
