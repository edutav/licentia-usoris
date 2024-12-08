package usecases

import (
	"context"
	"time"

	"github.com/edutav/licentia-usoris/infrastructure/email"
	otpapp "github.com/edutav/licentia-usoris/infrastructure/otp_app"
	"github.com/edutav/licentia-usoris/internal/domain/entity"
	"github.com/edutav/licentia-usoris/internal/domain/reporitory"
	"github.com/edutav/licentia-usoris/internal/presentation/schemas"
	"github.com/edutav/licentia-usoris/internal/usecases/validator"
	"github.com/edutav/licentia-usoris/internal/utils"
	"golang.org/x/crypto/bcrypt"
)

type UserUseCase interface {
	// Pre-registration new user
	PreRegisterUser(ctx context.Context, preRegistration *schemas.PreRegistrationInput) error

	// Verify OTP code
	VerifyOTPCode(ctx context.Context, email, code string) error
}

type userUseCase struct {
	userRepository       reporitory.UserRepository
	emailSender          email.EmailSender
	validateUserPassword validator.ValidatePasswordFunc
}

// NewUserUseCase creates a new user use case
func NewUserUseCase(userRepository reporitory.UserRepository, emailSender email.EmailSender, validatePassword validator.ValidatePasswordFunc) UserUseCase {
	if validatePassword == nil {
		validatePassword = validator.ValidateUserPassword
	}
	return &userUseCase{
		userRepository:       userRepository,
		emailSender:          emailSender,
		validateUserPassword: validatePassword,
	}
}

// PreRegisterUser implements UserUseCase.
func (u *userUseCase) PreRegisterUser(ctx context.Context, preRegistration *schemas.PreRegistrationInput) error {
	// Check if the email is already registered
	existingUser, err := u.userRepository.GetUserByEmail(ctx, preRegistration.Email)
	if err == nil && existingUser != nil && existingUser.IsEmailVerified {
		return utils.ErrDuplicateEmail
	} else if err != nil && err != utils.ErrUserNotFound {
		return err
	}

	// Validate password
	err = u.validateUserPassword(preRegistration.Password)
	if err != nil {
		switch err {
		case utils.ErrPasswordTooShort:
			return utils.ErrPasswordTooShort
		case utils.ErrPasswordTooLong:
			return utils.ErrPasswordTooLong
		case utils.ErrPasswordSecurity:
			return utils.ErrPasswordSecurity
		}
	}

	// Password hashing
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(preRegistration.Password), bcrypt.DefaultCost)
	if err != nil {
		return utils.ErrHashingPassword
	}
	preRegistration.Password = ""

	err = validator.ValidateUserName(preRegistration.Name)
	if err != nil {
		switch err {
		case utils.ErrUserNameTooShort:
			return utils.ErrUserNameTooShort
		case utils.ErrUserNameTooLong:
			return utils.ErrUserNameTooLong
		case utils.ErrUserNameWithNumericVals:
			return utils.ErrUserNameWithNumericVals
		}
	}

	key, otp, err := otpapp.GenerateOTP()
	if err != nil {
		return utils.ErrGenerateOTP
	}

	// OTP Code expiration 60 minutes
	expiresAt := time.Now().UTC().Add(time.Minute * 60)

	// Parse date of birth
	var dob time.Time
	if preRegistration.DateOfBirth != "" {
		parsedDOB, err := time.Parse("2006-01-02", preRegistration.DateOfBirth)
		if err != nil {
			return utils.ErrDOBFormat
		}
		dob = parsedDOB
	}

	// Create new user
	newUser := &entity.User{
		Name:         preRegistration.Name,
		Email:        preRegistration.Email,
		DOB:          dob,
		PasswordHash: string(passwordHash),
		PhoneNumber:  preRegistration.PhoneNumber,
	}

	preRegistrationEntity := &entity.PreRegistration{
		Email:        preRegistration.Email,
		PasswordHash: string(passwordHash),
		CodeOTP:      key,
		UserData:     newUser,
		ExpiresAt:    expiresAt,
		CreatedAt:    time.Now().UTC(),
	}

	// Save pre registration to database
	err = u.userRepository.PreRegisterUser(ctx, preRegistrationEntity)
	if err != nil {
		return utils.ErrCreateVericationEntry
	}

	// Send OTP to user email
	err = u.emailSender.SendOTP(preRegistration.Email, otp)
	if err != nil {
		return utils.ErrSMTPServerIssue
	}

	return err
}

// VerifyOTPCode implements UserUseCase.
func (u *userUseCase) VerifyOTPCode(ctx context.Context, email string, code string) error {
	// Get user by email
	userRegistred, err := u.userRepository.GetPreRegisteredByEmailAndOTPCode(ctx, email, code)
	if err != nil {
		if err == utils.ErrPreRegistredUserNotFound {
			return utils.ErrPreRegistredUserNotFound
		}
	}

	// Check if the OTP code has expired
	if userRegistred.ExpiresAt.Before(time.Now().UTC()) {
		return utils.ErrOTPExpired
	}

	// Check if user is already verified
	if userRegistred.IsVerified {
		return utils.ErrOTPAlreadyVerified
	}

	// Check if the OTP code is valid
	if otpapp.ValidateOTP(code, userRegistred.CodeOTP) {
		return utils.ErrInvalidOTP
	}

	// Check user existing in database
	existingUser, _ := u.userRepository.GetUserByEmail(ctx, email)
	if existingUser != nil {
		return utils.ErrDuplicateEmail
	}

	// Set user as pre registered
	newUser := entity.User{
		Name:            userRegistred.UserData.Name,
		Email:           userRegistred.UserData.Email,
		PasswordHash:    userRegistred.UserData.PasswordHash,
		DOB:             userRegistred.UserData.DOB,
		PhoneNumber:     userRegistred.UserData.PhoneNumber,
		IsBlocked:       false,
		IsEmailVerified: true,
		CreatedAt:       time.Now().UTC(),
		UpdatedAt:       time.Now().UTC(),
		DeletedAt:       time.Time{},
		IsDeleted:       false,
		LastLogin:       time.Time{},
	}

	// Save user to database
	err = u.userRepository.CreateUser(ctx, &newUser)
	if err != nil {
		return utils.ErrDuplicateEmail
	}

	err = u.userRepository.UpdateUserIsVerified(ctx, userRegistred.Email)
	if err != nil {
		return err
	}

	return err
}
