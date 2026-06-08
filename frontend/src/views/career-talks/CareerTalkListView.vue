<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import PageHeader from '@/components/business/PageHeader.vue'
import CareerTalkFilter from '@/components/business/CareerTalkFilter.vue'
import CareerTalkCard from '@/components/business/CareerTalkCard.vue'
import CareerTalkSidebar from '@/components/business/CareerTalkSidebar.vue'
import AppPagination from '@/components/common/AppPagination.vue'
import { getCareerTalkListApi } from '@/api/modules/careerTalk'
import { addCalendarEventApi } from '@/api/modules/calendar'
import { useAuthStore } from '@/stores'
import { usePagination } from '@/composables/usePagination'
import { normalizeCareerTalk } from '@/utils/careerTalk'
import type { CareerTalk, CareerTalkQuery } from '@/types'

const router = useRouter()
const authStore = useAuthStore()
const { page, pageSize, total, totalPages, setTotal } = usePagination()

const loading = ref(false)
const list = ref<CareerTalk[]>([])
const calendarLoading = ref<Record<number, boolean>>({})

const filters = reactive({
  keyword: '',
  dateRange: 'all',
  campus: 'all',
  industry: 'all',
  company: '',
  sortBy: 'time_asc',
})

async function fetchList() {
  loading.value = true
  try {
    const params: CareerTalkQuery = {
      page: page.value,
      pageSize: pageSize.value,
    }
    if (filters.keyword) params.keyword = filters.keyword
    if (filters.dateRange !== 'all') params.dateRange = filters.dateRange
    if (filters.campus !== 'all') params.campus = filters.campus
    if (filters.industry !== 'all') params.industry = filters.industry
    if (filters.company) params.company = filters.company
    if (filters.sortBy) params.sortBy = filters.sortBy

    const res = await getCareerTalkListApi(params)
    list.value = res.list.map(normalizeCareerTalk)
    setTotal(Number(res.total))
  } catch {
    // 错误由拦截器处理
  } finally {
    loading.value = false
  }
}

function handleSearch() {
  filters.company = ''
  page.value = 1
  fetchList()
}

function handleFilterByCompany(company: string) {
  filters.keyword = ''
  filters.company = company
  page.value = 1
  fetchList()
}

async function handleAddToCalendar(id: number) {
  if (!authStore.isLoggedIn) {
    router.push({ name: 'Login', query: { redirect: '/career-talks' } })
    return
  }
  const item = list.value.find((t) => t.id === id)
  if (!item || item.inCalendar) {
    if (item?.inCalendar) ElMessage.info('该活动已在您的日历中')
    return
  }
  calendarLoading.value[id] = true
  try {
    await addCalendarEventApi({ eventType: 'career_talk', refId: id })
    item.inCalendar = true
    ElMessage.success('已加入日历')
  } catch (e: any) {
    if (e?.response?.data?.code === 40901) {
      item.inCalendar = true
      ElMessage.info('该活动已在您的日历中')
    }
  } finally {
    calendarLoading.value[id] = false
  }
}

function handleToggleFavorite(id: number) {
  const item = list.value.find((t) => t.id === id)
  if (item) item.favorited = !item.favorited
}

function handleViewDetail(id: number) {
  router.push(`/career-talks/${id}`)
}

function handlePageChange(newPage: number) {
  page.value = newPage
  fetchList()
}

onMounted(fetchList)
</script>

<template>
  <div>
    <PageHeader
      title="宣讲会"
      description="发现优质校招机会，把握每一次与心仪企业面对面交流的机会"
      :count="total"
      count-label="场宣讲会"
    />

    <CareerTalkFilter
      v-model:keyword="filters.keyword"
      v-model:date-range="filters.dateRange"
      v-model:campus="filters.campus"
      v-model:industry="filters.industry"
      v-model:sort-by="filters.sortBy"
      @search="handleSearch"
    />

    <div class="flex gap-8">
      <div class="flex-1">
        <div v-if="loading" class="grid gap-5">
          <div v-for="i in 3" :key="i" class="h-40 animate-pulse rounded-2xl bg-ink-100" />
        </div>

        <div
          v-else-if="list.length === 0"
          class="flex flex-col items-center justify-center rounded-2xl border border-dashed border-ink-200 bg-white py-20"
        >
          <svg class="mb-4 h-12 w-12 text-ink-300" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M19 21V5a2 2 0 00-2-2H7a2 2 0 00-2 2v16m14 0h2m-2 0h-5m-9 0H3m2 0h5M9 7h6m-6 4h6m-6 4h6" />
          </svg>
          <p class="text-sm text-ink-500">暂无宣讲会信息</p>
        </div>

        <div v-else class="grid gap-5">
          <CareerTalkCard
            v-for="item in list"
            :key="item.id"
            :item="item"
            @add-to-calendar="handleAddToCalendar"
            @toggle-favorite="handleToggleFavorite"
            @view-detail="handleViewDetail"
          />
        </div>
        <AppPagination :page="page" :total-pages="totalPages" @change="handlePageChange" />
      </div>
      <CareerTalkSidebar @filter-company="handleFilterByCompany" />
    </div>
  </div>
</template>
