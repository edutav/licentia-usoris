package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

// Error response is the response object for the API
type ErrorResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message,omitempty"`
	Error   string `json:"error,omitempty"`
}

// Single response object for the API
type SingleResponse struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// List response object for the API
type ListResponse struct {
	Status int           `json:"status"`
	Data   []interface{} `json:"data,omitempty"`
	Meta   Meta          `json:"meta,omitempty"`
	Links  Links         `json:"links,omitempty"`
}

// Meta response object for the API
type Meta struct {
	TotalItems   int `json:"total_items"`
	TotalPages   int `json:"total_pages"`
	CurrentPage  int `json:"current_page"`
	ItemsPerPage int `json:"items_per_page"`
}

// Links represents the pagination links in the response
type Links struct {
	Self     string `json:"self"`
	Next     string `json:"next"`
	Previous string `json:"previous"`
	First    string `json:"first"`
	Last     string `json:"last"`
}

// PaginationParams holds the parameters for pagination
type PaginationParams struct {
	Page     int
	PageSize int
}

// SendErrorResponse sends an error response
func SendErrorResponse(w http.ResponseWriter, status int, message string, err string) {
	response := ErrorResponse{
		Status:  status,
		Message: message,
		Error:   err,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(response)
}

// SendSingleResponse sends a single object response
func SendSingleResponse(w http.ResponseWriter, status int, message string, data interface{}) {
	response := SingleResponse{
		Status:  status,
		Message: message,
		Data:    data,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(response)
}

// SendListResponse sends a paginated list response
func SendPaginatedResponse(w http.ResponseWriter, status int, data []interface{}, totalItems int, page, pageSize int, baseURL string) {
	meta, links := BuildPagination(totalItems, page, pageSize, baseURL)
	response := ListResponse{
		Status: status,
		Data:   data,
		Meta:   meta,
		Links:  links,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(response)
}

// GetPaginationParams extracts pagination parameters from the request
func GetPaginationParams(r *http.Request) PaginationParams {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}

	pageSize, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if pageSize < 1 {
		pageSize = 10
	}

	return PaginationParams{
		Page:     page,
		PageSize: pageSize,
	}
}

// BuildPagination constructs the pagination metadata and links
func BuildPagination(totalItems, page, pageSize int, baseURL string) (Meta, Links) {
	totalPages := (totalItems + pageSize - 1) / pageSize

	meta := Meta{
		TotalItems:   totalItems,
		TotalPages:   totalPages,
		CurrentPage:  page,
		ItemsPerPage: pageSize,
	}

	links := Links{
		Self:  fmt.Sprintf("%s?page=%d&limit=%d", baseURL, page, pageSize),
		First: fmt.Sprintf("%s?page=1&limit=%d", baseURL, pageSize),
		Last:  fmt.Sprintf("%s?page=%d&limit=%d", baseURL, totalPages, pageSize),
	}

	if page > 1 {
		links.Previous = fmt.Sprintf("%s?page=%d&limit=%d", baseURL, page-1, pageSize)
	}

	if page < totalPages {
		links.Next = fmt.Sprintf("%s?page=%d&limit=%d", baseURL, page+1, pageSize)
	}

	return meta, links
}
