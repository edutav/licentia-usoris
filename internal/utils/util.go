package utils

import "errors"

const (
	InternalServerErrorString = "Internal server error"
)

var (
	// email errors
	ErrInvalidEmail    = errors.New("invalid email format")
	ErrMissingEmail    = errors.New("no email input given")
	ErrDuplicateEmail  = errors.New("email already exists")
	ErrSMTPServerIssue = errors.New("SMTP server issue")

	// users errors
	ErrUserNotFound            = errors.New("user not found")
	ErrPasswordTooShort        = errors.New("password too short")
	ErrPasswordTooLong         = errors.New("password too long")
	ErrPasswordSecurity        = errors.New("password not safe")
	ErrUserNameTooShort        = errors.New("user name too short")
	ErrUserNameTooLong         = errors.New("user name too long")
	ErrUserNameWithNumericVals = errors.New("user name with numeric characters")
	ErrHashingPassword         = errors.New("error hashing password")
	ErrDOBFormat               = errors.New("invalid dob format")
	ErrInvalidName             = errors.New("invalid name")
	ErrPasswordInvalid         = errors.New("invalid password")
	ErrInvalidPhoneNumber      = errors.New("invalid phone number")
	ErrMissingPhoneNumber      = errors.New("no phone number input given")

	// Pre-registration errors
	ErrCreateVericationEntry    = errors.New("error creating verification entry")
	ErrPreRegistredUserNotFound = errors.New("pre-registred user not found")

	// login errors
	ErrGenerateJWTTokenWithRole = errors.New("error generate jwt token with role")

	// otp errors
	ErrGenerateOTP        = errors.New("error generating otp")
	ErrMissingOTP         = errors.New("need valid otp input")
	ErrOtpLength          = errors.New("otp must have 6 digits")
	ErrOtpNums            = errors.New("non digits in otp")
	ErrOTPExpired         = errors.New("OTP has expired")
	ErrOTPAlreadyVerified = errors.New("email already verified")
	ErrInvalidOTP         = errors.New("invalid OTP")
)
