package controllers

import (
	"github.com/bingemate/keycloak-service/internal/features"
	"github.com/gin-gonic/gin"
	"strconv"
)

func InitUserAdminController(engine *gin.RouterGroup, userEditService *features.UserEditService, userInfoService *features.UserInfoService) {
	engine.GET("/search", func(c *gin.Context) {
		adminGetUsers(c, userInfoService)
	})
	engine.PUT("/edit/:userID", func(c *gin.Context) {
		adminUpdateUser(c, userEditService)
	})
	engine.PUT("/edit/:userID/password", func(c *gin.Context) {
		adminUpdateUserPassword(c, userEditService)
	})
	engine.PUT("/roles/:userID", func(c *gin.Context) {
		addRoleToUser(c, userEditService)
	})
	engine.DELETE("/roles/:userID", func(c *gin.Context) {
		removeRoleFromUser(c, userEditService)
	})
	engine.GET("/roles", func(c *gin.Context) {
		getRoles(c, userInfoService)
	})
	engine.DELETE("/delete/:userID", func(c *gin.Context) {
		adminDeleteUser(c, userEditService)
	})
}

// @Summary Search users
// @Description Search users
// @Tags User Admin
// @Param query query string false "Username"
// @Param page query int false "Page size"
// @Param limit query int false "Page number"
// @Produce json
// @Success 200 {object} userResults
// @Failure 400 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /user-admin/search [get]
func adminGetUsers(c *gin.Context, userInfoService *features.UserInfoService) {
	query := c.Query("query")
	pageSize, err := strconv.Atoi(c.Query("limit"))
	if err != nil {
		pageSize = 10
	}
	pageNumber, err := strconv.Atoi(c.Query("page"))
	if err != nil {
		pageNumber = 1
	}

	users, count, err := userInfoService.GetUsers(query, pageNumber, pageSize)
	if err != nil {
		c.JSON(500, errorResponse{Error: err.Error()})
		return
	}
	c.JSON(200, userResults{Results: toUsersResponse(users), Total: count})
}

// @Summary Update user
// @Description Update user
// @Tags User Admin
// @Param userID path string true "User ID"
// @Param user body userUpdateRequest true "User"
// @Produce json
// @Success 200 {object} userResponse
// @Failure 400 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /user-admin/edit/{userID} [put]
func adminUpdateUser(c *gin.Context, userEditService *features.UserEditService) {
	userID := c.Param("userID")
	if userID == "" {
		c.JSON(400, errorResponse{Error: "userID must not be empty"})
		return
	}
	var user userUpdateRequest
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(400, errorResponse{Error: err.Error()})
		return
	}
	updatedUser, err := userEditService.UpdateUser(
		userID,
		user.Username,
		user.FirstName,
		user.LastName,
		user.Email,
	)
	if err != nil {
		c.JSON(500, errorResponse{Error: err.Error()})
		return
	}
	c.JSON(200, toUserResponse(updatedUser))
}

// @Summary Update user password
// @Description Update user password
// @Tags User Admin
// @Param userID path string true "User ID"
// @Param password body updatePasswordRequest true "Password"
// @Produce json
// @Success 200 {string} string "Password updated"
// @Failure 400 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /user-admin/edit/{userID}/password [put]
func adminUpdateUserPassword(c *gin.Context, userEditService *features.UserEditService) {
	userID := c.Param("userID")
	if userID == "" {
		c.JSON(400, errorResponse{Error: "userID must not be empty"})
		return
	}
	var password updatePasswordRequest
	if err := c.ShouldBindJSON(&password); err != nil {
		c.JSON(400, errorResponse{Error: err.Error()})
		return
	}
	if err := userEditService.UpdateUserPassword(userID, password.Password); err != nil {
		c.JSON(500, errorResponse{Error: err.Error()})
		return
	}
	c.JSON(200, "Password updated")
}

// @Summary Add role to user
// @Description Add role to user
// @Tags User Admin
// @Param userID path string true "User ID"
// @Param role body roleRequest true "Role"
// @Produce json
// @Success 200 {string} string "Role added"
// @Failure 400 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /user-admin/roles/{userID} [put]
func addRoleToUser(c *gin.Context, userEditService *features.UserEditService) {
	userID := c.Param("userID")
	if userID == "" {
		c.JSON(400, errorResponse{Error: "userID must not be empty"})
		return
	}
	var role roleRequest
	if err := c.ShouldBindJSON(&role); err != nil {
		c.JSON(400, errorResponse{Error: err.Error()})
		return
	}
	if err := userEditService.AddUserRole(userID, role.Role); err != nil {
		c.JSON(500, errorResponse{Error: err.Error()})
		return
	}
	c.JSON(200, "Role added")
}

// @Summary Remove role from user
// @Description Remove role from user
// @Tags User Admin
// @Param userID path string true "User ID"
// @Param role body roleRequest true "Role"
// @Produce json
// @Success 200 {string} string "Role removed"
// @Failure 400 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /user-admin/roles/{userID} [delete]
func removeRoleFromUser(c *gin.Context, userEditService *features.UserEditService) {
	userID := c.Param("userID")
	if userID == "" {
		c.JSON(400, errorResponse{Error: "userID must not be empty"})
		return
	}
	var role roleRequest
	if err := c.ShouldBindJSON(&role); err != nil {
		c.JSON(400, errorResponse{Error: err.Error()})
		return
	}
	if err := userEditService.RemoveUserRole(userID, role.Role); err != nil {
		c.JSON(500, errorResponse{Error: err.Error()})
		return
	}
	c.JSON(200, "Role removed")
}

// @Summary Get roles
// @Description Get roles
// @Tags User Admin
// @Produce json
// @Success 200 {array} string
// @Failure 500 {object} errorResponse
// @Router /user-admin/roles [get]
func getRoles(c *gin.Context, userInfoService *features.UserInfoService) {
	roles, err := userInfoService.GetAvailableRoles()
	if err != nil {
		c.JSON(500, errorResponse{Error: err.Error()})
		return
	}
	c.JSON(200, roles)
}

// @Summary Delete user
// @Description Delete user
// @Tags User Admin
// @Param userID path string true "User ID"
// @Produce json
// @Success 200 {string} string "User deleted"
// @Failure 400 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /user-admin/delete/{userID} [delete]
func adminDeleteUser(c *gin.Context, userEditService *features.UserEditService) {
	userID := c.Param("userID")
	if userID == "" {
		c.JSON(400, errorResponse{Error: "userID must not be empty"})
		return
	}
	if err := userEditService.DeleteUser(userID); err != nil {
		c.JSON(500, errorResponse{Error: err.Error()})
		return
	}
	c.JSON(200, "User deleted")
}
