package routes

import (
	"net/http"
	"todo-list/internal/api/handlers"
	"todo-list/internal/services"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine, todoService services.TodoServiceInterface, listService *services.ListService) {
	todoHandler := handlers.NewTodoHandler(todoService)
	listHandler := handlers.NewListHandler(listService)

	api := r.Group("/api")
	{
		api.GET("/health", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "Server is running"})
		})
		// Todo routes
		api.GET("/todos", todoHandler.GetTodos)
		api.GET("/todos/:id", todoHandler.GetTodo)
		api.POST("/todos", todoHandler.CreateTodo)
		api.PUT("/todos/:id", todoHandler.UpdateTodo)
		api.DELETE("/todos/:id", todoHandler.DeleteTodo)

		// List routes
		api.GET("/lists", listHandler.GetLists)
		api.GET("/lists/:id", listHandler.GetList)
		api.POST("/lists", listHandler.CreateList)
		api.PUT("/lists/:id", listHandler.UpdateList)
		api.DELETE("/lists/:id", listHandler.DeleteList)
	}
}
