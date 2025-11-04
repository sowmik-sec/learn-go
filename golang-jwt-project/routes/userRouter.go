package routes

import (
	controller "golang-jwt-project/controllers"
	"golang-jwt-project/middleware"

	"github.com/gin-gonic/gin"
)

func UserRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.GET("/users", middleware.Authenticate(), controller.GetUsers())
	incomingRoutes.GET("/users/:user_id", middleware.Authenticate(), controller.GetUser())
}
