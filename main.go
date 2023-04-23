package main

import (
	"attendance/controllers"
	"fmt"
	"net/http"

	docs "attendance/docs"
	seed "attendance/seeder"

	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
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
	// err := godotenv.Load()
	// if err != nil {
	// 	fmt.Println("Error loading .env file")
	// 	os.Exit(1)
	// }

	//migrate and seeder
	seed.CreateMigration()
	seed.SeedUsers()

	router := echo.New()
	// Serve Swagger UI
	router.GET("/swagger/*", echoSwagger.WrapHandler)

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

	// new endpoint to check if service is running
	router.GET("/", func(c echo.Context) error {
		return c.HTML(http.StatusOK, fmt.Sprintf(`Attendance management system is running! <br/><a href="http://localhost:8080/swagger/index.html">View Swagger UI</a>`))
	})

	router.Logger.Fatal(router.Start(":8080"))
}
