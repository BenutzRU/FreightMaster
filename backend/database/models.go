package database

type User struct {
	ID       uint   `gorm:"primaryKey"`
	Email    string `gorm:"unique;not null"`
	Password string `gorm:"not null"`
	Role     string `gorm:"type:varchar(10);default:'user'"` // user / admin
}

type Shipment struct {
	ID          uint   `gorm:"primaryKey"`
	Tracking    string `gorm:"unique" json:"tracking"`
	Status      string `json:"status"`
	OrderNumber string `gorm:"type:varchar(100);not null;unique" json:"order_number"`
	Destination string `gorm:"size:255" json:"destination"`
	UserID      uint   `json:"user_id"`
}
