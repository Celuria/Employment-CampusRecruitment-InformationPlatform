<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import PageHeader from '@/components/business/PageHeader.vue'
import CareerTalkFilter from '@/components/business/CareerTalkFilter.vue'
import CareerTalkCard from '@/components/business/CareerTalkCard.vue'
import CareerTalkSidebar from '@/components/business/CareerTalkSidebar.vue'
import AppPagination from '@/components/common/AppPagination.vue'
import { usePagination } from '@/composables/usePagination'
import type { CareerTalk } from '@/types'

const router = useRouter()
const { page } = usePagination()

/** 静态 Mock 数据，后续替换为 API 调用 */
const mockList = ref<CareerTalk[]>([
  {
    id: 1,
    title: '字节跳动2025校园招聘技术专场宣讲会',
    company: '字节跳动',
    industry: '互联网',
    companySize: '10000人以上',
    startTime: '2025年1月15日 周三 14:00',
    location: '本部校区 · 学术报告厅A301',
    format: 'hybrid',
    positions: ['研发工程师', '产品经理', '算法'],
    status: 'upcoming',
  },
  {
    id: 2,
    title: '阿里巴巴2025校招「星耀计划」全球宣讲会',
    company: '阿里巴巴集团',
    industry: '互联网',
    companySize: '10000人以上',
    startTime: '2025年1月16日 周四 19:00',
    location: '沙河校区 · 教学楼B205',
    format: 'offline',
    positions: ['Java开发', '数据分析师', '运营'],
    favorited: true,
    status: 'upcoming',
  },
  {
    id: 3,
    title: '腾讯2025校园招聘「技术大咖」面对面',
    company: '腾讯科技',
    industry: '互联网',
    companySize: '10000人以上',
    startTime: '2025年1月17日 周五 15:30',
    location: '本部校区 · 图书馆报告厅',
    format: 'online',
    positions: ['后端开发', 'AI工程师', '游戏开发'],
    inCalendar: true,
    status: 'upcoming',
  },
  {
    id: 4,
    title: '西门子2025「数字化工业」校园宣讲',
    company: '西门子中国',
    industry: '制造/自动化',
    companySize: '10000人以上',
    startTime: '2025年1月10日 周五 14:00',
    location: '本部校区 · 机械楼报告厅',
    format: 'offline',
    positions: ['自动化', '工业软件'],
    status: 'ended',
  },
])

function handleSearch() {
  // TODO: 接入 getCareerTalkListApi
}

function handleAddToCalendar(_id: number) {
  // TODO: 接入 addCalendarEventApi，需登录校验
}

function handleToggleFavorite(_id: number) {
  // TODO: 收藏功能（可选）
}

function handleViewDetail(id: number) {
  router.push(`/career-talks/${id}`)
}
</script>

<template>
  <div>
    <PageHeader
      title="宣讲会"
      description="发现优质校招机会，把握每一次与心仪企业面对面交流的机会"
      :count="128"
      count-label="场宣讲会"
    />

    <CareerTalkFilter @search="handleSearch" />

    <div class="flex gap-8">
      <div class="flex-1">
        <div class="grid gap-5">
          <CareerTalkCard
            v-for="item in mockList"
            :key="item.id"
            :item="item"
            @add-to-calendar="handleAddToCalendar"
            @toggle-favorite="handleToggleFavorite"
            @view-detail="handleViewDetail"
          />
        </div>
        <AppPagination :page="page" :total-pages="13" @change="page = $event" />
      </div>
      <CareerTalkSidebar />
    </div>
  </div>
</template>
