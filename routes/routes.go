package routes

import (
	"github.com/SomtochiAma/smartvac-api/controllers"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Init() *gin.Engine{
	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Welcome to Smartvac API",
		})
	})

	r.POST("/signup", controllers.CreateUser)
	r.POST("signin", controllers.Signin)

	r.POST("/data", controllers.PostReading)
	r.GET("/data", controllers.GetReading)

	r.GET("/user/:id", controllers.GetUser)
	r.PUT("/user/:id", controllers.UpdateUser)

	return r
}
