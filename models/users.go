package models

import (
	"errors"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

var (
	// Generic error when resource not found
	ErrNotFound = errors.New("models: resource not found")
)

// User is the model for a user
type User struct {
	gorm.Model
	Name  string
	Email string `gorm:"not null;unique_index"`
}

// UserService is the type to connect to a DB
type UserService struct {
	db *gorm.DB
}

// NewUserService creates a new DB connection
func NewUserService(connectionInfo string) (*UserService, error) {
	db, err := gorm.Open("postgres", connectionInfo)
	if err != nil {
		return nil, err
	}

	db.LogMode(true)

	return &UserService{db}, nil
}

// ByID will look a user with theprovided ID
func (us *UserService) ByID(id uint) (*User, error) {
	var user User

	err := us.db.Where("id = ?", id).First(&user).Error

	switch err {
	case nil:
		return &user, nil
	case gorm.ErrRecordNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

// Create creates the provided user
func (us *UserService) Create(user *User) error {
	return us.db.Create(user).Error
}

// Close closes the user service database connection
func (us *UserService) Close() error {
	return us.db.Close()
}

// ResetDB drop and recreate all database tables
func (us *UserService) ResetDB() {
	us.db.DropTableIfExists(&User{})
	us.db.AutoMigrate(&User{})

}
