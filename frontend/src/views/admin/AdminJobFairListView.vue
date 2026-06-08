<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  getAdminJobFairsApi,
  createAdminJobFairApi,
  updateAdminJobFairApi,
  deleteAdminJobFairApi,
  batchAdminJobFairStatusApi,
} from '@/api/modules/admin/jobFair'
import { CAMPUS_OPTIONS } from '@/constants'
import type { AdminJobFairForm, AdminJobFairVO, PublishStatus } from '@/types/admin'

const loading = ref(false)
const list = ref<AdminJobFairVO[]>([])
const total = ref(0)
const page = ref(1)
const pageSize = ref(10)
const selectedIds = ref<number[]>([])

const filters = ref({
  keyword: '',
  campus: '',
  publishStatus: '' as PublishStatus | '',
  sourceType: '',
})

const dialogVisible = ref(false)
const editing = ref<AdminJobFairVO | null>(null)
const saving = ref(false)

const emptyForm = (): AdminJobFairForm => ({
  title: '',
  startDate: '',
  endDate: '',
  startTime: '',
  location: '',
  campus: 'main',
  companyCount: undefined,
  targetAudience: '',
  targetMajors: [],
  deadline: '',
  detailUrl: '',
  sourceUrl: '',
  description: '',
  publishStatus: 'draft',
})

const form = ref<AdminJobFairForm>(emptyForm())
const majorsText = ref('')

const PUBLISH_OPTIONS = [
  { label: '草稿', value: 'draft' },
  { label: '已发布', value: 'published' },
  { label: '已归档', value: 'archived' },
]

const STATUS_LABEL: Record<string, string> = {
  draft: '草稿',
  published: '已发布',
  archived: '已归档',
}

async function fetchList() {
  loading.value = true
  try {
    const res = await getAdminJobFairsApi({
      keyword: filters.value.keyword || undefined,
      campus: filters.value.campus || undefined,
      publishStatus: filters.value.publishStatus || undefined,
      sourceType: filters.value.sourceType || undefined,
      page: page.value,
      pageSize: pageSize.value,
    })
    list.value = res.list
    total.value = res.total
  } finally {
    loading.value = false
  }
}

function openCreate() {
  editing.value = null
  form.value = emptyForm()
  majorsText.value = ''
  dialogVisible.value = true
}

function openEdit(row: AdminJobFairVO) {
  editing.value = row
  form.value = {
    title: row.title,
    startDate: row.startDate,
    endDate: row.endDate,
    startTime: row.startTime?.slice(0, 16),
    location: row.location,
    campus: row.campus,
    companyCount: row.companyCount,
    targetAudience: row.targetAudience,
    targetMajors: row.targetMajors || [],
    deadline: row.deadline?.slice(0, 16),
    detailUrl: row.detailUrl,
    sourceUrl: row.sourceUrl,
    description: row.description,
    publishStatus: row.publishStatus,
  }
  majorsText.value = (row.targetMajors || []).join(', ')
  dialogVisible.value = true
}

function buildPayload(): AdminJobFairForm {
  return {
    ...form.value,
    targetMajors: majorsText.value.split(/[,，]/).map((s) => s.trim()).filter(Boolean),
    startTime: form.value.startTime?.replace('T', ' '),
    deadline: form.value.deadline?.replace('T', ' '),
  }
}

async function handleSave() {
  saving.value = true
  try {
    const payload = buildPayload()
    if (editing.value) {
      await updateAdminJobFairApi(editing.value.id, payload)
      ElMessage.success('更新成功')
    } else {
      await createAdminJobFairApi(payload)
      ElMessage.success('创建成功')
    }
    dialogVisible.value = false
    fetchList()
  } finally {
    saving.value = false
  }
}

async function handleDelete(row: AdminJobFairVO) {
  await ElMessageBox.confirm(`确定归档「${row.title}」？`, '确认删除', { type: 'warning' })
  await deleteAdminJobFairApi(row.id)
  ElMessage.success('已归档')
  fetchList()
}

async function handleBatchStatus(status: PublishStatus) {
  if (selectedIds.value.length === 0) {
    ElMessage.warning('请先选择记录')
    return
  }
  await batchAdminJobFairStatusApi(selectedIds.value, status)
  ElMessage.success('批量更新成功')
  selectedIds.value = []
  fetchList()
}

function handleSelectionChange(rows: AdminJobFairVO[]) {
  selectedIds.value = rows.map((r) => r.id)
}

onMounted(fetchList)
</script>

<template>
  <div>
    <div class="mb-4 flex flex-wrap items-center justify-between gap-3">
      <h2 class="text-lg font-bold text-ink-900">双选会管理</h2>
      <button type="button" class="btn-primary rounded-xl px-4 py-2 text-sm text-white" @click="openCreate">
        新建双选会
      </button>
    </div>

    <div class="card-shadow mb-4 flex flex-wrap gap-3 rounded-2xl bg-white p-4">
      <el-input v-model="filters.keyword" placeholder="搜索名称/地点" clearable class="!w-48" @keyup.enter="fetchList" />
      <el-select v-model="filters.campus" placeholder="校区" clearable class="!w-32">
        <el-option v-for="o in CAMPUS_OPTIONS.filter((c) => c.value !== 'all')" :key="o.value" :label="o.label" :value="o.value" />
      </el-select>
      <el-select v-model="filters.publishStatus" placeholder="发布状态" clearable class="!w-28">
        <el-option v-for="o in PUBLISH_OPTIONS" :key="o.value" :label="o.label" :value="o.value" />
      </el-select>
      <el-select v-model="filters.sourceType" placeholder="来源" clearable class="!w-28">
        <el-option label="手动" value="manual" />
        <el-option label="同步" value="sync" />
      </el-select>
      <button type="button" class="rounded-xl border border-ink-200 px-4 py-2 text-sm" @click="fetchList">搜索</button>
    </div>

    <div class="mb-3 flex gap-2">
      <button type="button" class="rounded-lg border border-brand-200 px-3 py-1.5 text-xs text-brand-600" @click="handleBatchStatus('published')">批量发布</button>
      <button type="button" class="rounded-lg border border-ink-200 px-3 py-1.5 text-xs text-ink-600" @click="handleBatchStatus('draft')">批量下架</button>
    </div>

    <div class="card-shadow rounded-2xl bg-white p-4">
      <el-table v-loading="loading" :data="list" @selection-change="handleSelectionChange">
        <el-table-column type="selection" width="40" />
        <el-table-column prop="title" label="名称" min-width="160" show-overflow-tooltip />
        <el-table-column label="日期" width="180">
          <template #default="{ row }">{{ row.startDate }}<template v-if="row.endDate"> ~ {{ row.endDate }}</template></template>
        </el-table-column>
        <el-table-column prop="location" label="地点" min-width="120" show-overflow-tooltip />
        <el-table-column prop="companyCount" label="企业数" width="80" />
        <el-table-column label="状态" width="80">
          <template #default="{ row }">
            <span class="rounded-md px-2 py-0.5 text-xs" :class="row.publishStatus === 'published' ? 'bg-brand-100 text-brand-700' : 'bg-ink-100 text-ink-500'">
              {{ STATUS_LABEL[row.publishStatus] }}
            </span>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="140" fixed="right">
          <template #default="{ row }">
            <button type="button" class="mr-2 text-xs text-brand-600" @click="openEdit(row)">编辑</button>
            <button type="button" class="text-xs text-red-500" @click="handleDelete(row)">删除</button>
          </template>
        </el-table-column>
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

    <el-dialog v-model="dialogVisible" :title="editing ? '编辑双选会' : '新建双选会'" width="640px" :close-on-click-modal="false">
      <div class="space-y-3">
        <el-input v-model="form.title" placeholder="名称 *" />
        <div class="grid grid-cols-2 gap-3">
          <el-date-picker v-model="form.startDate" type="date" placeholder="开始日期 *" value-format="YYYY-MM-DD" class="!w-full" />
          <el-date-picker v-model="form.endDate" type="date" placeholder="结束日期" value-format="YYYY-MM-DD" class="!w-full" />
        </div>
        <el-input v-model="form.location" placeholder="举办地点 *" />
        <div class="grid grid-cols-3 gap-3">
          <el-select v-model="form.campus" placeholder="校区">
            <el-option v-for="o in CAMPUS_OPTIONS.filter((c) => c.value !== 'all')" :key="o.value" :label="o.label" :value="o.value" />
          </el-select>
          <el-input-number v-model="form.companyCount" :min="0" placeholder="企业数" class="!w-full" />
          <el-select v-model="form.publishStatus" placeholder="发布状态">
            <el-option v-for="o in PUBLISH_OPTIONS" :key="o.value" :label="o.label" :value="o.value" />
          </el-select>
        </div>
        <el-input v-model="form.targetAudience" placeholder="面向对象" />
        <el-input v-model="majorsText" placeholder="面向专业（逗号分隔）" />
        <el-date-picker v-model="form.deadline" type="datetime" placeholder="报名截止" value-format="YYYY-MM-DDTHH:mm:ss" class="!w-full" />
        <el-input v-model="form.description" type="textarea" :rows="3" placeholder="详细描述" />
      </div>
      <template #footer>
        <button type="button" class="rounded-xl border px-4 py-2 text-sm" @click="dialogVisible = false">取消</button>
        <button type="button" class="btn-primary ml-2 rounded-xl px-4 py-2 text-sm text-white" :disabled="saving" @click="handleSave">
          {{ saving ? '保存中…' : '保存' }}
        </button>
      </template>
    </el-dialog>
  </div>
</template>
