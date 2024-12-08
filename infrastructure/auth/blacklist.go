package auth

import (
	"context"
	"database/sql"
	"log"
	"time"
)

// Blacklist struct
type Blacklist struct {
	db *sql.DB
}

// NewBlacklist creates a new blacklist
func NewBlacklist(db *sql.DB) *Blacklist {
	return &Blacklist{
		db: db,
	}
}

// Add adds a token to the blacklist
func (b *Blacklist) Add(ctx context.Context, token string, expiresAt time.Time) (bool, error) {
	query := `SELECT EXISTS (SELECT 1 FROM blacklisted_tokens WHERE token = $1 AND expires_at > NOW()`
	var exists bool
	err := b.db.QueryRowContext(ctx, query, token).Scan(&exists)
	if err != nil {
		log.Printf("Erro ao verificar se o token está na blacklist: %v", err)
	}
	return exists, err
}

// IsBlacklisted checks if a token is blacklisted
func (b *Blacklist) IsBlacklisted(ctx context.Context, token string) (bool, error) {
	query := `SELECT COUNT(*) FROM blacklist WHERE token = $1 AND expires_at > NOW()`
	var count int
	err := b.db.QueryRowContext(ctx, query, token).Scan(&count)
	if err != nil {
		log.Printf("Error checking if token is blacklisted: %v", err)
		return false, err
	}

	return count > 0, nil
}

// Remove removes a token from the blacklist
func (b *Blacklist) Remove(ctx context.Context) error {
	query := `DELETE FROM blacklisted_tokens WHERE expires_at <= NOW()`
	_, err := b.db.ExecContext(ctx, query)
	if err != nil {
		log.Printf("Error removing token from blacklist: %v", err)
	}
	return err
}
