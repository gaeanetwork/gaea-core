package models

// Response return struct of the swagger api
type Response struct {
	Code    string      `json:"code,omitempty"`
	Message interface{} `json:"message,omitempty"`
	Result  interface{} `json:"result,omitempty"`
}
