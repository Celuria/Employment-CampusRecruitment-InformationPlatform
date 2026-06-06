<script setup lang="ts">
import { computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { useAuthStore } from '@/stores'
import { APP_NAME } from '@/constants'

const route = useRoute()
const router = useRouter()
const authStore = useAuthStore()

const navItems = [
  { name: 'Home', label: '首页', path: '/' },
  { name: 'Recommendations', label: '推荐', path: '/recommendations', requiresAuth: true },
  { name: 'CareerTalkList', label: '宣讲会', path: '/career-talks' },
  { name: 'JobFairList', label: '双选会', path: '/job-fairs' },
  { name: 'Calendar', label: '我的日历', path: '/calendar', requiresAuth: true },
]

const activeNav = computed(() => {
  if (route.path.startsWith('/career-talks')) return 'CareerTalkList'
  if (route.path.startsWith('/job-fairs')) return 'JobFairList'
  if (route.path.startsWith('/profile')) return 'Profile'
  return route.name as string
})

function isActive(name: string) {
  return activeNav.value === name
}

function handleNavClick(item: (typeof navItems)[number]) {
  if (item.requiresAuth && !authStore.isLoggedIn) {
    router.push({ name: 'Login', query: { redirect: item.path } })
    return
  }
  router.push(item.path)
}

function goProfile() {
  if (!authStore.isLoggedIn) {
    router.push({ name: 'Login' })
    return
  }
  router.push('/profile')
}

async function handleLogout() {
  await authStore.logout()
  ElMessage.success('已退出登录')
  router.push('/login')
}

function handleUserCommand(command: string) {
  if (command === 'profile') goProfile()
  if (command === 'logout') handleLogout()
}

function displayName() {
  const info = authStore.userInfo
  if (!info) return '用户'
  return info.name || info.username
}
</script>

<template>
  <header class="sticky top-0 z-50 border-b border-ink-200/60 bg-white/95 backdrop-blur-sm">
    <div class="mx-auto flex h-16 max-w-[1440px] items-center justify-between px-8">
      <!-- Logo -->
      <RouterLink to="/" class="flex items-center gap-3">
        <div class="btn-primary flex h-9 w-9 items-center justify-center rounded-xl">
          <svg class="h-5 w-5 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path
              stroke-linecap="round"
              stroke-linejoin="round"
              stroke-width="2"
              d="M19 21V5a2 2 0 00-2-2H7a2 2 0 00-2 2v16m14 0h2m-2 0h-5m-9 0H3m2 0h5M9 7h1m-1 4h1m4-4h1m-1 4h1m-5 10v-5a1 1 0 011-1h2a1 1 0 011 1v5m-4 0h4"
            />
          </svg>
        </div>
        <span class="text-lg font-bold tracking-tight text-ink-900">{{ APP_NAME }}</span>
      </RouterLink>

      <!-- 主导航 -->
      <nav class="flex items-center gap-1">
        <button
          v-for="item in navItems"
          :key="item.name"
          type="button"
          class="rounded-lg px-4 py-2 text-sm transition-colors"
          :class="
            isActive(item.name)
              ? 'bg-brand-50 font-semibold text-brand-600'
              : 'font-medium text-ink-600 hover:bg-ink-100'
          "
          @click="handleNavClick(item)"
        >
          {{ item.label }}
        </button>
      </nav>

      <!-- 用户区域 -->
      <div class="flex items-center gap-4">
        <template v-if="authStore.isLoggedIn">
          <button
            type="button"
            class="relative rounded-lg p-2 transition-colors hover:bg-ink-100"
            title="通知"
          >
            <svg class="h-5 w-5 text-ink-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path
                stroke-linecap="round"
                stroke-linejoin="round"
                stroke-width="2"
                d="M15 17h5l-1.405-1.405A2.032 2.032 0 0118 14.158V11a6.002 6.002 0 00-4-5.659V5a2 2 0 10-4 0v.341C7.67 6.165 6 8.388 6 11v3.159c0 .538-.214 1.055-.595 1.436L4 17h5m6 0v1a3 3 0 11-6 0v-1m6 0H9"
              />
            </svg>
            <span class="absolute right-1.5 top-1.5 h-2 w-2 rounded-full bg-red-500" />
          </button>
          <el-dropdown trigger="click" @command="handleUserCommand">
            <button
              type="button"
              class="flex items-center gap-3 border-l border-ink-200 pl-4 outline-none"
            >
              <div class="text-right">
                <p class="text-sm font-semibold text-ink-800">{{ displayName() }}</p>
                <p class="text-xs text-ink-500">
                  {{ authStore.userInfo?.college || '完善资料' }}
                  <template v-if="authStore.userInfo?.grade"> · {{ authStore.userInfo.grade }}</template>
                </p>
              </div>
              <div
                class="flex h-10 w-10 items-center justify-center rounded-full border-2 border-brand-200 bg-brand-100 text-sm font-semibold text-brand-700"
              >
                {{ displayName().charAt(0).toUpperCase() }}
              </div>
            </button>
            <template #dropdown>
              <el-dropdown-menu>
                <el-dropdown-item command="profile">个人中心</el-dropdown-item>
                <el-dropdown-item divided command="logout">退出登录</el-dropdown-item>
              </el-dropdown-menu>
            </template>
          </el-dropdown>
        </template>
        <template v-else>
          <RouterLink
            to="/login"
            class="rounded-lg px-4 py-2 text-sm font-medium text-ink-600 transition-colors hover:bg-ink-100"
          >
            登录
          </RouterLink>
          <RouterLink
            to="/register"
            class="btn-primary rounded-xl px-4 py-2 text-sm font-medium text-white"
          >
            注册
          </RouterLink>
        </template>
      </div>
    </div>
  </header>
</template>
