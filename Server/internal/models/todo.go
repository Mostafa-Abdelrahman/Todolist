package models

import (
	"time"

	"gorm.io/gorm"
)

type Todo struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Title     string         `json:"title" gorm:"not null"`
	ListID    uint           `json:"listId" gorm:"not null"`
	List      List           `json:"list" gorm:"foreignKey:ListID"`
	Done      bool           `json:"done" gorm:"default:false"`
	Due       string         `json:"due" gorm:"not null"`
	CreatedAt time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}
