<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { getAuditLogsApi } from '@/api/modules/admin/auditLog'
import type { AuditLogVO } from '@/types/admin'

const loading = ref(false)
const list = ref<AuditLogVO[]>([])
const total = ref(0)
const page = ref(1)
const pageSize = ref(20)

const filters = ref({
  action: '',
  resourceType: '',
  startDate: '',
  endDate: '',
})

const ACTION_OPTIONS = ['CREATE', 'UPDATE', 'DELETE', 'SYNC']
const RESOURCE_OPTIONS = ['career_talk', 'job_fair', 'user', 'sync']

async function fetchList() {
  loading.value = true
  try {
    const res = await getAuditLogsApi({
      action: filters.value.action || undefined,
      resourceType: filters.value.resourceType || undefined,
      startDate: filters.value.startDate || undefined,
      endDate: filters.value.endDate || undefined,
      page: page.value,
      pageSize: pageSize.value,
    })
    list.value = res.list
    total.value = res.total
  } finally {
    loading.value = false
  }
}

function formatTime(s: string) {
  return new Date(s).toLocaleString('zh-CN')
}

onMounted(fetchList)
</script>

<template>
  <div>
    <h2 class="mb-4 text-lg font-bold text-ink-900">审计日志</h2>

    <div class="card-shadow mb-4 flex flex-wrap gap-3 rounded-2xl bg-white p-4">
      <el-select v-model="filters.action" placeholder="操作类型" clearable class="!w-32">
        <el-option v-for="a in ACTION_OPTIONS" :key="a" :label="a" :value="a" />
      </el-select>
      <el-select v-model="filters.resourceType" placeholder="资源类型" clearable class="!w-36">
        <el-option v-for="r in RESOURCE_OPTIONS" :key="r" :label="r" :value="r" />
      </el-select>
      <el-date-picker v-model="filters.startDate" type="date" placeholder="起始日期" value-format="YYYY-MM-DD" class="!w-36" />
      <el-date-picker v-model="filters.endDate" type="date" placeholder="结束日期" value-format="YYYY-MM-DD" class="!w-36" />
      <button type="button" class="rounded-xl border border-ink-200 px-4 py-2 text-sm" @click="fetchList">筛选</button>
    </div>

    <div class="card-shadow rounded-2xl bg-white p-4">
      <el-table v-loading="loading" :data="list">
        <el-table-column label="时间" width="160">
          <template #default="{ row }">{{ formatTime(row.createdAt) }}</template>
        </el-table-column>
        <el-table-column prop="operatorName" label="操作人" width="100" />
        <el-table-column prop="action" label="动作" width="80" />
        <el-table-column prop="resourceType" label="资源类型" width="100" />
        <el-table-column prop="resourceId" label="资源ID" width="80" />
        <el-table-column prop="ip" label="IP" width="120" />
        <el-table-column prop="detail" label="详情" min-width="200" show-overflow-tooltip />
      </el-table>
      <div class="mt-4 flex justify-end">
        <el-pagination
          v-model:current-page="page"
          v-model:page-size="pageSize"
          :total="total"
          layout="total, prev, pager, next"
          @current-change="fetchList"
          @size-change="fetchList"
        />
      </div>
    </div>
  </div>
</template>
