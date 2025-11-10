import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

// 路由組件懶載入
const Home = () => import('@/views/Home.vue')
const Login = () => import('@/views/Login.vue')
const Matches = () => import('@/views/Matches.vue')
const CreateMatch = () => import('@/views/CreateMatch.vue')
const MyMatches = () => import('@/views/MyMatches.vue')
const Admin = () => import('@/views/Admin.vue')
const Profile = () => import('@/views/Profile.vue')
const MatchDetails = () => import('@/views/MatchDetails.vue')

const routes = [
  {
    path: '/',
    name: 'Home',
    component: Home,
    meta: { title: '首頁' }
  },
  {
    path: '/login',
    name: 'Login',
    component: Login,
    meta: { title: '登入' }
  },
  {
    path: '/matches',
    name: 'Matches',
    component: Matches,
    meta: { title: '配對列表', requiresAuth: true }
  },
  {
    path: '/matches/create',
    name: 'CreateMatch',
    component: CreateMatch,
    meta: { title: '創建配對', requiresAuth: true }
  },
  {
    path: '/matches/:id',
    name: 'MatchDetails',
    component: MatchDetails,
    meta: { title: '配對詳情', requiresAuth: true }
  },
  {
    path: '/my-matches',
    name: 'MyMatches',
    component: MyMatches,
    meta: { title: '我的配對', requiresAuth: true }
  },
  {
    path: '/profile',
    name: 'Profile',
    component: Profile,
    meta: { title: '個人資料', requiresAuth: true }
  },
  {
    path: '/admin',
    name: 'Admin',
    component: Admin,
    meta: { title: '管理後台', requiresAuth: true, requiresAdmin: true }
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

// 路由守衛
router.beforeEach((to, from, next) => {
  const authStore = useAuthStore()
  
  // 設置頁面標題
  if (to.meta.title) {
    document.title = `${to.meta.title} - 買一送一配對網站`
  }
  
  // 檢查是否需要登入
  if (to.meta.requiresAuth && !authStore.isAuthenticated) {
    next({ name: 'Login', query: { redirect: to.fullPath } })
    return
  }
  
  // 檢查是否需要管理員權限
  if (to.meta.requiresAdmin && !authStore.isAdmin) {
    next({ name: 'Home' })
    return
  }
  
  // 如果已登入且訪問登入頁面，重定向到首頁
  if (to.name === 'Login' && authStore.isAuthenticated) {
    next({ name: 'Home' })
    return
  }
  
  next()
})

export default router