<script setup>
import { ref, nextTick, watch, onMounted, computed } from 'vue'
import apiService from '@/services/api.js'

// Reactive state
const isCollapsed = ref(false)
const showNewListInput = ref(false)
const newListName = ref('')
const newListInput = ref(null)
const loading = ref(false)
const error = ref(null)

const lists = ref([])
const todos = ref([])

// Computed: lists with task counts
const listsWithCounts = computed(() => {
  return lists.value.map(list => {
    const count = todos.value.filter(todo => todo.listId === list.id).length
    return { ...list, count }
  })
})

// Methods
const toggleSidebar = () => {
  isCollapsed.value = !isCollapsed.value
}

async function fetchData() {
  loading.value = true
  error.value = null
  
  try {
    const [listsData, todosData] = await Promise.all([
      apiService.getLists(),
      apiService.getTodos()
    ])
    
    lists.value = listsData
    todos.value = todosData
  } catch (err) {
    error.value = err.message
    console.error('Failed to fetch data:', err)
  } finally {
    loading.value = false
  }
}

const addNewList = async () => {
  if (newListName.value.trim()) {
    try {
      const newList = await apiService.createList({
        name: newListName.value.trim()
      })
      lists.value.push(newList)
      newListName.value = ''
      showNewListInput.value = false
    } catch (err) {
      error.value = err.message
      console.error('Failed to create list:', err)
    }
  }
}

const cancelNewList = () => {
  newListName.value = ''
  showNewListInput.value = false
}

// Auto-focus input when shown
watch(showNewListInput, (newValue) => {
  if (newValue) {
    nextTick(() => {
      if (newListInput.value) {
        newListInput.value.focus()
      }
    })
  }
})

onMounted(() => {
  fetchData()
})
</script>

<template>
  <div class="flex h-screen bg-gray-100">
    <!-- Sidebar -->
    <div 
      :class="[
        'bg-white shadow-lg transition-all duration-300 ease-in-out',
        isCollapsed ? 'w-16' : 'w-64'
      ]"
    >
      <!-- Header with toggle button -->
      <div class="flex items-center justify-between p-4 border-b border-gray-200">
        <h2 v-if="!isCollapsed" class="text-lg font-semibold text-gray-800">
          Menu
        </h2>
        <button 
          @click="toggleSidebar"
          class="p-2 rounded-md hover:bg-gray-100 transition-colors"
          :title="isCollapsed ? 'Expand Sidebar' : 'Collapse Sidebar'"
        >
          <svg class="w-5 h-5 text-gray-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path v-if="!isCollapsed" stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 19l-7-7 7-7m8 14l-7-7 7-7"></path>
            <path v-else stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 5l7 7-7 7M5 5l7 7-7 7"></path>
          </svg>
        </button>
      </div>

      <!-- Home Link -->
      <nav class="p-4 border-b border-gray-100">
        <RouterLink to="/" class="flex items-center gap-2 p-2 rounded-md hover:bg-blue-50 hover:text-blue-600 transition-colors font-medium text-gray-700">
          <span v-if="!isCollapsed">Home</span>
        </RouterLink>
      </nav>

      <!-- Tasks Section -->
      <div class="p-4 border-b border-gray-100">
        <div class="flex items-center mb-2">
          <span v-if="!isCollapsed" class="font-medium text-gray-700">Tasks</span>
        </div>
        <nav class="space-y-1">
          <RouterLink to="/today" class="flex items-center gap-2 p-2 rounded-md hover:bg-blue-50 hover:text-blue-600 transition-colors font-medium text-gray-700">
            <span v-if="!isCollapsed">Today</span>
          </RouterLink>
          
          <RouterLink to="/Upcoming" class="flex items-center gap-2 p-2 rounded-md hover:bg-blue-50 hover:text-blue-600 transition-colors font-medium text-gray-700">
            <span v-if="!isCollapsed">Upcoming</span>
          </RouterLink>
        </nav>
      </div>

      <!-- Sidebar Content -->
      <div class="p-4 space-y-6">
        <!-- Lists Section -->
        <div>
          <div class="flex items-center justify-between mb-3">
            <div class="flex items-center">
              <h3 v-if="!isCollapsed" class="font-medium text-gray-700">Lists</h3>
            </div>
            <button 
              v-if="!isCollapsed"
              @click="showNewListInput = true"
              class="text-blue-600 hover:text-blue-700 text-sm font-medium"
            >
              + New
            </button>
          </div>

          <!-- Error Message -->
          <div v-if="error" class="mb-3 p-2 bg-red-100 border border-red-400 text-red-700 rounded text-xs">
            {{ error }}
          </div>

          <!-- Loading State -->
          <div v-if="loading" class="mb-3 text-center">
            <div class="inline-block animate-spin rounded-full h-4 w-4 border-b-2 border-blue-600"></div>
            <span class="text-xs text-gray-500 ml-2">Loading lists...</span>
          </div>

          <!-- New List Input -->
          <div v-if="showNewListInput && !isCollapsed" class="mb-4 p-3 bg-gray-50 rounded-lg">
            <div class="space-y-3">
              <input 
                v-model="newListName"
                @keyup.enter="addNewList"
                @keyup.esc="cancelNewList"
                type="text"
                placeholder="Enter list name"
                class="w-full px-3 py-2 text-sm border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                ref="newListInput"
              />
              <div class="flex space-x-2">
                <button 
                  @click="addNewList"
                  class="flex-1 px-3 py-1 text-xs bg-blue-500 text-white rounded hover:bg-blue-600 transition-colors font-medium"
                >
                  Create
                </button>
                <button 
                  @click="cancelNewList"
                  class="flex-1 px-3 py-1 text-xs bg-gray-300 text-gray-700 rounded hover:bg-gray-400 transition-colors"
                >
                  Cancel
                </button>
              </div>
            </div>
          </div>

          <!-- Lists -->
          <ul class="space-y-2">
            <li v-for="list in listsWithCounts" :key="list.id">
              <RouterLink
                :to="`/list/${list.id}`" 
                class="flex items-center justify-between p-2 rounded-md hover:bg-blue-50 hover:text-blue-600 transition-colors"
                :title="isCollapsed ? list.name : ''"
              >
                <div class="flex items-center">
                  <span class="text-sm mr-2"></span>
                  <span v-if="!isCollapsed" class="text-sm font-medium">{{ list.name }}</span>
                </div>
                <span v-if="!isCollapsed && list.count > 0" class="text-xs bg-gray-200 text-gray-600 px-2 py-1 rounded-full">
                  {{ list.count }}
                </span>
              </RouterLink>
            </li>
          </ul>
          
          <!-- Empty State -->
          <div v-if="!loading && listsWithCounts.length === 0" class="text-center py-4">
            <p class="text-xs text-gray-500">No lists yet</p>
            <p class="text-xs text-gray-400">Create your first list!</p>
          </div>
        </div>
      </div>
    </div>

    <!-- Main Content Area -->
    <main class="flex-1 overflow-auto">
      <slot />
    </main>
  </div>
</template>

<style scoped>
/* Additional custom styles if needed */
</style>

