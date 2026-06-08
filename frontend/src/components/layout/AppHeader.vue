<script setup lang="ts">
import { computed, ref, watch, onUnmounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { useAuthStore } from '@/stores'
import { APP_NAME } from '@/constants'
import { getReminderLogsApi } from '@/api'
import type { ReminderLog } from '@/types'

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

function goAdmin() {
  router.push('/admin')
}

function handleUserCommand(command: string) {
  if (command === 'profile') goProfile()
  if (command === 'admin') goAdmin()
  if (command === 'logout') handleLogout()
}

function displayName() {
  const info = authStore.userInfo
  if (!info) return '用户'
  return info.name || info.username
}

// ============ 消息通知 ============
const LS_KEY = 'reminder_last_seen_at'

const reminders = ref<ReminderLog[]>([])
const notificationVisible = ref(false)
const notificationLoading = ref(false)
const lastSeen = ref(localStorage.getItem(LS_KEY) || '')

function persistLastSeen(time: string) {
  localStorage.setItem(LS_KEY, time)
  lastSeen.value = time
}

const hasNewReminders = computed(() => {
  if (!lastSeen.value) return reminders.value.length > 0
  return reminders.value.some((r) => r.sentTime && r.sentTime > lastSeen.value)
})

async function fetchReminders() {
  if (!authStore.isLoggedIn) return
  notificationLoading.value = true
  try {
    const result = await getReminderLogsApi({ page: 1, pageSize: 5 })
    reminders.value = result.list
  } catch {
    reminders.value = []
  } finally {
    notificationLoading.value = false
  }
}

async function onNotificationShow() {
  await fetchReminders()
  // 打开弹窗时，将最新 sentTime 记为已读时间
  const latest = reminders.value
    .filter((r) => r.sentTime)
    .map((r) => r.sentTime!)
    .sort()
    .reverse()[0]
  if (latest) {
    persistLastSeen(latest)
  } else if (reminders.value.length === 0) {
    // 如果没有消息，也更新时间戳，避免下次误报
    persistLastSeen(new Date().toISOString())
  }
}

function formatTime(iso: string): string {
  if (!iso) return ''
  const d = new Date(iso)
  const pad = (n: number) => String(n).padStart(2, '0')
  return `${d.getFullYear()}-${pad(d.getMonth() + 1)}-${pad(d.getDate())} ${pad(d.getHours())}:${pad(d.getMinutes())}`
}

function goReminders() {
  notificationVisible.value = false
  router.push('/profile/reminders')
}

// 登录后自动拉取；每 30 秒轮询，确保调度器处理后的提醒能及时出现红点
let pollTimer: ReturnType<typeof setInterval> | null = null

watch(() => authStore.isLoggedIn, (val) => {
  if (val) {
    fetchReminders()
    if (!pollTimer) {
      pollTimer = setInterval(fetchReminders, 30_000)
    }
  } else {
    if (pollTimer) {
      clearInterval(pollTimer)
      pollTimer = null
    }
  }
}, { immediate: true })

onUnmounted(() => {
  if (pollTimer) {
    clearInterval(pollTimer)
    pollTimer = null
  }
})
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
          <el-popover
            v-model:visible="notificationVisible"
            trigger="click"
            placement="bottom-end"
            :width="360"
            :offset="8"
            popper-class="notification-popover"
            @show="onNotificationShow"
          >
            <template #reference>
              <button
                type="button"
                class="relative rounded-lg p-2 transition-colors hover:bg-ink-100"
                title="消息通知"
              >
                <svg class="h-5 w-5 text-ink-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path
                    stroke-linecap="round"
                    stroke-linejoin="round"
                    stroke-width="2"
                    d="M15 17h5l-1.405-1.405A2.032 2.032 0 0118 14.158V11a6.002 6.002 0 00-4-5.659V5a2 2 0 10-4 0v.341C7.67 6.165 6 8.388 6 11v3.159c0 .538-.214 1.055-.595 1.436L4 17h5m6 0v1a3 3 0 11-6 0v-1m6 0H9"
                  />
                </svg>
                <span
                  v-if="hasNewReminders"
                  class="absolute right-1.5 top-1.5 h-2.5 w-2.5 rounded-full bg-red-500 ring-2 ring-white"
                />
              </button>
            </template>

            <!-- 弹窗内容 -->
            <div class="flex flex-col">
              <div class="flex items-center justify-between border-b border-ink-100 pb-3">
                <span class="text-sm font-semibold text-ink-800">消息通知</span>
                <button
                  type="button"
                  class="text-xs text-brand-600 hover:text-brand-700"
                  @click="goReminders"
                >
                  查看全部
                </button>
              </div>

              <div v-if="notificationLoading" class="flex items-center justify-center py-10">
                <span class="text-sm text-ink-400">加载中...</span>
              </div>

              <div v-else-if="reminders.length === 0" class="flex flex-col items-center py-10">
                <svg class="mb-3 h-10 w-10 text-ink-300" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path
                    stroke-linecap="round"
                    stroke-linejoin="round"
                    stroke-width="1.5"
                    d="M15 17h5l-1.405-1.405A2.032 2.032 0 0118 14.158V11a6.002 6.002 0 00-4-5.659V5a2 2 0 10-4 0v.341C7.67 6.165 6 8.388 6 11v3.159c0 .538-.214 1.055-.595 1.436L4 17h5m6 0v1a3 3 0 11-6 0v-1m6 0H9"
                  />
                </svg>
                <span class="text-sm text-ink-400">暂无消息</span>
              </div>

              <div v-else class="-mx-3 max-h-80 overflow-y-auto">
                <button
                  v-for="reminder in reminders"
                  :key="reminder.id"
                  type="button"
                  class="flex w-full items-start gap-3 px-3 py-3 text-left transition-colors hover:bg-ink-50"
                  @click="goReminders"
                >
                  <div
                    class="mt-0.5 flex h-8 w-8 shrink-0 items-center justify-center rounded-full"
                    :class="reminder.status === 'success' ? 'bg-green-100 text-green-600' : 'bg-red-100 text-red-500'"
                  >
                    <svg v-if="reminder.status === 'sent'" class="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
                    </svg>
                    <svg v-else-if="reminder.status === 'pending'" class="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
                    </svg>
                    <svg v-else class="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
                    </svg>
                  </div>
                  <div class="min-w-0 flex-1">
                    <p class="text-sm font-medium text-ink-700 truncate">{{ reminder.eventTitle }}</p>
                    <p class="mt-0.5 text-xs text-ink-400">
                      {{ reminder.scheduledTime ? formatTime(reminder.scheduledTime) : '' }}
                    </p>
                  </div>
                  <span
                    class="mt-1 shrink-0 text-xs"
                    :class="{
                      'text-green-500': reminder.status === 'sent',
                      'text-yellow-500': reminder.status === 'pending',
                      'text-red-400': reminder.status === 'failed' || reminder.status === 'cancelled',
                    }"
                  >
                    {{ reminder.status === 'sent' ? '已提醒' : reminder.status === 'pending' ? '待提醒' : reminder.status === 'cancelled' ? '已取消' : '失败' }}
                  </span>
                </button>
              </div>
            </div>
          </el-popover>
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
                <el-dropdown-item v-if="authStore.isAdmin" command="admin">管理后台</el-dropdown-item>
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
