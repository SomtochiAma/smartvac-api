package main

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"

	"github.com/SomtochiAma/smartvac-api/models"
	"github.com/SomtochiAma/smartvac-api/routes"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Warnf("unable to load .env file: %s", err)
	}

	db, err := models.Init()
	if err != nil {
		log.Fatalf("unable to initialize database: %s", err)
	}

	postgresDB, err := db.DB()
	if err != nil {
		fmt.Println(err)
	}
	defer postgresDB.Close()
	log.Println("Successfully connected to the database.")

	r := routes.Init()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	if err := r.Run(":" + port); err != nil {
		log.Fatal(err)
	}
}
