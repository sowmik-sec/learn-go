package routes

import (
	controller "golang-jwt-project/controllers"

	"github.com/gin-gonic/gin"
)

func AuthRoutes(router *gin.Engine) {
	router.POST("users/signup", controller.Signup())
	router.POST("users/login", controller.Login())
	router.POST("users/refresh", controller.RefreshToken())
}
