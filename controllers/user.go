package controllers

import (
	"errors"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"net/http"
	"strings"

	"github.com/SomtochiAma/smartvac-api/models"
)

type ReturnUser struct {
	ID        uint   `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	MinUnit   int    `json:"min_unit"`
	UsedUnit   int    `json:"used_unit"`
	TotalUnit int    `json:"total_unit"`
}

func Signin(c *gin.Context) {
	type UserCred struct {
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	var userCred UserCred
	if err := c.ShouldBindJSON(&userCred); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	var user models.User
	err := models.DB.Where("email = ?", userCred.Email).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) || user.Email == "" || !checkPasswordHash(userCred.Password, user.Password) {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "invalid username or password",
		})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
	})
}

func CreateUser(c *gin.Context) {
	var newUser models.User
	if err := c.ShouldBindJSON(&newUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	var existingUser ReturnUser
	_ = models.DB.Model(&models.User{}).Where("email = ?", newUser.Email).First(&existingUser)
	if existingUser.Email != "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "user already exists",
		})
		return
	}

	hashedPassword, err := HashPassword(newUser.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	newUser.Password = hashedPassword

	if err := models.DB.Create(&newUser).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "signup successful",
		"data": existingUser,
	})
}

func GetUser(c *gin.Context) {
	var user ReturnUser
	id, _ := c.Params.Get("id")
	result := models.DB.Model(&models.User{}).Select("id", "name", "email", "min_unit").First(&user, id)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": result.Error.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": user,
	})
}

func UpdateUser(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil && !strings.Contains(err.Error(), "Field validation") {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	id, _ := c.Params.Get("id")
	result := models.DB.Model(models.User{}).Where("id = ?", id).Updates(&user)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": result.Error.Error(),
		})
		return
	}

	var retUser ReturnUser
	err := models.DB.Model(&models.User{}).Where("id = ?", id).First(&retUser).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "user doesn't exist",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "update successful",
	})
}


func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		log.Errorf("error decrypting password: %s", err)
	}
	return err == nil
}
