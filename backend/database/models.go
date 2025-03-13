package database

type User struct {
	ID    uint   `gorm:"primaryKey"`
	Name  string `gorm:"size:100"`
	Email string `gorm:"unique"`
}

type Shipment struct {
	ID          uint   `gorm:"primaryKey"`
	Tracking    string `gorm:"unique"`
	Status      string
	OrderNumber string `gorm:"type:varchar(100);not null;unique" json:"order_number"`
	Destination string `gorm:"size:255" json:"destination"`
}
