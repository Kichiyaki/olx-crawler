package models

type Response struct {
	Data  interface{} `json:"data,omitempty"`
	Error error       `json:"error,omitempty"`
}
