package database

import (
	"gorm.io/gorm"
	"time"
)

type User struct {
	gorm.Model
	Username string `json:"username" gorm:"unique"`
	Password string `json:"password"`
	Role     string `json:"role" gorm:"default:'user'"` // Роль пользователя: user или admin
}

type Client struct {
	gorm.Model
	Name      string     `json:"name"`
	Contact   string     `json:"contact"`
	Shipments []Shipment `gorm:"foreignKey:ClientID"` // Один клиент может иметь много отправлений
}

type Shipment struct {
	gorm.Model
	UserID           uint       `json:"user_id"`
	ClientID         uint       `json:"client_id"`
	Description      string     `json:"description"`
	Status           string     `json:"status"`
	Cost             float64    `json:"cost"` // Меняем тип на float64 для согласованности
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`
	DeletedAt        *time.Time `json:"deleted_at"`         // Используем указатель, так как поле может быть NULL
	DeliveryMethodID uint       `json:"delivery_method_id"` // Связь с методом доставки
	BranchID         uint       `json:"branch_id"`
}

type DeliveryMethod struct {
	gorm.Model
	MethodName string     `json:"method_name"`                 // Название метода доставки (например, "Air", "Sea")
	City       string     `json:"city"`                        // Город доставки
	Rate       float64    `json:"rate"`                        // Расценка за доставку
	Shipments  []Shipment `gorm:"foreignKey:DeliveryMethodID"` // Один метод доставки может использоваться многими отправлениями
}

type Branch struct {
	gorm.Model
	Location  string     `json:"location"`            // Местоположение отделения
	Shipments []Shipment `gorm:"foreignKey:BranchID"` // Связь с отправлениями (опционально)
}
