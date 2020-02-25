package models

import (
	"errors"
	"regexp"
	"strings"

	"github.com/eduardopeixe/instago/hash"
	"github.com/eduardopeixe/instago/rand"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"golang.org/x/crypto/bcrypt"
)

const (
	userPwPepper  = "secret-instago-string"
	hmacSecretKey = "secret-hmac-key-12345678901234567890"
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

// ensure that userValidator implements UserDB
var _ UserDB = &userValidator{}

type userValidator struct {
	UserDB
	hmac       hash.HMAC
	emailRegex *regexp.Regexp
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
}

var _ UserDB = &userGorm{}

func newUserValidator(udb UserDB, hmac hash.HMAC) *userValidator {
	return &userValidator{
		UserDB: udb,
		hmac:   hmac,
		// This is not the best design because the program could panic
		// when creating a new UserValidator
		// having this in the global scope would make panic at starting it
		// the program
		emailRegex: regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,16}$`),
	}
}

// NewUserService creates a new DB connection
func NewUserService(db *gorm.DB) UserService {
	ug := &userGorm{db}

	hmac := hash.NewHMAC(hmacSecretKey)
	uv := newUserValidator(ug, hmac)

	return &userService{UserDB: uv}
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

func (uv *userValidator) passwordMinLength(user *User) error {
	if user.Password != "" && len(user.Password) < 8 {
		return ErrPasswordTooShort
	}
	return nil
}

func (uv *userValidator) passwordRequired(user *User) error {
	if user.Password == "" {
		return ErrPasswordRequired
	}
	return nil
}

func (uv *userValidator) passwordHashRequired(user *User) error {
	if user.PasswordHash == "" {
		return ErrPasswordHashRequired
	}
	return nil
}

func (uv *userValidator) hmacRemember(user *User) error {
	if user.Remember == "" {
		return nil
	}
	user.RememberHash = uv.hmac.Hash(user.Remember)
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

func (uv *userValidator) normalizeEmail(user *User) error {
	user.Email = strings.ToLower(user.Email)
	user.Email = strings.TrimSpace(user.Email)
	return nil
}

func (uv *userValidator) requireEmail(user *User) error {
	if user.Email == "" {
		return ErrEmailRequired
	}
	return nil
}

func (uv *userValidator) emailFormat(user *User) error {
	if user.Email != "" && !uv.emailRegex.MatchString(user.Email) {
		return ErrEmailInvalid
	}
	return nil
}

func (uv *userValidator) emailIsAvailable(user *User) error {
	existing, err := uv.ByEmail(user.Email)
	if err == ErrNotFound {
		return nil
	}
	if err != nil {
		return err
	}

	// In case the user is found and the userID is a match, that means
	// it is an update, so return nil
	if user.ID == existing.ID {
		return nil
	} else {
		return ErrEmailTaken
	}
}

func (uv *userValidator) ByEmail(email string) (*User, error) {
	user := User{
		Email: email,
	}
	err := runUserValFuncs(&user,
		uv.normalizeEmail,
		uv.emailFormat,
	)
	if err != nil {
		return nil, err
	}
	return uv.UserDB.ByEmail(email)
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
		uv.normalizeEmail,
		uv.requireEmail,
		uv.emailFormat,
		uv.emailIsAvailable,
		uv.passwordRequired,
		uv.passwordMinLength,
		uv.bcryptPassword,
		uv.passwordHashRequired,
		uv.SetRemember,
		uv.hmacRemember,
	)
	if err != nil {
		return err
	}

	return uv.UserDB.Create(user)
}

// Update will hash a token, if provided
func (uv *userValidator) Update(user *User) error {
	err := runUserValFuncs(user,
		uv.bcryptPassword,
		uv.passwordHashRequired,
		uv.normalizeEmail,
		uv.requireEmail,
		uv.hmacRemember,
		uv.emailIsAvailable)
	if err != nil {
		return err
	}

	return uv.UserDB.Update(user)
}

// Delete updates a provided user with user data provided
func (uv *userValidator) Delete(id uint) error {
	var user User
	user.ID = id
	err := runUserValFuncs(&user,
		uv.IDGreaterThanZero,
		uv.normalizeEmail,
	)
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
