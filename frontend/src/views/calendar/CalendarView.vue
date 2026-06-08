<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import PageHeader from '@/components/business/PageHeader.vue'
import CalendarDayPanel from '@/components/business/CalendarDayPanel.vue'
import { getCalendarEventsApi, deleteCalendarEventApi, updateCalendarEventApi } from '@/api/modules/calendar'
import {
  buildCalendarDays,
  formatEventChipTime,
  getEventsForDay,
  getVisibleGridRange,
} from '@/utils/calendar'
import type { CalendarEvent } from '@/types'

const router = useRouter()

const loading = ref(true)
const events = ref<CalendarEvent[]>([])
const viewMode = ref<'list' | 'month'>('list')
const currentMonth = ref(new Date())
const selectedDate = ref<Date | null>(null)

const editDialogVisible = ref(false)
const editTarget = ref<CalendarEvent | null>(null)
const editForm = ref({ customNote: '', remindBefore: [] as string[] })
const editSaving = ref(false)

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

const calendarDays = computed(() => buildCalendarDays(currentMonth.value, events.value))

const selectedDayEvents = computed(() => {
  if (!selectedDate.value) return []
  return getEventsForDay(selectedDate.value, events.value)
})

const monthHasEvents = computed(() =>
  calendarDays.value.some((d) => isCurrentMonth(d.date) && d.events.length > 0),
)

async function fetchEvents() {
  loading.value = true
  try {
    const params = viewMode.value === 'month'
      ? { ...getVisibleGridRange(currentMonth.value), view: 'month' as const }
      : { view: 'list' as const }
    const res = await getCalendarEventsApi(params)
    events.value = Array.isArray(res) ? res : []
  } catch {
    events.value = []
  } finally {
    loading.value = false
  }
}

function handleViewDetail(event: CalendarEvent) {
  const path = event.eventType === 'career_talk'
    ? `/career-talks/${event.refId}`
    : `/job-fairs/${event.refId}`
  router.push(path)
}

function switchViewMode(mode: 'list' | 'month') {
  if (viewMode.value === mode) return
  viewMode.value = mode
  selectedDate.value = null
  fetchEvents()
}

function prevMonth() {
  currentMonth.value = new Date(currentMonth.value.getFullYear(), currentMonth.value.getMonth() - 1, 1)
  selectedDate.value = null
  if (viewMode.value === 'month') fetchEvents()
}

function nextMonth() {
  currentMonth.value = new Date(currentMonth.value.getFullYear(), currentMonth.value.getMonth() + 1, 1)
  selectedDate.value = null
  if (viewMode.value === 'month') fetchEvents()
}

function goToToday() {
  currentMonth.value = new Date()
  selectedDate.value = new Date()
  if (viewMode.value === 'month') fetchEvents()
}

function isCurrentMonth(date: Date) {
  return date.getMonth() === currentMonth.value.getMonth()
    && date.getFullYear() === currentMonth.value.getFullYear()
}

function isToday(date: Date) {
  const t = new Date()
  return date.getDate() === t.getDate()
    && date.getMonth() === t.getMonth()
    && date.getFullYear() === t.getFullYear()
}

function isSelectedDay(date: Date) {
  if (!selectedDate.value) return false
  return date.getDate() === selectedDate.value.getDate()
    && date.getMonth() === selectedDate.value.getMonth()
    && date.getFullYear() === selectedDate.value.getFullYear()
}

function selectDay(date: Date) {
  selectedDate.value = new Date(date.getFullYear(), date.getMonth(), date.getDate())
}

function eventTooltip(event: CalendarEvent) {
  const time = formatDateTime(event.startTime)
  return event.location ? `${event.title}\n${time}\n${event.location}` : `${event.title}\n${time}`
}

const groupedEvents = computed(() => {
  const groups: { label: string; events: CalendarEvent[] }[] = []
  const now = new Date()
  const upcoming = events.value.filter((e) => new Date(e.startTime) >= now)
  const past = events.value.filter((e) => new Date(e.startTime) < now)
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

function openEdit(event: CalendarEvent) {
  editTarget.value = event
  editForm.value = {
    customNote: event.customNote || '',
    remindBefore: [...(event.remindBefore ?? [])],
  }
  editDialogVisible.value = true
}

async function handleSaveEdit() {
  if (!editTarget.value) return
  editSaving.value = true
  try {
    const updated = await updateCalendarEventApi(editTarget.value.id, {
      customNote: editForm.value.customNote,
      remindBefore: editForm.value.remindBefore,
    })
    const idx = events.value.findIndex((e) => e.id === editTarget.value!.id)
    if (idx !== -1) events.value[idx] = { ...events.value[idx], ...updated }
    editDialogVisible.value = false
    ElMessage.success('日程已更新')
  } catch {
    // 错误由拦截器处理
  } finally {
    editSaving.value = false
  }
}

async function handleDelete(event: CalendarEvent) {
  try {
    await ElMessageBox.confirm(
      `确认删除「${event.title}」？删除后将取消未发送的提醒。`,
      '删除确认',
      { confirmButtonText: '确认删除', cancelButtonText: '取消', type: 'warning' },
    )
  } catch {
    return
  }
  deleteTarget.value = event
  deleteLoading.value = true
  try {
    await deleteCalendarEventApi(event.id)
    events.value = events.value.filter((e) => e.id !== event.id)
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
          @click="switchViewMode('list')"
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
          @click="switchViewMode('month')"
        >
          <svg class="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z" />
          </svg>
          月视图
        </button>
      </div>
      <p class="text-sm text-ink-500">共 {{ events.length }} 个日程</p>
    </div>

    <!-- 列表视图：加载中 -->
    <div v-if="loading && viewMode === 'list'" class="space-y-4">
      <div v-for="i in 3" :key="i" class="h-24 animate-pulse rounded-2xl bg-ink-100" />
    </div>

    <!-- 列表视图：空状态 -->
    <div
      v-else-if="viewMode === 'list' && events.length === 0"
      class="flex flex-col items-center justify-center rounded-2xl border border-dashed border-ink-200 bg-white py-20"
    >
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
            class="card-shadow group cursor-pointer rounded-2xl bg-white p-5 transition-all hover:shadow-md"
            @click="handleViewDetail(event)"
          >
            <div class="flex items-start gap-4">
              <div
                class="shrink-0 rounded-xl px-3 py-2 text-center text-xs font-bold"
                :class="event.eventType === 'career_talk' ? 'bg-brand-100 text-brand-700' : 'bg-blue-100 text-blue-700'"
              >
                {{ EVENT_TYPE_LABEL[event.eventType] }}
              </div>
              <div class="min-w-0 flex-1">
                <div class="flex items-start justify-between gap-3">
                  <h3 class="font-semibold text-ink-900 transition-colors group-hover:text-brand-600">{{ event.title }}</h3>
                  <div class="flex shrink-0 gap-2" @click.stop>
                    <button
                      type="button"
                      class="rounded-lg border border-ink-200 px-3 py-1.5 text-xs text-ink-600 transition-colors hover:border-brand-300 hover:text-brand-600"
                      @click="openEdit(event)"
                    >
                      编辑
                    </button>
                    <button
                      type="button"
                      class="rounded-lg border border-red-200 px-3 py-1.5 text-xs text-red-500 transition-colors hover:bg-red-50 disabled:opacity-50"
                      :disabled="deleteTarget?.id === event.id && deleteLoading"
                      @click="handleDelete(event)"
                    >
                      删除
                    </button>
                  </div>
                </div>
                <p class="mt-1 text-sm text-ink-500">{{ formatDateTime(event.startTime) }}</p>
                <p v-if="event.location" class="mt-0.5 text-sm text-ink-500">{{ event.location }}</p>
                <p v-if="event.customNote" class="mt-1 rounded-lg bg-warm-50 px-3 py-2 text-xs text-ink-600">
                  备注：{{ event.customNote }}
                </p>
                <div class="mt-2 flex flex-wrap gap-2">
                  <span
                    v-for="rb in event.remindBefore"
                    :key="rb"
                    class="rounded-md bg-ink-100 px-2 py-0.5 text-xs text-ink-500"
                  >
                    {{ REMIND_OPTIONS.find((o) => o.value === rb)?.label || rb }}
                  </span>
                  <span v-if="event.sourceUpdated" class="rounded-md bg-warm-50 px-2 py-0.5 text-xs text-warm-300">
                    活动信息已更新
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
    <div v-else class="flex gap-6">
      <div class="min-w-0 flex-1">
        <div class="card-shadow rounded-2xl bg-white p-6">
          <!-- 月导航 -->
          <div class="mb-4 flex flex-wrap items-center justify-between gap-3">
            <div class="flex items-center gap-2">
              <button type="button" class="rounded-lg p-2 hover:bg-ink-50" @click="prevMonth">
                <svg class="h-5 w-5 text-ink-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7" />
                </svg>
              </button>
              <h3 class="min-w-[8rem] text-center text-base font-bold text-ink-900">
                {{ currentMonth.toLocaleDateString('zh-CN', { year: 'numeric', month: 'long' }) }}
              </h3>
              <button type="button" class="rounded-lg p-2 hover:bg-ink-50" @click="nextMonth">
                <svg class="h-5 w-5 text-ink-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7" />
                </svg>
              </button>
              <button
                type="button"
                class="rounded-lg border border-ink-200 px-3 py-1.5 text-xs font-medium text-ink-600 transition-colors hover:border-brand-300 hover:text-brand-600"
                @click="goToToday"
              >
                今天
              </button>
            </div>
            <div class="flex items-center gap-3 text-xs text-ink-500">
              <span class="flex items-center gap-1.5">
                <span class="h-2.5 w-2.5 rounded-sm bg-brand-400" />宣讲会
              </span>
              <span class="flex items-center gap-1.5">
                <span class="h-2.5 w-2.5 rounded-sm bg-blue-400" />双选会
              </span>
            </div>
          </div>

          <p v-if="!loading && !monthHasEvents" class="mb-3 text-center text-xs text-ink-400">
            本月暂无日程
          </p>

          <!-- 星期标题 -->
          <div class="mb-2 grid grid-cols-7 text-center">
            <div v-for="dow in ['日','一','二','三','四','五','六']" :key="dow" class="py-2 text-xs font-medium text-ink-400">
              {{ dow }}
            </div>
          </div>

          <!-- 加载骨架 -->
          <div v-if="loading" class="grid grid-cols-7 gap-1">
            <div v-for="i in 42" :key="i" class="min-h-[88px] animate-pulse rounded-xl bg-ink-100" />
          </div>

          <!-- 日格 -->
          <div v-else class="grid grid-cols-7 gap-1">
            <div
              v-for="day in calendarDays"
              :key="day.date.toISOString()"
              class="min-h-[88px] cursor-pointer rounded-xl p-1.5 transition-colors hover:bg-ink-50"
              :class="{
                'bg-ink-50/60': !isCurrentMonth(day.date),
                'ring-2 ring-brand-500': isToday(day.date),
                'bg-brand-50 ring-1 ring-brand-200': isSelectedDay(day.date) && !isToday(day.date),
              }"
              @click="selectDay(day.date)"
            >
              <p
                class="mb-1 text-right text-xs font-medium"
                :class="isCurrentMonth(day.date) ? 'text-ink-700' : 'text-ink-300'"
              >
                {{ day.date.getDate() }}
              </p>
              <div class="space-y-0.5" @click.stop>
                <div
                  v-for="ev in day.events.slice(0, 2)"
                  :key="ev.id"
                  class="cursor-pointer truncate rounded px-1.5 py-0.5 text-[10px] font-medium leading-tight"
                  :class="ev.eventType === 'career_talk' ? 'bg-brand-100 text-brand-700 hover:bg-brand-200' : 'bg-blue-100 text-blue-700 hover:bg-blue-200'"
                  :title="eventTooltip(ev)"
                  @click="handleViewDetail(ev)"
                >
                  <span v-if="formatEventChipTime(ev.startTime)" class="opacity-80">
                    {{ formatEventChipTime(ev.startTime) }}
                  </span>
                  {{ ev.title }}
                </div>
                <button
                  v-if="day.events.length > 2"
                  type="button"
                  class="w-full rounded px-1.5 py-0.5 text-left text-[10px] text-ink-400 hover:bg-ink-100 hover:text-brand-600"
                  @click="selectDay(day.date)"
                >
                  +{{ day.events.length - 2 }} 更多
                </button>
              </div>
            </div>
          </div>
        </div>
      </div>

      <CalendarDayPanel
        v-if="selectedDate"
        :date="selectedDate"
        :events="selectedDayEvents"
        :delete-loading-id="deleteLoading && deleteTarget ? deleteTarget.id : null"
        @close="selectedDate = null"
        @view-detail="handleViewDetail"
        @edit="openEdit"
        @delete="handleDelete"
      />
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
          <button
            type="button"
            class="flex-1 rounded-xl border border-ink-200 py-2.5 text-sm font-medium text-ink-600 hover:bg-ink-50"
            @click="editDialogVisible = false"
          >
            取消
          </button>
          <button
            type="button"
            class="btn-primary flex-1 rounded-xl py-2.5 text-sm font-medium text-white disabled:opacity-60"
            :disabled="editSaving"
            @click="handleSaveEdit"
          >
            {{ editSaving ? '保存中…' : '保存' }}
          </button>
        </div>
      </template>
    </el-dialog>
  </div>
</template>
