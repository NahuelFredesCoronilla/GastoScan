package main

import (
	"log"
	"net/http"

	"gastoscan-api/config"
	"gastoscan-api/internal/database"
	"gastoscan-api/internal/expense"
	"gastoscan-api/internal/router"
)

func main() {
	cfg := config.Load()

	db, err := database.Connect(cfg)
	if err != nil {
		log.Fatal("error connecting to database:", err)
	}

	log.Println("database connected")

	if err := db.AutoMigrate(&expense.Expense{}); err != nil {
		log.Fatal("error migrating database:", err)
	}

	r := router.New(db)

	log.Printf("server running on port %s", cfg.Port)

	err = http.ListenAndServe(":"+cfg.Port, r)
	if err != nil {
		log.Fatal(err)
	}
}
