package models

import (
	"gorm.io/gorm"
)

// Shipment представляет модель отправления
type Shipment struct {
	gorm.Model
	OrderNumber string `gorm:"type:varchar(100);not null;unique" json:"order_number"`
	Destination string `gorm:"size:255" json:"destination"`
	Status      string `gorm:"size:255" json:"status"`
}

// CreateShipment добавляет новое отправление в базу
func CreateShipment(db *gorm.DB, shipment *Shipment) error {
	return db.Create(shipment).Error
}

// GetAllShipments получает все отправления
func GetAllShipments(db *gorm.DB, shipments *[]Shipment) error {
	return db.Find(shipments).Error
}

// GetShipmentByID получает конкретную запись по ID
func GetShipmentByID(db *gorm.DB, id uint, shipment *Shipment) error {
	return db.First(shipment, id).Error // ✅ Теперь корректно
}

// UpdateShipment обновляет информацию об отправлении
func UpdateShipment(db *gorm.DB, shipment *Shipment) error {
	return db.Save(shipment).Error
}

// DeleteShipment удаляет отправление по ID
func DeleteShipment(db *gorm.DB, id uint) error {
	return db.Delete(&Shipment{}, id).Error
}
