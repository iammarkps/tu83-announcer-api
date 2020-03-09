package middlewares

import (
	"github.com/jinzhu/gorm"
)

// Middleware inject DB
type Middleware struct {
	DB *gorm.DB
}
