<script setup>
import { ref, computed, onMounted } from 'vue'
import apiService from '@/services/api.js'

// Reactive data
const lists = ref([])
const tasks = ref([])
const loading = ref(false)
const error = ref(null)

const today = new Date().toISOString().slice(0, 10)
const todayDisplay = new Date().toLocaleDateString(undefined, { weekday: 'long', year: 'numeric', month: 'long', day: 'numeric' })

// Computed: tasks for today
const todayTasks = computed(() => tasks.value.filter(t => t.due === today))

// Add new task for today
const newTaskTitle = ref('')
const newTaskList = ref(null)

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
    
    // Set default list for new task form
    if (lists.value.length > 0 && !newTaskList.value) {
      newTaskList.value = lists.value[0].id
    }
  } catch (err) {
    error.value = err.message
    console.error('Failed to fetch data:', err)
  } finally {
    loading.value = false
  }
}

async function addTask() {
  if (!newTaskTitle.value || !newTaskList.value) return
  
  try {
    const newTask = {
      title: newTaskTitle.value,
      listId: newTaskList.value,
      due: today,
      done: false
    }
    
    const createdTask = await apiService.createTodo(newTask)
    tasks.value.push(createdTask)
    
    // Reset form
    newTaskTitle.value = ''
    newTaskList.value = lists.value[0]?.id || null
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
    const task = todayTasks.value[idx]
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
  const task = todayTasks.value[idx]
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
    <div class="flex flex-col md:flex-row md:items-center md:justify-between mb-2 gap-2">
      <h1 class="text-2xl md:text-3xl font-bold text-gray-800">Today's Tasks</h1>
      <div class="w-full md:w-auto text-right md:text-right">
        <span class="text-xl md:text-2xl font-bold text-blue-700">{{ todayDisplay }}</span>
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

    <!-- Add Task for Today -->
    <div v-if="!loading" class="bg-white rounded-xl shadow-sm mb-8 w-full px-4 py-5 md:px-8 flex flex-col md:flex-row md:items-end gap-4">
      <form @submit.prevent="addTask" class="flex flex-col md:flex-row w-full gap-3 md:gap-4">
        <input
          v-model="newTaskTitle"
          type="text"
          placeholder="New Task for Today"
          class="border border-gray-200 rounded px-3 py-2 flex-1 focus:ring-2 focus:ring-blue-200 focus:outline-none"
          required
        />
        <select v-model="newTaskList" class="border border-gray-200 rounded px-3 py-2 flex-1 focus:ring-2 focus:ring-blue-200 focus:outline-none" required>
          <option v-for="list in lists" :key="list.id" :value="list.id">
            {{ list.name }}
          </option>
        </select>
        <button
          type="submit"
          class="bg-blue-600 text-white px-4 py-2 rounded-lg hover:bg-blue-700 transition font-semibold shadow-sm min-w-[120px]"
        >
          Add Task
        </button>
      </form>
    </div>

    <!-- Today's Tasks List -->
    <div v-if="!loading" class="bg-white rounded-xl shadow-sm p-6 w-full">
      <h2 class="text-lg font-semibold text-gray-700 mb-4">Tasks List</h2>
      <ul class="space-y-3">
        <li v-for="(task, idx) in todayTasks" :key="task.id" class="flex items-center gap-3">
          <input type="checkbox" :checked="task.done" @change="toggleDone(idx)" class="h-5 w-5 text-blue-600 rounded border-gray-300 focus:ring-blue-500" />
          <span class="text-xs bg-gray-100 text-gray-600 px-2 py-1 rounded" v-if="lists.find(l => l.id === task.listId)">
            {{ lists.find(l => l.id === task.listId).name }}
          </span>
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
        <li v-if="todayTasks.length === 0" class="text-gray-400 px-3 py-2 text-center">
          No tasks for today.
        </li>
      </ul>
    </div>
  </div>
</template>
