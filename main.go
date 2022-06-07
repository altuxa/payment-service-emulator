package main

import (
	"log"

	"github.com/altuxa/payment-service-emulator/internal/repository"
	"github.com/altuxa/payment-service-emulator/internal/service"
)

func main() {
	db, err := repository.NewSqliteDB()
	if err != nil {
		log.Fatalf("failed to initialize db %s", err)
	}
	repository := repository.NewRepository(db)
	service := service.NewService(repository)
}
