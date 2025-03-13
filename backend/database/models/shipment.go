package models

import (
	"FreightMaster/backend/database" // Импортируем пакет database
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// CreateShipment добавляет новое отправление в базу
func CreateShipment(shipment *database.Shipment) error {
	return database.DB.Create(shipment).Error
}

// GetAllShipments получает все отправления
func GetAllShipments(shipments *[]database.Shipment) error {
	return database.DB.Find(shipments).Error
}

// GetShipmentByID получает конкретное отправление по ID
func GetShipmentByID(id uint, shipment *database.Shipment) error {
	return database.DB.First(shipment, id).Error
}

// UpdateShipment обновляет отправление
func UpdateShipment(shipment *database.Shipment) error {
	return database.DB.Save(shipment).Error
}

// DeleteShipment удаляет отправление по ID
func DeleteShipment(id uint) error {
	return database.DB.Delete(&database.Shipment{}, id).Error
}

// Обработчик создания отправления
func CreateShipmentHandler(c *gin.Context) {
	var shipment database.Shipment
	if err := c.ShouldBindJSON(&shipment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверные данные"})
		return
	}

	if err := CreateShipment(&shipment); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при создании отправления"})
		return
	}

	c.JSON(http.StatusCreated, shipment)
}

// Обработчик получения всех отправлений
func GetAllShipmentsHandler(c *gin.Context) {
	var shipments []database.Shipment
	if err := GetAllShipments(&shipments); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при получении отправлений"})
		return
	}

	c.JSON(http.StatusOK, shipments)
}

// Обработчик получения отправления по ID
func GetShipmentByIDHandler(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректный ID"})
		return
	}

	var shipment database.Shipment
	if err := GetShipmentByID(uint(id), &shipment); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Отправление не найдено"})
		return
	}

	c.JSON(http.StatusOK, shipment)
}

// Обработчик обновления отправления
func UpdateShipmentHandler(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректный ID"})
		return
	}

	var shipment database.Shipment
	if err := GetShipmentByID(uint(id), &shipment); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Отправление не найдено"})
		return
	}

	if err := c.ShouldBindJSON(&shipment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверные данные"})
		return
	}

	if err := UpdateShipment(&shipment); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при обновлении отправления"})
		return
	}

	c.JSON(http.StatusOK, shipment)
}

// Обработчик удаления отправления
func DeleteShipmentHandler(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректный ID"})
		return
	}

	if err := DeleteShipment(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при удалении отправления"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Отправление удалено"})
}
