package models

import (
	"errors"
	"fmt"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Init() (*gorm.DB, error) {
	requiredVar := []string{"DB_USER", "DB_PASS", "DB_HOST", "DB_NAME"}
	for _, val := range requiredVar {
		if _, exist := os.LookupEnv(val); !exist {
			return nil, errors.New(fmt.Sprintf("env variable %s not set", val))
		}
	}

	user, pass, host, dbname := os.Getenv("DB_USER"), os.Getenv("DB_PASS"), os.Getenv("DB_HOST"), os.Getenv("DB_NAME")
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=5432 sslmode=disable", host, user, pass, dbname)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	DB = db
	DB.AutoMigrate(&Reading{})
	DB.AutoMigrate(&User{})
	DB.AutoMigrate(&Payment{})

	payments := Payment{
		UserID: 1,
		Amount: 1000,
		Units:  250,
		Day:    time.Now(),
	}
	db.Create(&payments)

	return db, nil
}
