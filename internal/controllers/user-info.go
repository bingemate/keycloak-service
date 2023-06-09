package controllers

import (
	"github.com/bingemate/keycloak-service/internal/features"
	"github.com/gin-gonic/gin"
)

func InitUserInfoController(engine *gin.RouterGroup, userInfoService *features.UserInfoService) {
	engine.GET("/username/:userID", func(c *gin.Context) {
		getUsername(c, userInfoService)
	})
	engine.GET("/search", func(c *gin.Context) {
		searchUsers(c, userInfoService)
	})
	engine.GET("/:userID", func(c *gin.Context) {
		getUser(c, userInfoService)
	})
	engine.GET("/me", func(c *gin.Context) {
		getMe(c, userInfoService)
	})
}

// @Summary Get user's username
// @Description Get user's username
// @Tags User Info
// @Param userID path string true "User ID"
// @Produce json
// @Success 200 {object} usernameResponse
// @Failure 400 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /user-info/username/{userID} [get]
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

// @Summary Search users
// @Description Search users
// @Tags User Info
// @Param query query string true "Username"
// @Param includeRoles query bool false "Include roles"
// @Produce json
// @Success 200 {array} userResponse
// @Failure 400 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /user-info/search [get]
func searchUsers(c *gin.Context, userInfoService *features.UserInfoService) {
	query := c.Query("query")
	if query == "" {
		c.JSON(400, errorResponse{Error: "query must not be empty"})
		return
	}
	includeRoles := c.Query("includeRoles") == "true"

	users, err := userInfoService.SearchUsers(query, includeRoles)
	if err != nil {
		c.JSON(500, errorResponse{Error: err.Error()})
		return
	}
	c.JSON(200, toUsersResponse(users))
}

// @Summary Get user
// @Description Get user
// @Tags User Info
// @Param userID path string true "User ID"
// @Produce json
// @Success 200 {object} userResponse
// @Failure 400 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /user-info/{userID} [get]
func getUser(c *gin.Context, userInfoService *features.UserInfoService) {
	userID := c.Param("userID")
	if userID == "" {
		c.JSON(400, errorResponse{Error: "userID must not be empty"})
		return
	}
	user, err := userInfoService.GetUser(userID)
	if err != nil {
		c.JSON(500, errorResponse{Error: err.Error()})
		return
	}
	c.JSON(200, toUserResponse(user))
}

// @Summary Get current user
// @Description Get current user
// @Tags User Info
// @Param user-id header string true "User ID"
// @Produce json
// @Success 200 {object} userResponse
// @Failure 500 {object} errorResponse
// @Router /user-info/me [get]
func getMe(c *gin.Context, userInfoService *features.UserInfoService) {
	userID := c.GetHeader("user-id")
	if userID == "" {
		c.JSON(500, errorResponse{Error: "user-id must not be empty"})
		return
	}
	user, err := userInfoService.GetUser(userID)
	if err != nil {
		c.JSON(500, errorResponse{Error: err.Error()})
		return
	}
	c.JSON(200, toUserResponse(user))
}
