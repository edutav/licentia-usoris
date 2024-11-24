package handlers

import "net/http"

// IndexHandler is the handler for index related operations
type IndexHandler struct {
	version string
}

// NewIndexHandler creates a new index handler
func NewIndexHandler() *IndexHandler {
	return &IndexHandler{}
}

// Index function is the handler for the index route
func Index(w http.ResponseWriter, r *http.Request) {
	version := IndexHandler{
		version: "v0.0.1",
	}
	// Return the version formatted as JSON
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"version": "` + version.version + `"}`))
}
