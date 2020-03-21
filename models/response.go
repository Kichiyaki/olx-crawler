package models

type Response struct {
	Data   interface{} `json:"data,omitempty"`
	Errors []error     `json:"errors,omitempty"`
}
