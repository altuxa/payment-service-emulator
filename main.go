package main

import (
	"log"

	"github.com/altuxa/payment-service-emulator/internal/handlers"
	"github.com/altuxa/payment-service-emulator/internal/repository"
	"github.com/altuxa/payment-service-emulator/internal/service"
)

func main() {
	db, err := repository.NewSqliteDB()
	if err != nil {
		log.Fatalf("failed to initialize db %s", err)
	}
	repository.CreateTable(db)
	repository := repository.NewRepository(db)
	service := service.NewService(repository)
	handler := handlers.NewHandler(service)
	handler.Server()
}
