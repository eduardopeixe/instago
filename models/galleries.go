package models

import "github.com/jinzhu/gorm"

// Gallery is image container resources
type Gallery struct {
	gorm.Model
	UserID uint   `gorm:"not_null;index"`
	Title  string `gorm:"not_null"`
}

type GalleryService interface{}

type GalleryDB interface {
	Create(gallery *Gallery) error
}

type GalleryService struct {
	GalleryDB
}

func NewGalleryService(db *gorm.DB) GalleryService {
	return &galleryService{
		GalleryDB: &galleryValidator{&galleryGorm{db}},
	}
}

type GalleryValidator struct {
	GalleryDB
}

var _ GalleryDB = &galleryGorm{}

func (gg *galleryGorm) Create(gallery *Gallery) error {
	return gg.db.Craete(gallery).Error
}
