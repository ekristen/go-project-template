package api

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/swaggest/openapi-go"
)

type OpenAPIHandler interface {
	GetMethod() string
	GetPath() string
	SetupOpenAPIOperation(oc openapi.OperationContext) error
	ServeHTTP(http.ResponseWriter, *http.Request)
}

func Register(r *chi.Mux, handler OpenAPIHandler) {
	r.Method(handler.GetMethod(), handler.GetPath(), handler)
}

type Request struct {
}

type PaginationRequest struct {
	Page  int `json:"page" query:"page" description:"The page number you want to retrieve" default:"1"`
	Limit int `json:"limit" query:"limit" description:"The number of items per page" default:"20"`
}

type PaginationResponse struct {
	PaginationRequest
	TotalItems  int  `json:"total_items" query:"-" description:"The total number of items"`
	TotalPages  int  `json:"total_pages" query:"-" description:"The total number of pages"`
	HasNextPage bool `json:"has_next_page" query:"-" description:"Whether there is a next page"`
}

type ErrorResponse struct {
	Type     string `json:"type,omitempty"`
	Title    string `json:"title,omitempty"`
	Detail   string `json:"detail,omitempty"`
	Status   int    `json:"status,omitempty"`
	Instance string `json:"instance,omitempty"`
}

type MetadataResponse struct {
	Count int `json:"count"`
}

type DataResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
}

type Response struct {
	Pagination *PaginationResponse `json:"pagination,omitempty"`
	Metadata   *MetadataResponse   `json:"metadata,omitempty"`
}

func (r *Response) Respond(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(r)
	if err != nil {
		r.HandleError(w, http.StatusInternalServerError, err.Error())
		return
	}
}

func (r *Response) HandleError(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/problem+json")
	w.WriteHeader(statusCode)
	err := json.NewEncoder(w).Encode(ErrorResponse{
		Detail: message,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
