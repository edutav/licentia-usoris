package handlers

import (
	"net/http"

	"github.com/edutav/licentia-usoris/infrastructure/server/api"
)

// IndexHandler is the handler for index related operations
// @Description IndexHandler is the handler for index related operations
// @Param version body string true "version"
type IndexHandler struct{}

// NewIndexHandler creates a new index handler
func NewIndexHandler() *IndexHandler {
	return &IndexHandler{}
}

// Index function is the handler for the index route
// @Summary Get the API version
// @Description Get the API version
// @Tags index
// @Accept json
// @Produce json
// @Success 200 {object} api.SingleResponse  "API version"
// @Router /index [get]
func Index(w http.ResponseWriter, r *http.Request) {
	api.SendSingleResponse(
		w,
		http.StatusOK,
		"API version",
		map[string]interface{}{
			"version": "v0.1.0",
		},
	)
}
