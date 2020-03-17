package models

type PaginatedResponse struct {
	Total int         `json:"total"`
	Items interface{} `json:"items"`
}
