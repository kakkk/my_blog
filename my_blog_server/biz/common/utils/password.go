package utils

import "golang.org/x/crypto/bcrypt"

// CompareHashAndPassword 比较Hash和密码
func CompareHashAndPassword(hash string, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		return false
	}
	return true
}

// GenerateHashFromPassword 通过密码生成Hash
func GenerateHashFromPassword(password string) (hash string) {
	hashBytes, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	hash = string(hashBytes)
	return
}
