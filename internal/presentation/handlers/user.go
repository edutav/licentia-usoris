package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/edutav/licentia-usoris/infrastructure/server/api"
	"github.com/edutav/licentia-usoris/internal/presentation/schemas"
	"github.com/edutav/licentia-usoris/internal/usecases"
	"github.com/edutav/licentia-usoris/internal/utils"
	"github.com/edutav/licentia-usoris/internal/utils/helpers"
)

// UserHandler is the handler for user related operations
type UserHandler struct {
	userUseCase usecases.UserUseCase
}

// NewUserHandler creates a new user handler
func NewUserHandler(userUseCase usecases.UserUseCase) *UserHandler {
	return &UserHandler{
		userUseCase: userUseCase,
	}
}

func (h *UserHandler) PreRegister(w http.ResponseWriter, r *http.Request) {
	// Check content type
	if r.Header.Get("Content-Type") != "application/json" {
		api.SendErrorResponse(w, http.StatusBadRequest, "Invalid content type", "Content type must be application/json")
		return
	}

	var input *schemas.PreRegistrationInput
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		api.SendErrorResponse(w, http.StatusBadRequest, "Invalid request body", err.Error())
		return
	}

	input.Name = strings.TrimSpace(input.Name)
	input.Email = strings.ToLower(strings.TrimSpace(input.Email))
	input.PhoneNumber = strings.TrimSpace(input.PhoneNumber)
	input.DateOfBirth = strings.TrimSpace(input.DateOfBirth)

	// Validate input name
	err = helpers.ValidateName(input.Name)
	if err != nil {
		api.SendErrorResponse(w, http.StatusBadRequest, "Invalid name", "Please provide a valid name")
		return
	}

	// Validate input email
	err = helpers.ValidateEmail(input.Email)
	if err != nil {
		switch err {
		case utils.ErrMissingEmail:
			api.SendErrorResponse(w, http.StatusBadRequest, "Missing email", "Please provide an email address")
		case utils.ErrInvalidEmail:
			api.SendErrorResponse(w, http.StatusBadRequest, "Invalid email", "Please provide a valid email address")
		default:
			api.SendErrorResponse(w, http.StatusInternalServerError, "Internal server error", err.Error())
		}

		return
	}

	// Validate input phone number
	if input.PhoneNumber != "" {
		err = helpers.ValidatePhoneNumber(input.PhoneNumber)
		if err != nil {
			switch err {
			case utils.ErrMissingPhoneNumber:
				api.SendErrorResponse(w, http.StatusBadRequest, "Missing phone number", "Please provide a phone number")
			case utils.ErrInvalidPhoneNumber:
				api.SendErrorResponse(w, http.StatusBadRequest, "Invalid phone number", "Please provide a valid phone number")
			default:
				api.SendErrorResponse(w, http.StatusInternalServerError, "Internal server error", err.Error())
			}

			return
		}
	}

	// Validate input date of birth
	if input.DateOfBirth != "" {
		err = helpers.ValidateDOB(input.DateOfBirth)
		if err != nil {
			api.SendErrorResponse(w, http.StatusBadRequest, "Invalid date of birth", "Please provide a valid date of birth")
			return
		}
	}

	// Validate input password
	err = helpers.ValidatePassword(input.Password)
	if err != nil {
		api.SendErrorResponse(w, http.StatusBadRequest, "Missing password", "Please provide a password")
		return
	}

	// Pre-register user
	err = h.userUseCase.PreRegisterUser(r.Context(), input)
	if err != nil {
		switch err {
		case utils.ErrUserNotFound:
			api.SendErrorResponse(w, http.StatusNotFound, "User not found", "User not found")
		case utils.ErrUserNotFound:
			api.SendErrorResponse(w, http.StatusNotFound, "User not found", "User not found")
		case utils.ErrPasswordTooShort:
			api.SendErrorResponse(w, http.StatusBadRequest, "Password too short", "Password must be at least 8 characters")
		case utils.ErrPasswordTooLong:
			api.SendErrorResponse(w, http.StatusBadRequest, "Password too long", "Password must be at most 72 characters")
		case utils.ErrPasswordSecurity:
			api.SendErrorResponse(w, http.StatusBadRequest, "Password not safe", "Password must contain at least one uppercase letter, one lowercase letter, one number, and one special character")
		case utils.ErrHashingPassword:
			api.SendErrorResponse(w, http.StatusInternalServerError, "Error hashing password", "Error hashing password")
		case utils.ErrUserNameTooShort:
			api.SendErrorResponse(w, http.StatusBadRequest, "User name too short", "User name must be at least 3 characters")
		case utils.ErrUserNameTooLong:
			api.SendErrorResponse(w, http.StatusBadRequest, "User name too long", "User name must be at most 200 characters")
		case utils.ErrUserNameWithNumericVals:
			api.SendErrorResponse(w, http.StatusBadRequest, "User name with numeric characters", "User name must not contain numeric characters")
		case utils.ErrDOBFormat:
			api.SendErrorResponse(w, http.StatusBadRequest, "Invalid date of birth format", "Invalid date of birth format")

		case utils.ErrDuplicateEmail:
			api.SendErrorResponse(w, http.StatusConflict, "Email already exists", "Email already exists")
		case utils.ErrGenerateOTP:
			api.SendErrorResponse(w, http.StatusInternalServerError, "Error generating OTP", "Error generating OTP")
		case utils.ErrCreateVericationEntry:
			api.SendErrorResponse(
				w,
				http.StatusInternalServerError,
				"Error creating verification entry",
				"Error creating verification entry",
			)
		case utils.ErrSMTPServerIssue:
			api.SendErrorResponse(w, http.StatusInternalServerError, "SMTP server issue", "SMTP server issue")
		default:
			api.SendErrorResponse(w, http.StatusInternalServerError, "Internal server error", err.Error())
		}

		return
	}

	api.SendSingleResponse(w, http.StatusCreated, "User pre-registered successfully", nil)
}

func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	// Check content type
	if r.Header.Get("Content-Type") != "application/json" {
		api.SendErrorResponse(w, http.StatusBadRequest, "Invalid content type", "Content type must be application/json")
		return
	}

	var input *schemas.VerifyOTPInput
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		api.SendErrorResponse(w, http.StatusBadRequest, "Invalid request body", err.Error())
		return
	}

	input.Email = strings.TrimSpace(input.Email)

	// Validate input email
	err = helpers.ValidateEmail(input.Email)
	if err != nil {
		switch err {
		case utils.ErrMissingEmail:
			api.SendErrorResponse(w, http.StatusBadRequest, "Missing email", "Please provide an email address")
		case utils.ErrInvalidEmail:
			api.SendErrorResponse(w, http.StatusBadRequest, "Invalid email", "Please provide a valid email address")
		default:
			api.SendErrorResponse(w, http.StatusInternalServerError, "Internal server error", err.Error())
		}

		return
	}

	// Verify OTP code
	err = h.userUseCase.VerifyOTPCode(r.Context(), input.Email, input.OTP)
	if err != nil {
		switch err {
		case utils.ErrPreRegistredUserNotFound:
			api.SendErrorResponse(w, http.StatusNotFound, "User not found", "User not found")
		case utils.ErrOTPExpired:
			api.SendErrorResponse(w, http.StatusBadRequest, "OTP code expired", "OTP code expired")
		case utils.ErrOTPAlreadyVerified:
			api.SendErrorResponse(w, http.StatusBadRequest, "OTP code already verified", "OTP code already verified")
		case utils.ErrInvalidOTP:
			api.SendErrorResponse(w, http.StatusBadRequest, "Invalid OTP code", "Invalid OTP code")
		case utils.ErrUserNotFound:
			api.SendErrorResponse(w, http.StatusNotFound, "User not found", "User not found")
		case utils.ErrDuplicateEmail:
			api.SendErrorResponse(w, http.StatusConflict, "Email already exists", "Email already exists")
		default:
			api.SendErrorResponse(w, http.StatusInternalServerError, "Internal server error", err.Error())
		}

		return
	}

	api.SendSingleResponse(w, http.StatusCreated, "User created successfully", nil)
}
