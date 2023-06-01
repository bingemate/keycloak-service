package controllers

import "github.com/Nerzal/gocloak/v13"

type errorResponse struct {
	Error string `json:"error" example:"error message"`
}

type usernameResponse struct {
	Username string `json:"username" example:"user1"`
}

type userResponse struct {
	ID               string    `json:"id" example:"123e4567-e89b-12d3-a456-426614174000"`
	CreatedTimestamp int64     `json:"createdTimestamp" example:"1600000000000"`
	Username         string    `json:"username" example:"user1"`
	FirstName        string    `json:"firstName" example:"John"`
	LastName         string    `json:"lastName" example:"Doe"`
	Email            string    `json:"email" example:"example@email.com"`
	Roles            *[]string `json:"roles" example:"[\"admin\", \"user\"]"`
}

func toUserResponse(user *gocloak.User) *userResponse {
	response := &userResponse{
		ID:               *user.ID,
		CreatedTimestamp: *user.CreatedTimestamp,
		Username:         *user.Username,
		FirstName:        *user.FirstName,
		LastName:         *user.LastName,
		Email:            *user.Email,
	}
	if user.RealmRoles != nil {
		response.Roles = user.RealmRoles
	} else {
		response.Roles = &[]string{}
	}
	return response
}

func toUsersResponse(users []*gocloak.User) []*userResponse {
	var usersResponse []*userResponse
	for _, user := range users {
		usersResponse = append(usersResponse, toUserResponse(user))
	}
	return usersResponse
}
