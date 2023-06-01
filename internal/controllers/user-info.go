package controllers

import (
	"github.com/bingemate/keycloak-service/internal/features"
	"github.com/gin-gonic/gin"
)

func InitUserInfoController(engine *gin.RouterGroup, userInfoService *features.UserInfoService) {
	engine.GET("/:userID/username", func(c *gin.Context) {
		getUsername(c, userInfoService)
	})
}

// @Summary Get user's username
// @Description Get user's username
// @Tags User
// @Param userID path string true "User ID"
// @Produce json
// @Success 200 {object} usernameResponse
// @Failure 400 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /user-info/{userID}/username [get]
func getUsername(c *gin.Context, userInfoService *features.UserInfoService) {
	userID := c.Param("userID")
	if userID == "" {
		c.JSON(400, errorResponse{Error: "userID must not be empty"})
		return
	}
	username, err := userInfoService.GetUsername(userID)
	if err != nil {
		c.JSON(500, errorResponse{Error: err.Error()})
		return
	}
	c.JSON(200, usernameResponse{Username: username})
}
