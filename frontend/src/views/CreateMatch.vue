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
    <main class="max-w-3xl mx-auto py-6 px-4 sm:px-6 lg:px-8">
      <!-- 頁面標題 -->
      <div class="mb-6">
        <h1 class="text-2xl font-bold text-gray-900">創建配對</h1>
        <p class="text-gray-600">發起一個新的配對活動</p>
      </div>

      <!-- 創建表單 -->
      <div class="card">
        <form @submit.prevent="createMatch">
          <!-- 活動選擇 -->
          <div class="mb-6">
            <label class="label">選擇活動 *</label>
            <select v-model="formData.activity_id" class="input" required>
              <option value="">請選擇活動</option>
              <option v-for="activity in activities" :key="activity.id" :value="activity.id">
                {{ activity.title }} - {{ activity.location.name }}
              </option>
            </select>
            <p class="text-sm text-gray-500 mt-1">
              如果沒有適合的活動，請先聯繫管理員創建。
            </p>
          </div>

          <!-- 配對時間 -->
          <div class="mb-6">
            <label class="label">配對時間 *</label>
            <input 
              v-model="formData.match_time" 
              type="datetime-local" 
              class="input" 
              required
              :min="minDateTime"
            />
            <p class="text-sm text-gray-500 mt-1">
              請選擇未來的時間點
            </p>
          </div>

          <!-- 預覽 -->
          <div v-if="selectedActivity" class="mb-6 p-4 bg-gray-50 rounded-lg">
            <h3 class="font-semibold text-gray-900 mb-2">配對預覽</h3>
            <div class="space-y-2 text-sm">
              <p><span class="font-medium">活動：</span>{{ selectedActivity.title }}</p>
              <p><span class="font-medium">地點：</span>{{ selectedActivity.location.name }}</p>
              <p><span class="font-medium">地址：</span>{{ selectedActivity.location.address }}</p>
              <p><span class="font-medium">目標人數：</span>{{ selectedActivity.target_count }}人</p>
              <p><span class="font-medium">描述：</span>{{ selectedActivity.description }}</p>
            </div>
          </div>

          <!-- 錯誤訊息 -->
          <div v-if="errorMessage" class="mb-4 p-3 bg-red-100 border border-red-400 text-red-700 rounded">
            {{ errorMessage }}
          </div>

          <!-- 成功訊息 -->
          <div v-if="successMessage" class="mb-4 p-3 bg-green-100 border border-green-400 text-green-700 rounded">
            {{ successMessage }}
          </div>

          <!-- 提交按鈕 -->
          <div class="flex space-x-4">
            <button 
              type="submit" 
              :disabled="isSubmitting"
              class="flex-1 btn-primary"
            >
              {{ isSubmitting ? '創建中...' : '創建配對' }}
            </button>
            <router-link to="/matches" class="flex-1 btn-secondary text-center">
              取消
            </router-link>
          </div>
        </form>
      </div>
    </main>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import ApiService from '@/services/api'
import { useToast } from 'vue-toastification'

const router = useRouter()
const authStore = useAuthStore()
const toast = useToast()

const activities = ref<any[]>([])
const isLoading = ref(false)
const isSubmitting = ref(false)
const errorMessage = ref('')
const successMessage = ref('')

const formData = ref({
  activity_id: '',
  match_time: ''
})

// 最小可選擇的時間
const minDateTime = computed(() => {
  const now = new Date()
  now.setMinutes(now.getMinutes() - now.getTimezoneOffset())
  return now.toISOString().slice(0, 16)
})

// 選中的活動
const selectedActivity = computed(() => {
  return activities.value.find(activity => activity.id === parseInt(formData.value.activity_id))
})

// 載入活動列表
const loadActivities = async () => {
  try {
    isLoading.value = true
    const response = await ApiService.getActivities()
    activities.value = response.data
  } catch (error) {
    console.error('載入活動失敗:', error)
    toast.error('載入活動失敗')
  } finally {
    isLoading.value = false
  }
}

// 創建配對
const createMatch = async () => {
  try {
    isSubmitting.value = true
    errorMessage.value = ''
    successMessage.value = ''

    // 驗證表單
    if (!formData.value.activity_id || !formData.value.match_time) {
      errorMessage.value = '請填寫所有必填欄位'
      return
    }

    // 驗證時間
    const matchTime = new Date(formData.value.match_time)
    const now = new Date()
    if (matchTime <= now) {
      errorMessage.value = '配對時間必須是未來時間'
      return
    }

    const matchData = {
      activity_id: parseInt(formData.value.activity_id),
      match_time: formData.value.match_time
    }

    await ApiService.createMatch(matchData)
    
    successMessage.value = '配對創建成功！'
    toast.success('配對創建成功！')
    
    // 延遲跳轉，讓用戶看到成功訊息
    setTimeout(() => {
      router.push('/matches')
    }, 2000)
    
  } catch (error) {
    console.error('創建配對失敗:', error)
    errorMessage.value = '創建配對失敗，請重試'
    toast.error('創建配對失敗')
  } finally {
    isSubmitting.value = false
  }
}

onMounted(() => {
  loadActivities()
})
</script>