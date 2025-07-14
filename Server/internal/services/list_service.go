package services

import (
	"todo-list/internal/models"
	"gorm.io/gorm"
)

type ListService struct {
	db *gorm.DB
}

func NewListService(db *gorm.DB) *ListService {	
	return &ListService{db: db}
}

func (s *ListService) GetAll() ([]models.List, error) {
	var lists []models.List
	err := s.db.Find(&lists).Error
	return lists, err
}

func (s *ListService) GetByID(id uint) (*models.List, error) {
	var list models.List
	err := s.db.First(&list, id).Error
	return &list, err
}

func (s *ListService) Create(list *models.List) error {
	return s.db.Create(list).Error
}

func (s *ListService) Update(list *models.List) error {
	return s.db.Save(list).Error	
}

func (s *ListService) Delete(id uint) error {
	return s.db.Delete(&models.List{}, id).Error
}	
