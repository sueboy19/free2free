<template>
  <div class="min-h-screen bg-gray-50">
    <!-- 導航列 -->
    <nav class="bg-white shadow-sm border-b">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div class="flex justify-between h-16">
          <div class="flex items-center">
            <router-link to="/" class="text-xl font-bold text-gray-900">買一送一配對</router-link>
          </div>
          
          <div class="hidden md:flex items-center space-x-4">
            <router-link to="/" class="text-gray-700 hover:text-gray-900 px-3 py-2 rounded-md text-sm font-medium">
              首頁
            </router-link>
            <router-link to="/matches" class="text-gray-700 hover:text-gray-900 px-3 py-2 rounded-md text-sm font-medium">
              配對列表
            </router-link>
            <router-link to="/my-matches" class="text-gray-700 hover:text-gray-900 px-3 py-2 rounded-md text-sm font-medium">
              我的配對
            </router-link>
            <router-link to="/profile" class="text-primary-600 px-3 py-2 rounded-md text-sm font-medium">
              個人資料
            </router-link>
            <router-link v-if="authStore.isAdmin" to="/admin" class="text-gray-700 hover:text-gray-900 px-3 py-2 rounded-md text-sm font-medium">
              管理後台
            </router-link>
            <button @click="authStore.logout" class="btn-secondary">
              登出
            </button>
          </div>
        </div>
      </div>
    </nav>

    <!-- 主要內容 -->
    <main class="max-w-4xl mx-auto py-6 px-4 sm:px-6 lg:px-8">
      <!-- 頁面標題 -->
      <div class="mb-6">
        <h1 class="text-2xl font-bold text-gray-900">個人資料</h1>
        <p class="text-gray-600">管理您的個人資訊</p>
      </div>

      <!-- 用戶資訊卡片 -->
      <div class="card mb-6">
        <div class="flex items-center space-x-6">
          <!-- 頭像 -->
          <div class="flex-shrink-0">
            <img 
              :src="authStore.user?.avatar_url || '/default-avatar.png'" 
              :alt="authStore.user?.name"
              class="h-24 w-24 rounded-full object-cover"
            />
          </div>
          
          <!-- 基本資訊 -->
          <div class="flex-1">
            <h2 class="text-xl font-semibold text-gray-900">{{ authStore.user?.name }}</h2>
            <p class="text-gray-600">{{ authStore.user?.email }}</p>
            <p class="text-sm text-gray-500">
              登入方式：{{ authStore.user?.social_provider === 'facebook' ? 'Facebook' : 'Instagram' }}
            </p>
            <div class="mt-2">
              <span 
                v-if="authStore.isAdmin"
                class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-blue-100 text-blue-800"
              >
                管理員
              </span>
            </div>
          </div>
        </div>
      </div>

      <!-- 統計資訊 -->
      <div class="grid grid-cols-1 md:grid-cols-3 gap-6 mb-6">
        <div class="card text-center">
          <div class="text-2xl font-bold text-primary-600">{{ stats.organizedCount }}</div>
          <div class="text-sm text-gray-500">開局數量</div>
        </div>
        <div class="card text-center">
          <div class="text-2xl font-bold text-primary-600">{{ stats.participatedCount }}</div>
          <div class="text-sm text-gray-500">參與數量</div>
        </div>
        <div class="card text-center">
          <div class="text-2xl font-bold text-primary-600">{{ stats.completedCount }}</div>
          <div class="text-sm text-gray-500">完成數量</div>
        </div>
      </div>

      <!-- 最近活動 -->
      <div class="card">
        <h3 class="text-lg font-semibold text-gray-900 mb-4">最近活動</h3>
        <div v-if="isLoading" class="text-center py-8">
          <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-primary-600 mx-auto"></div>
          <p class="text-gray-500 mt-2">載入中...</p>
        </div>
        <div v-else-if="recentActivities.length === 0" class="text-center py-8">
          <p class="text-gray-500">暫無活動記錄</p>
        </div>
        <div v-else class="space-y-4">
          <div 
            v-for="activity in recentActivities" 
            :key="activity.id"
            class="flex items-center justify-between p-3 border rounded-lg hover:bg-gray-50"
          >
            <div class="flex items-center space-x-3">
              <div class="flex-shrink-0">
                <div 
                  :class="[
                    'w-2 h-2 rounded-full',
                    activity.type === 'organize' ? 'bg-blue-500' : 'bg-green-500'
                  ]"
                ></div>
              </div>
              <div>
                <p class="text-sm font-medium text-gray-900">
                  {{ activity.type === 'organize' ? '開局' : '參與' }}: {{ activity.activity_title }}
                </p>
                <p class="text-xs text-gray-500">{{ formatDate(activity.match_time) }}</p>
              </div>
            </div>
            <span 
              :class="[
                'px-2 py-1 text-xs font-semibold rounded-full',
                activity.status === 'open' ? 'bg-green-100 text-green-800' :
                activity.status === 'completed' ? 'bg-blue-100 text-blue-800' :
                'bg-red-100 text-red-800'
              ]"
            >
              {{ activity.status === 'open' ? '進行中' : 
                 activity.status === 'completed' ? '已完成' : '已取消' }}
            </span>
          </div>
        </div>
      </div>
    </main>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useAuthStore } from '@/stores/auth'
import ApiService from '@/services/api'
import { format } from 'date-fns'
import { zhTW } from 'date-fns/locale'

const authStore = useAuthStore()

const isLoading = ref(false)
const recentActivities = ref<any[]>([])
const stats = ref({
  organizedCount: 0,
  participatedCount: 0,
  completedCount: 0
})

// 格式化日期
const formatDate = (date: string) => {
  return format(new Date(date), 'MM月dd日 HH:mm', { locale: zhTW })
}

// 載入統計資料
const loadStats = async () => {
  try {
    const [matchesResponse, pastMatchesResponse] = await Promise.all([
      ApiService.getMatches(),
      ApiService.getPastMatches()
    ])
    
    const allMatches = [...matchesResponse.data, ...pastMatchesResponse.data]
    
    stats.value.organizedCount = allMatches.filter((match: any) => 
      match.organizer_id === authStore.user?.id
    ).length
    
    stats.value.participatedCount = allMatches.filter((match: any) => 
      match.organizer_id !== authStore.user?.id
    ).length
    
    stats.value.completedCount = allMatches.filter((match: any) => 
      match.status === 'completed'
    ).length
    
  } catch (error) {
    console.error('載入統計資料失敗:', error)
  }
}

// 載入最近活動
const loadRecentActivities = async () => {
  try {
    isLoading.value = true
    
    // 載入最近的配對活動
    const response = await ApiService.getMatches()
    const matches = response.data
    
    // 轉換為活動記錄格式
    recentActivities.value = matches.map((match: any) => ({
      id: match.id,
      type: match.organizer_id === authStore.user?.id ? 'organize' : 'participate',
      activity_title: match.activity?.title,
      match_time: match.match_time,
      status: match.status
    })).slice(0, 5) // 只顯示最近5筆
    
  } catch (error) {
    console.error('載入最近活動失敗:', error)
  } finally {
    isLoading.value = false
  }
}

onMounted(() => {
  loadStats()
  loadRecentActivities()
})
</script>