import type { RouteRecordRaw } from 'vue-router'

declare module 'vue-router' {
  interface RouteMeta {
    title?: string
    requiresAuth?: boolean
    guestOnly?: boolean
    requiresAdmin?: boolean
    layout?: 'default' | 'auth' | 'admin' | 'blank'
  }
}

export const routes: RouteRecordRaw[] = [
  {
    path: '/',
    name: 'Home',
    component: () => import('@/views/home/HomeView.vue'),
    meta: { title: '首页', layout: 'default' },
  },
  {
    path: '/recommendations',
    name: 'Recommendations',
    component: () => import('@/views/recommendations/RecommendationView.vue'),
    meta: { title: '个性化推荐', requiresAuth: true, layout: 'default' },
  },
  {
    path: '/career-talks',
    name: 'CareerTalkList',
    component: () => import('@/views/career-talks/CareerTalkListView.vue'),
    meta: { title: '宣讲会', layout: 'default' },
  },
  {
    path: '/career-talks/:id',
    name: 'CareerTalkDetail',
    component: () => import('@/views/career-talks/CareerTalkDetailView.vue'),
    meta: { title: '宣讲会详情', layout: 'default' },
  },
  {
    path: '/job-fairs',
    name: 'JobFairList',
    component: () => import('@/views/job-fairs/JobFairListView.vue'),
    meta: { title: '双选会', layout: 'default' },
  },
  {
    path: '/job-fairs/:id',
    name: 'JobFairDetail',
    component: () => import('@/views/job-fairs/JobFairDetailView.vue'),
    meta: { title: '双选会详情', layout: 'default' },
  },
  {
    path: '/calendar',
    name: 'Calendar',
    component: () => import('@/views/calendar/CalendarView.vue'),
    meta: { title: '我的日历', requiresAuth: true, layout: 'default' },
  },
  {
    path: '/profile',
    component: () => import('@/views/profile/ProfileView.vue'),
    meta: { title: '个人中心', requiresAuth: true, layout: 'default' },
    redirect: '/profile/info',
    children: [
      {
        path: 'info',
        name: 'ProfileInfo',
        component: () => import('@/views/profile/ProfileInfoView.vue'),
        meta: { title: '基本资料' },
      },
      {
        path: 'preferences',
        name: 'ProfilePreferences',
        component: () => import('@/views/profile/PreferencesView.vue'),
        meta: { title: '偏好设置' },
      },
      {
        path: 'reminders',
        name: 'ProfileReminders',
        component: () => import('@/views/profile/RemindersView.vue'),
        meta: { title: '提醒记录' },
      },
    ],
  },
  {
    path: '/admin',
    component: () => import('@/views/admin/AdminView.vue'),
    meta: { title: '管理后台', requiresAuth: true, requiresAdmin: true, layout: 'admin' },
    redirect: '/admin/career-talks',
    children: [
      {
        path: 'career-talks',
        name: 'AdminCareerTalks',
        component: () => import('@/views/admin/AdminCareerTalkListView.vue'),
        meta: { title: '宣讲会管理' },
      },
      {
        path: 'job-fairs',
        name: 'AdminJobFairs',
        component: () => import('@/views/admin/AdminJobFairListView.vue'),
        meta: { title: '双选会管理' },
      },
      {
        path: 'users',
        name: 'AdminUsers',
        component: () => import('@/views/admin/AdminUserListView.vue'),
        meta: { title: '用户管理' },
      },
      {
        path: 'sync',
        name: 'AdminSync',
        component: () => import('@/views/admin/AdminSyncView.vue'),
        meta: { title: '信息同步' },
      },
      {
        path: 'audit-logs',
        name: 'AdminAuditLogs',
        component: () => import('@/views/admin/AdminAuditLogView.vue'),
        meta: { title: '审计日志' },
      },
    ],
  },
  {
    path: '/login',
    name: 'Login',
    component: () => import('@/views/auth/LoginView.vue'),
    meta: { title: '登录', guestOnly: true, layout: 'auth' },
  },
  {
    path: '/register',
    name: 'Register',
    component: () => import('@/views/auth/RegisterView.vue'),
    meta: { title: '注册', guestOnly: true, layout: 'auth' },
  },
  {
    path: '/:pathMatch(.*)*',
    name: 'NotFound',
    component: () => import('@/views/error/NotFoundView.vue'),
    meta: { title: '页面不存在', layout: 'blank' },
  },
]
