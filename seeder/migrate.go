package seeder

import (
	"attendance/models"
	"attendance/utils"
)

func CreateMigration() {
	db, err := utils.Connect()
	if err != nil {
		panic("Failed to connect to database!")
	}

	// Auto migrate all entities
	db.AutoMigrate(&models.Employee{}, &models.ClockIn{}, &models.ClockOut{}, &models.WorkingHours{})
}
