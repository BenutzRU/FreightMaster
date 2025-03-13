package models

import (
	"FreightMaster/backend/database" // Импортируем пакет database
	"gorm.io/gorm"
)

// CreateShipment добавляет новое отправление в базу
func CreateShipment(db *gorm.DB, shipment *database.Shipment) error {
	return db.Create(shipment).Error
}

// GetAllShipments получает все отправления
func GetAllShipments(db *gorm.DB, shipments *[]database.Shipment) error {
	return db.Find(shipments).Error
}

// GetShipmentByID получает конкретное отправление по ID
func GetShipmentByID(db *gorm.DB, id uint, shipment *database.Shipment) error {
	return db.First(shipment, id).Error
}

// UpdateShipment обновляет отправление
func UpdateShipment(db *gorm.DB, shipment *database.Shipment) error {
	return db.Save(shipment).Error
}

// DeleteShipment удаляет отправление по ID
func DeleteShipment(db *gorm.DB, id uint) error {
	return db.Delete(&database.Shipment{}, id).Error
}
