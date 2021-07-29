package controllers

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/SomtochiAma/smartvac-api/models"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

func PostReading(c *gin.Context) {
	var newReading models.Reading
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

func GetTotalReading(c *gin.Context) {
	//time := c.Query("time")
	var latestPayment models.Payment
	res := models.DB.Order("day desc").First(&latestPayment)
	if res.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": res.Error.Error(),
		})
		return
	}

	var sum int
	query := fmt.Sprintf("time > '%s'", latestPayment.Day.Format("2006-01-02T15:04:05-0700"))
	res = models.DB.Model(&models.Reading{}).Select("sum(total_power)").Where(query).Take(&sum)
	if res.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": res.Error.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"used":  sum,
		"total": latestPayment.Units,
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
	fmt.Println(len(values))
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
		var readings []models.Reading
		res := models.DB.Model(&models.Reading{}).
			Select("date_trunc('hour', time) as date, sum(total_power)").Group("date").
			Order("1").
			Find(&readings)
		if res.Error != nil {
			log.Printf("error writing message: %s", res.Error.Error())
			break
		}

		err = ws.WriteJSON(readings)
		if err != nil {
			log.Printf("error writing message: %s", err.Error())
			break
		}
		time.Sleep(5 * time.Minute)
	}
}
