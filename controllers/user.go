package controllers

import (
	"fmt"
	"github.com/SomtochiAma/smartvac-api/models"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"strings"
)

func Signin(c *gin.Context) {
	type UserCred struct {
		Email string `json:"email" binding:"required"`
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
	result := models.DB.Where("email = ?", userCred.Email).First(&user)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": result.Error.Error(),
		})
		return
	}

	if user.Email == "" {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "user not found",
		})
		return
	}

	if !checkPasswordHash(userCred.Password, user.Password) {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "invalid username or password",
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

	var existingUser models.User
	_ = models.DB.Where("email = ?", newUser.Email).First(&existingUser)
	if existingUser.Email != "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "user already exists",
		})
		return
	};

	hashedPassword, err := HashPassword(newUser.Password);
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	newUser.Password = hashedPassword

 	result := models.DB.Create(&newUser)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": result.Error.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "signin successful",
		"id": newUser.ID,
	})
}

func GetUser(c *gin.Context) {
	var user models.User
	id, _ := c.Params.Get("id")
	result := models.DB.Select("id","name", "email", "min_unit", "address").First(&user, id)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": result.Error.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  user,
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

	c.JSON(http.StatusOK, gin.H{
		"message": "update successful",
	})
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return "", err
	}

	fmt.Println(string(bytes))
	return string(bytes), nil
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	fmt.Println(err)
	return err == nil
}
