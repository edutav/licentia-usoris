package entity

import "time"

type PreRegistration struct {
	UUID         string
	Email        string
	PasswordHash string
	CodeOTP      string
	UserData     *User
	ExpiresAt    time.Time
	IsVerified   bool
	CreatedAt    time.Time
}
