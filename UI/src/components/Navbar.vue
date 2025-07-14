<script setup>
import { ref, onMounted, computed } from 'vue'
import apiService from '@/services/api.js'

const stats = ref({
  totalTasks: 0,
  completedTasks: 0,
  todayTasks: 0,
  totalLists: 0
})
const loading = ref(false)
const error = ref(null)

const completionRate = computed(() => {
  if (stats.value.totalTasks === 0) return 0
  return Math.round((stats.value.completedTasks / stats.value.totalTasks) * 100)
})

async function fetchStats() {
  loading.value = true
  error.value = null
  
  try {
    const [todos, lists] = await Promise.all([
      apiService.getTodos(),
      apiService.getLists()
    ])
    
    const today = new Date().toISOString().slice(0, 10)
    const todayTasks = todos.filter(todo => todo.due === today)
    const completedTasks = todos.filter(todo => todo.done)
    
    stats.value = {
      totalTasks: todos.length,
      completedTasks: completedTasks.length,
      todayTasks: todayTasks.length,
      totalLists: lists.length
    }
  } catch (err) {
    error.value = err.message
    console.error('Failed to fetch stats:', err)
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  fetchStats()
})
</script>

<template>
  <nav class="bg-gray-800 p-4">
    <div class="container mx-auto flex justify-between items-center">
      <div class="text-white text-2xl font-bold">Todo List</div>
      
      <!-- Stats -->
      <div v-if="!loading" class="flex items-center space-x-6 text-white">
        <div class="text-center">
          <div class="text-sm text-gray-300">Total Tasks</div>
          <div class="text-xl font-bold">{{ stats.totalTasks }}</div>
        </div>
        <div class="text-center">
          <div class="text-sm text-gray-300">Completed</div>
          <div class="text-xl font-bold text-green-400">{{ stats.completedTasks }}</div>
        </div>
        <div class="text-center">
          <div class="text-sm text-gray-300">Today</div>
          <div class="text-xl font-bold text-blue-400">{{ stats.todayTasks }}</div>
        </div>
        <div class="text-center">
          <div class="text-sm text-gray-300">Lists</div>
          <div class="text-xl font-bold text-purple-400">{{ stats.totalLists }}</div>
        </div>
        <div class="text-center">
          <div class="text-sm text-gray-300">Progress</div>
          <div class="text-xl font-bold text-yellow-400">{{ completionRate }}%</div>
        </div>
      </div>
      
      <!-- Loading State -->
      <div v-else class="flex items-center space-x-4">
        <div class="animate-spin rounded-full h-6 w-6 border-b-2 border-white"></div>
        <span class="text-white text-sm">Loading stats...</span>
      </div>
    </div>
  </nav>
</template>

