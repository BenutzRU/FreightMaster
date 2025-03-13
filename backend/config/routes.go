package config

import (
	"FreightMaster/backend"
	"FreightMaster/backend/database"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

func SetupRoutes(r *gin.Engine, db *gorm.DB) {
	for _, route := range r.Routes() {
		println("Маршрут:", route.Method, route.Path)
	}

	shipments := r.Group("/shipments")
	{
		shipments.GET("/", func(c *gin.Context) {
			var shipments []database.Shipment
			if err := db.Find(&shipments).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при получении отправлений"})
				return
			}
			c.JSON(http.StatusOK, shipments)
		})

		shipments.GET("/:id", func(c *gin.Context) {
			var shipment database.Shipment
			id, err := strconv.Atoi(c.Param("id"))
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректный ID"})
				return
			}
			if err := db.First(&shipment, id).Error; err != nil {
				if err == gorm.ErrRecordNotFound {
					c.JSON(http.StatusNotFound, gin.H{"error": "Отправление не найдено"})
				} else {
					c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка базы данных"})
				}
				return
			}
			c.JSON(http.StatusOK, shipment)
		})

		shipments.POST("/", func(c *gin.Context) {
			var shipment database.Shipment
			if err := c.ShouldBindJSON(&shipment); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректный JSON"})
				return
			}
			if err := db.Create(&shipment).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка создания отправления"})
				return
			}
			c.JSON(http.StatusCreated, shipment)
		})

		shipments.DELETE("/:id", func(c *gin.Context) {
			id, err := strconv.Atoi(c.Param("id"))
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректный ID"})
				return
			}
			if err := db.Delete(&database.Shipment{}, id).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка удаления"})
				return
			}
			c.JSON(http.StatusOK, gin.H{"message": "Отправление удалено"})
		})
	}

	// Маршруты пользователей
	users := r.Group("/users")
	{
		users.GET("/", func(c *gin.Context) {
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

		users.POST("/", func(c *gin.Context) {
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

	// Вынесенные маршруты авторизации
	r.POST("/register", backend.Register)
	r.POST("/login", backend.Login)

	// Фильтрация отправлений
	r.GET("/shipments", func(c *gin.Context) {
		status := c.Query("status") // Фильтр по статусу

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
