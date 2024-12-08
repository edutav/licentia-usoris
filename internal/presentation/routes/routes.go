package routes

import (
	"log"
	"net/http"
	"time"

	"github.com/edutav/licentia-usoris/internal/presentation/handlers"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	httpSwagger "github.com/swaggo/http-swagger"
)

// statusRecorder struct to record the status of the response
type statusRecorder struct {
	http.ResponseWriter
	statusCode int
}

// loggingMiddleware logs the request and response
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		logrus.WithFields(logrus.Fields{
			"method": r.Method,
			"url":    r.URL.Path,
			"remote": r.RemoteAddr,
			"agent":  r.UserAgent(),
		}).Info("Incoming HTTP request")

		rec := statusRecorder{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}

		next.ServeHTTP(&rec, r)

		logrus.WithFields(logrus.Fields{
			"status":   rec.statusCode,
			"duration": time.Since(start).String(),
		}).Info("Completed HTTP request")
	})
}

// NewRouter creates a new router
func NewRouter(
	indexHandler *handlers.IndexHandler,
	userHandler *handlers.UserHandler,
) http.Handler {
	log.Println("Settings up router...")

	r := mux.NewRouter()

	r.Use(loggingMiddleware) // logging

	r.PathPrefix("/docs/").Handler(httpSwagger.WrapHandler)

	// Routes v1
	prefixRouteV1 := r.PathPrefix("/api/v1").Subrouter()

	indexRouter := prefixRouteV1.PathPrefix("/").Subrouter()
	indexRouter.HandleFunc("/index", handlers.Index).Methods(http.MethodGet)

	// Routes for users
	userRouter := prefixRouteV1.PathPrefix("/user").Subrouter()
	userRouter.HandleFunc("/pre-register", userHandler.PreRegister).Methods(http.MethodPost)
	userRouter.HandleFunc("/register", userHandler.Register).Methods(http.MethodPost)

	log.Println("List all routes:")
	r.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		pathTemplate, err := route.GetPathTemplate()
		if err == nil {
			methods, _ := route.GetMethods()
			if len(methods) > 0 {
				log.Printf("Route found: %s - Methods: %v", pathTemplate, methods)
			}
		}
		return nil
	})

	log.Println("Setting up router... Done!")

	return r
}
