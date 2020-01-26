package models

import (
	"errors"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

var (
	// ErrNotFound is a generic error for resource not found
	ErrNotFound = errors.New("models: resource not found")
	// ErrInvalidID returns error when a invalid ID is provided
	ErrInvalidID = errors.New("models: invalid ID")
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

	db := us.db.Where("id = ?", id)
	err := first(db, &user)

	return &user, err
}

// ByEmail will look for a user with the provided email
func (us *UserService) ByEmail(email string) (*User, error) {
	var user User

	db := us.db.Where("email = ?", email)
	err := first(db, &user)

	return &user, err
}

// first will query provided gorm.db and get the first item
// returned an place it into dst interface
func first(db *gorm.DB, dst interface{}) error {
	err := db.First(dst).Error
	if err == gorm.ErrRecordNotFound {
		return ErrNotFound
	}
	return err

}

// Create creates the provided user
func (us *UserService) Create(user *User) error {
	return us.db.Create(user).Error
}

// Update updates a provided user with user data provided
func (us *UserService) Update(user *User) error {
	return us.db.Save(user).Error
}

// Delete updates a provided user with user data provided
func (us *UserService) Delete(id uint) error {
	if id <= 0 {
		return ErrInvalidID
	}

	user := User{Model: gorm.Model{ID: id}}

	return us.db.Delete(&user).Error
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
