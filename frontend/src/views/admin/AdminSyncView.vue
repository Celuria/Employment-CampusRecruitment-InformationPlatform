<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { triggerSyncApi, getSyncLogsApi } from '@/api/modules/admin/sync'
import type { SyncLogVO } from '@/types/admin'

const syncing = ref(false)
const sourceType = ref('all')
const force = ref(false)

const loading = ref(false)
const list = ref<SyncLogVO[]>([])
const total = ref(0)
const page = ref(1)
const pageSize = ref(10)

const STATUS_LABEL: Record<string, string> = {
  pending: '等待中',
  running: '运行中',
  success: '成功',
  failed: '失败',
}

async function fetchLogs() {
  loading.value = true
  try {
    const res = await getSyncLogsApi(page.value, pageSize.value)
    list.value = res.list
    total.value = res.total
  } finally {
    loading.value = false
  }
}

async function handleSync() {
  syncing.value = true
  try {
    const res = await triggerSyncApi(sourceType.value, force.value)
    ElMessage.success(res.message || '同步任务已提交')
    fetchLogs()
  } finally {
    syncing.value = false
  }
}

function formatTime(s?: string) {
  if (!s) return '-'
  return new Date(s).toLocaleString('zh-CN')
}

onMounted(fetchLogs)
</script>

<template>
  <div>
    <h2 class="mb-4 text-lg font-bold text-ink-900">信息同步</h2>

    <div class="card-shadow mb-6 rounded-2xl bg-white p-6">
      <h3 class="mb-4 text-sm font-semibold text-ink-700">触发同步</h3>
      <div class="flex flex-wrap items-end gap-4">
        <div>
          <p class="mb-1 text-xs text-ink-500">同步类型</p>
          <el-select v-model="sourceType" class="!w-40">
            <el-option label="全部" value="all" />
            <el-option label="宣讲会" value="career_talk" />
            <el-option label="双选会" value="job_fair" />
          </el-select>
        </div>
        <el-checkbox v-model="force">强制全量同步</el-checkbox>
        <button
          type="button"
          class="btn-primary rounded-xl px-5 py-2 text-sm text-white disabled:opacity-60"
          :disabled="syncing"
          @click="handleSync"
        >
          {{ syncing ? '同步中…' : '开始同步' }}
        </button>
      </div>
      <p class="mt-3 text-xs text-ink-400">当前为轻量同步实现：刷新来源为 sync 且含 source_url 的条目 synced_at 时间戳。</p>
    </div>

    <div class="card-shadow rounded-2xl bg-white p-4">
      <h3 class="mb-4 text-sm font-semibold text-ink-700">同步记录</h3>
      <el-table v-loading="loading" :data="list">
        <el-table-column prop="taskId" label="任务ID" min-width="160" show-overflow-tooltip />
        <el-table-column prop="sourceType" label="类型" width="100" />
        <el-table-column label="状态" width="80">
          <template #default="{ row }">
            <span :class="row.status === 'success' ? 'text-green-600' : row.status === 'failed' ? 'text-red-500' : ''">
              {{ STATUS_LABEL[row.status] || row.status }}
            </span>
          </template>
        </el-table-column>
        <el-table-column prop="updatedCount" label="更新数" width="80" />
        <el-table-column label="开始时间" width="160">
          <template #default="{ row }">{{ formatTime(row.startedAt) }}</template>
        </el-table-column>
        <el-table-column label="结束时间" width="160">
          <template #default="{ row }">{{ formatTime(row.finishedAt) }}</template>
        </el-table-column>
        <el-table-column prop="errorMessage" label="错误信息" min-width="120" show-overflow-tooltip />
      </el-table>
      <div class="mt-4 flex justify-end">
        <el-pagination
          v-model:current-page="page"
          v-model:page-size="pageSize"
          :total="total"
          layout="total, prev, pager, next"
          @current-change="fetchLogs"
          @size-change="fetchLogs"
        />
      </div>
    </div>
  </div>
</template>
