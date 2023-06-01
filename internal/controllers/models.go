package controllers

type errorResponse struct {
	Error string `json:"error" example:"error message"`
}

type usernameResponse struct {
	Username string `json:"username" example:"user1"`
}
