package validator

import (
	"regexp"
	"unicode"

	"github.com/edutav/licentia-usoris/internal/utils"
)

type ValidatePasswordFunc func(password string) error

var ValidateUserPassword ValidatePasswordFunc = defaultValidateUserPassword

func defaultValidateUserPassword(password string) error {
	if len(password) < 8 {
		return utils.ErrPasswordTooShort
	}

	if len(password) > 64 {
		return utils.ErrPasswordTooLong
	}

	var (
		hasUppercase   bool
		hasLowercase   bool
		hasNumber      bool
		hasSpecialChar bool
	)

	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUppercase = true
		case unicode.IsLower(char):
			hasLowercase = true
		case unicode.IsNumber(char):
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecialChar = true
		}
	}

	if !hasUppercase || !hasLowercase || !hasNumber || !hasSpecialChar {
		return utils.ErrPasswordSecurity
	}

	return nil
}

func ValidateUserName(name string) error {
	if len(name) < 2 {
		return utils.ErrUserNameTooShort
	}

	if len(name) > 200 {
		return utils.ErrUserNameTooLong
	}

	if matched, _ := regexp.MatchString(`[0-9]`, name); matched {
		return utils.ErrUserNameWithNumericVals
	}

	return nil
}

func ValidateOTP(otp string) error {
	if len(otp) == 0 {
		return utils.ErrMissingOTP
	}
	if len(otp) != 6 {
		return utils.ErrOtpLength
	}

	matched, _ := regexp.MatchString(`^\d{6}$`, otp)
	if !matched {
		return utils.ErrOtpNums
	}

	return nil
}
