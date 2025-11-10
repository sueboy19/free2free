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
            <router-link to="/my-matches" class="text-primary-600 px-3 py-2 rounded-md text-sm font-medium">
              我的配對
            </router-link>
            <router-link to="/profile" class="text-gray-700 hover:text-gray-900 px-3 py-2 rounded-md text-sm font-medium">
              個人資料
            </router-link>
            <button @click="authStore.logout" class="btn-secondary">
              登出
            </button>
          </div>
        </div>
      </div>
    </nav>

    <!-- 主要內容 -->
    <main class="max-w-7xl mx-auto py-6 px-4 sm:px-6 lg:px-8">
      <!-- 頁面標題 -->
      <div class="mb-6">
        <h1 class="text-2xl font-bold text-gray-900">我的配對</h1>
        <p class="text-gray-600">查看您參與的配對記錄</p>
      </div>

      <!-- 標籤頁 -->
      <div class="mb-6">
        <nav class="flex space-x-8" aria-label="Tabs">
          <button
            @click="activeTab = 'organizing'"
            :class="[
              activeTab === 'organizing'
                ? 'border-primary-500 text-primary-600'
                : 'border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300',
              'whitespace-nowrap py-2 px-1 border-b-2 font-medium text-sm'
            ]"
          >
            我開局的 ({{ organizedMatches.length }})
          </button>
          <button
            @click="activeTab = 'participating'"
            :class="[
              activeTab === 'participating'
                ? 'border-primary-500 text-primary-600'
                : 'border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300',
              'whitespace-nowrap py-2 px-1 border-b-2 font-medium text-sm'
            ]"
          >
            我參與的 ({{ participatingMatches.length }})
          </button>
          <button
            @click="activeTab = 'past'"
            :class="[
              activeTab === 'past'
                ? 'border-primary-500 text-primary-600'
                : 'border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300',
              'whitespace-nowrap py-2 px-1 border-b-2 font-medium text-sm'
            ]"
          >
            歷史配對 ({{ pastMatches.length }})
          </button>
        </nav>
      </div>

      <!-- 載入狀態 -->
      <div v-if="isLoading" class="text-center py-12">
        <div class="animate-spin rounded-full h-12 w-12 border-b-2 border-primary-600 mx-auto"></div>
        <p class="text-gray-500 mt-4">載入中...</p>
      </div>

      <!-- 標籤內容 -->
      <div v-else>
        <!-- 我開局的配對 -->
        <div v-if="activeTab === 'organizing'">
          <div v-if="organizedMatches.length === 0" class="text-center py-12">
            <p class="text-gray-500 text-lg">您還沒有開局過配對</p>
            <router-link to="/matches/create" class="btn-primary mt-4 inline-block">
              創建第一個配對
            </router-link>
          </div>
          <div v-else class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
            <div 
              v-for="match in organizedMatches" 
              :key="match.id" 
              class="card hover:shadow-lg transition-shadow duration-200"
            >
              <div class="mb-4">
                <h3 class="text-lg font-semibold text-gray-900 mb-2">
                  {{ match.activity?.title }}
                </h3>
                <p class="text-gray-600 text-sm mb-3">
                  {{ match.activity?.description }}
                </p>
                
                <div class="flex items-center text-sm text-gray-500 mb-2">
                  <svg class="w-4 h-4 mr-1" fill="currentColor" viewBox="0 0 20 20">
                    <path fill-rule="evenodd" d="M5.05 4.05a7 7 0 119.9 9.9L10 18.9l-4.95-4.95a7 7 0 010-9.9zM10 11a2 2 0 100-4 2 2 0 000 4z" clip-rule="evenodd" />
                  </svg>
                  {{ match.activity?.location?.name }}
                </div>
                
                <div class="flex items-center text-sm text-gray-500 mb-2">
                  <svg class="w-4 h-4 mr-1" fill="currentColor" viewBox="0 0 20 20">
                    <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm1-12a1 1 0 10-2 0v4a1 1 0 00.293.707l2.828 2.829a1 1 0 101.415-1.415L11 9.586V6z" clip-rule="evenodd" />
                  </svg>
                  {{ formatDate(match.match_time) }}
                </div>
                
                <span 
                  :class="[
                    'inline-flex px-2 py-1 text-xs font-semibold rounded-full',
                    match.status === 'open' ? 'bg-green-100 text-green-800' :
                    match.status === 'completed' ? 'bg-blue-100 text-blue-800' :
                    'bg-red-100 text-red-800'
                  ]"
                >
                  {{ match.status === 'open' ? '進行中' : 
                     match.status === 'completed' ? '已完成' : '已取消' }}
                </span>
              </div>
              
              <div class="flex space-x-2">
                <router-link 
                  :to="`/matches/${match.id}`" 
                  class="flex-1 btn-secondary text-center text-sm"
                >
                  管理配對
                </router-link>
              </div>
            </div>
          </div>
        </div>

        <!-- 我參與的配對 -->
        <div v-if="activeTab === 'participating'">
          <div v-if="participatingMatches.length === 0" class="text-center py-12">
            <p class="text-gray-500 text-lg">您還沒有參與任何配對</p>
            <router-link to="/matches" class="btn-primary mt-4 inline-block">
              瀏覽配對
            </router-link>
          </div>
          <div v-else class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
            <div 
              v-for="match in participatingMatches" 
              :key="match.id" 
              class="card hover:shadow-lg transition-shadow duration-200"
            >
              <div class="mb-4">
                <h3 class="text-lg font-semibold text-gray-900 mb-2">
                  {{ match.activity?.title }}
                </h3>
                <p class="text-gray-600 text-sm mb-3">
                  {{ match.activity?.description }}
                </p>
                
                <div class="flex items-center text-sm text-gray-500 mb-2">
                  <svg class="w-4 h-4 mr-1" fill="currentColor" viewBox="0 0 20 20">
                    <path fill-rule="evenodd" d="M5.05 4.05a7 7 0 119.9 9.9L10 18.9l-4.95-4.95a7 7 0 010-9.9zM10 11a2 2 0 100-4 2 2 0 000 4z" clip-rule="evenodd" />
                  </svg>
                  {{ match.activity?.location?.name }}
                </div>
                
                <div class="flex items-center text-sm text-gray-500 mb-2">
                  <svg class="w-4 h-4 mr-1" fill="currentColor" viewBox="0 0 20 20">
                    <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm1-12a1 1 0 10-2 0v4a1 1 0 00.293.707l2.828 2.829a1 1 0 101.415-1.415L11 9.586V6z" clip-rule="evenodd" />
                  </svg>
                  {{ formatDate(match.match_time) }}
                </div>
                
                <div class="flex items-center text-sm text-gray-500 mb-2">
                  <svg class="w-4 h-4 mr-1" fill="currentColor" viewBox="0 0 20 20">
                    <path d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
                  </svg>
                  開局者: {{ match.organizer?.name }}
                </div>
                
                <span 
                  :class="[
                    'inline-flex px-2 py-1 text-xs font-semibold rounded-full',
                    match.status === 'open' ? 'bg-green-100 text-green-800' :
                    match.status === 'completed' ? 'bg-blue-100 text-blue-800' :
                    'bg-red-100 text-red-800'
                  ]"
                >
                  {{ match.status === 'open' ? '進行中' : 
                     match.status === 'completed' ? '已完成' : '已取消' }}
                </span>
              </div>
              
              <div class="flex space-x-2">
                <router-link 
                  :to="`/matches/${match.id}`" 
                  class="flex-1 btn-primary text-center text-sm"
                >
                  查看詳情
                </router-link>
              </div>
            </div>
          </div>
        </div>

        <!-- 歷史配對 -->
        <div v-if="activeTab === 'past'">
          <div v-if="pastMatches.length === 0" class="text-center py-12">
            <p class="text-gray-500 text-lg">沒有歷史配對記錄</p>
          </div>
          <div v-else class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
            <div 
              v-for="match in pastMatches" 
              :key="match.id" 
              class="card hover:shadow-lg transition-shadow duration-200"
            >
              <div class="mb-4">
                <h3 class="text-lg font-semibold text-gray-900 mb-2">
                  {{ match.activity?.title }}
                </h3>
                <p class="text-gray-600 text-sm mb-3">
                  {{ match.activity?.description }}
                </p>
                
                <div class="flex items-center text-sm text-gray-500 mb-2">
                  <svg class="w-4 h-4 mr-1" fill="currentColor" viewBox="0 0 20 20">
                    <path fill-rule="evenodd" d="M5.05 4.05a7 7 0 119.9 9.9L10 18.9l-4.95-4.95a7 7 0 010-9.9zM10 11a2 2 0 100-4 2 2 0 000 4z" clip-rule="evenodd" />
                  </svg>
                  {{ match.activity?.location?.name }}
                </div>
                
                <div class="flex items-center text-sm text-gray-500 mb-2">
                  <svg class="w-4 h-4 mr-1" fill="currentColor" viewBox="0 0 20 20">
                    <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm1-12a1 1 0 10-2 0v4a1 1 0 00.293.707l2.828 2.829a1 1 0 101.415-1.415L11 9.586V6z" clip-rule="evenodd" />
                  </svg>
                  {{ formatDate(match.match_time) }}
                </div>
                
                <span class="inline-flex px-2 py-1 text-xs font-semibold rounded-full bg-gray-100 text-gray-800">
                  已完成
                </span>
              </div>
              
              <div class="flex space-x-2">
                <router-link 
                  :to="`/matches/${match.id}`" 
                  class="flex-1 btn-secondary text-center text-sm"
                >
                  查看詳情
                </router-link>
              </div>
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
const activeTab = ref('organizing')
const organizedMatches = ref<any[]>([])
const participatingMatches = ref<any[]>([])
const pastMatches = ref<any[]>([])

// 格式化日期
const formatDate = (date: string) => {
  return format(new Date(date), 'MM月dd日 HH:mm', { locale: zhTW })
}

// 載入配對數據
const loadMatches = async () => {
  try {
    isLoading.value = true
    
    // 載入所有配對，然後分類
    const response = await ApiService.getMatches()
    const allMatches = response.data
    
    // 分類為我開局的配對
    organizedMatches.value = allMatches.filter((match: any) => 
      match.organizer_id === authStore.user?.id
    )
    
    // 分類為我參與的配對
    participatingMatches.value = allMatches.filter((match: any) => 
      match.organizer_id !== authStore.user?.id
    )
    
    // 載入歷史配對
    const pastResponse = await ApiService.getPastMatches()
    pastMatches.value = pastResponse.data
    
  } catch (error) {
    console.error('載入配對失敗:', error)
  } finally {
    isLoading.value = false
  }
}

onMounted(() => {
  loadMatches()
})
</script>