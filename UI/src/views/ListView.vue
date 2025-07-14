<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import apiService from '@/services/api.js'

// Reactive data
const lists = ref([])
const tasks = ref([])
const loading = ref(false)
const error = ref(null)

const route = useRoute()
const listId = computed(() => {
  // Support both /list/1 and /list/Work%20Projects
  const param = route.params.id
  const found = lists.value.find(
    l => l.id === Number(param) || l.name === param || l.name.replace(/ /g, '%20') === param
  )
  return found ? found.id : null
})

const listName = computed(() => {
  const found = lists.value.find(l => l.id === listId.value)
  return found ? found.name : 'List'
})

// Filter tasks for this list
const listTasks = computed(() => tasks.value.filter(t => t.listId === listId.value))

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

// Add new task
const newTaskTitle = ref('')
const newTaskDue = ref('')

async function addTask() {
  if (!newTaskTitle.value || !newTaskDue.value || !listId.value) return
  
  try {
    const newTask = {
      title: newTaskTitle.value,
      listId: listId.value,
      due: newTaskDue.value,
      done: false
    }
    
    const createdTask = await apiService.createTodo(newTask)
    tasks.value.push(createdTask)
    
    // Reset form
    newTaskTitle.value = ''
    newTaskDue.value = ''
  } catch (err) {
    error.value = err.message
    console.error('Failed to create task:', err)
  }
}

// Edit logic
const editingIndex = ref(null)
const editingTitle = ref('')

function startEdit(idx, task) {
  editingIndex.value = idx
  editingTitle.value = task.title
}

async function saveEdit(idx) {
  if (editingTitle.value.trim()) {
    const task = listTasks.value[idx]
    try {
      const updatedTask = await apiService.updateTodo(task.id, {
        ...task,
        title: editingTitle.value
      })
      
      // Update in local state
      const taskIndex = tasks.value.findIndex(t => t.id === task.id)
      if (taskIndex !== -1) {
        tasks.value[taskIndex] = updatedTask
      }
    } catch (err) {
      error.value = err.message
      console.error('Failed to update task:', err)
      return
    }
  }
  editingIndex.value = null
  editingTitle.value = ''
}

function cancelEdit() {
  editingIndex.value = null
  editingTitle.value = ''
}

// Toggle done
async function toggleDone(idx) {
  const task = listTasks.value[idx]
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
    <h1 class="text-2xl md:text-3xl font-bold text-gray-800 mb-6">{{ listName }}</h1>

    <!-- Error Message -->
    <div v-if="error" class="bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded mb-6">
      {{ error }}
    </div>

    <!-- Loading State -->
    <div v-if="loading" class="text-center py-8">
      <div class="inline-block animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600"></div>
      <p class="mt-2 text-gray-600">Loading...</p>
    </div>

    <!-- Add Task Form -->
    <div v-if="!loading && listId" class="bg-white rounded-xl shadow-sm mb-8 w-full px-4 py-5 md:px-8 flex flex-col md:flex-row md:items-end gap-4">
      <form @submit.prevent="addTask" class="flex flex-col md:flex-row w-full gap-3 md:gap-4">
        <input
          v-model="newTaskTitle"
          type="text"
          placeholder="New Task"
          class="border border-gray-200 rounded px-3 py-2 flex-1 focus:ring-2 focus:ring-blue-200 focus:outline-none"
          required
        />
        <input
          v-model="newTaskDue"
          type="date"
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

    <!-- List Not Found -->
    <div v-if="!loading && !listId" class="bg-white rounded-xl shadow-sm p-8 text-center">
      <div class="text-gray-400 mb-4">
        <svg class="w-16 h-16 mx-auto mb-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"></path>
        </svg>
      </div>
      <h3 class="text-lg font-semibold text-gray-600 mb-2">List not found</h3>
      <p class="text-gray-500">The list you're looking for doesn't exist.</p>
    </div>

    <!-- List Tasks -->
    <div v-if="!loading && listId" class="bg-white rounded-xl shadow-sm p-6 w-full">
      <h2 class="text-lg font-semibold text-gray-700 mb-4">Tasks ({{ listTasks.length }})</h2>
      <ul class="space-y-3">
        <li v-for="(task, idx) in listTasks" :key="task.id" class="flex items-center gap-3">
          <input type="checkbox" :checked="task.done" @change="toggleDone(idx)" class="h-5 w-5 text-blue-600 rounded border-gray-300 focus:ring-blue-500" />
          <span class="text-xs bg-gray-100 text-gray-600 px-2 py-1 rounded">{{ task.due }}</span>
          <template v-if="editingIndex === idx">
            <input v-model="editingTitle" @keyup.enter="saveEdit(idx)" @keyup.esc="cancelEdit" class="border border-gray-300 rounded px-2 py-1 flex-1" />
            <button @click="saveEdit(idx)" class="text-green-600 hover:underline text-xs ml-2">Save</button>
            <button @click="cancelEdit" class="text-gray-400 hover:underline text-xs ml-1">Cancel</button>
          </template>
          <template v-else>
            <span :class="['flex-1', 'truncate', task.done ? 'line-through text-gray-400' : 'text-gray-800']">{{ task.title }}</span>
            <button @click="startEdit(idx, task)" class="text-blue-500 hover:underline text-xs ml-2">Edit</button>
          </template>
        </li>
        <li v-if="listTasks.length === 0" class="text-gray-400 px-3 py-2 text-center">
          No tasks in this list.
        </li>
      </ul>
    </div>
  </div>
</template>