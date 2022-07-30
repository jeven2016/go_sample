package dto

type Result struct {
	Code    int    `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
	Payload any    `json:"payload,omitempty"`
	Errors  any    `json:"errors,omitempty"`
}
