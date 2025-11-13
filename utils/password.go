package utils

import "golang.org/x/crypto/bcrypt"

var MockCheckPassword func(pw, hash string) bool

func HashPassword(pw string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(pw), bcrypt.DefaultCost)
	return string(b), err
}

func CheckPassword(pw, hash string) bool {
	if MockCheckPassword != nil {
		return MockCheckPassword(pw, hash)
	}
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(pw)) == nil
}
