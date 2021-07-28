package controllers

import (
	"fmt"
	"github.com/SomtochiAma/smartvac-api/models"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func PostReading(c *gin.Context) {
	var newReading models.CurrentReading
	if err := c.ShouldBindJSON(&newReading); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	result := models.DB.Create(&newReading)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": result.Error.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": newReading,
	})
}

func GetReading(c *gin.Context) {
	type Reading struct {
		Date time.Time `json:"date"`
		Sum uint `json:"sum"`
	}
	var values []Reading
	frequency := c.DefaultQuery("freq", "hour")
	id := c.Query("id")

	query := fmt.Sprintf("date_trunc('%s', time) as date, sum(value)", frequency)
	res := models.DB.Table("current_readings").
		Where("user_id = ?", id).
		Select(query).Group("date").
		Order("1").
		Find(&values)
	if res.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
	 		"error": "unable to retrieve data",
		})
	 	return
	}
	fmt.Println(len(values))
	fmt.Println(values)

	 c.JSON(http.StatusOK, gin.H{
		"data": values,
		"message": "values retrieved successfully",
	 })
}
