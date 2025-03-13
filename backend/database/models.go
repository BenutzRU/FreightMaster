package database

type User struct {
	ID    uint   `gorm:"primaryKey"`
	Name  string `gorm:"size:100"`
	Email string `gorm:"unique"`
}

type Shipment struct {
	ID       uint   `gorm:"primaryKey"`
	Tracking string `gorm:"unique"`
	Status   string
}
