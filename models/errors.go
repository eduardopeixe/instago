package models

import (
	"strings"
)

var (
	// ErrNotFound is a generic error for resource not found
	ErrNotFound modelError = "models: resource not found"
	// ErrInvalidEmail when an invalid email is provided
	ErrInvalidEmail modelError = "models: invalid email address"
	// ErrIncorrectPassword when a invalid password is provided
	ErrIncorrectPassword modelError = "models: incorrect password provided"
	//ErrEmailRequired when an email is not present
	ErrEmailRequired modelError = "models: email address is required"
	//ErrEmailInvalid when email is not in the correct format
	ErrEmailInvalid modelError = "models: email address is invalid"
	//ErrEmailTaken when a user with same email already exists
	ErrEmailTaken modelError = "models: email address is already taken"
	//ErrPasswordTooShort when a password's length is not enough
	ErrPasswordTooShort modelError = "models: password must be at least 8 characters long"
	//ErrPasswordRequired when a password is not present
	ErrPasswordRequired modelError = "models: password is required"
	//ErrTitleRequired when a title is not provided
	ErrTitleRequired modelError = "models: Title is required"
	// ErrInvalidID when an invalid ID is provided
	ErrInvalidID privateError = "models: invalid ID"
	//ErrUserIDRequired when a userID is not present
	ErrUserIDRequired privateError = "models: UserID is required"
	//ErrPasswordHashRequired when a password is not present
	ErrPasswordHashRequired privateError = "models: password hash not present"
)

type modelError string

func (e modelError) Error() string {
	return string(e)
}

func (e modelError) Public() string {
	s := strings.Replace(string(e), "models: ", "", 1)
	split := strings.Split(s, " ")
	split[0] = strings.Title(split[0])
	return strings.Join(split, " ")
}

type privateError string

func (e privateError) Error() string {
	return string(e)
}
