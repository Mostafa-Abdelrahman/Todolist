<script setup>
import { ref, computed } from 'vue'

const lists = ref([
  { name: 'Work Projects', id: 1 },
  { name: 'Personal Goals', id: 2 },
  { name: 'Shopping List', id: 3 }
])

const tasks = ref([
  { title: 'Finish quarterly report', listId: 1, due: '2026-06-15', done: false },
  { title: 'Team meeting', listId: 1, due: '2026-06-15', done: false },
  { title: 'Buy birthday gift', listId: 3, due: '2026-06-18', done: false },
  { title: 'Gym session', listId: 2, due: '2026-06-20', done: false },
  { title: 'Project presentation', listId: 1, due: '2026-06-25', done: false },
  { title: 'Read new book', listId: 2, due: '2026-06-30', done: false }
])

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

const newTaskTitle = ref('')
const newTaskList = ref(lists.value[0]?.id || null)
const newTaskDue = ref('')

function addTask() {
  if (!newTaskTitle.value || !newTaskList.value || !newTaskDue.value) return
  if (newTaskDue.value < today) return 
  
  tasks.value.push({
    title: newTaskTitle.value,
    listId: newTaskList.value,
    due: newTaskDue.value,
    done: false
  })
  newTaskTitle.value = ''
  newTaskList.value = lists.value[0]?.id || null
  newTaskDue.value = ''
}

const editingTask = ref(null)
const editingTitle = ref('')

function startEdit(task) {
  editingTask.value = task
  editingTitle.value = task.title
}

function saveEdit() {
  if (editingTitle.value.trim() && editingTask.value) {
    editingTask.value.title = editingTitle.value
  }
  editingTask.value = null
  editingTitle.value = ''
}

function cancelEdit() {
  editingTask.value = null
  editingTitle.value = ''
}

// Toggle done
function toggleDone(task) {
  task.done = !task.done
}
</script>

<template>
  <div class="p-6 md:p-10 bg-gray-50 min-h-screen">
    <div class="flex flex-col md:flex-row md:items-center md:justify-between mb-6 gap-2">
      <h1 class="text-2xl md:text-3xl font-bold text-gray-800">Upcoming Tasks</h1>
      <div class="w-full md:w-auto text-right md:text-right">
        <span class="text-lg md:text-xl font-bold text-blue-700">{{ upcomingTasks.length }} upcoming tasks</span>
      </div>
    </div>

    <!-- Add Task for Upcoming -->
    <div class="bg-white rounded-xl shadow-sm mb-8 w-full px-4 py-5 md:px-8 flex flex-col md:flex-row md:items-end gap-4">
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
    <div v-if="Object.keys(groupedTasks).length > 0" class="space-y-6">
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
              :key="task.title + task.listId"
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
    <div v-else class="bg-white rounded-xl shadow-sm p-8 text-center">
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
