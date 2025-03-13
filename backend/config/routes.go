package config

import (
	"FreightMaster/backend/database"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
)

func SetupRoutes(r *gin.Engine, db *gorm.DB) {
	shipments := r.Group("/shipments")
	{
		shipments.GET("/", func(c *gin.Context) {
			var shipments []database.Shipment
			if err := db.Find(&shipments).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось загрузить список отправлений"})
				return
			}
			c.JSON(http.StatusOK, shipments)
		})

		// Получение одного отправления
		shipments.GET("/:id", func(c *gin.Context) {
			var shipment database.Shipment
			id := c.Param("id")
			if err := db.First(&shipment, id).Error; err != nil {
				c.JSON(http.StatusNotFound, gin.H{"error": "Отправление не найдено"})
				return
			}
			c.JSON(http.StatusOK, shipment)
		})

		// Создание нового отправления
		shipments.POST("/", func(c *gin.Context) {
			var shipment database.Shipment
			if err := c.ShouldBindJSON(&shipment); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат JSON"})
				return
			}
			if err := db.Create(&shipment).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка создания отправления"})
				return
			}
			c.JSON(http.StatusCreated, shipment)
		})

		// Обновление отправления
		shipments.PUT("/:id", func(c *gin.Context) {
			id := c.Param("id")
			var shipment database.Shipment
			if err := db.First(&shipment, id).Error; err != nil {
				c.JSON(http.StatusNotFound, gin.H{"error": "Отправление не найдено"})
				return
			}

			// Создаем новую структуру для обновления
			var updatedShipment database.Shipment
			if err := c.ShouldBindJSON(&updatedShipment); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат JSON"})
				return
			}

			// Обновляем только переданные поля
			db.Model(&shipment).Updates(updatedShipment)

			c.JSON(http.StatusOK, shipment)
		})

		shipments.DELETE("/:id", func(c *gin.Context) {
			id := c.Param("id")
			var shipment database.Shipment
			if err := db.First(&shipment, id).Error; err != nil {
				c.JSON(http.StatusNotFound, gin.H{"error": "Отправление не найдено"})
				return
			}

			if err := db.Delete(&shipment).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка удаления"})
				return
			}

			c.JSON(http.StatusOK, gin.H{"message": "Отправление удалено"})
		})

	}
}
