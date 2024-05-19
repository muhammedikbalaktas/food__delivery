package controllers

import "golang.org/x/crypto/bcrypt"

func hashPassword(password string) ([]byte, error) {

	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

func comparePasswords(hashedPassword []byte, password []byte) error {

	return bcrypt.CompareHashAndPassword(hashedPassword, password)
}
