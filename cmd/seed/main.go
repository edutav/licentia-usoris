package main

import (
	"context"
	"database/sql"
	"log"
	"strings"
	"time"

	"github.com/edutav/licentia-usoris/infrastructure/database"
	"github.com/edutav/licentia-usoris/internal/config"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	initSeed()
}

func initSeed() {
	println("init")
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Database connection
	db, err := database.NewConnectionPostgres(
		cfg.Database.Host, cfg.Database.Port, cfg.Database.User, cfg.Database.Password, cfg.Database.Name,
	)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	root := seedRoot(context.Background(), db)
	permissions := seedPermissionsDefault(context.Background(), db)
	roles := seedRolesDefault(context.Background(), db)
	seedRolesPermissionsDefault(context.Background(), db, roles, permissions)
	seedUserRolesDefault(context.Background(), db, roles, root)
}

type Root struct {
	uuid            string
	name            string
	email           string
	password        string
	isEmailVerified bool
	createedAt      time.Time
	updatedAt       time.Time
	isBlocked       bool
	isDeleted       bool
}

func seedRoot(ctx context.Context, db *sql.DB) *Root {
	// Hash the password
	// TODO: Move password to environment variable
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("Secret@123"), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("failed to hash password: %v", err)
	}

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		log.Fatalf("Failed to begin transaction: %v", err)
	}
	defer tx.Rollback()

	root := Root{
		name:            "root",
		email:           "root@mail.com",
		password:        string(hashedPassword),
		isEmailVerified: true,
		createedAt:      time.Now().UTC(),
		updatedAt:       time.Now().UTC(),
		isBlocked:       false,
		isDeleted:       false,
	}

	query := `
		INSERT INTO users (
			name, 
			email, 
			password_hash, 
			is_email_verified,
			created_at,
			updated_at,
			is_blocked,
			is_deleted
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING uuid
	`

	err = tx.QueryRowContext(ctx, query,
		root.name,
		root.email,
		root.password,
		root.isEmailVerified,
		root.createedAt,
		root.updatedAt,
		root.isBlocked,
		root.isDeleted,
	).Scan(&root.uuid)

	if err != nil {
		log.Fatalf("Failed to insert root user: %v", err)
	}

	if err := tx.Commit(); err != nil {
		log.Fatalf("Failed to commit transaction: %v", err)
	}

	log.Println("Root user created successfully")

	return &root
}

type Permission struct {
	name        string
	description string
}

func buildListPermissions() []*Permission {
	return []*Permission{
		{
			name:        "create:user",
			description: "Create user information",
		},
		{
			name:        "read:user",
			description: "Read user information",
		},
		{
			name:        "update:user",
			description: "Update user information",
		},
		{
			name:        "delete:user",
			description: "Delete a user",
		},
		{
			name:        "create:group",
			description: "Create a group information",
		},
		{
			name:        "read:group",
			description: "Read group information",
		},
		{
			name:        "update:group",
			description: "Update group information",
		},
		{
			name:        "delete:group",
			description: "Delete a group",
		},
		{
			name:        "create:role",
			description: "Create a role information",
		},
		{
			name:        "read:role",
			description: "Read role information",
		},
		{
			name:        "update:role",
			description: "Update role information",
		},
		{
			name:        "delete:role",
			description: "Delete a role information",
		},
		{
			name:        "create:permission",
			description: "Create a permission information",
		},
		{
			name:        "read:permission",
			description: "Read permission information",
		},
		{
			name:        "update:permission",
			description: "Update permission information",
		},
		{
			name:        "delete:permission",
			description: "Delete a permission information",
		},
	}
}

func seedPermissionsDefault(ctx context.Context, db *sql.DB) []*Permission {
	println("seedPermissionsDefault")

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		log.Fatalf("Failed to begin transaction: %v", err)
	}
	defer tx.Rollback()

	permissions := buildListPermissions()

	for _, permission := range permissions {
		query := `
			INSERT INTO permissions (
				name,
				description
			)
			VALUES ($1, $2)
		`

		_, err := tx.ExecContext(ctx, query, permission.name, permission.description)
		if err != nil {
			log.Fatalf("Failed to insert permission: %v", err)
		}
	}

	if err := tx.Commit(); err != nil {
		log.Fatalf("Failed to commit transaction: %v", err)
	}

	log.Println("Permissions created successfully")

	return permissions
}

type Role struct {
	name        string
	description string
}

func buildListRoles() []*Role {
	return []*Role{
		{
			name:        "admin",
			description: "Administrator",
		},
		{
			name:        "user",
			description: "User",
		},
	}
}

func seedRolesDefault(ctx context.Context, db *sql.DB) []*Role {
	println("seedRolesDefault")
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		log.Fatalf("Failed to begin transaction: %v", err)
	}
	defer tx.Rollback()

	roles := buildListRoles()

	for _, role := range roles {
		query := `
			INSERT INTO roles (
				name,
				description
			)
			VALUES ($1, $2)
		`

		_, err := tx.ExecContext(ctx, query, role.name, role.description)
		if err != nil {
			log.Fatalf("Failed to insert role: %v", err)
		}
	}

	if err := tx.Commit(); err != nil {
		log.Fatalf("Failed to commit transaction: %v", err)
	}

	log.Println("Roles created successfully")

	return roles
}

func seedRolesPermissionsDefault(ctx context.Context, db *sql.DB, roles []*Role, permissions []*Permission) {
	println("seedRolesPermissionsDefault")
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		log.Fatalf("Failed to begin transaction: %v", err)
	}
	defer tx.Rollback()

	for _, role := range roles {
		for _, permission := range permissions {
			if role.name == "user" && !strings.HasPrefix(permission.name, "read:") {
				continue
			}
			query := `
				INSERT INTO role_permissions (
					role_id,
					permission_id
				)
				VALUES (
					(SELECT uuid FROM roles WHERE name = $1),
					(SELECT uuid FROM permissions WHERE name = $2)
				)
			`

			_, err := tx.ExecContext(ctx, query, role.name, permission.name)
			if err != nil {
				log.Fatalf("Failed to insert role permission: %v", err)
			}
		}
	}

	if err := tx.Commit(); err != nil {
		log.Fatalf("Failed to commit transaction: %v", err)
	}

	log.Println("Roles permissions created successfully")
}

func seedUserRolesDefault(ctx context.Context, db *sql.DB, roles []*Role, root *Root) {
	println("seedUserRolesDefault")
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		log.Fatalf("Failed to begin transaction: %v", err)
	}
	defer tx.Rollback()

	for _, role := range roles {
		if root.name == "root" && role.name != "admin" {
			continue
		}
		query := `
			INSERT INTO user_roles (
				user_id,
				role_id
			)
			VALUES (
				(SELECT uuid FROM users WHERE email = $1),
				(SELECT uuid FROM roles WHERE name = $2)
			)
		`

		_, err := tx.ExecContext(ctx, query, root.email, role.name)
		if err != nil {
			log.Fatalf("Failed to insert user role: %v", err)
		}
	}

	if err := tx.Commit(); err != nil {
		log.Fatalf("Failed to commit transaction: %v", err)
	}

	log.Println("User roles created successfully")
}
