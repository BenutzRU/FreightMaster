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

	// Проверяем, что подключение успешно
	if database.DB == nil {
		log.Fatal("Не удалось подключиться к базе данных")
	}

	// Миграция базы данных
	if err := database.DB.AutoMigrate(
		&database.User{},
		&database.Client{},
		&database.Branch{},
		&database.DeliveryMethod{},
		&database.Shipment{},
	); err != nil {
		log.Fatal("Ошибка миграции базы данных:", err)
	}

	log.Println("✅ Подключение к базе данных успешно. Миграция выполнена.")

	// Создание маршрутизатора
	r := gin.Default()

	// Инициализация маршрутов
	config.SetupRoutes(r, database.DB)

	// Запуск сервера
	if err := r.Run(":8080"); err != nil {
		log.Fatal("Ошибка запуска сервера:", err)
	}
}
