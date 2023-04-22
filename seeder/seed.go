package seeder

import (
	"log"

	"attendance/models"
	"attendance/utils"
)

func SeedUsers() {
	db, err := utils.Connect()
	if err != nil {
		log.Fatalf("failed to connect database: %s", err.Error())
	}

	// check if any user already exists in the database
	var user models.Employee
	if db.First(&user).Error == nil {
		log.Println("users already seeded")
		return
	}

	// migrate the user table
	db.AutoMigrate(&models.Employee{})

	// create some users
	users := []models.Employee{
		{Username: "john_doe", Password: "password1", Role: "user", Email: "jhon@gmail.com", Fullname: "Jhon Doe", Address: "", PhoneNumber: ""},
		{Username: "jane_doe", Password: "password2", Role: "admin", Email: "adm@gmail.com", Fullname: "Jane Doe", Address: "", PhoneNumber: ""},
		{Username: "bob_smith", Password: "password3", Role: "user", Email: "bob@gmail.com", Fullname: "Bob Smith", Address: "", PhoneNumber: ""},
	}

	for i := range users {
		users[i].ID = uint(i) + 1
		err = db.Create(&users[i]).Error
		if err != nil {
			log.Fatalf("failed to seed users: %s", err.Error())
		}
	}

	log.Println("users seeded")
}
