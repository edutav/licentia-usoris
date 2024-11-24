package server

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/edutav/licentia-usoris/infrastructure/email"
	"github.com/edutav/licentia-usoris/internal/config"
	"github.com/edutav/licentia-usoris/internal/presentation/handlers"
	"github.com/edutav/licentia-usoris/internal/presentation/routes"
)

type Server struct {
	router http.Handler
}

func NewServer(db *sql.DB, emailSender *email.Sender, cfg *config.Config) *Server {
	log.Println("Initializing components for server")

	indexHandler := handlers.NewIndexHandler()

	// Create router
	router := routes.NewRouter(indexHandler)
	log.Println("Router created")

	return &Server{
		router: router,
	}
}

// ServerHttp starts the server
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}
