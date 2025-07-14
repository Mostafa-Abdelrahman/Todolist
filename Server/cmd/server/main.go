package main

import (

	"log"
	"todo-list/internal/config"
	"todo-list/internal/database"
	"todo-list/internal/api"

)

func main() {
	cfg := config.Load()
	
	db,err:= database.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	server := api.NewServer(cfg,db)
	log.Printf("Server is running on port %s", cfg.ServerPort)
	if err := server.Start(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}