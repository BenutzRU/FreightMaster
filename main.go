package main

import (
	"FreightMaster/backend/config"
	"FreightMaster/backend/database"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	// Подключение к базе данных
	database.ConnectDatabase()
	db := database.DB
	if db == nil {
		log.Fatal("Не удалось подключиться к базе данных")
	}

	log.Println("✅ Подключение к базе данных успешно")

	// Автоматическая миграция моделей
	db.AutoMigrate(&database.User{}, &database.Shipment{})

	// Создание маршрутизатора
	r := gin.Default()

	// Инициализация маршрутов
	config.SetupRoutes(r, db)

	// Запуск сервера
	r.Run(":8080")
}
