package models

import (
	"errors"

	"github.com/eduardopeixe/instago/hash"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"golang.org/x/crypto/bcrypt"
)

const (
	userPwPepper  = "secret-instago-string"
	hmacSecretKey = "secret-hmac-key"
)

var (
	// ErrNotFound is a generic error for resource not found
	ErrNotFound = errors.New("models: resource not found")
	// ErrInvalidID returns error when an invalid ID is provided
	ErrInvalidID = errors.New("models: invalid ID")
	// ErrInvalidEmail returns error when an invalid email is provided
	ErrInvalidEmail = errors.New("models: invalid email address")
	// ErrIncorrectPassword returns error when a invalid password is provided
	ErrIncorrectPassword = errors.New("models: incorrect password provided")
)

// User is the model for a user
type User struct {
	gorm.Model
	Name         string
	Email        string `gorm:"not null;unique_index"`
	Password     string `gorm:"-"`
	PasswordHash string `gorm:"not null"`
	Remember     string `gorm:"-"`
	RememberHash string `gorm:"not null,unique_index"`
}

// UserService is the type to connect to a DB
type UserService struct {
	db   *gorm.DB
	hmac hash.HMAC
}

// NewUserService creates a new DB connection
func NewUserService(connectionInfo string) (*UserService, error) {
	db, err := gorm.Open("postgres", connectionInfo)
	if err != nil {
		return nil, err
	}

	db.LogMode(true)

	hmac := hash.NewHMAC(hmacSecretKey)

	return &UserService{db, hmac}, nil
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

// ByRemember looks up a user with the given remember token
// and return that user
func (us *UserService) ByRemember(token string) (*User, error){
	var user User
	rememberHash := us.hmac.Hash(token)
	err := first(us.db.Where("remember_hash" = ?, rememberHash), &user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

//Authenticate a user with provided email and password
func (us *UserService) Authenticate(email, password string) (*User, error) {
	foundUser, err := us.ByEmail(email)
	if err != nil {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(foundUser.PasswordHash), []byte(password+userPwPepper))
	if err != nil {
		switch err {
		case bcrypt.ErrMismatchedHashAndPassword:
			return nil, ErrIncorrectPassword
		default:
			return nil, err
		}
	}
	return foundUser, nil
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
	pwBytes := []byte(user.Password + userPwPepper)
	hashedBytes, err := bcrypt.GenerateFromPassword(pwBytes, bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user.PasswordHash = string(hashedBytes)
	user.Password = ""

	if user.Remember == "" {
		token, err := rand.RememberToken()
		if err != nil {
			return err
		}
		user.Remember = token
	}

	user.RememberHash = us.hmac.Hash(user.Remember)
	return us.db.Create(user).Error

	if user.Remember != "" {
		user.RememberHash = us.hmac.Hash(user.Remember)
	}

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

// AutoMigrate migrates users table
func (us *UserService) AutoMigrate() error {
	if err := us.db.AutoMigrate(&User{}).Error; err != nil {
		return err
	}
	return nil
}
