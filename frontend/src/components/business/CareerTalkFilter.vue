<script setup lang="ts">
import { CAMPUS_OPTIONS, DATE_FILTER_OPTIONS, INDUSTRY_OPTIONS, SORT_OPTIONS } from '@/constants'

const keyword = defineModel<string>('keyword', { default: '' })
const dateRange = defineModel<string>('dateRange', { default: 'all' })
const campus = defineModel<string>('campus', { default: 'all' })
const industry = defineModel<string>('industry', { default: 'all' })
const sortBy = defineModel<string>('sortBy', { default: 'time_asc' })

const emit = defineEmits<{
  search: []
}>()

function setDateFilter(value: string) {
  dateRange.value = value
  emit('search')
}

function setIndustry(value: string) {
  industry.value = value
  emit('search')
}

function handleSearch() {
  emit('search')
}
</script>

<template>
  <div class="card-shadow mb-8 rounded-2xl bg-white p-6">
    <div class="mb-5 flex gap-3">
      <div class="relative flex-1">
        <svg
          class="absolute left-4 top-1/2 h-5 w-5 -translate-y-1/2 text-ink-400"
          fill="none"
          stroke="currentColor"
          viewBox="0 0 24 24"
        >
          <path
            stroke-linecap="round"
            stroke-linejoin="round"
            stroke-width="2"
            d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"
          />
        </svg>
        <input
          v-model="keyword"
          type="text"
          placeholder="搜索公司、岗位、关键词..."
          class="w-full rounded-xl border border-ink-200 bg-ink-50 py-3 pl-12 pr-4 text-sm transition-all focus:border-brand-400 focus:outline-none focus:ring-2 focus:ring-brand-300"
          @keyup.enter="handleSearch"
        />
      </div>
      <button
        type="button"
        class="btn-primary rounded-xl px-8 py-3 text-sm font-medium text-white transition-all hover:shadow-lg"
        @click="handleSearch"
      >
        搜索
      </button>
    </div>

    <div class="flex flex-wrap items-center gap-6 text-sm">
      <div class="flex items-center gap-2">
        <span class="font-medium text-ink-500">日期：</span>
        <div class="flex gap-1.5">
          <button
            v-for="opt in DATE_FILTER_OPTIONS"
            :key="opt.value"
            type="button"
            class="filter-pill rounded-lg px-3 py-1.5 text-xs font-medium"
            :class="
              dateRange === opt.value
                ? 'active'
                : 'bg-ink-50 text-ink-600'
            "
            @click="setDateFilter(opt.value)"
          >
            {{ opt.label }}
          </button>
        </div>
      </div>

      <div class="flex items-center gap-2">
        <span class="font-medium text-ink-500">地点：</span>
        <select
          v-model="campus"
          class="cursor-pointer rounded-lg border-none bg-ink-50 px-3 py-1.5 text-xs font-medium text-ink-600 focus:ring-2 focus:ring-brand-300"
          @change="handleSearch"
        >
          <option v-for="opt in CAMPUS_OPTIONS" :key="opt.value" :value="opt.value">
            {{ opt.label }}
          </option>
        </select>
      </div>

      <div class="flex items-center gap-2">
        <span class="font-medium text-ink-500">行业：</span>
        <div class="flex gap-1.5">
          <button
            v-for="opt in INDUSTRY_OPTIONS"
            :key="opt.value"
            type="button"
            class="filter-pill rounded-lg px-3 py-1.5 text-xs font-medium"
            :class="
              industry === opt.value
                ? 'active'
                : 'bg-ink-50 text-ink-600'
            "
            @click="setIndustry(opt.value)"
          >
            {{ opt.label }}
          </button>
        </div>
      </div>

      <div class="ml-auto flex items-center gap-2">
        <span class="font-medium text-ink-500">排序：</span>
        <select
          v-model="sortBy"
          class="cursor-pointer rounded-lg border-none bg-ink-50 px-3 py-1.5 text-xs font-medium text-ink-600 focus:ring-2 focus:ring-brand-300"
          @change="handleSearch"
        >
          <option v-for="opt in SORT_OPTIONS" :key="opt.value" :value="opt.value">
            {{ opt.label }}
          </option>
        </select>
      </div>
    </div>
  </div>
</template>
