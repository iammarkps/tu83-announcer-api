package handler

import "github.com/jinzhu/gorm"

// Handler provide db connection to all handler
type Handler struct {
	DB *gorm.DB
}
