const API_BASE_URL = 'http://localhost:8080/api'

class ApiService {
  constructor() {
    this.baseURL = API_BASE_URL
  }

  // Generic request method
  async request(endpoint, options = {}) {
    const url = `${this.baseURL}${endpoint}`
    const config = {
      headers: {
        'Content-Type': 'application/json',
        ...options.headers
      },
      ...options
    }

    try {
      const response = await fetch(url, config)
      
      console.log(`API Response: ${response.status} ${response.statusText}`)
      
      if (!response.ok) {
        const errorData = await response.json().catch(() => ({}))
        console.error('API Error:', errorData)
        throw new Error(errorData.error || `HTTP error! status: ${response.status}`)
      }

      // Handle empty responses
      const contentType = response.headers.get('content-type')
      if (contentType && contentType.includes('application/json')) {
        const data = await response.json()
        console.log('API Response data:', data)
        return data
      }
      
      return null
    } catch (error) {
      console.error('API request failed:', error)
      throw error
    }
  }

  // Todo API methods
  async getTodos() {
    return this.request('/todos')
  }

  async getTodo(id) {
    return this.request(`/todos/${id}`)
  }

  async createTodo(todo) {
    return this.request('/todos', {
      method: 'POST',
      body: JSON.stringify(todo)
    })
  }

  async updateTodo(id, todo) {
    return this.request(`/todos/${id}`, {
      method: 'PUT',
      body: JSON.stringify(todo)
    })
  }

  async deleteTodo(id) {
    return this.request(`/todos/${id}`, {
      method: 'DELETE'
    })
  }

  // List API methods
  async getLists() {
    return this.request('/lists')
  }

  async getList(id) {
    return this.request(`/lists/${id}`)
  }

  async createList(list) {
    return this.request('/lists', {
      method: 'POST',
      body: JSON.stringify(list)
    })
  }

  async updateList(id, list) {
    return this.request(`/lists/${id}`, {
      method: 'PUT',
      body: JSON.stringify(list)
    })
  }

  async deleteList(id) {
    return this.request(`/lists/${id}`, {
      method: 'DELETE'
    })
  }
}

// Create and export a singleton instance
const apiService = new ApiService()
export default apiService 