package main

import (
	"golang.org/x/crypto/bcrypt"
	"regexp"
)

// Validate password format (lowercase + digits, min 8 chars)
func isValidPassword(password string) bool {
	re := regexp.MustCompile(`^[a-z\d]{8,}$`)
	return re.MatchString(password)
}

// HashPassword hashes a plain text password using bcrypt
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// CheckPasswordHash compares a hashed password with its plain text counterpart
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
