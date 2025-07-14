package unit

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"todo-list/internal/api/handlers"
	"todo-list/internal/models"
	"todo-list/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock TodoService
type MockTodoService struct {
	mock.Mock
}

// Ensure MockTodoService implements TodoServiceInterface
var _ services.TodoServiceInterface = (*MockTodoService)(nil)

func (m *MockTodoService) GetAll() ([]models.Todo, error) {
	args := m.Called()
	return args.Get(0).([]models.Todo), args.Error(1)
}

func (m *MockTodoService) GetByID(id uint) (*models.Todo, error) {
	args := m.Called(id)
	return args.Get(0).(*models.Todo), args.Error(1)
}

func (m *MockTodoService) Create(todo *models.Todo) error {
	args := m.Called(todo)
	return args.Error(0)
}

func (m *MockTodoService) Update(todo *models.Todo) error {
	args := m.Called(todo)
	return args.Error(0)
}

func (m *MockTodoService) Delete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func TestGetTodos(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	mockService := new(MockTodoService)
	handler := handlers.NewTodoHandler(mockService)

	// Mock data
	todos := []models.Todo{
		{ID: 1, Title: "Test Todo 1", ListID: 1, Due: "2024-12-25", Done: false},
		{ID: 2, Title: "Test Todo 2", ListID: 1, Due: "2024-12-26", Done: true},
	}

	// Expectations
	mockService.On("GetAll").Return(todos, nil)

	// Create request
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/todos", nil)
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	// Execute
	handler.GetTodos(c)

	// Assertions
	assert.Equal(t, http.StatusOK, w.Code)

	var response []models.Todo
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Len(t, response, 2)
	assert.Equal(t, "Test Todo 1", response[0].Title)

	mockService.AssertExpectations(t)
}

func TestCreateTodo(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	mockService := new(MockTodoService)
	handler := handlers.NewTodoHandler(mockService)

	// Test data
	todoData := models.Todo{
		Title:  "New Todo",
		ListID: 1,
		Due:    "2024-12-25",
		Done:   false,
	}

	// Expectations
	mockService.On("Create", mock.AnythingOfType("*models.Todo")).Return(nil)

	// Create request
	w := httptest.NewRecorder()
	jsonData, _ := json.Marshal(todoData)
	req, _ := http.NewRequest("POST", "/api/todos", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	// Execute
	handler.CreateTodo(c)

	// Assertions
	assert.Equal(t, http.StatusCreated, w.Code)

	var response models.Todo
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "New Todo", response.Title)

	mockService.AssertExpectations(t)
}
