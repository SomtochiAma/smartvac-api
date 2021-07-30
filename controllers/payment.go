package controllers

import (
	"fmt"
	"gorm.io/gorm"
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

	if err := models.DB.Transaction(CreatePaymentTx(payment)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "payment successful",
	})
}

func GetPaymentHistory(c *gin.Context) {
	id := c.Param("id")
	var history []models.Payment
	result := models.DB.Where("user_id = ?", id).Find(&history)
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

func CreatePaymentTx(payment models.Payment) func(tx *gorm.DB) error {
	return func(tx *gorm.DB) error {
		if err := tx.Create(&payment).Error; err != nil {
			return err
		}

		var usage struct {
			UsedUnit  int
			TotalUnit int
		}

		if err := models.DB.Model(&models.User{}).
			Where("id = ?", payment.UserID).
			Select("id", "total_unit", "used_unit").
			Take(&usage).Error; err != nil {
			return err
		}
		fmt.Println(usage)

		newTotal := payment.Units + (usage.TotalUnit - usage.UsedUnit)
		if err := models.DB.Model(&models.User{}).
			Where("id = ?", payment.UserID).
			Updates(models.User{TotalUnit: newTotal, UsedUnit: 0}).Error; err != nil {
			return err
		}

		return nil
	}
}
