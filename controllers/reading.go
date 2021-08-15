package controllers

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"

	"github.com/SomtochiAma/smartvac-api/models"
)

func PostReading(c *gin.Context) {
	var newReading models.Reading
	if err := c.ShouldBindJSON(&newReading); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	newReading.Time = time.Now()

	result := models.DB.Create(&newReading)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": result.Error.Error(),
		})
		return
	}
	if err := models.DB.Model(&models.User{}).Where("id = ?", newReading.UserID).
		UpdateColumn("used_unit",
			gorm.Expr("used_unit + ?", newReading.TotalPower)).
		Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": result.Error.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": newReading,
	})
}

func GetTotalReading(c *gin.Context) {
	var user struct {
		ID        uint
		UsedUnit  int
		TotalUnit int
	}
	id, _ := c.Params.Get("id")
	fmt.Println(id)
	err := models.DB.Model(&models.User{}).Where("id = ?", id).
		Select("id", "used_unit", "total_unit").
		First(&user).Error
	if err != nil {
		message := err.Error()
		logrus.Errorf("error getting payments: %s", err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			message = "no such user"
		}

		c.JSON(http.StatusBadRequest, gin.H{
			"message": message,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "usage summary retrieved successfully",
		"data":    user,
	})
}

func GetPostReading(c *gin.Context) {
	userId, err := strconv.Atoi(c.Query("user_id"))
	power, err := strconv.ParseFloat(c.Query("power"), 32)
	current, err := strconv.ParseFloat(c.Query("current"), 32)
	totalPower, err := strconv.ParseFloat(c.Query("total_power"), 32)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"error": fmt.Sprintf("unable to convert to int: %s", err),
		})
		return
	}

	newReading := models.Reading{
		UserID:     userId,
		Current:    float32(current),
		Power:      float32(power),
		TotalPower: float32(totalPower),
		Time:       time.Now(),
	}
	fmt.Printf("%v", newReading)

	result := models.DB.Create(&newReading)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": result.Error.Error(),
		})
		return
	}
	if err := models.DB.Model(&models.User{}).Where("id = ?", newReading.UserID).
		UpdateColumn("used_unit",
			gorm.Expr("used_unit + ?", newReading.TotalPower)).
		Error; err != nil {
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
		Sum  uint      `json:"sum"`
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
	//fmt.Println(len(values))
	fmt.Println(values)

	c.JSON(http.StatusOK, gin.H{
		"data":    values,
		"message": "values retrieved successfully",
	})
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func WebSocket(c *gin.Context) {
	id := c.Query("id")
	fmt.Println(id)
	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("error getting socket connection: %s", err)
		return
	}
	defer ws.Close()

	for {
		type Readings []struct {
			Sum  float32       `json:"sum"`
			Date time.Time `json:"date"`
		}
		var readings Readings
		//res := models.DB.Find(&readings)
		res := models.DB.Model(&models.Reading{}).
			Select("date_trunc('hour', time) as date, sum(total_power)").Group("date").
			Order("1").
			Find(&readings)
		if res.Error != nil {
			log.Printf("error writing message: %s", res.Error.Error())
			break
		}

		 type Units struct{
			UsedUnit float32 `json:"used_unit"`
			TotalUnit float32 `json:"total_unit"`
		}
		var units Units
		//res := models.DB.Find(&readings)
		total := models.DB.Model(&models.User{}).
			Select("used_unit, total_unit").Where("id = ?", 1).
			Take(&units)
		if total.Error != nil {
			logrus.Warnf("error writing message: %s", total.Error)
		}

		data := struct{
			Units Units `json:"units"`
			Readings Readings `json:"readings"`
		}{
			Units: units,
			Readings: readings,
		}

		err = ws.WriteJSON(data)
		if err != nil {
			log.Printf("error writing message: %s", err.Error())
			break
		}
		time.Sleep(5 * time.Minute)
	}
}

func GetAllReading(c *gin.Context) {
	var readings []models.Reading
	res := models.DB.Find(&readings)
	if res.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "unable to retrieve data",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":    readings,
		"message": "readings retrieved successfully",
	})
}
