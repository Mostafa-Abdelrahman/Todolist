<script setup>
import { ref, nextTick, watch } from 'vue'

// Reactive state
const isCollapsed = ref(false)
const showNewListInput = ref(false)
const newListName = ref('')
const newListInput = ref(null)

const lists = ref([
  { name: 'Work Projects', id: 1, count: 8 },
  { name: 'Personal Goals', id: 2, count: 4 },
  { name: 'Shopping List', id: 3, count: 6 }
])

// Methods
const toggleSidebar = () => {
  isCollapsed.value = !isCollapsed.value
}

const addNewList = () => {
  if (newListName.value.trim()) {
    lists.value.push({
      name: newListName.value,
      count: 0
    })
    newListName.value = ''
    showNewListInput.value = false
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
              class="text-blue-600 hover:text-blue-700 text-sm"
            >
              + New
            </button>
          </div>

          <!-- New List Input -->
          <div v-if="showNewListInput && !isCollapsed" class="mb-3">
            <div class="flex flex-col space-y-2">
              <input 
                v-model="newListName"
                @keyup.enter="addNewList"
                @keyup.esc="cancelNewList"
                type="text"
                placeholder="List name"
                class="w-full px-3 py-2 text-sm border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                ref="newListInput"
              />
              <div class="flex space-x-2">
                <button 
                  @click="addNewList"
                  class="flex-1 px-3 py-1 text-xs bg-blue-500 text-white rounded hover:bg-blue-600 transition-colors"
                >
                  Add
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
            <li v-for="list in lists" :key="list.name">
              <RouterLink
                :to="`/list/${list.id}`" 
                class="flex items-center justify-between p-2 rounded-md hover:bg-blue-50 hover:text-blue-600 transition-colors"
                :title="isCollapsed ? list.name : ''"
              >
                <div class="flex items-center">
                  <span class="text-sm mr-2"></span>
                  <span v-if="!isCollapsed" class="text-sm">{{ list.name }}</span>
                </div>
                <span v-if="!isCollapsed && list.count > 0" class="text-xs bg-gray-200 text-gray-600 px-2 py-1 rounded-full">
                  {{ list.count }}
                </span>
              </RouterLink>
            </li>
          </ul>
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

