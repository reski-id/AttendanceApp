package utils

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func Connect() (*gorm.DB, error) {
	// err := godotenv.Load()
	// if err != nil {
	// 	fmt.Println("Error loading .env file")
	// 	os.Exit(1)
	// }

	dbUser := "root"
	dbPassword := "reski123"
	dbHost := "localhost"
	dbPort := "3306"
	dbName := "attendancedb"

	dsn := dbUser + ":" + dbPassword + "@tcp(" + dbHost + ":" + dbPort + ")/" + dbName + "?charset=utf8mb4&parseTime=True&loc=Local"

	return gorm.Open(mysql.Open(dsn), &gorm.Config{})
}
