package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/SomtochiAma/smartvac-api/models"
)

func MakePayment(c *gin.Context) {
	var payment models.Payment

	err := c.ShouldBindJSON(&payment)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	models.DB.Create(&payment)
	c.JSON(http.StatusOK, gin.H{
		"message": "payment successful",
	})
}

func GetPaymentHistory(c *gin.Context) {
	var history []models.Payment
	result := models.DB.Find(&history)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": result.Error.Error(),
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"data":    history,
		"message": "payment history retrieved successfully",
	})
}
