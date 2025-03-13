package config

import (
	"FreightMaster/backend"
	"FreightMaster/backend/database"
	"FreightMaster/backend/database/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

func SetupRoutes(r *gin.Engine, db *gorm.DB) {
	// Основные маршруты shipments
	shipments := r.Group("/shipments")
	{
		shipments.GET("", models.GetAllShipmentsHandler)
		shipments.POST("", models.CreateShipmentHandler)
		shipments.GET("/:id", models.GetShipmentByIDHandler)
		shipments.PUT("/:id", models.UpdateShipmentHandler)
		shipments.DELETE("/:id", models.DeleteShipmentHandler)
	}

	// Маршруты пользователей
	users := r.Group("/users")
	{
		users.GET("", func(c *gin.Context) {
			var users []database.User
			if err := db.Find(&users).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка получения пользователей"})
				return
			}
			c.JSON(http.StatusOK, users)
		})

		users.GET("/:id", func(c *gin.Context) {
			var user database.User
			id, err := strconv.Atoi(c.Param("id"))
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректный ID"})
				return
			}
			if err := db.First(&user, id).Error; err != nil {
				if err == gorm.ErrRecordNotFound {
					c.JSON(http.StatusNotFound, gin.H{"error": "Пользователь не найден"})
				} else {
					c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка базы данных"})
				}
				return
			}
			c.JSON(http.StatusOK, user)
		})

		users.POST("", func(c *gin.Context) {
			var user database.User
			if err := c.ShouldBindJSON(&user); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректный JSON"})
				return
			}
			if err := db.Create(&user).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка создания пользователя"})
				return
			}
			c.JSON(http.StatusCreated, user)
		})

		users.DELETE("/:id", func(c *gin.Context) {
			id, err := strconv.Atoi(c.Param("id"))
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректный ID"})
				return
			}
			if err := db.Delete(&database.User{}, id).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка удаления"})
				return
			}
			c.JSON(http.StatusOK, gin.H{"message": "Пользователь удален"})
		})
	}

	// Маршруты авторизации
	r.POST("/register", backend.Register)
	r.POST("/login", backend.Login)

	// Фильтрация отправлений
	r.GET("/shipments/filter", func(c *gin.Context) {
		status := c.Query("status")

		var shipments []database.Shipment
		query := db

		if status != "" {
			query = query.Where("status = ?", status)
		}

		if err := query.Find(&shipments).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка загрузки данных"})
			return
		}

		c.JSON(http.StatusOK, shipments)
	})
}
