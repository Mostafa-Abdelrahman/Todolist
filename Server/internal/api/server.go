package api
import (
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"todo-list/internal/api/routes"
	"todo-list/internal/config"
	"todo-list/internal/services"
)
 type Server struct {
	config *config.Config
	db *gorm.DB
	router *gin.Engine
 }

func NewServer(cfg *config.Config, db *gorm.DB) *Server{
	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173", "http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))
	 
	todoService := services.NewTodoService(db)
	listService := services.NewListService(db)

	routes.SetupRoutes(router, todoService, listService)

	return &Server{
		config: cfg,
		db: db,
		router: router,
	}
}

func (s *Server) Start() error {
	return s.router.Run(fmt.Sprintf(":%s", s.config.ServerPort))
}