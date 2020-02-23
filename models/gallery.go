package models

import "github.com/jinzhu/gorm"

// Gallery is image container resources
type Gallery struct {
	gorm.Model
	UserID uint   `gorm:"not_null;index"`
	Title  string `gorm:"not_null"`
}
