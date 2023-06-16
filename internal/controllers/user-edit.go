package controllers

import (
	"github.com/bingemate/keycloak-service/internal/features"
	"github.com/gin-gonic/gin"
)

func InitUserEditController(engine *gin.RouterGroup, userEditService *features.UserEditService) {
	engine.PUT("", func(c *gin.Context) {
		updateUser(c, userEditService)
	})
	engine.PUT("/password", func(c *gin.Context) {
		updateUserPassword(c, userEditService)
	})
	engine.DELETE("", func(c *gin.Context) {
		deleteUser(c, userEditService)
	})
}

// @Summary Update user
// @Description Update user infos
// @Tags User Edit
// @Param user-id header string true "User ID"
// @Param userUpdateRequest body userUpdateRequest true "User update request"
// @Produce json
// @Success 200 {object} userResponse
// @Failure 400 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /user-edit [put]
func updateUser(c *gin.Context, userEditService *features.UserEditService) {
	userID := c.GetHeader("user-id")
	if userID == "" {
		c.JSON(400, errorResponse{Error: "user-id must not be empty"})
		return
	}
	var userUpdateRequest userUpdateRequest
	err := c.BindJSON(&userUpdateRequest)
	if err != nil {
		c.JSON(400, errorResponse{Error: err.Error()})
		return
	}
	user, err := userEditService.UpdateUser(userID, userUpdateRequest.Username, userUpdateRequest.FirstName, userUpdateRequest.LastName, userUpdateRequest.Email)
	if err != nil {
		c.JSON(500, errorResponse{Error: err.Error()})
		return
	}
	c.JSON(200, toUserResponse(user))
}

// @Summary Update user password
// @Description Update user password
// @Tags User Edit
// @Param user-id header string true "User ID"
// @Param userPasswordUpdateRequest body updatePasswordRequest true "User password update request"
// @Produce json
// @Success 200 {string} string "Password updated"
// @Failure 400 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /user-edit/password [put]
func updateUserPassword(c *gin.Context, userEditService *features.UserEditService) {
	userID := c.GetHeader("user-id")
	if userID == "" {
		c.JSON(400, errorResponse{Error: "user-id must not be empty"})
		return
	}
	var updatePasswordRequest updatePasswordRequest
	err := c.BindJSON(&updatePasswordRequest)
	if err != nil {
		c.JSON(400, errorResponse{Error: err.Error()})
		return
	}
	err = userEditService.UpdateUserPassword(userID, updatePasswordRequest.Password)
	if err != nil {
		c.JSON(500, errorResponse{Error: err.Error()})
		return
	}
	c.JSON(200, "Password updated")
}

// @Summary Delete user
// @Description Delete user
// @Tags User Edit
// @Param user-id header string true "User ID"
// @Produce json
// @Success 200 {string} string "User deleted"
// @Failure 400 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /user-edit [delete]
func deleteUser(c *gin.Context, userEditService *features.UserEditService) {
	userID := c.GetHeader("user-id")
	if userID == "" {
		c.JSON(400, errorResponse{Error: "user-id must not be empty"})
		return
	}
	err := userEditService.DeleteUser(userID)
	if err != nil {
		c.JSON(500, errorResponse{Error: err.Error()})
		return
	}
	c.JSON(200, "User deleted")
}
