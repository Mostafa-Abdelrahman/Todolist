<script setup>
import { ref, computed } from 'vue'
import { useRoute } from 'vue-router'

// Example lists and tasks (replace with global store or props as needed)
const lists = ref([
  { name: 'Work Projects', id: 1 },
  { name: 'Personal Goals', id: 2 },
  { name: 'Shopping List', id: 3 }
])

const tasks = ref([
  { title: 'Finish report', listId: 1, due: '2024-06-10', done: false },
  { title: 'Buy groceries', listId: 3, due: '2024-06-11', done: false },
  { title: 'Go for a run', listId: 2, due: '2024-06-10', done: false }
])

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

// Add new task
const newTaskTitle = ref('')
const newTaskDue = ref('')
function addTask() {
  if (!newTaskTitle.value || !newTaskDue.value) return
  tasks.value.push({
    title: newTaskTitle.value,
    listId: listId.value,
    due: newTaskDue.value,
    done: false
  })
  newTaskTitle.value = ''
  newTaskDue.value = ''
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
    listTasks.value[idx].title = editingTitle.value
  }
  editingIndex.value = null
  editingTitle.value = ''
}
function cancelEdit() {
  editingIndex.value = null
  editingTitle.value = ''
}

// Toggle done
function toggleDone(idx) {
  listTasks.value[idx].done = !listTasks.value[idx].done
}
</script>

<template>
  <div class="p-6 md:p-10 bg-gray-50 min-h-screen">
    <h1 class="text-2xl md:text-3xl font-bold text-gray-800 mb-6">{{ listName }}</h1>

    <!-- Add Task Form -->
    <div class="bg-white rounded-xl shadow-sm mb-8 w-full px-4 py-5 md:px-8 flex flex-col md:flex-row md:items-end gap-4">
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

    <!-- List Tasks -->
    <div class="bg-white rounded-xl shadow-sm p-6 w-full">
      <h2 class="text-lg font-semibold text-gray-700 mb-4">Tasks</h2>
      <ul class="space-y-3">
        <li v-for="(task, idx) in listTasks" :key="task.title + task.due" class="flex items-center gap-3">
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