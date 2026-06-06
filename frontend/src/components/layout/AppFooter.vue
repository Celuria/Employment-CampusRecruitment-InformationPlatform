<script setup lang="ts">
import { APP_NAME } from '@/constants'

const serviceLinks = [
  { label: '宣讲会信息', to: '/career-talks' },
  { label: '双选会信息', to: '/job-fairs' },
  { label: '个性化推荐', to: '/recommendations' },
  { label: '日程管理', to: '/calendar' },
]

const helpLinks = [
  { label: '使用指南', to: '#' },
  { label: '常见问题', to: '#' },
  { label: '意见反馈', to: '#' },
  { label: '联系客服', to: '#' },
]

const relatedLinks = [
  { label: '学校就业中心', to: '#' },
  { label: '学生事务大厅', to: '#' },
  { label: '教务处', to: '#' },
  { label: '研究生院', to: '#' },
]
</script>

<template>
  <footer class="mt-16 border-t border-ink-200 bg-white">
    <div class="mx-auto max-w-[1440px] px-8 py-10">
      <div class="flex items-start justify-between">
        <div>
          <div class="mb-3 flex items-center gap-3">
            <div class="btn-primary flex h-8 w-8 items-center justify-center rounded-lg">
              <svg class="h-4 w-4 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path
                  stroke-linecap="round"
                  stroke-linejoin="round"
                  stroke-width="2"
                  d="M19 21V5a2 2 0 00-2-2H7a2 2 0 00-2 2v16m14 0h2m-2 0h-5m-9 0H3m2 0h5M9 7h1m-1 4h1m4-4h1m-1 4h1m-5 10v-5a1 1 0 011-1h2a1 1 0 011 1v5m-4 0h4"
                />
              </svg>
            </div>
            <span class="font-bold text-ink-900">{{ APP_NAME }}</span>
          </div>
          <p class="max-w-sm text-xs leading-relaxed text-ink-500">
            为高校学生提供校招信息浏览、个性化推荐、日程管理与邮件提醒的一站式就业信息服务平台
          </p>
        </div>
        <div class="flex gap-12">
          <div>
            <h4 class="mb-3 text-sm font-semibold text-ink-800">平台服务</h4>
            <ul class="space-y-2">
              <li v-for="link in serviceLinks" :key="link.label">
                <RouterLink
                  :to="link.to"
                  class="text-xs text-ink-500 transition-colors hover:text-brand-600"
                >
                  {{ link.label }}
                </RouterLink>
              </li>
            </ul>
          </div>
          <div>
            <h4 class="mb-3 text-sm font-semibold text-ink-800">帮助支持</h4>
            <ul class="space-y-2">
              <li v-for="link in helpLinks" :key="link.label">
                <a href="#" class="text-xs text-ink-500 transition-colors hover:text-brand-600">
                  {{ link.label }}
                </a>
              </li>
            </ul>
          </div>
          <div>
            <h4 class="mb-3 text-sm font-semibold text-ink-800">相关链接</h4>
            <ul class="space-y-2">
              <li v-for="link in relatedLinks" :key="link.label">
                <a href="#" class="text-xs text-ink-500 transition-colors hover:text-brand-600">
                  {{ link.label }}
                </a>
              </li>
            </ul>
          </div>
        </div>
      </div>
      <div class="mt-8 flex items-center justify-between border-t border-ink-100 pt-6">
        <p class="text-xs text-ink-400">{{ APP_NAME }} 版权所有</p>
        <div class="flex items-center gap-4">
          <a href="#" class="text-xs text-ink-400 transition-colors hover:text-ink-600">隐私政策</a>
          <a href="#" class="text-xs text-ink-400 transition-colors hover:text-ink-600">用户协议</a>
          <span class="text-xs text-ink-400">京ICP备XXXXXXXX号</span>
        </div>
      </div>
    </div>
  </footer>
</template>
