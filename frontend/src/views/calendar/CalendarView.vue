<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import PageHeader from '@/components/business/PageHeader.vue'
import { getCalendarEventsApi, deleteCalendarEventApi, updateCalendarEventApi } from '@/api/modules/calendar'
import type { CalendarEvent } from '@/types'

const loading = ref(true)
const events = ref<CalendarEvent[]>([])
const viewMode = ref<'list' | 'month'>('list')

// 编辑弹窗
const editDialogVisible = ref(false)
const editTarget = ref<CalendarEvent | null>(null)
const editForm = ref({ customNote: '', remindBefore: [] as string[] })
const editSaving = ref(false)

// 删除确认
const deleteTarget = ref<CalendarEvent | null>(null)
const deleteLoading = ref(false)

const REMIND_OPTIONS = [
  { label: '提前 1 小时', value: '1h' },
  { label: '提前 1 天', value: '1d' },
  { label: '提前 3 天', value: '3d' },
]

const EVENT_TYPE_LABEL: Record<string, string> = {
  career_talk: '宣讲会',
  job_fair: '双选会',
}

const REMINDER_STATUS_MAP: Record<string, { label: string; class: string }> = {
  pending: { label: '待发送', class: 'bg-warm-100 text-warm-300' },
  sent: { label: '已发送', class: 'bg-brand-100 text-brand-700' },
  failed: { label: '发送失败', class: 'bg-red-100 text-red-600' },
  cancelled: { label: '已取消', class: 'bg-ink-100 text-ink-500' },
}

async function fetchEvents() {
  loading.value = true
  try {
    const res = await getCalendarEventsApi()
    events.value = Array.isArray(res) ? res : []
  } catch {
    events.value = []
  } finally {
    loading.value = false
  }
}

// ---- 月视图 ----
const currentMonth = ref(new Date())

const calendarDays = computed(() => {
  const year = currentMonth.value.getFullYear()
  const month = currentMonth.value.getMonth()
  const firstDay = new Date(year, month, 1)
  const lastDay = new Date(year, month + 1, 0)
  const startDow = firstDay.getDay() // 0=Sun
  const days: { date: Date; events: CalendarEvent[] }[] = []

  // 上月补齐
  for (let i = startDow - 1; i >= 0; i--) {
    days.push({ date: new Date(year, month, -i), events: [] })
  }
  // 本月
  for (let d = 1; d <= lastDay.getDate(); d++) {
    const date = new Date(year, month, d)
    const dayEvents = events.value.filter(e => {
      const ed = new Date(e.startTime)
      return ed.getFullYear() === year && ed.getMonth() === month && ed.getDate() === d
    })
    days.push({ date, events: dayEvents })
  }
  // 补齐末尾
  while (days.length % 7 !== 0) {
    const last = days[days.length - 1].date
    days.push({ date: new Date(last.getFullYear(), last.getMonth(), last.getDate() + 1), events: [] })
  }
  return days
})

function prevMonth() {
  currentMonth.value = new Date(currentMonth.value.getFullYear(), currentMonth.value.getMonth() - 1, 1)
}
function nextMonth() {
  currentMonth.value = new Date(currentMonth.value.getFullYear(), currentMonth.value.getMonth() + 1, 1)
}
function isCurrentMonth(date: Date) {
  return date.getMonth() === currentMonth.value.getMonth()
}
function isToday(date: Date) {
  const t = new Date()
  return date.getDate() === t.getDate() && date.getMonth() === t.getMonth() && date.getFullYear() === t.getFullYear()
}

// ---- 列表分组 ----
const groupedEvents = computed(() => {
  const groups: { label: string; events: CalendarEvent[] }[] = []
  const now = new Date()
  const upcoming = events.value.filter(e => new Date(e.startTime) >= now)
  const past = events.value.filter(e => new Date(e.startTime) < now)
  if (upcoming.length) groups.push({ label: '即将到来', events: upcoming })
  if (past.length) groups.push({ label: '已结束', events: past })
  return groups
})

function formatDateTime(str: string) {
  return new Date(str).toLocaleString('zh-CN', {
    year: 'numeric', month: 'short', day: 'numeric',
    weekday: 'short', hour: '2-digit', minute: '2-digit',
  })
}

// ---- 编辑 ----
function openEdit(event: CalendarEvent) {
  editTarget.value = event
  editForm.value = {
    customNote: event.customNote || '',
    remindBefore: Array.isArray(event.remindBefore) ? [...(event.remindBefore as string[])] : [event.remindBefore as string],
  }
  editDialogVisible.value = true
}

async function handleSaveEdit() {
  if (!editTarget.value) return
  editSaving.value = true
  try {
    const updated = await updateCalendarEventApi(editTarget.value.id, {
      customNote: editForm.value.customNote,
      remindBefore: editForm.value.remindBefore as any,
    })
    const idx = events.value.findIndex(e => e.id === editTarget.value!.id)
    if (idx !== -1) events.value[idx] = { ...events.value[idx], ...updated }
    editDialogVisible.value = false
    ElMessage.success('日程已更新')
  } catch {
    // 错误由拦截器处理
  } finally {
    editSaving.value = false
  }
}

// ---- 删除 ----
async function handleDelete(event: CalendarEvent) {
  try {
    await ElMessageBox.confirm(
      `确认删除「${event.title}」？删除后将取消未发送的提醒。`,
      '删除确认',
      { confirmButtonText: '确认删除', cancelButtonText: '取消', type: 'warning' }
    )
  } catch {
    return
  }
  deleteTarget.value = event
  deleteLoading.value = true
  try {
    await deleteCalendarEventApi(event.id)
    events.value = events.value.filter(e => e.id !== event.id)
    ElMessage.success('日程已删除')
  } catch {
    // 错误由拦截器处理
  } finally {
    deleteLoading.value = false
    deleteTarget.value = null
  }
}

onMounted(fetchEvents)
</script>

<template>
  <div>
    <PageHeader title="我的日历" description="统一管理宣讲会与双选会日程，按时收到邮件提醒" />

    <!-- 视图切换 -->
    <div class="mb-6 flex items-center justify-between">
      <div class="flex rounded-xl border border-ink-200 bg-white p-1">
        <button
          type="button"
          class="flex items-center gap-1.5 rounded-lg px-4 py-2 text-sm font-medium transition-colors"
          :class="viewMode === 'list' ? 'bg-brand-500 text-white' : 'text-ink-600 hover:bg-ink-50'"
          @click="viewMode = 'list'"
        >
          <svg class="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M4 10h16M4 14h16M4 18h16" />
          </svg>
          列表
        </button>
        <button
          type="button"
          class="flex items-center gap-1.5 rounded-lg px-4 py-2 text-sm font-medium transition-colors"
          :class="viewMode === 'month' ? 'bg-brand-500 text-white' : 'text-ink-600 hover:bg-ink-50'"
          @click="viewMode = 'month'"
        >
          <svg class="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z" />
          </svg>
          月视图
        </button>
      </div>
      <p class="text-sm text-ink-500">共 {{ events.length }} 个日程</p>
    </div>

    <!-- 加载中 -->
    <div v-if="loading" class="space-y-4">
      <div v-for="i in 3" :key="i" class="h-24 animate-pulse rounded-2xl bg-ink-100" />
    </div>

    <!-- 空状态 -->
    <div v-else-if="events.length === 0" class="flex flex-col items-center justify-center rounded-2xl border border-dashed border-ink-200 bg-white py-20">
      <svg class="mb-4 h-12 w-12 text-ink-300" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z" />
      </svg>
      <p class="mb-2 text-base font-medium text-ink-600">暂无日程</p>
      <p class="mb-4 text-sm text-ink-400">在宣讲会或双选会页面点击「加入日历」添加日程</p>
      <RouterLink to="/career-talks" class="btn-primary rounded-xl px-5 py-2.5 text-sm font-medium text-white">
        浏览宣讲会
      </RouterLink>
    </div>

    <!-- 列表视图 -->
    <template v-else-if="viewMode === 'list'">
      <div v-for="group in groupedEvents" :key="group.label" class="mb-8">
        <h2 class="mb-4 text-base font-bold text-ink-700">{{ group.label }}</h2>
        <div class="space-y-4">
          <div
            v-for="event in group.events"
            :key="event.id"
            class="card-shadow rounded-2xl bg-white p-5"
          >
            <div class="flex items-start gap-4">
              <!-- 类型徽章 -->
              <div class="shrink-0 rounded-xl px-3 py-2 text-center text-xs font-bold"
                :class="event.eventType === 'career_talk' ? 'bg-brand-100 text-brand-700' : 'bg-blue-100 text-blue-700'"
              >
                {{ EVENT_TYPE_LABEL[event.eventType] }}
              </div>
              <div class="min-w-0 flex-1">
                <div class="flex items-start justify-between gap-3">
                  <h3 class="font-semibold text-ink-900">{{ event.title }}</h3>
                  <div class="flex shrink-0 gap-2">
                    <button type="button" class="rounded-lg border border-ink-200 px-3 py-1.5 text-xs text-ink-600 transition-colors hover:border-brand-300 hover:text-brand-600"
                      @click="openEdit(event)">
                      编辑
                    </button>
                    <button type="button"
                      class="rounded-lg border border-red-200 px-3 py-1.5 text-xs text-red-500 transition-colors hover:bg-red-50 disabled:opacity-50"
                      :disabled="deleteTarget?.id === event.id && deleteLoading"
                      @click="handleDelete(event)">
                      删除
                    </button>
                  </div>
                </div>
                <p class="mt-1 text-sm text-ink-500">{{ formatDateTime(event.startTime) }}</p>
                <p v-if="event.location" class="mt-0.5 text-sm text-ink-500">📍 {{ event.location }}</p>
                <p v-if="event.customNote" class="mt-1 rounded-lg bg-warm-50 px-3 py-2 text-xs text-ink-600">
                  备注：{{ event.customNote }}
                </p>
                <div class="mt-2 flex flex-wrap gap-2">
                  <span
                    v-for="rb in (Array.isArray(event.remindBefore) ? event.remindBefore : [event.remindBefore])"
                    :key="rb"
                    class="rounded-md bg-ink-100 px-2 py-0.5 text-xs text-ink-500"
                  >
                    {{ REMIND_OPTIONS.find(o => o.value === rb)?.label || rb }}
                  </span>
                  <span
                    v-if="event.reminderStatus"
                    class="rounded-md px-2 py-0.5 text-xs"
                    :class="REMINDER_STATUS_MAP[event.reminderStatus]?.class"
                  >
                    {{ REMINDER_STATUS_MAP[event.reminderStatus]?.label }}
                  </span>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </template>

    <!-- 月视图 -->
    <div v-else class="card-shadow rounded-2xl bg-white p-6">
      <!-- 月导航 -->
      <div class="mb-4 flex items-center justify-between">
        <button type="button" class="rounded-lg p-2 hover:bg-ink-50" @click="prevMonth">
          <svg class="h-5 w-5 text-ink-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7" />
          </svg>
        </button>
        <h3 class="text-base font-bold text-ink-900">
          {{ currentMonth.toLocaleDateString('zh-CN', { year: 'numeric', month: 'long' }) }}
        </h3>
        <button type="button" class="rounded-lg p-2 hover:bg-ink-50" @click="nextMonth">
          <svg class="h-5 w-5 text-ink-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7" />
          </svg>
        </button>
      </div>
      <!-- 星期标题 -->
      <div class="mb-2 grid grid-cols-7 text-center">
        <div v-for="dow in ['日','一','二','三','四','五','六']" :key="dow" class="py-2 text-xs font-medium text-ink-400">
          {{ dow }}
        </div>
      </div>
      <!-- 日格 -->
      <div class="grid grid-cols-7 gap-1">
        <div
          v-for="day in calendarDays"
          :key="day.date.toISOString()"
          class="min-h-[80px] rounded-xl p-1.5"
          :class="{
            'bg-ink-50': !isCurrentMonth(day.date),
            'ring-2 ring-brand-500': isToday(day.date),
          }"
        >
          <p class="mb-1 text-right text-xs font-medium"
            :class="isCurrentMonth(day.date) ? 'text-ink-700' : 'text-ink-300'"
          >
            {{ day.date.getDate() }}
          </p>
          <div class="space-y-0.5">
            <div
              v-for="ev in day.events.slice(0, 2)"
              :key="ev.id"
              class="cursor-pointer truncate rounded px-1.5 py-0.5 text-[10px] font-medium"
              :class="ev.eventType === 'career_talk' ? 'bg-brand-100 text-brand-700' : 'bg-blue-100 text-blue-700'"
              @click="openEdit(ev)"
            >
              {{ ev.title }}
            </div>
            <p v-if="day.events.length > 2" class="px-1 text-[10px] text-ink-400">
              +{{ day.events.length - 2 }} 更多
            </p>
          </div>
        </div>
      </div>
    </div>

    <!-- 编辑弹窗 -->
    <el-dialog v-model="editDialogVisible" title="编辑日程" width="480px" :close-on-click-modal="false">
      <div v-if="editTarget" class="space-y-4">
        <div>
          <p class="mb-1 text-sm font-medium text-ink-700">活动</p>
          <p class="text-sm text-ink-800">{{ editTarget.title }}</p>
        </div>
        <div>
          <p class="mb-1 text-sm font-medium text-ink-700">时间</p>
          <p class="text-sm text-ink-600">{{ formatDateTime(editTarget.startTime) }}</p>
        </div>
        <div>
          <p class="mb-2 text-sm font-medium text-ink-700">个人备注</p>
          <el-input
            v-model="editForm.customNote"
            type="textarea"
            :rows="3"
            placeholder="添加备注（可选）"
            maxlength="200"
            show-word-limit
          />
        </div>
        <div>
          <p class="mb-2 text-sm font-medium text-ink-700">提醒时间</p>
          <el-checkbox-group v-model="editForm.remindBefore">
            <el-checkbox v-for="opt in REMIND_OPTIONS" :key="opt.value" :label="opt.value">
              {{ opt.label }}
            </el-checkbox>
          </el-checkbox-group>
        </div>
      </div>
      <template #footer>
        <div class="flex gap-3">
          <button type="button" class="flex-1 rounded-xl border border-ink-200 py-2.5 text-sm font-medium text-ink-600 hover:bg-ink-50" @click="editDialogVisible = false">
            取消
          </button>
          <button type="button" class="btn-primary flex-1 rounded-xl py-2.5 text-sm font-medium text-white disabled:opacity-60" :disabled="editSaving" @click="handleSaveEdit">
            {{ editSaving ? '保存中…' : '保存' }}
          </button>
        </div>
      </template>
    </el-dialog>
  </div>
</template>