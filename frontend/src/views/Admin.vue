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
            <router-link to="/admin" class="text-primary-600 px-3 py-2 rounded-md text-sm font-medium">
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
    <main class="max-w-7xl mx-auto py-6 px-4 sm:px-6 lg:px-8">
      <!-- 頁面標題 -->
      <div class="mb-6">
        <h1 class="text-2xl font-bold text-gray-900">管理後台</h1>
        <p class="text-gray-600">管理配對活動和地點資訊</p>
      </div>

      <!-- 標籤頁 -->
      <div class="mb-6">
        <nav class="flex space-x-8" aria-label="Tabs">
          <button
            @click="activeTab = 'activities'"
            :class="[
              activeTab === 'activities'
                ? 'border-primary-500 text-primary-600'
                : 'border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300',
              'whitespace-nowrap py-2 px-1 border-b-2 font-medium text-sm'
            ]"
          >
            配對活動管理
          </button>
          <button
            @click="activeTab = 'locations'"
            :class="[
              activeTab === 'locations'
                ? 'border-primary-500 text-primary-600'
                : 'border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300',
              'whitespace-nowrap py-2 px-1 border-b-2 font-medium text-sm'
            ]"
          >
            地點管理
          </button>
        </nav>
      </div>

      <!-- 載入狀態 -->
      <div v-if="isLoading" class="text-center py-12">
        <div class="animate-spin rounded-full h-12 w-12 border-b-2 border-primary-600 mx-auto"></div>
        <p class="text-gray-500 mt-4">載入中...</p>
      </div>

      <!-- 配對活動管理 -->
      <div v-if="activeTab === 'activities' && !isLoading">
        <div class="mb-4 flex justify-between items-center">
          <h2 class="text-lg font-semibold">配對活動列表</h2>
          <button @click="showActivityForm = !showActivityForm" class="btn-primary">
            {{ showActivityForm ? '取消' : '新增活動' }}
          </button>
        </div>

        <!-- 新增活動表單 -->
        <div v-if="showActivityForm" class="card mb-6">
          <h3 class="text-lg font-semibold mb-4">新增配對活動</h3>
          <form @submit.prevent="createActivity">
            <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
              <div>
                <label class="label">活動標題 *</label>
                <input v-model="activityForm.title" type="text" class="input" required />
              </div>
              <div>
                <label class="label">目標人數 *</label>
                <input v-model="activityForm.target_count" type="number" min="1" class="input" required />
              </div>
              <div>
                <label class="label">選擇地點 *</label>
                <select v-model="activityForm.location_id" class="input" required>
                  <option value="">請選擇地點</option>
                  <option v-for="location in locations" :key="location.id" :value="location.id">
                    {{ location.name }}
                  </option>
                </select>
              </div>
              <div class="md:col-span-2">
                <label class="label">活動描述</label>
                <textarea v-model="activityForm.description" class="input" rows="3"></textarea>
              </div>
            </div>
            <div class="flex space-x-4 mt-4">
              <button type="submit" :disabled="isSubmitting" class="btn-primary">
                {{ isSubmitting ? '建立中...' : '建立活動' }}
              </button>
              <button type="button" @click="resetActivityForm" class="btn-secondary">
                重置
              </button>
            </div>
          </form>
        </div>

        <!-- 活動列表 -->
        <div class="card">
          <div v-if="activities.length === 0" class="text-center py-8">
            <p class="text-gray-500">暫無配對活動</p>
          </div>
          <div v-else class="overflow-x-auto">
            <table class="min-w-full divide-y divide-gray-200">
              <thead class="bg-gray-50">
                <tr>
                  <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                    活動標題
                  </th>
                  <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                    地點
                  </th>
                  <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                    目標人數
                  </th>
                  <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                    創建者
                  </th>
                  <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                    操作
                  </th>
                </tr>
              </thead>
              <tbody class="bg-white divide-y divide-gray-200">
                <tr v-for="activity in activities" :key="activity.id">
                  <td class="px-6 py-4 whitespace-nowrap">
                    <div class="text-sm font-medium text-gray-900">{{ activity.title }}</div>
                    <div class="text-sm text-gray-500">{{ activity.description }}</div>
                  </td>
                  <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                    {{ activity.location?.name }}
                  </td>
                  <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                    {{ activity.target_count }}
                  </td>
                  <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                    {{ activity.created_by }}
                  </td>
                  <td class="px-6 py-4 whitespace-nowrap text-sm font-medium">
                    <button @click="editActivity(activity)" class="text-primary-600 hover:text-primary-900 mr-4">
                      編輯
                    </button>
                    <button @click="deleteActivity(activity.id)" class="text-red-600 hover:text-red-900">
                      刪除
                    </button>
                  </td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>
      </div>

      <!-- 地點管理 -->
      <div v-if="activeTab === 'locations' && !isLoading">
        <div class="mb-4 flex justify-between items-center">
          <h2 class="text-lg font-semibold">地點列表</h2>
          <button @click="showLocationForm = !showLocationForm" class="btn-primary">
            {{ showLocationForm ? '取消' : '新增地點' }}
          </button>
        </div>

        <!-- 新增地點表單 -->
        <div v-if="showLocationForm" class="card mb-6">
          <h3 class="text-lg font-semibold mb-4">新增地點</h3>
          <form @submit.prevent="createLocation">
            <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
              <div>
                <label class="label">地點名稱 *</label>
                <input v-model="locationForm.name" type="text" class="input" required />
              </div>
              <div class="md:col-span-2">
                <label class="label">地址 *</label>
                <input v-model="locationForm.address" type="text" class="input" required />
              </div>
              <div>
                <label class="label">緯度 *</label>
                <input v-model="locationForm.latitude" type="number" step="0.000001" class="input" required />
              </div>
              <div>
                <label class="label">經度 *</label>
                <input v-model="locationForm.longitude" type="number" step="0.000001" class="input" required />
              </div>
            </div>
            <div class="flex space-x-4 mt-4">
              <button type="submit" :disabled="isSubmitting" class="btn-primary">
                {{ isSubmitting ? '建立中...' : '建立地點' }}
              </button>
              <button type="button" @click="resetLocationForm" class="btn-secondary">
                重置
              </button>
            </div>
          </form>
        </div>

        <!-- 地點列表 -->
        <div class="card">
          <div v-if="locations.length === 0" class="text-center py-8">
            <p class="text-gray-500">暫無地點</p>
          </div>
          <div v-else class="overflow-x-auto">
            <table class="min-w-full divide-y divide-gray-200">
              <thead class="bg-gray-50">
                <tr>
                  <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                    地點名稱
                  </th>
                  <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                    地址
                  </th>
                  <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                    座標
                  </th>
                  <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                    操作
                  </th>
                </tr>
              </thead>
              <tbody class="bg-white divide-y divide-gray-200">
                <tr v-for="location in locations" :key="location.id">
                  <td class="px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900">
                    {{ location.name }}
                  </td>
                  <td class="px-6 py-4 text-sm text-gray-500">
                    {{ location.address }}
                  </td>
                  <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                    {{ location.latitude }}, {{ location.longitude }}
                  </td>
                  <td class="px-6 py-4 whitespace-nowrap text-sm font-medium">
                    <button @click="editLocation(location)" class="text-primary-600 hover:text-primary-900 mr-4">
                      編輯
                    </button>
                    <button @click="deleteLocation(location.id)" class="text-red-600 hover:text-red-900">
                      刪除
                    </button>
                  </td>
                </tr>
              </tbody>
            </table>
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
import { useToast } from 'vue-toastification'

const authStore = useAuthStore()
const toast = useToast()

const isLoading = ref(false)
const isSubmitting = ref(false)
const activeTab = ref('activities')
const showActivityForm = ref(false)
const showLocationForm = ref(false)

const activities = ref<any[]>([])
const locations = ref<any[]>([])

const activityForm = ref({
  title: '',
  target_count: 2,
  location_id: '',
  description: ''
})

const locationForm = ref({
  name: '',
  address: '',
  latitude: 0,
  longitude: 0
})

// 載入活動列表
const loadActivities = async () => {
  try {
    const response = await ApiService.getActivities()
    activities.value = response.data
  } catch (error) {
    console.error('載入活動失敗:', error)
    toast.error('載入活動失敗')
  }
}

// 載入地點列表
const loadLocations = async () => {
  try {
    const response = await ApiService.getLocations()
    locations.value = response.data
  } catch (error) {
    console.error('載入地點失敗:', error)
    toast.error('載入地點失敗')
  }
}

// 創建活動
const createActivity = async () => {
  try {
    isSubmitting.value = true
    await ApiService.createActivity(activityForm.value)
    toast.success('活動建立成功')
    await loadActivities()
    resetActivityForm()
  } catch (error) {
    console.error('建立活動失敗:', error)
    toast.error('建立活動失敗')
  } finally {
    isSubmitting.value = false
  }
}

// 重置活動表單
const resetActivityForm = () => {
  activityForm.value = {
    title: '',
    target_count: 2,
    location_id: '',
    description: ''
  }
  showActivityForm.value = false
}

// 創建地點
const createLocation = async () => {
  try {
    isSubmitting.value = true
    await ApiService.createLocation(locationForm.value)
    toast.success('地點建立成功')
    await loadLocations()
    resetLocationForm()
  } catch (error) {
    console.error('建立地點失敗:', error)
    toast.error('建立地點失敗')
  } finally {
    isSubmitting.value = false
  }
}

// 重置地點表單
const resetLocationForm = () => {
  locationForm.value = {
    name: '',
    address: '',
    latitude: 0,
    longitude: 0
  }
  showLocationForm.value = false
}

// 刪除活動
const deleteActivity = async (id: number) => {
  if (confirm('確定要刪除這個活動嗎？')) {
    try {
      await ApiService.deleteActivity(id)
      toast.success('活動刪除成功')
      await loadActivities()
    } catch (error) {
      console.error('刪除活動失敗:', error)
      toast.error('刪除活動失敗')
    }
  }
}

// 刪除地點
const deleteLocation = async (id: number) => {
  if (confirm('確定要刪除這個地點嗎？')) {
    try {
      await ApiService.deleteLocation(id)
      toast.success('地點刪除成功')
      await loadLocations()
    } catch (error) {
      console.error('刪除地點失敗:', error)
      toast.error('刪除地點失敗')
    }
  }
}

// 編輯活動
const editActivity = (activity: any) => {
  activityForm.value = { ...activity }
  showActivityForm.value = true
}

// 編輯地點
const editLocation = (location: any) => {
  locationForm.value = { ...location }
  showLocationForm.value = true
}

onMounted(async () => {
  isLoading.value = true
  await Promise.all([loadActivities(), loadLocations()])
  isLoading.value = false
})
</script>