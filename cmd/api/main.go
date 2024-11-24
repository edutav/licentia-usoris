package main

import (
	"log"
	"net/http"
	"time"

	"github.com/edutav/licentia-usoris/infrastructure/database"
	"github.com/edutav/licentia-usoris/infrastructure/email"
	"github.com/edutav/licentia-usoris/infrastructure/server"
	"github.com/edutav/licentia-usoris/internal/config"
)

func main() {
	// Define timezone default to UTC for the application
	time.Local = time.UTC

	log.Println("Starting Aplication")

	cfg, err := config.Load()
	if err != nil {
		log.Printf("Error loading config: %s", err)
	}

	db, err := database.NewConnectionPostgres(
		cfg.Database.Host, cfg.Database.Port, cfg.Database.User, cfg.Database.Password, cfg.Database.Name,
	)
	if err != nil {
		log.Fatalf("Error connecting to database: %s", err)
	}
	defer db.Close()

	// Initialize email sender
	emailSender := email.NewEmailSender(cfg.SMTP.Host, cfg.SMTP.Port, cfg.SMTP.Username, cfg.SMTP.Password)

	server := server.NewServer(db, emailSender, cfg)

	log.Println("Server is running on port", cfg.Server.Port)
	if err := http.ListenAndServe(":"+cfg.Server.Port, server); err != nil {
		log.Fatalf("Error starting server: %s", err)
	}
}
