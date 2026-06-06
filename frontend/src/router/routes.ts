import type { RouteRecordRaw } from 'vue-router'

declare module 'vue-router' {
  interface RouteMeta {
    title?: string
    requiresAuth?: boolean
    guestOnly?: boolean
    layout?: 'default' | 'auth' | 'blank'
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
