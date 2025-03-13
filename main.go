package main

import (
	"FreightMaster/backend/config"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	// Инициализация базы данных
	db := config.ConnectDatabase() // ✅ Теперь правильно вызывается из config

	if db == nil {
		log.Fatal("Не удалось подключиться к базе данных")
	}

	fmt.Println("✅ Подключение к базе данных успешно")

	// Создаем маршрутизатор Gin
	r := gin.Default()

	// Инициализация маршрутов
	config.SetupRoutes(r, db) // ✅ Передаем Gin и базу данных

	// Запуск сервера
	r.Run(":8080")
}
