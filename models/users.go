package models

import (
	"errors"

	"github.com/eduardopeixe/instago/hash"
	"github.com/eduardopeixe/instago/rand"
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

type userValidator struct {
	UserDB
	hmac hash.HMAC
}

type userGorm struct {
	db *gorm.DB
}

type userService struct {
	UserDB
}

// UserService is a set of methods used to manipulate and work with the user model
type UserService interface {
	// Authenticate will verify the provided email and password and
	// if correct the user corresponding to that email will be return.
	// Otherwise returns an error.
	Authenticate(email, password string) (*User, error)
	UserDB
}

// UserDB is used to interact with the users database
//
// For pretty much all single user quiries:
// If user is found, we will return a nil error
// If user is not found, we will return ErrNotFound
// If there is another error, we will return an aeror with more information
// about what went wrong.  This mau not be an error generated by the models package
// As a general rule, any error but ErrNotFound should result in a 500 error
type UserDB interface {
	// Medhods for querying for single user
	ByID(id uint) (*User, error)
	ByEmail(email string) (*User, error)
	ByRemember(token string) (*User, error)

	// Methods for altering users
	Create(user *User) error
	Update(user *User) error
	Delete(id uint) error

	//Used to close DB connection
	Close() error

	//Migration helpers
	AutoMigrate() error
	DestructiveReset()
}

var _ UserDB = &userGorm{}

func newuserGorm(connectionInfo string) (*userGorm, error) {
	db, err := gorm.Open("postgres", connectionInfo)
	if err != nil {
		return nil, err
	}

	db.LogMode(true)

	return &userGorm{
		db: db,
	}, nil
}

// NewUserService creates a new DB connection
func NewUserService(connectionInfo string) (UserService, error) {
	ug, err := newuserGorm(connectionInfo)
	if err != nil {
		return nil, err
	}

	hmac := hash.NewHMAC(hmacSecretKey)
	uv := &userValidator{
		hmac:   hmac,
		UserDB: ug,
	}

	return &userService{UserDB: uv}, nil
}

type userValFunc func(*User) error

func runUserValFuncs(user *User, fns ...userValFunc) error {
	for _, fn := range fns {
		if err := fn(user); err != nil {
			return err
		}
	}
	return nil
}

// bcryptPassword will hash a user password with a predefined pepper
// and bcrypt
func (uv *userValidator) bcryptPassword(user *User) error {
	if user.Password == "" {
		return nil
	}
	pwBytes := []byte(user.Password + userPwPepper)
	hashedBytes, err := bcrypt.GenerateFromPassword(pwBytes, bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.PasswordHash = string(hashedBytes)
	user.Password = ""
	return nil
}

func (uv *userValidator) hmacRemember(user *User) error {
	if user.Remember == "" {
		return nil
	}
	user.RememberHash = uv.hmac.Hash(user.Remember)
	user.Remember = ""
	return nil
}

func (uv *userValidator) SetRemember(user *User) error {
	if user.Remember != "" {

		return nil
	}
	token, err := rand.RememberToken()
	if err != nil {
		return err
	}
	user.Remember = token
	return nil
}

func (uv *userValidator) IDGreaterThanZero(user *User) error {
	if user.ID <= 0 {
		return ErrInvalidID
	}
	return nil
}

// ByRemember will hash the remeber token and then call ByRemember
// on the subsequent UserDB layer
func (uv *userValidator) ByRemember(token string) (*User, error) {
	user := User{
		Remember: token,
	}

	err := runUserValFuncs(&user, uv.hmacRemember)
	if err != nil {
		return nil, err
	}
	return uv.UserDB.ByRemember(user.RememberHash)

}

// Create creates the provided user
func (uv *userValidator) Create(user *User) error {

	err := runUserValFuncs(user,
		uv.bcryptPassword,
		uv.SetRemember,
		uv.hmacRemember)
	if err != nil {
		return err
	}

	return uv.UserDB.Create(user)
}

// Update will hash a token, if provided
func (uv *userValidator) Update(user *User) error {
	err := runUserValFuncs(user,
		uv.bcryptPassword,
		uv.hmacRemember)
	if err != nil {
		return err
	}

	return uv.UserDB.Update(user)
}

// Delete updates a provided user with user data provided
func (uv *userValidator) Delete(id uint) error {
	var user User
	user.ID = id
	err := runUserValFuncs(&user, uv.IDGreaterThanZero)
	if err != nil {
		return err
	}

	return uv.UserDB.Delete(id)
}

func (uv *userValidator) ByID(id uint) (*User, error) {
	//validate the ID
	if id <= 0 {
		return nil, errors.New("Invalid ID")
	}
	return uv.UserDB.ByID(id)
}

// ByID will look a user with theprovided ID
func (ug *userGorm) ByID(id uint) (*User, error) {
	var user User

	db := ug.db.Where("id = ?", id)
	err := first(db, &user)

	return &user, err
}

// ByEmail will look for a user with the provided email
func (ug *userGorm) ByEmail(email string) (*User, error) {
	var user User

	db := ug.db.Where("email = ?", email)
	err := first(db, &user)

	return &user, err
}

// ByRemember looks up a user with the given remember token and returns
// that user. This method expects the remember token to already be hashed
func (ug *userGorm) ByRemember(rememberHash string) (*User, error) {
	var user User
	err := first(ug.db.Where("remember_hash = ?", rememberHash), &user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

//Authenticate a user with provided email and password
func (us *userService) Authenticate(email, password string) (*User, error) {
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
func (ug *userGorm) Create(user *User) error {

	return ug.db.Create(user).Error
}

// Update updates a provided user with user data provided
func (ug *userGorm) Update(user *User) error {

	return ug.db.Save(user).Error
}

// Delete updates a provided user with user data provided
func (ug *userGorm) Delete(id uint) error {

	user := User{Model: gorm.Model{ID: id}}

	return ug.db.Delete(&user).Error
}

// Close closes the user service database connection
func (ug *userGorm) Close() error {
	return ug.db.Close()
}

// DestructiveReset drop and recreate all database tables
func (ug *userGorm) DestructiveReset() {
	ug.db.DropTableIfExists(&User{})
	ug.db.AutoMigrate(&User{})

}

// AutoMigrate migrates users table
func (ug *userGorm) AutoMigrate() error {
	if err := ug.db.AutoMigrate(&User{}).Error; err != nil {
		return err
	}
	return nil
}
