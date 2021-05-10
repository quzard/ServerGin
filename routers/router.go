package routers

import (
	"ServerGin/controllers"
	"github.com/gin-gonic/gin"
)

func Init(router *gin.Engine) {
	router.POST("/login", controllers.LoginController)
}
