<script setup>
import { ref, computed, onMounted } from 'vue'
import apiService from '@/services/api.js'

// Reactive data
const lists = ref([])
const tasks = ref([])
const loading = ref(false)
const error = ref(null)

const today = new Date().toISOString().slice(0, 10)

const upcomingTasks = computed(() => 
  tasks.value.filter(t => t.due >= today).sort((a, b) => a.due.localeCompare(b.due))
)

const groupedTasks = computed(() => {
  const groups = {}
  upcomingTasks.value.forEach(task => {
    if (!groups[task.due]) {
      groups[task.due] = []
    }
    groups[task.due].push(task)
  })
  return groups
})

function formatDate(dateString) {
  const date = new Date(dateString)
  const today = new Date()
  const tomorrow = new Date(today)
  tomorrow.setDate(tomorrow.getDate() + 1)
  
  if (dateString === tomorrow.toISOString().slice(0, 10)) {
    return 'Tomorrow'
  }
  
  return date.toLocaleDateString(undefined, { 
    weekday: 'long', 
    month: 'short', 
    day: 'numeric' 
  })
}

// Fetch data from API
async function fetchData() {
  loading.value = true
  error.value = null
  
  try {
    const [listsData, tasksData] = await Promise.all([
      apiService.getLists(),
      apiService.getTodos()
    ])
    
    lists.value = listsData
    tasks.value = tasksData
  } catch (err) {
    error.value = err.message
    console.error('Failed to fetch data:', err)
  } finally {
    loading.value = false
  }
}

const newTaskTitle = ref('')
const newTaskList = ref(null)
const newTaskDue = ref('')

async function addTask() {
  if (!newTaskTitle.value || !newTaskList.value || !newTaskDue.value) return
  if (newTaskDue.value < today) return 
  
  try {
    const newTask = {
      title: newTaskTitle.value,
      listId: newTaskList.value,
      due: newTaskDue.value,
      done: false
    }
    
    const createdTask = await apiService.createTodo(newTask)
    tasks.value.push(createdTask)
    
    // Reset form
    newTaskTitle.value = ''
    newTaskDue.value = ''
    newTaskList.value = lists.value[0]?.id || null
  } catch (err) {
    error.value = err.message
    console.error('Failed to create task:', err)
  }
}

const editingTask = ref(null)
const editingTitle = ref('')

function startEdit(task) {
  editingTask.value = task
  editingTitle.value = task.title
}

async function saveEdit() {
  if (editingTitle.value.trim() && editingTask.value) {
    try {
      const updatedTask = await apiService.updateTodo(editingTask.value.id, {
        ...editingTask.value,
        title: editingTitle.value
      })
      
      // Update in local state
      const taskIndex = tasks.value.findIndex(t => t.id === editingTask.value.id)
      if (taskIndex !== -1) {
        tasks.value[taskIndex] = updatedTask
      }
    } catch (err) {
      error.value = err.message
      console.error('Failed to update task:', err)
      return
    }
  }
  editingTask.value = null
  editingTitle.value = ''
}

function cancelEdit() {
  editingTask.value = null
  editingTitle.value = ''
}

// Toggle done
async function toggleDone(task) {
  try {
    const updatedTask = await apiService.updateTodo(task.id, {
      ...task,
      done: !task.done
    })
    
    // Update in local state
    const taskIndex = tasks.value.findIndex(t => t.id === task.id)
    if (taskIndex !== -1) {
      tasks.value[taskIndex] = updatedTask
    }
  } catch (err) {
    error.value = err.message
    console.error('Failed to update task:', err)
  }
}

// Load data on component mount
onMounted(() => {
  fetchData()
})
</script>

<template>
  <div class="p-6 md:p-10 bg-gray-50 min-h-screen">
    <div class="flex flex-col md:flex-row md:items-center md:justify-between mb-6 gap-2">
      <h1 class="text-2xl md:text-3xl font-bold text-gray-800">Upcoming Tasks</h1>
      <div class="w-full md:w-auto text-right md:text-right">
        <span class="text-lg md:text-xl font-bold text-blue-700">{{ upcomingTasks.length }} upcoming tasks</span>
      </div>
    </div>

    <!-- Error Message -->
    <div v-if="error" class="bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded mb-6">
      {{ error }}
    </div>

    <!-- Loading State -->
    <div v-if="loading" class="text-center py-8">
      <div class="inline-block animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600"></div>
      <p class="mt-2 text-gray-600">Loading...</p>
    </div>

    <!-- Add Task for Upcoming -->
    <div v-if="!loading" class="bg-white rounded-xl shadow-sm mb-8 w-full px-4 py-5 md:px-8 flex flex-col md:flex-row md:items-end gap-4">
      <form @submit.prevent="addTask" class="flex flex-col md:flex-row w-full gap-3 md:gap-4">
        <input
          v-model="newTaskTitle"
          type="text"
          placeholder="New upcoming task"
          class="border border-gray-200 rounded px-3 py-2 flex-1 focus:ring-2 focus:ring-blue-200 focus:outline-none"
          required
        />
        <select v-model="newTaskList" class="border border-gray-200 rounded px-3 py-2 flex-1 focus:ring-2 focus:ring-blue-200 focus:outline-none" required>
          <option v-for="list in lists" :key="list.id" :value="list.id">
            {{ list.name }}
          </option>
        </select>
        <input
          v-model="newTaskDue"
          type="date"
          :min="new Date().toISOString().slice(0, 10)"
          class="border border-gray-200 rounded px-3 py-2 flex-1 focus:ring-2 focus:ring-blue-200 focus:outline-none"
          required
        />
        <button
          type="submit"
          class="bg-blue-600 text-white px-4 py-2 rounded-lg hover:bg-blue-700 transition font-semibold shadow-sm min-w-[120px]"
        >
          Add Task
        </button>
      </form>
    </div>

    <!-- Upcoming Tasks by Date -->
    <div v-if="!loading && Object.keys(groupedTasks).length > 0" class="space-y-6">
      <div 
        v-for="(tasksForDate, date) in groupedTasks" 
        :key="date"
        class="bg-white rounded-xl shadow-sm overflow-hidden"
      >
        <!-- Date Header -->
        <div class="bg-gradient-to-r from-blue-50 to-blue-100 px-6 py-4 border-b border-gray-200">
          <h2 class="text-lg font-semibold text-blue-800">{{ formatDate(date) }}</h2>
          <p class="text-sm text-blue-600">{{ tasksForDate.length }} task{{ tasksForDate.length !== 1 ? 's' : '' }}</p>
        </div>
        
        <!-- Tasks for this date -->
        <div class="p-6">
          <ul class="space-y-3">
            <li 
              v-for="task in tasksForDate" 
              :key="task.id"
              class="flex items-center gap-3 p-3 bg-gray-50 rounded-lg hover:bg-gray-100 transition-colors"
            >
              <input 
                type="checkbox" 
                :checked="task.done" 
                @change="toggleDone(task)" 
                class="h-5 w-5 text-blue-600 rounded border-gray-300 focus:ring-blue-500" 
              />
              <span class="text-xs bg-blue-100 text-blue-700 px-2 py-1 rounded-full">
                {{ lists.find(l => l.id === task.listId)?.name }}
              </span>
              
              <!-- Edit Mode -->
              <template v-if="editingTask === task">
                <input 
                  v-model="editingTitle" 
                  @keyup.enter="saveEdit" 
                  @keyup.esc="cancelEdit" 
                  class="border border-gray-300 rounded px-2 py-1 flex-1" 
                />
                <button @click="saveEdit" class="text-green-600 hover:underline text-xs ml-2">Save</button>
                <button @click="cancelEdit" class="text-gray-400 hover:underline text-xs ml-1">Cancel</button>
              </template>
              
              <!-- View Mode -->
              <template v-else>
                <span :class="['flex-1', 'truncate', task.done ? 'line-through text-gray-400' : 'text-gray-800']">
                  {{ task.title }}
                </span>
                <button @click="startEdit(task)" class="text-blue-500 hover:underline text-xs ml-2">Edit</button>
              </template>
            </li>
          </ul>
        </div>
      </div>
    </div>

    <!-- Empty State -->
    <div v-if="!loading && Object.keys(groupedTasks).length === 0" class="bg-white rounded-xl shadow-sm p-8 text-center">
      <div class="text-gray-400 mb-4">
        <svg class="w-16 h-16 mx-auto mb-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z"></path>
        </svg>
      </div>
      <h3 class="text-lg font-semibold text-gray-600 mb-2">No upcoming tasks</h3>
      <p class="text-gray-500">Add some tasks with future dates to see them here.</p>
    </div>
  </div>
</template>
