package config

import (
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() *gorm.DB {
	dsn := "host=localhost user=postgres password=student dbname=freight port=5432 sslmode=disable TimeZone=UTC"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Ошибка подключения к базе данных: %v", err)
	}

	fmt.Println("Подключено к базе данных!")

	// Можно добавить миграции (раскомментировать при необходимости)
	// db.AutoMigrate(&models.Shipment{})

	return db
}
