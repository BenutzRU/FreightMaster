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
	// Маршруты авторизации (без middleware)
	r.POST("/register", backend.Register)
	r.POST("/login", backend.Login)

	// Группа маршрутов, доступных всем авторизованным пользователям
	auth := r.Group("/api")
	auth.Use(backend.AuthMiddleware())
	{
		// Отправления
		shipments := auth.Group("/shipments")
		{
			shipments.GET("", models.GetAllShipmentsHandler)
			shipments.GET("/:id", models.GetShipmentByIDHandler)
		}

		// Клиенты
		clients := auth.Group("/clients")
		{
			clients.GET("", func(c *gin.Context) {
				var clients []database.Client
				if err := db.Preload("Shipments").Find(&clients).Error; err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка получения клиентов"})
					return
				}
				c.JSON(http.StatusOK, clients)
			})
			clients.GET("/:id", func(c *gin.Context) {
				id, err := strconv.Atoi(c.Param("id"))
				if err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректный ID"})
					return
				}
				var client database.Client
				if err := db.Preload("Shipments").First(&client, id).Error; err != nil {
					if err == gorm.ErrRecordNotFound {
						c.JSON(http.StatusNotFound, gin.H{"error": "Клиент не найден"})
					} else {
						c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка базы данных"})
					}
					return
				}
				c.JSON(http.StatusOK, client)
			})
		}

		// Методы доставки
		deliveryMethods := auth.Group("/delivery-methods")
		{
			deliveryMethods.GET("", func(c *gin.Context) {
				var deliveryMethods []database.DeliveryMethod
				if err := db.Preload("Shipments").Find(&deliveryMethods).Error; err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка получения методов доставки"})
					return
				}
				c.JSON(http.StatusOK, deliveryMethods)
			})
			deliveryMethods.GET("/:id", func(c *gin.Context) {
				id, err := strconv.Atoi(c.Param("id"))
				if err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректный ID"})
					return
				}
				var deliveryMethod database.DeliveryMethod
				if err := db.Preload("Shipments").First(&deliveryMethod, id).Error; err != nil {
					if err == gorm.ErrRecordNotFound {
						c.JSON(http.StatusNotFound, gin.H{"error": "Метод доставки не найден"})
					} else {
						c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка базы данных"})
					}
					return
				}
				c.JSON(http.StatusOK, deliveryMethod)
			})
		}

		// Отделения
		branches := auth.Group("/branches")
		{
			branches.GET("", func(c *gin.Context) {
				var branches []database.Branch
				if err := db.Preload("Shipments").Find(&branches).Error; err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка получения отделений"})
					return
				}
				c.JSON(http.StatusOK, branches)
			})
			branches.GET("/:id", func(c *gin.Context) {
				id, err := strconv.Atoi(c.Param("id"))
				if err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректный ID"})
					return
				}
				var branch database.Branch
				if err := db.Preload("Shipments").First(&branch, id).Error; err != nil {
					if err == gorm.ErrRecordNotFound {
						c.JSON(http.StatusNotFound, gin.H{"error": "Отделение не найдено"})
					} else {
						c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка базы данных"})
					}
					return
				}
				c.JSON(http.StatusOK, branch)
			})
		}

		// Фильтрация отправлений
		auth.GET("/shipments/filter", func(c *gin.Context) {
			status := c.Query("status")
			var shipments []database.Shipment
			query := db.Preload("Client").Preload("DeliveryMethod")

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

	// Группа маршрутов, доступных только администраторам
	admin := r.Group("/admin")
	admin.Use(backend.AuthMiddleware(), backend.AdminMiddleware())
	{
		// Отправления
		shipments := admin.Group("/shipments")
		{
			shipments.POST("", models.CreateShipmentHandler)
			shipments.PUT("/:id", models.UpdateShipmentHandler)
			shipments.DELETE("/:id", models.DeleteShipmentHandler)
		}

		// Клиенты
		clients := admin.Group("/clients")
		{
			clients.POST("", func(c *gin.Context) {
				var client database.Client
				if err := c.ShouldBindJSON(&client); err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректный JSON"})
					return
				}
				if err := db.Create(&client).Error; err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка создания клиента"})
					return
				}
				c.JSON(http.StatusCreated, client)
			})
			clients.PUT("/:id", func(c *gin.Context) {
				id, err := strconv.Atoi(c.Param("id"))
				if err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректный ID"})
					return
				}
				var client database.Client
				if err := db.First(&client, id).Error; err != nil {
					if err == gorm.ErrRecordNotFound {
						c.JSON(http.StatusNotFound, gin.H{"error": "Клиент не найден"})
						return
					}
					c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка базы данных"})
					return
				}
				if err := c.ShouldBindJSON(&client); err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректный JSON"})
					return
				}
				if err := db.Save(&client).Error; err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка обновления"})
					return
				}
				c.JSON(http.StatusOK, client)
			})
			clients.DELETE("/:id", func(c *gin.Context) {
				id, err := strconv.Atoi(c.Param("id"))
				if err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректный ID"})
					return
				}
				if err := db.Delete(&database.Client{}, id).Error; err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка удаления"})
					return
				}
				c.JSON(http.StatusOK, gin.H{"message": "Клиент удалён"})
			})
		}

		// Методы доставки
		deliveryMethods := admin.Group("/delivery-methods")
		{
			deliveryMethods.POST("", func(c *gin.Context) {
				var deliveryMethod database.DeliveryMethod
				if err := c.ShouldBindJSON(&deliveryMethod); err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректный JSON"})
					return
				}
				if err := db.Create(&deliveryMethod).Error; err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка создания метода доставки"})
					return
				}
				c.JSON(http.StatusCreated, deliveryMethod)
			})
			deliveryMethods.PUT("/:id", func(c *gin.Context) {
				id, err := strconv.Atoi(c.Param("id"))
				if err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректный ID"})
					return
				}
				var deliveryMethod database.DeliveryMethod
				if err := db.First(&deliveryMethod, id).Error; err != nil {
					if err == gorm.ErrRecordNotFound {
						c.JSON(http.StatusNotFound, gin.H{"error": "Метод доставки не найден"})
						return
					}
					c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка базы данных"})
					return
				}
				if err := c.ShouldBindJSON(&deliveryMethod); err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректный JSON"})
					return
				}
				if err := db.Save(&deliveryMethod).Error; err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка обновления"})
					return
				}
				c.JSON(http.StatusOK, deliveryMethod)
			})
			deliveryMethods.DELETE("/:id", func(c *gin.Context) {
				id, err := strconv.Atoi(c.Param("id"))
				if err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректный ID"})
					return
				}
				if err := db.Delete(&database.DeliveryMethod{}, id).Error; err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка удаления"})
					return
				}
				c.JSON(http.StatusOK, gin.H{"message": "Метод доставки удалён"})
			})
		}

		// Отделения
		branches := admin.Group("/branches")
		{
			branches.POST("", func(c *gin.Context) {
				var branch database.Branch
				if err := c.ShouldBindJSON(&branch); err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректный JSON"})
					return
				}
				if err := db.Create(&branch).Error; err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка создания отделения"})
					return
				}
				c.JSON(http.StatusCreated, branch)
			})
			branches.PUT("/:id", func(c *gin.Context) {
				id, err := strconv.Atoi(c.Param("id"))
				if err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректный ID"})
					return
				}
				var branch database.Branch
				if err := db.First(&branch, id).Error; err != nil {
					if err == gorm.ErrRecordNotFound {
						c.JSON(http.StatusNotFound, gin.H{"error": "Отделение не найдено"})
						return
					}
					c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка базы данных"})
					return
				}
				if err := c.ShouldBindJSON(&branch); err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректный JSON"})
					return
				}
				if err := db.Save(&branch).Error; err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка обновления"})
					return
				}
				c.JSON(http.StatusOK, branch)
			})
			branches.DELETE("/:id", func(c *gin.Context) {
				id, err := strconv.Atoi(c.Param("id"))
				if err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректный ID"})
					return
				}
				if err := db.Delete(&database.Branch{}, id).Error; err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка удаления"})
					return
				}
				c.JSON(http.StatusOK, gin.H{"message": "Отделение удалено"})
			})
		}

		// Пользователи
		users := admin.Group("/users")
		{
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
				c.JSON(http.StatusOK, gin.H{"message": "Пользователь удалён"})
			})
		}
	}
}
