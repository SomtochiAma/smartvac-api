package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/SomtochiAma/smartvac-api/controllers"
)

func Init() *gin.Engine {
	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Welcome to Smartvac API",
		})
	})

	r.POST("/signup", controllers.CreateUser)
	r.POST("signin", controllers.Signin)

	r.GET("/user/:id", controllers.GetUser)
	r.PUT("/user/:id", controllers.UpdateUser)

	r.GET("/history/:id", controllers.GetPaymentHistory)
	r.GET("/summary/:id", controllers.GetTotalReading)
	r.POST("/pay", controllers.MakePayment)
	r.POST("/data", controllers.PostReading)

	r.GET("/ws", controllers.WebSocket)

	return r
}
