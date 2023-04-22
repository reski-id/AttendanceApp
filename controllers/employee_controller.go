package controllers

import (
	"attendance/models"
	"attendance/utils"
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

type EmployeeController struct{}

func (controller EmployeeController) GetEmployees(c echo.Context) error {
	_, role, err := utils.ExtractData(c)

	if role != "admin" {
		return c.JSON(http.StatusUnauthorized, models.ErrorResponse{Error: "Only admin can Access"})
	}

	db, err := utils.Connect()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Database connection error"})
	}

	var employees []models.Employee
	page, _ := strconv.Atoi(c.QueryParam("page"))
	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	offset := (page - 1) * limit

	result := db.Offset(offset).Limit(limit).Find(&employees)
	if result.Error != nil {
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: result.Error.Error()})
	}

	return c.JSON(http.StatusOK, employees)
}

// @Summary Get a employee
// @Description Get a single employee by ID
// @Tags Employees
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer {token}"
// @Param id path int true "Employee ID"
// @Success 200 {object} models.Employee
// @Failure 401 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /employees/{id} [get]
func (controller EmployeeController) GetEmployee(c echo.Context) error {
	_, role, err := utils.ExtractData(c)

	if role != "admin" {
		return c.JSON(http.StatusUnauthorized, models.ErrorResponse{Error: "Only admin can Access"})
	}
	fmt.Println(role)
	db, err := utils.Connect()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: err.Error()})
	}
	var employee models.Employee
	result := db.First(&employee, c.Param("id"))
	if result.Error != nil {
		return c.JSON(http.StatusNotFound, models.ErrorResponse{Error: "Employee not found"})
	}

	return c.JSON(http.StatusOK, employee)
}

// @Summary Create a employee
// @Description Create a new employee
// @Tags Employees
// @Accept json
// @Produce json
// @Param employee body models.Employee true "Employee object"
// @Success 200 {object} models.CreateEmployeeResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /employees [post]
func (controller EmployeeController) CreateEmployee(c echo.Context) error {
	// all user can access
	db, err := utils.Connect()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: err.Error()})
	}

	var employee models.Employee
	err = c.Bind(&employee)
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
	}

	//cek
	var existingUser models.Employee
	checkUsername := db.Where("username = ?", employee.Username).First(&existingUser)
	if checkUsername.Error == nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Username already exists"})
	}
	checkEmail := db.Where("email = ?", employee.Email).First(&existingUser)
	if checkEmail.Error == nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Email already exists"})
	}

	employee.Role = "user"

	hash, err := bcrypt.GenerateFromPassword([]byte(employee.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Password hashing error"})
	}

	newEmployee := models.Employee{
		Fullname:    employee.Fullname,
		Username:    employee.Username,
		Password:    string(hash),
		Email:       employee.Email,
		PhoneNumber: employee.PhoneNumber,
		Address:     employee.Address,
		Role:        employee.Role}

	result := db.Create(&newEmployee)
	if result.Error != nil {
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: result.Error.Error()})
	}

	response := models.CreateEmployeeResponse{
		Fullname:    newEmployee.Fullname,
		Username:    newEmployee.Username,
		Password:    employee.Password,
		Email:       newEmployee.Email,
		Role:        newEmployee.Role,
		PhoneNumber: newEmployee.PhoneNumber,
		Address:     newEmployee.Address,
	}

	return c.JSON(http.StatusOK, response)
}

// UpdateEmployee godoc
// @Summary Update a employee by ID
// @Description Update a employee by ID
// @Tags Employees
// @Param id path int true "Employee ID"
// @Accept json
// @Produce json
// @Param employee body models.Employee true "Employee data"
// @Success 200 {object} models.Employee
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /employees/{id} [put]
func (controller EmployeeController) UpdateEmployee(c echo.Context) error {
	// all user can access
	db, err := utils.Connect()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: err.Error()})
	}

	var employee models.Employee
	result := db.First(&employee, c.Param("id"))
	if result.Error != nil {
		return c.JSON(http.StatusNotFound, models.ErrorResponse{Error: "Employee not found"})
	}

	err = c.Bind(&employee)
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
	}

	result = db.Save(&employee)
	if result.Error != nil {
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: result.Error.Error()})
	}

	return c.JSON(http.StatusOK, employee)
}

// DeleteEmployee godoc
// @Summary Delete a employee by ID
// @Description Delete a employee by ID
// @Tags Employees
// @Security ApiKeyAuth
// @Param Authorization header string true "Bearer {token}"
// @Param id path int true "Employee ID"
// @Produce json
// @Success 200 {object} models.MessageResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /employees/{id} [delete]
func (controller EmployeeController) DeleteEmployee(c echo.Context) error {
	_, role, err := utils.ExtractData(c)

	if role != "admin" {
		return c.JSON(http.StatusUnauthorized, models.ErrorResponse{Error: "Only admin can Access"})
	}
	db, err := utils.Connect()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: err.Error()})
	}

	var employee models.Employee
	result := db.First(&employee, c.Param("id"))
	if result.Error != nil {
		return c.JSON(http.StatusNotFound, models.ErrorResponse{Error: "Employee not found"})
	}

	result = db.Delete(&employee)
	if result.Error != nil {
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: result.Error.Error()})
	}

	return c.JSON(http.StatusOK, models.MessageResponse{Message: "Employee Deleted Succesfully"})
}

// SearchEmployees godoc
// @Summary Search employees by name
// @Description Search employees by name
// @Tags Employees
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer {token}"
// @Param query query string true "Search query"
// @Success 200 {array} models.Employee
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /employees/search [get]
func (controller EmployeeController) SearchEmployees(c echo.Context) error {
	_, role, err := utils.ExtractData(c)

	if role != "admin" {
		return c.JSON(http.StatusUnauthorized, models.ErrorResponse{Error: "Only admin can access"})
	}
	db, err := utils.Connect()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: err.Error()})
	}

	var employees []models.Employee
	query := "%" + c.QueryParam("query") + "%"

	result := db.Where("name LIKE ?", query).Find(&employees)
	if result.Error != nil {
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: result.Error.Error()})
	}

	return c.JSON(http.StatusOK, employees)
}
