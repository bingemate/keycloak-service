package controllers

import (
	"github.com/bingemate/keycloak-service/initializers"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func InitRouter(engine *gin.Engine, db *gorm.DB, env initializers.Env) {
}
