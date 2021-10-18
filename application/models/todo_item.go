package models

import "github.com/jinzhu/gorm"

type ToDoItem struct {
	gorm.Model
	Description string
	IsCompleted bool
	ImageURL    string
}
