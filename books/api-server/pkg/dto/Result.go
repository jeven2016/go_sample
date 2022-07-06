package dto

type Result struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Payload any    `json:"payload"`
	Errors  any    `json:"errors"`
}
