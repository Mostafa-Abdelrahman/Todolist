<script setup>
import { ref } from 'vue'

// Example data (replace with real data or fetch from backend/store)
const lists = ref([
  { name: 'Work Projects', id: 1 },
  { name: 'Personal Goals', id: 2 },
  { name: 'Shopping List', id: 3 }
])

const tasks = ref([
  { title: 'Finish report', listId: 1, due: '2024-06-10' },
  { title: 'Buy groceries', listId: 3, due: '2024-06-11' },
  { title: 'Go for a run', listId: 2, due: '2024-06-10' }
])

// Form state
const newTaskTitle = ref('')
const newTaskList = ref(lists.value[0]?.id || null)
const newTaskDue = ref('')

// Add new task
function addTask() {
  if (!newTaskTitle.value || !newTaskList.value || !newTaskDue.value) return
  tasks.value.push({
    title: newTaskTitle.value,
    listId: newTaskList.value,
    due: newTaskDue.value
  })
  newTaskTitle.value = ''
  newTaskList.value = lists.value[0]?.id || null
  newTaskDue.value = ''
}
</script>

<template>
  <div class="p-6 md:p-10 bg-gray-50 min-h-screen">
    <h1 class="text-2xl md:text-3xl font-bold text-gray-800 mb-6">Todo List Dashboard</h1>

    <!-- Add Task Card (Full Width, Simple) -->
    <div class="bg-white rounded-xl shadow-sm mb-8 w-full max-w-full px-4 py-5 md:px-8 flex flex-col md:flex-row md:items-end gap-4">
      <form @submit.prevent="addTask" class="flex flex-col md:flex-row w-full gap-3 md:gap-4">
        <input
          v-model="newTaskTitle"
          type="text"
          placeholder="New Task "
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

    <!-- Lists as Cards -->
    <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
      <div
        v-for="list in lists"
        :key="list.id"
        class="bg-white rounded-xl shadow-sm p-5 flex flex-col"
      >
        <div class="flex items-center justify-between mb-3">
          <h3 class="text-lg font-bold text-gray-700">{{ list.name }}</h3>
          <span class="text-xs bg-blue-100 text-blue-700 px-2 py-1 rounded-full">{{ tasks.filter(t => t.listId === list.id).length }} tasks</span>
        </div>
        <ul class="flex-1 space-y-2">
          <li
            v-for="task in tasks.filter(t => t.listId === list.id)"
            :key="task.title + task.due"
            class="flex justify-between items-center bg-gray-50 rounded px-3 py-2"
          >
            <span class="truncate">{{ task.title }}</span>
            <span class="text-xs text-gray-500 ml-2">{{ task.due }}</span>
          </li>
          <li v-if="tasks.filter(t => t.listId === list.id).length === 0" class="text-gray-400 px-3 py-2">
            No tasks in this list.
          </li>
        </ul>
      </div>
    </div>
  </div>
</template>