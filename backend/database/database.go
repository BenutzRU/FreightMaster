// database/database.go
package database

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

var DB *gorm.DB

func ConnectDatabase() {
	dsn := "host=localhost user=postgres password=student dbname=freight port=5432 sslmode=disable TimeZone=UTC"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Error connecting to database: ", err)
	}

	DB = db
	if DB == nil {
		log.Fatal("❌ DB is nil after initialization!")
	}

	fmt.Println("Database connected!")

	// Миграция всех моделей
	err = DB.AutoMigrate(&User{}, &Client{}, &Shipment{})
	if err != nil {
		log.Fatal("Failed to migrate database: ", err)
	}
	fmt.Println("Database migrated!")

	// Создаём пользователя-администратора, если он не существует
	var admin User
	if err := DB.Where("username = ?", "admin").First(&admin).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			hashedPassword, err := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
			if err != nil {
				log.Fatal("Failed to hash admin password: ", err)
			}
			admin = User{
				Username: "admin",
				Password: string(hashedPassword),
				Role:     "admin",
			}
			if err := DB.Create(&admin).Error; err != nil {
				log.Fatal("Failed to create admin user: ", err)
			}
			fmt.Println("Created admin user")
		} else {
			log.Fatal("Error checking admin user: ", err)
		}
	} else {
		fmt.Println("Admin user already exists")
	}
}
