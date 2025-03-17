// config/db.go
package config

import (
	"FreightMaster/backend/database"
)

var DB = database.DB

func ConnectDatabase() {
	database.ConnectDatabase()
}
