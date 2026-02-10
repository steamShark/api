package model_response

/*
Response of the API when it's successfully
*/
type ApiResponse[T any] struct {
	Success bool   `json:"success"`
	Data    T      `json:"data,omitempty"`
	Error   string `json:"error,omitempty"`
}
