package controllers

import (
	"github.com/bingemate/keycloak-service/initializers"
	"github.com/bingemate/keycloak-service/internal/features"
	"github.com/gin-gonic/gin"
)

// func InitRouter(engine *gin.Engine, db *gorm.DB, env initializers.Env) {
func InitRouter(engine *gin.Engine, keycloakClient *initializers.KeycloakClient) {
	var keycloakServiceGroup = engine.Group("/keycloak-service")
	var userInfoService = features.NewUserInfoService(keycloakClient)
	InitPingController(keycloakServiceGroup.Group("/ping"))
	InitUserInfoController(keycloakServiceGroup.Group("/user-info"), userInfoService)
}
