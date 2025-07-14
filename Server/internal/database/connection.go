package database

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"todo-list/internal/models"
)

func Connect(databaseURL string) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(databaseURL), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	if err := db.AutoMigrate(&models.Todo{}, &models.List{}); err != nil {
		return nil, err
	}

	return db, nil
}