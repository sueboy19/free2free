<template>
  <div class="min-h-screen bg-gray-50">
    <!-- 導航列 -->
    <nav class="bg-white shadow-sm border-b">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div class="flex justify-between h-16">
          <div class="flex items-center">
            <h1 class="text-xl font-bold text-gray-900">買一送一配對</h1>
          </div>
          
          <!-- 桌面版導航 -->
          <div class="hidden md:flex items-center space-x-4" v-if="authStore.isAuthenticated">
            <router-link to="/" class="text-gray-700 hover:text-gray-900 px-3 py-2 rounded-md text-sm font-medium">
              首頁
            </router-link>
            <router-link to="/matches" class="text-gray-700 hover:text-gray-900 px-3 py-2 rounded-md text-sm font-medium">
              配對列表
            </router-link>
            <router-link to="/my-matches" class="text-gray-700 hover:text-gray-900 px-3 py-2 rounded-md text-sm font-medium">
              我的配對
            </router-link>
            <router-link to="/profile" class="text-gray-700 hover:text-gray-900 px-3 py-2 rounded-md text-sm font-medium">
              個人資料
            </router-link>
            <router-link v-if="authStore.isAdmin" to="/admin" class="text-gray-700 hover:text-gray-900 px-3 py-2 rounded-md text-sm font-medium">
              管理後台
            </router-link>
            <button @click="authStore.logout" class="btn-secondary">
              登出
            </button>
          </div>
          
          <!-- 登入按鈕 -->
          <div class="flex items-center" v-else>
            <router-link to="/login" class="btn-primary">
              登入
            </router-link>
          </div>
        </div>
      </div>
      
      <!-- 手機版導航 -->
      <div class="md:hidden" v-if="authStore.isAuthenticated">
        <div class="px-2 pt-2 pb-3 space-y-1 sm:px-3 bg-gray-50">
          <router-link to="/" class="block px-3 py-2 text-gray-700 hover:text-gray-900">
            首頁
          </router-link>
          <router-link to="/matches" class="block px-3 py-2 text-gray-700 hover:text-gray-900">
            配對列表
          </router-link>
          <router-link to="/my-matches" class="block px-3 py-2 text-gray-700 hover:text-gray-900">
            我的配對
          </router-link>
          <router-link to="/profile" class="block px-3 py-2 text-gray-700 hover:text-gray-900">
            個人資料
          </router-link>
          <router-link v-if="authStore.isAdmin" to="/admin" class="block px-3 py-2 text-gray-700 hover:text-gray-900">
            管理後台
          </router-link>
          <button @click="authStore.logout" class="block w-full text-left px-3 py-2 text-gray-700 hover:text-gray-900">
            登出
          </button>
        </div>
      </div>
    </nav>

    <!-- 主要內容 -->
    <main class="max-w-7xl mx-auto py-6 px-4 sm:px-6 lg:px-8">
      <!-- 未登入狀態 -->
      <div v-if="!authStore.isAuthenticated" class="text-center">
        <div class="max-w-2xl mx-auto">
          <h2 class="text-3xl font-bold text-gray-900 mb-4">
            歡迎來到買一送一配對網站
          </h2>
          <p class="text-lg text-gray-600 mb-8">
            找尋夥伴一起享受買一送一的優惠！登入後即可開始瀏覽和創建配對。
          </p>
          <router-link to="/login" class="btn-primary text-lg px-8 py-3">
            立即登入
          </router-link>
        </div>
      </div>

      <!-- 已登入狀態 -->
      <div v-else>
        <!-- 歡迎訊息 -->
        <div class="mb-8">
          <h2 class="text-2xl font-bold text-gray-900 mb-2">
            歡迎回來，{{ authStore.user?.name }}！
          </h2>
          <p class="text-gray-600">
            今日已有 {{ todayMatchesCount }} 個新的配對機會
          </p>
        </div>

        <!-- 快速操作 -->
        <div class="grid grid-cols-1 md:grid-cols-3 gap-6 mb-8">
          <div class="card">
            <h3 class="text-lg font-semibold mb-2">瀏覽配對</h3>
            <p class="text-gray-600 mb-4">查看可參與的配對機會</p>
            <router-link to="/matches" class="btn-primary w-full text-center block">
              立即瀏覽
            </router-link>
          </div>
          
          <div class="card">
            <h3 class="text-lg font-semibold mb-2">創建配對</h3>
            <p class="text-gray-600 mb-4">發起一個新的配對活動</p>
            <router-link to="/matches/create" class="btn-primary w-full text-center block">
              創建配對
            </router-link>
          </div>
          
          <div class="card">
            <h3 class="text-lg font-semibold mb-2">我的配對</h3>
            <p class="text-gray-600 mb-4">查看我參與的配對</p>
            <router-link to="/my-matches" class="btn-primary w-full text-center block">
              查看配對
            </router-link>
          </div>
        </div>

        <!-- 今日配對預覽 -->
        <div class="card">
          <h3 class="text-lg font-semibold mb-4">今日推薦配對</h3>
          <div v-if="isLoading" class="text-center py-8">
            <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-primary-600 mx-auto"></div>
            <p class="text-gray-500 mt-2">載入中...</p>
          </div>
          <div v-else-if="featuredMatches.length === 0" class="text-center py-8">
            <p class="text-gray-500">暫無可用的配對</p>
          </div>
          <div v-else class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
            <div v-for="match in featuredMatches" :key="match.id" class="border rounded-lg p-4 hover:shadow-md transition-shadow">
              <h4 class="font-semibold text-gray-900">{{ match.activity?.title }}</h4>
              <p class="text-sm text-gray-600 mt-1">{{ match.activity?.description }}</p>
              <div class="mt-3 text-xs text-gray-500">
                <p>時間: {{ formatDate(match.match_time) }}</p>
                <p>地點: {{ match.activity?.location?.name }}</p>
              </div>
              <router-link :to="`/matches/${match.id}`" class="mt-3 btn-primary w-full text-center block text-sm">
                查看詳情
              </router-link>
            </div>
          </div>
        </div>
      </div>
    </main>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { useAuthStore } from '@/stores/auth'
import ApiService from '@/services/api'
import { format } from 'date-fns'
import { zhTW } from 'date-fns/locale'

const authStore = useAuthStore()
const isLoading = ref(false)
const featuredMatches = ref<any[]>([])

// 今日配對數量
const todayMatchesCount = computed(() => {
  const today = new Date().toDateString()
  return featuredMatches.value.filter(match => 
    new Date(match.match_time).toDateString() === today
  ).length
})

// 格式化日期
const formatDate = (date: string) => {
  return format(new Date(date), 'MM月dd日 HH:mm', { locale: zhTW })
}

// 載入精選配對
const loadFeaturedMatches = async () => {
  try {
    isLoading.value = true
    const response = await ApiService.getMatches()
    // 取前 6 個配對作為精選
    featuredMatches.value = response.data.slice(0, 6)
  } catch (error) {
    console.error('載入配對失敗:', error)
  } finally {
    isLoading.value = false
  }
}

onMounted(() => {
  if (authStore.isAuthenticated) {
    loadFeaturedMatches()
  }
})
</script>