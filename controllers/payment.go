package controllers

import (
	"context"
	"fmt"
	"github.com/mailgun/mailgun-go/v4"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"

	"github.com/SomtochiAma/smartvac-api/models"
)

type EmailVars struct {
	email  string
	amount int
	units  int
	date   string
	name   string
}

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

	var userDetails struct {
		ID    int
		Name  string
		Email string `json:"email"`
	}
	models.DB.Model(&models.User{}).Select("id, name, email").Where("id = ?", payment.UserID).First(&userDetails)
	var domain = "sandbox869c162826e04f39b909a85aeedcfc65.mailgun.org"
	apiKey := os.Getenv("EMAIL_API_KEY")
	_, err = SendSimpleMessage(domain, apiKey, EmailVars{
		email:  userDetails.Email,
		amount: payment.Amount,
		units:  payment.Units,
		date:   payment.Day.Format("02-01-2006"),
		name:   userDetails.Name,
	})
	if err != nil {
		logrus.Errorf("error while sending email: %s", err)
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

func SendSimpleMessage(domain, apiKey string, msg EmailVars) (string, error) {
	mg := mailgun.NewMailgun(domain, apiKey)
	m := mg.NewMessage(
		"Excited User mailgun@sandbox869c162826e04f39b909a85aeedcfc65.mailgun.org",
		"Receipts for your Smartvac Meter",
		"Testing some Mailgun awesomeness!",
		msg.email,
	)
	m.SetTemplate("smartvac-api")
	m.AddTemplateVariable("name", msg.name)
	m.AddTemplateVariable("amount", msg.amount)
	m.AddTemplateVariable("units", msg.units)
	m.AddTemplateVariable("date", msg.date)

	_, id, err := mg.Send(context.Background(), m)
	return id, err
}
