package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"

	"github.com/edutav/licentia-usoris/internal/domain/entity"
	"github.com/edutav/licentia-usoris/internal/domain/reporitory"
	"github.com/edutav/licentia-usoris/internal/utils"
	"github.com/lib/pq"
)

type userRepository struct {
	db *sql.DB
}

// NewUserRepository creates a new instance of UserRepository
func NewUserRepository(db *sql.DB) reporitory.UserRepository {
	return &userRepository{
		db: db,
	}
}

// PreRegisterUser pre-registers a new user
func (repo *userRepository) PreRegisterUser(
	ctx context.Context, preRegistration *entity.PreRegistration,
) error {
	tx, err := repo.db.BeginTx(ctx, nil)
	if err != nil {
		log.Printf("Error beginning transaction: %v", err)
		return err
	}
	defer tx.Rollback()

	query := `
		INSERT INTO pre_registrations (
			email, 
			password_hash, 
			code_otp, 
			user_data, 
			expires_at, 
			is_verified, 
			created_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	userDataJSON, err := json.Marshal(preRegistration.UserData)
	if err != nil {
		log.Printf("Error marshalling user data: %v", err)
		return err
	}

	err = tx.QueryRowContext(ctx, query,
		preRegistration.Email,
		preRegistration.PasswordHash,
		preRegistration.CodeOTP,
		userDataJSON,
		preRegistration.ExpiresAt,
		preRegistration.IsVerified,
		preRegistration.CreatedAt,
	).Err()

	if err != nil {
		pqErr, ok := err.(*pq.Error)
		if ok {
			switch pqErr.Code {
			case "23505":
				if pqErr.Constraint == "pre_registrations_email_key" {
					return utils.ErrDuplicateEmail
				}
			}
		}

		log.Printf("Error inserting pre-registration: %v", err)
		return err
	}

	if err := tx.Commit(); err != nil {
		log.Printf("Error committing transaction: %v", err)
		return err
	}

	return err
}

// GetUserByEmail gets a user by email
func (repo *userRepository) GetUserByEmail(ctx context.Context, email string) (*entity.User, error) {
	query := `
		SELECT 
			uuid, 
			name, 
			email, 
			password_hash, 
			date_of_birth, 
			phone_number,
			is_blocked, 
			is_email_verified, 
			created_at, 
			updated_at, 
			deleted_at,
			is_deleted, 
			last_login
		FROM
			users
		WHERE
			email = $1
		LIMIT 1
	`

	user := &entity.User{}

	err := repo.db.QueryRowContext(ctx, query, email).Scan(
		&user.UUID,
		&user.Name,
		&user.Email,
		&user.PasswordHash,
		&user.DOB,
		&user.PhoneNumber,
		&user.IsBlocked,
		&user.IsEmailVerified,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.DeletedAt,
		&user.IsDeleted,
		&user.LastLogin,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, utils.ErrUserNotFound
		}

		log.Printf("Error getting user by email: %v", err)
		return nil, err
	}

	return user, err
}

// CreateUser creates a new user
func (repo *userRepository) CreateUser(ctx context.Context, user *entity.User) error {
	tx, err := repo.db.BeginTx(ctx, nil)
	if err != nil {
		log.Printf("Error beginning transaction: %v", err)
		return err
	}
	defer tx.Rollback()

	query := `
		INSERT INTO users (
			name,
			email,
			password_hash,
			date_of_birth,
			phone_number,
			is_blocked,
			is_email_verified,
			created_at,
			updated_at,
			deleted_at,
			is_deleted,
			last_login
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`

	err = tx.QueryRowContext(ctx, query,
		user.Name,
		user.Email,
		user.PasswordHash,
		user.DOB,
		user.PhoneNumber,
		user.IsBlocked,
		user.IsEmailVerified,
		user.CreatedAt,
		user.UpdatedAt,
		user.DeletedAt,
		user.IsDeleted,
		user.LastLogin,
	).Err()

	if err != nil {
		pqErr, ok := err.(*pq.Error)
		if ok {
			switch pqErr.Code {
			case "23505":
				if pqErr.Constraint == "users_email_key" {
					return utils.ErrDuplicateEmail
				}
			}
		}

		log.Printf("Error inserting user: %v", err)
		return err
	}

	if err := tx.Commit(); err != nil {
		log.Printf("Error committing transaction: %v", err)
		return err
	}

	return nil
}

// GetPreRegisteredByEmailAndOTPCode gets a pre-registered user by email and OTP code
func (repo *userRepository) GetPreRegisteredByEmailAndOTPCode(ctx context.Context, email, otpCode string) (*entity.PreRegistration, error) {
	query := `
		SELECT
			uuid,
			email,
			password_hash,
			code_otp,
			user_data,
			expires_at,
			is_verified,
			created_at
		FROM
			pre_registrations
		WHERE
			email = $1
		LIMIT 1`

	preRegistration := &entity.PreRegistration{}
	var userDataJSON []byte

	err := repo.db.QueryRowContext(ctx, query, email).Scan(
		&preRegistration.UUID,
		&preRegistration.Email,
		&preRegistration.PasswordHash,
		&preRegistration.CodeOTP,
		&userDataJSON,
		&preRegistration.ExpiresAt,
		&preRegistration.IsVerified,
		&preRegistration.CreatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, utils.ErrPreRegistredUserNotFound
		}
		log.Printf("Error getting pre-registered user by email and OTP code: %v", err)
		return nil, err
	}

	err = json.Unmarshal(userDataJSON, &preRegistration.UserData)
	if err != nil {
		log.Printf("Error unmarshalling user data: %v", err)
		return nil, err
	}

	return preRegistration, nil

}

// UpdateUserIsVerified updates the user is verified status
func (repo *userRepository) UpdateUserIsVerified(ctx context.Context, email string) error {
	tx, err := repo.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("error beginning transaction: %w", err)
	}

	query := `
		UPDATE
			pre_registrations
		SET
			is_verified = true
		WHERE
			email = $1`

	_, err = tx.ExecContext(ctx, query, email)
	if err != nil {
		tx.Rollback()
		log.Printf("Error updating user is verified: %v", err)
		return err
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("error committing transaction: %w", err)
	}

	return err
}
