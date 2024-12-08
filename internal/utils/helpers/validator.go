package helpers

import (
	"regexp"
	"time"

	"github.com/edutav/licentia-usoris/internal/utils"
)

func ValidateName(name string) error {
	if name == "" {
		return utils.ErrInvalidName
	}

	return nil
}

func ValidateEmail(email string) error {
	if email == "" {
		return utils.ErrMissingEmail
	}

	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(email) {
		return utils.ErrInvalidEmail
	}

	return nil
}

func ValidatePassword(password string) error {
	if password == "" {
		return utils.ErrPasswordInvalid
	}

	return nil
}

func ValidateDOB(dob string) error {
	if _, err := time.Parse("2006-01-02", dob); err != nil {
		return utils.ErrDOBFormat
	}

	return nil
}

func ValidatePhoneNumber(phoneNumber string) error {
	phoneRegex := regexp.MustCompile(`^\d{11}$`)
	if !phoneRegex.MatchString(phoneNumber) {
		return utils.ErrInvalidPhoneNumber
	}

	return nil
}
