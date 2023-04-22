package main

import (
	"attendance/controllers"
	"attendance/utils"
	"fmt"
	"os"

	"github.com/joho/godotenv"

	docs "attendance/docs"
	seed "attendance/seeder"

	"github.com/labstack/echo/v4"
)

// @title           Swagger Costumer APP
// @version         2.0
// @description     This is a swagger documentation for Costumer APP.
// @BasePath        /api/v1
// @host            localhost:8080
// @schemes         http https
// @SecurityDefinition  jwt
// @Security        jwt
func main() {
	//setting env
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
		os.Exit(1)
	}

	//test email
	to := "reski.devel@gmail.com"
	subject := "Test Email dari main,go"
	body := "This is a test email sent from Go Attenddance in main.go!"

	errem := utils.SendEmail(to, subject, body)
	if errem != nil {
		fmt.Println("Error sending email:", errem)
		return
	}

	fmt.Println("Email sent successfully!")

	//migrate and seeder
	seed.CreateMigration()
	seed.SeedUsers()

	router := echo.New()

	docs.SwaggerInfo.BasePath = "/api/v1"

	employeesController := &controllers.EmployeeController{}
	authController := &controllers.AuthController{}
	attendanceController := &controllers.AttendanceController{}

	v1 := router.Group("/api/v1")

	v1.POST("/login", authController.Login)
	v1.POST("/register", authController.Register)

	v1.POST("/employees", employeesController.CreateEmployee)
	v1.PUT("/employees/:id", employeesController.UpdateEmployee)
	v1.DELETE("/employees/:id", employeesController.DeleteEmployee)
	v1.GET("/employees", employeesController.GetEmployees)
	v1.GET("/employees/:id", employeesController.GetEmployee)
	v1.GET("/employees/search", employeesController.SearchEmployees)

	// attendance endpoints
	v1.POST("/attendance/clock-in/:id", attendanceController.ClockIn)
	v1.POST("/attendance/clock-out/:id", attendanceController.ClockOut)
	v1.GET("/attendance/work-hours/:id", attendanceController.GetWorkHours)

	router.Logger.Fatal(router.Start(":8080"))
}
