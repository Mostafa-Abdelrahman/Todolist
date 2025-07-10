<script setup>
import { ref, computed } from 'vue'

// Example lists (could be imported or shared via store)
const lists = ref([
  { name: 'Work Projects', id: 1 },
  { name: 'Personal Goals', id: 2 },
  { name: 'Shopping List', id: 3 }
])

// Example tasks (could be imported or shared via store)
const tasks = ref([
  { title: 'Finish report', listId: 1, due: new Date().toISOString().slice(0, 10), done: false },
  { title: 'Buy groceries', listId: 3, due: new Date().toISOString().slice(0, 10), done: false },
  { title: 'Go for a run', listId: 2, due: '2024-06-10', done: false }
])

const today = new Date().toISOString().slice(0, 10)
const todayDisplay = new Date().toLocaleDateString(undefined, { weekday: 'long', year: 'numeric', month: 'long', day: 'numeric' })

// Computed: tasks for today
const todayTasks = computed(() => tasks.value.filter(t => t.due === today))

// Add new task for today
const newTaskTitle = ref('')
const newTaskList = ref(lists.value[0]?.id || null)
function addTask() {
  if (!newTaskTitle.value || !newTaskList.value) return
  tasks.value.push({
    title: newTaskTitle.value,
    listId: newTaskList.value,
    due: today,
    done: false
  })
  newTaskTitle.value = ''
  newTaskList.value = lists.value[0]?.id || null
}

// Edit logic
const editingIndex = ref(null)
const editingTitle = ref('')
function startEdit(idx, task) {
  editingIndex.value = idx
  editingTitle.value = task.title
}
function saveEdit(idx) {
  if (editingTitle.value.trim()) {
    // Find the task in the main tasks array and update
    const task = todayTasks.value[idx]
    const mainIdx = tasks.value.findIndex(t => t === task)
    if (mainIdx !== -1) tasks.value[mainIdx].title = editingTitle.value
  }
  editingIndex.value = null
  editingTitle.value = ''
}
function cancelEdit() {
  editingIndex.value = null
  editingTitle.value = ''
}

// Toggle done (fix: update in main tasks array)
function toggleDone(idx) {
  const task = todayTasks.value[idx]
  const mainIdx = tasks.value.findIndex(t => t === task)
  if (mainIdx !== -1) tasks.value[mainIdx].done = !tasks.value[mainIdx].done
}
</script>

<template>
  <div class="p-6 md:p-10 bg-gray-50 min-h-screen">
    <div class="flex flex-col md:flex-row md:items-center md:justify-between mb-2 gap-2">
      <h1 class="text-2xl md:text-3xl font-bold text-gray-800">Today's Tasks</h1>
      <div class="w-full md:w-auto text-right md:text-right">
        <span class="text-xl md:text-2xl font-bold text-blue-700">{{ todayDisplay }}</span>
      </div>
    </div>

    <!-- Add Task for Today -->
    <div class="bg-white rounded-xl shadow-sm mb-8 w-full px-4 py-5 md:px-8 flex flex-col md:flex-row md:items-end gap-4">
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
    <div class="bg-white rounded-xl shadow-sm p-6 w-full">
      <h2 class="text-lg font-semibold text-gray-700 mb-4">Tasks List</h2>
      <ul class="space-y-3">
        <li v-for="(task, idx) in todayTasks" :key="task.title + task.listId" class="flex items-center gap-3">
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
