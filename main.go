package main

import (
	"fmt"
	"github.com/SomtochiAma/smartvac-api/routes"
	"log"

	"github.com/SomtochiAma/smartvac-api/models"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("unable to load .env file: %s", err)
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
	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
