package utils

type SuccessResponse struct {
	Message string `json:"message"`
}

type UserResponse struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}
