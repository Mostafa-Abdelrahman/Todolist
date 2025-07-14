<script setup>
import { ref, onMounted } from 'vue'
import apiService from '@/services/api.js'

// Reactive data
const lists = ref([])
const tasks = ref([])
const loading = ref(false)
const error = ref(null)

// Form state
const newTaskTitle = ref('')
const newTaskList = ref(null)
const newTaskDue = ref('')

// List management state
const showNewListForm = ref(false)
const newListName = ref('')
const editingListId = ref(null)
const editingListName = ref('')

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

// Add new task
async function addTask() {
  if (!newTaskTitle.value || !newTaskList.value || !newTaskDue.value) return
  
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

// Create new list
async function createList() {
  if (!newListName.value.trim()) return
  
  try {
    const newList = {
      name: newListName.value.trim()
    }
    
    const createdList = await apiService.createList(newList)
    lists.value.push(createdList)
    
    // Set as default if it's the first list
    if (lists.value.length === 1) {
      newTaskList.value = createdList.id
    }
    
    // Reset form
    newListName.value = ''
    showNewListForm.value = false
  } catch (err) {
    error.value = err.message
    console.error('Failed to create list:', err)
  }
}

// Start editing list
function startEditList(list) {
  editingListId.value = list.id
  editingListName.value = list.name
}

// Save list edit
async function saveListEdit() {
  if (!editingListName.value.trim()) return
  
  try {
    const updatedList = await apiService.updateList(editingListId.value, {
      name: editingListName.value.trim()
    })
    
    // Update in local state
    const listIndex = lists.value.findIndex(l => l.id === editingListId.value)
    if (listIndex !== -1) {
      lists.value[listIndex] = updatedList
    }
    
    // Reset edit state
    editingListId.value = null
    editingListName.value = ''
  } catch (err) {
    error.value = err.message
    console.error('Failed to update list:', err)
  }
}

// Cancel list edit
function cancelListEdit() {
  editingListId.value = null
  editingListName.value = ''
}

// Delete list
async function deleteList(listId) {
  if (!confirm('Are you sure you want to delete this list? All tasks in this list will also be deleted.')) return
  
  try {
    await apiService.deleteList(listId)
    
    // Remove from local state
    lists.value = lists.value.filter(l => l.id !== listId)
    
    // Update task list to remove tasks from deleted list
    tasks.value = tasks.value.filter(t => t.listId !== listId)
    
    // Update default list if needed
    if (newTaskList.value === listId) {
      newTaskList.value = lists.value[0]?.id || null
    }
  } catch (err) {
    error.value = err.message
    console.error('Failed to delete list:', err)
  }
}

// Load data on component mount
onMounted(() => {
  fetchData()
})
</script>

<template>
  <div class="p-6 md:p-10 bg-gray-50 min-h-screen">
    <h1 class="text-2xl md:text-3xl font-bold text-gray-800 mb-6">Todo List Dashboard</h1>

    <!-- Error Message -->
    <div v-if="error" class="bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded mb-6">
      {{ error }}
    </div>

    <!-- Loading State -->
    <div v-if="loading" class="text-center py-8">
      <div class="inline-block animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600"></div>
      <p class="mt-2 text-gray-600">Loading...</p>
    </div>

    <!-- Add Task Card (Full Width, Simple) -->
    <div v-if="!loading" class="bg-white rounded-xl shadow-sm mb-8 w-full max-w-full px-4 py-5 md:px-8 flex flex-col md:flex-row md:items-end gap-4">
      <form @submit.prevent="addTask" class="flex flex-col md:flex-row w-full gap-3 md:gap-4">
        <input
          v-model="newTaskTitle"
          type="text"
          placeholder="New Task"
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

    <!-- Lists Management Section -->
    <div v-if="!loading" class="mb-8">
      <div class="flex items-center justify-between mb-4">
        <h2 class="text-xl font-semibold text-gray-800">Lists Management</h2>
        <button
          @click="showNewListForm = !showNewListForm"
          class="bg-green-600 text-white px-4 py-2 rounded-lg hover:bg-green-700 transition font-semibold shadow-sm"
        >
          {{ showNewListForm ? 'Cancel' : 'Create New List' }}
        </button>
      </div>

      <!-- New List Form -->
      <div v-if="showNewListForm" class="bg-white rounded-xl shadow-sm p-6 mb-6">
        <h3 class="text-lg font-semibold text-gray-700 mb-4">Create New List</h3>
        <form @submit.prevent="createList" class="flex flex-col md:flex-row gap-4">
          <input
            v-model="newListName"
            type="text"
            placeholder="Enter list name"
            class="border border-gray-200 rounded px-3 py-2 flex-1 focus:ring-2 focus:ring-blue-200 focus:outline-none"
            required
          />
          <button
            type="submit"
            class="bg-blue-600 text-white px-6 py-2 rounded-lg hover:bg-blue-700 transition font-semibold shadow-sm"
          >
            Create List
          </button>
        </form>
      </div>

      <!-- Lists Grid -->
      <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
        
        <div
          v-for="list in lists"
          :key="list.id"
          class="bg-white rounded-xl shadow-sm p-5 flex flex-col"
        >
          <!-- List Header -->
          <div class="flex items-center justify-between mb-3">
            <div class="flex-1">
              <!-- Edit Mode -->
              <div v-if="editingListId === list.id" class="flex items-center gap-2">
                <input
                  v-model="editingListName"
                  @keyup.enter="saveListEdit"
                  @keyup.esc="cancelListEdit"
                  class="border border-gray-300 rounded px-2 py-1 flex-1 text-sm"
                  ref="editListInput"
                />
                <button @click="saveListEdit" class="text-green-600 hover:text-green-700 text-sm">
                  <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7"></path>
                  </svg>
                </button>
                <button @click="cancelListEdit" class="text-gray-400 hover:text-gray-600 text-sm">
                  <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"></path>
                  </svg>
                </button>
              </div>
              
              <!-- View Mode -->
              <div v-else class="flex items-center justify-between">
                <h3 class="text-lg font-bold text-gray-700">{{ list.name }}</h3>
                <div class="flex items-center gap-2">
                  <button @click="startEditList(list)" class="text-blue-500 hover:text-blue-700 text-sm">
                    <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z"></path>
                    </svg>
                  </button>
                  <button @click="deleteList(list.id)" class="text-red-500 hover:text-red-700 text-sm">
                    <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"></path>
                    </svg>
                  </button>
                </div>
              </div>
            </div>
          </div>
          
          <!-- Task Count -->
          <div class="mb-3">
            <span class="text-xs bg-blue-100 text-blue-700 px-2 py-1 rounded-full">
              {{ tasks.filter(t => t.listId === list.id).length }} tasks
            </span>
          </div>
          
          <!-- Tasks List -->
          <ul class="flex-1 space-y-2">
            <li
              v-for="task in tasks.filter(t => t.listId === list.id)"
              :key="task.id"
              class="flex justify-between items-center bg-gray-50 rounded px-3 py-2"
            >
              <span class="truncate">{{ task.title }}</span>
              <span class="text-xs text-gray-500 ml-2">{{ task.due }}</span>
            </li>
            <li v-if="tasks.filter(t => t.listId === list.id).length === 0" class="text-gray-400 px-3 py-2 text-center text-sm">
              No tasks in this list.
            </li>
          </ul>
        </div>
      </div>
    </div>
  </div>
</template>