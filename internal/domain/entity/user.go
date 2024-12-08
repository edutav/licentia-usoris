package entity

import (
	"time"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	UUID            string
	Name            string
	Email           string
	PasswordHash    string
	DOB             time.Time
	PhoneNumber     string
	IsBlocked       bool
	IsEmailVerified bool
	CreatedAt       time.Time
	UpdatedAt       time.Time
	DeletedAt       time.Time
	IsDeleted       bool
	LastLogin       time.Time
}

func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword(
		[]byte(u.PasswordHash),
		[]byte(password),
	)
	return err == nil
}
