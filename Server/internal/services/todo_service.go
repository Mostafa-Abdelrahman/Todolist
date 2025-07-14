package services

import (
	"gorm.io/gorm"
	"todo-list/internal/models"
)

// TodoServiceInterface defines the contract for todo services
type TodoServiceInterface interface {
	GetAll() ([]models.Todo, error)
	GetByID(id uint) (*models.Todo, error)
	Create(todo *models.Todo) error
	Update(todo *models.Todo) error
	Delete(id uint) error
}

type TodoService struct {
	db *gorm.DB
}

// Ensure TodoService implements TodoServiceInterface
var _ TodoServiceInterface = (*TodoService)(nil)

func NewTodoService(db *gorm.DB) *TodoService {
	return &TodoService{db: db}
}

func (s *TodoService) GetAll() ([]models.Todo, error) {
	var todos []models.Todo
	err := s.db.Preload("List").Find(&todos).Error
	return todos, err
}

func (s *TodoService) GetByID(id uint) (*models.Todo, error) {
	var todo models.Todo
	err := s.db.Preload("List").First(&todo, id).Error
	return &todo, err
}

func (s *TodoService) Create(todo *models.Todo) error {
	return s.db.Create(todo).Error
}

func (s *TodoService) Update(todo *models.Todo) error {
	return s.db.Save(todo).Error
}

func (s *TodoService) Delete(id uint) error {
	return s.db.Delete(&models.Todo{}, id).Error
}