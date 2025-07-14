package integration

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"todo-list/internal/api/handlers"
	"todo-list/internal/models"
	"todo-list/internal/services"
)

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(t, err)
	
	// Migrate test schema
	err = db.AutoMigrate(&models.Todo{}, &models.List{})
	assert.NoError(t, err)
	
	return db
}

func setupTestServer(t *testing.T) *gin.Engine {
	db := setupTestDB(t)
	todoService := services.NewTodoService(db)
	listService := services.NewListService(db)
	
	todoHandler := handlers.NewTodoHandler(todoService)
	listHandler := handlers.NewListHandler(listService)
	
	gin.SetMode(gin.TestMode)
	r := gin.New()
	
	api := r.Group("/api")
	{
		api.GET("/todos", todoHandler.GetTodos)
		api.GET("/todos/:id", todoHandler.GetTodo)
		api.POST("/todos", todoHandler.CreateTodo)
		api.PUT("/todos/:id", todoHandler.UpdateTodo)
		api.DELETE("/todos/:id", todoHandler.DeleteTodo)
		
		api.GET("/lists", listHandler.GetLists)
		api.GET("/lists/:id", listHandler.GetList)
		api.POST("/lists", listHandler.CreateList)
		api.PUT("/lists/:id", listHandler.UpdateList)
		api.DELETE("/lists/:id", listHandler.DeleteList)
	}
	
	return r
}

func TestTodoCRUD(t *testing.T) {
	router := setupTestServer(t)
	
	// Test 1: Create Todo
	todoData := models.Todo{
		Title:  "Integration Test Todo",
		ListID: 1,
		Due:    "2024-12-25",
		Done:   false,
	}
	
	jsonData, _ := json.Marshal(todoData)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/todos", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusCreated, w.Code)
	
	var createdTodo models.Todo
	err := json.Unmarshal(w.Body.Bytes(), &createdTodo)
	assert.NoError(t, err)
	assert.Equal(t, "Integration Test Todo", createdTodo.Title)
	assert.False(t, createdTodo.Done)
	
	// Test 2: Get All Todos
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/api/todos", nil)
	router.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusOK, w.Code)
	
	var todos []models.Todo
	err = json.Unmarshal(w.Body.Bytes(), &todos)
	assert.NoError(t, err)
	assert.Len(t, todos, 1)
}

func TestListCRUD(t *testing.T) {
	router := setupTestServer(t)
	
	// Test 1: Create List
	listData := models.List{
		Name: "Test List",
	}
	
	jsonData, _ := json.Marshal(listData)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/lists", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusCreated, w.Code)
	
	var createdList models.List
	err := json.Unmarshal(w.Body.Bytes(), &createdList)
	assert.NoError(t, err)
	assert.Equal(t, "Test List", createdList.Name)
}