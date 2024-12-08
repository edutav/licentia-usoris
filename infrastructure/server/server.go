package server

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/edutav/licentia-usoris/infrastructure/email"
	"github.com/edutav/licentia-usoris/internal/config"
	"github.com/edutav/licentia-usoris/internal/domain/reporitory/postgres"
	"github.com/edutav/licentia-usoris/internal/presentation/handlers"
	"github.com/edutav/licentia-usoris/internal/presentation/routes"
	"github.com/edutav/licentia-usoris/internal/usecases"
	"github.com/edutav/licentia-usoris/internal/usecases/validator"
)

type Server struct {
	router http.Handler
}

func NewServer(db *sql.DB, emailSender *email.Sender, cfg *config.Config) *Server {
	log.Println("Initializing components for server")

	indexHandler := handlers.NewIndexHandler()

	// Components the users
	userRepository := postgres.NewUserRepository(db)
	userUseCase := usecases.NewUserUseCase(userRepository, emailSender, validator.ValidateUserPassword)
	userHandler := handlers.NewUserHandler(userUseCase)

	// Create router
	router := routes.NewRouter(indexHandler, userHandler)
	log.Println("Router created")

	return &Server{
		router: router,
	}
}

// ServerHttp starts the server
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}
