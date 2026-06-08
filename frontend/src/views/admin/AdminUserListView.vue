<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  getAdminUsersApi,
  createAdminUserApi,
  updateAdminUserApi,
  updateAdminUserStatusApi,
  resetAdminUserPasswordApi,
} from '@/api/modules/admin/user'
import type { UserRole } from '@/types/auth'
import type { AdminUserCreateForm, AdminUserUpdateForm, AdminUserVO, UserStatus } from '@/types/admin'

const loading = ref(false)
const list = ref<AdminUserVO[]>([])
const total = ref(0)
const page = ref(1)
const pageSize = ref(10)

const filters = ref({
  keyword: '',
  role: '' as UserRole | '',
  status: '' as UserStatus | '',
})

const createDialogVisible = ref(false)
const editDialogVisible = ref(false)
const resetDialogVisible = ref(false)
const saving = ref(false)
const editing = ref<AdminUserVO | null>(null)

const createForm = ref<AdminUserCreateForm>({
  username: '',
  password: '',
  name: '',
  email: '',
  role: 'student',
})

const editForm = ref<AdminUserUpdateForm>({
  name: '',
  email: '',
  college: '',
  major: '',
  role: 'student',
})

const newPassword = ref('')

const ROLE_LABEL: Record<string, string> = { student: '学生', admin: '管理员' }
const STATUS_LABEL: Record<string, string> = { active: '正常', locked: '锁定', disabled: '禁用' }

async function fetchList() {
  loading.value = true
  try {
    const res = await getAdminUsersApi({
      keyword: filters.value.keyword || undefined,
      role: filters.value.role || undefined,
      status: filters.value.status || undefined,
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
  createForm.value = { username: '', password: '', name: '', email: '', role: 'student' }
  createDialogVisible.value = true
}

function openEdit(row: AdminUserVO) {
  editing.value = row
  editForm.value = {
    name: row.name,
    email: row.email,
    college: row.college,
    major: row.major,
    role: row.role,
  }
  editDialogVisible.value = true
}

function openReset(row: AdminUserVO) {
  editing.value = row
  newPassword.value = ''
  resetDialogVisible.value = true
}

async function handleCreate() {
  saving.value = true
  try {
    await createAdminUserApi(createForm.value)
    ElMessage.success('用户创建成功')
    createDialogVisible.value = false
    fetchList()
  } finally {
    saving.value = false
  }
}

async function handleEdit() {
  if (!editing.value) return
  saving.value = true
  try {
    await updateAdminUserApi(editing.value.id, editForm.value)
    ElMessage.success('更新成功')
    editDialogVisible.value = false
    fetchList()
  } finally {
    saving.value = false
  }
}

async function handleReset() {
  if (!editing.value) return
  saving.value = true
  try {
    await resetAdminUserPasswordApi(editing.value.id, newPassword.value)
    ElMessage.success('密码已重置')
    resetDialogVisible.value = false
  } finally {
    saving.value = false
  }
}

async function handleDisable(row: AdminUserVO) {
  await ElMessageBox.confirm(`确定禁用用户「${row.username}」？`, '确认禁用', { type: 'warning' })
  await updateAdminUserStatusApi(row.id, 'disabled')
  ElMessage.success('已禁用')
  fetchList()
}

function formatTime(s?: string) {
  if (!s) return '-'
  return new Date(s).toLocaleString('zh-CN')
}

onMounted(fetchList)
</script>

<template>
  <div>
    <div class="mb-4 flex flex-wrap items-center justify-between gap-3">
      <h2 class="text-lg font-bold text-ink-900">用户管理</h2>
      <button type="button" class="btn-primary rounded-xl px-4 py-2 text-sm text-white" @click="openCreate">
        创建用户
      </button>
    </div>

    <div class="card-shadow mb-4 flex flex-wrap gap-3 rounded-2xl bg-white p-4">
      <el-input v-model="filters.keyword" placeholder="搜索账号/姓名/邮箱" clearable class="!w-52" @keyup.enter="fetchList" />
      <el-select v-model="filters.role" placeholder="角色" clearable class="!w-28">
        <el-option label="学生" value="student" />
        <el-option label="管理员" value="admin" />
      </el-select>
      <el-select v-model="filters.status" placeholder="状态" clearable class="!w-28">
        <el-option label="正常" value="active" />
        <el-option label="锁定" value="locked" />
        <el-option label="禁用" value="disabled" />
      </el-select>
      <button type="button" class="rounded-xl border border-ink-200 px-4 py-2 text-sm" @click="fetchList">搜索</button>
    </div>

    <div class="card-shadow rounded-2xl bg-white p-4">
      <el-table v-loading="loading" :data="list">
        <el-table-column prop="username" label="账号" width="120" />
        <el-table-column prop="name" label="姓名" width="100" />
        <el-table-column prop="email" label="邮箱" min-width="160" show-overflow-tooltip />
        <el-table-column prop="college" label="学院" width="120" show-overflow-tooltip />
        <el-table-column label="角色" width="90">
          <template #default="{ row }">
            <span :class="row.role === 'admin' ? 'text-brand-600 font-medium' : ''">{{ ROLE_LABEL[row.role] }}</span>
          </template>
        </el-table-column>
        <el-table-column label="状态" width="80">
          <template #default="{ row }">
            <span :class="row.status === 'active' ? 'text-green-600' : 'text-ink-400'">{{ STATUS_LABEL[row.status] }}</span>
          </template>
        </el-table-column>
        <el-table-column label="最后登录" width="160">
          <template #default="{ row }">{{ formatTime(row.lastLoginAt) }}</template>
        </el-table-column>
        <el-table-column label="操作" width="200" fixed="right">
          <template #default="{ row }">
            <button type="button" class="mr-2 text-xs text-brand-600" @click="openEdit(row)">编辑</button>
            <button type="button" class="mr-2 text-xs text-ink-600" @click="openReset(row)">重置密码</button>
            <button v-if="row.status !== 'disabled'" type="button" class="text-xs text-red-500" @click="handleDisable(row)">禁用</button>
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

    <el-dialog v-model="createDialogVisible" title="创建用户" width="480px" :close-on-click-modal="false">
      <div class="space-y-3">
        <el-input v-model="createForm.username" placeholder="账号 *" />
        <el-input v-model="createForm.password" type="password" placeholder="初始密码 *" show-password />
        <el-input v-model="createForm.name" placeholder="姓名 *" />
        <el-input v-model="createForm.email" placeholder="邮箱 *" />
        <el-select v-model="createForm.role" placeholder="角色" class="!w-full">
          <el-option label="学生" value="student" />
          <el-option label="管理员" value="admin" />
        </el-select>
      </div>
      <template #footer>
        <button type="button" class="rounded-xl border px-4 py-2 text-sm" @click="createDialogVisible = false">取消</button>
        <button type="button" class="btn-primary ml-2 rounded-xl px-4 py-2 text-sm text-white" :disabled="saving" @click="handleCreate">创建</button>
      </template>
    </el-dialog>

    <el-dialog v-model="editDialogVisible" title="编辑用户" width="480px" :close-on-click-modal="false">
      <div class="space-y-3">
        <el-input v-model="editForm.name" placeholder="姓名" />
        <el-input v-model="editForm.email" placeholder="邮箱" />
        <el-input v-model="editForm.college" placeholder="学院" />
        <el-input v-model="editForm.major" placeholder="专业" />
        <el-select v-model="editForm.role" placeholder="角色（提权/降权）" class="!w-full">
          <el-option label="学生" value="student" />
          <el-option label="管理员" value="admin" />
        </el-select>
      </div>
      <template #footer>
        <button type="button" class="rounded-xl border px-4 py-2 text-sm" @click="editDialogVisible = false">取消</button>
        <button type="button" class="btn-primary ml-2 rounded-xl px-4 py-2 text-sm text-white" :disabled="saving" @click="handleEdit">保存</button>
      </template>
    </el-dialog>

    <el-dialog v-model="resetDialogVisible" title="重置密码" width="400px" :close-on-click-modal="false">
      <el-input v-model="newPassword" type="password" placeholder="新密码（至少8位）" show-password />
      <template #footer>
        <button type="button" class="rounded-xl border px-4 py-2 text-sm" @click="resetDialogVisible = false">取消</button>
        <button type="button" class="btn-primary ml-2 rounded-xl px-4 py-2 text-sm text-white" :disabled="saving" @click="handleReset">确认重置</button>
      </template>
    </el-dialog>
  </div>
</template>
