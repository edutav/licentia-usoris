package reporitory

import (
	"context"

	"github.com/edutav/licentia-usoris/internal/domain/entity"
)

type UserRepository interface {
	// Pre-registration new user
	PreRegisterUser(ctx context.Context, preRegistration *entity.PreRegistration) error

	// Get user by email
	GetUserByEmail(ctx context.Context, email string) (*entity.User, error)

	// Create user
	CreateUser(ctx context.Context, user *entity.User) error

	// Get pre-registered user by email and OTP code
	GetPreRegisteredByEmailAndOTPCode(ctx context.Context, email, otpCode string) (*entity.PreRegistration, error)

	// Update user is verified
	UpdateUserIsVerified(ctx context.Context, email string) error

	// TODO: Update last login
}
