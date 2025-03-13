package main

import (
	"FreightMaster/backend/config"
	"FreightMaster/backend/database"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	// Инициализация базы данных
	db := config.ConnectDatabase()
	if db == nil {
		log.Fatal("Не удалось подключиться к базе данных")
	}

	log.Println("✅ Подключение к базе данных успешно")

	// Выполняем миграции
	db.AutoMigrate(&database.User{}, &database.Shipment{}) // Теперь обе модели в database

	// Создаем маршрутизатор
	r := gin.Default()

	// Инициализация маршрутов
	config.SetupRoutes(r, db)

	// Запуск сервера
	r.Run(":8080")
}
