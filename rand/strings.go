package rand

import (
	"crypto/rand"
	"encoding/base64"
)

const rememberTokenBytes = 32

// Bytes generates n random bytes
func Bytes(n int) ([]byte, error) {
	b := make([]byte, n)

	_, err := rand.Read(b)

	if err != nil {
		return nil, err
	}

	return b, nil

}

// String will generate a byte slice of size nBytes and then
// return a string that is the base64 URL encoded version
func String(nBytes int) (string, error) {
	b, err := Bytes(nBytes)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

// RememberToken generates remembertoken nBytes size
func RememberToken() (string, error) {
	return String(rememberTokenBytes)
}
