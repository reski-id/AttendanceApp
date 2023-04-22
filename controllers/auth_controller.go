package controllers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"

	"attendance/models"
	"attendance/utils"
)

type AuthController struct{}

func (auth *AuthController) Login(c echo.Context) error {
	var loginData models.LoginData

	if err := c.Bind(&loginData); err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
	}

	db, err := utils.Connect()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Database connection error"})
	}

	var user models.Employee
	result := db.Where("username = ?", loginData.Username).First(&user)
	if result.Error != nil {
		return c.JSON(http.StatusUnauthorized, models.ErrorResponse{Error: "Invalid username or password"})
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginData.Password)); err != nil {
		return c.JSON(http.StatusUnauthorized, models.ErrorResponse{Error: "Invalid username or password"})
	}

	token, err := utils.GenerateToken(int(user.ID), user.Role)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Token generation error"})
	}

	response := models.TokenResponse{
		Username: user.Username,
		Email:    user.Email,
		Role:     user.Role,
		Token:    token,
	}

	return c.JSON(http.StatusOK, response)
}

// Register godoc
// @Summary Register to the system
// @Description Register to the system with username, password, email, and isAdmin flag
// @Tags Auth
// @Accept json
// @Produce json
// @Param registrationData body models.User true "Registration Data"
// @Success 200 {object} models.CreateEmployeeResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /register [post]
func (auth *AuthController) Register(c echo.Context) error {
	db, err := utils.Connect()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: err.Error()})
	}

	var customer models.Employee
	if err := c.Bind(&customer); err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
	}

	// Check if username and email already exist
	var existingUser models.Employee
	result := db.Where("username = ?", customer.Username).Or("email = ?", customer.Email).First(&existingUser)
	if result.Error == nil {
		if existingUser.Username == customer.Username {
			return c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Username already exists"})
		}
		if existingUser.Email == customer.Email {
			return c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Email already exists"})
		}
	}

	// Hash the password
	hash, err := bcrypt.GenerateFromPassword([]byte(customer.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Password hashing error"})
	}

	// Create a new user record
	newCustomer := models.Employee{
		Fullname:    customer.Fullname,
		Username:    customer.Username,
		Password:    string(hash),
		Email:       customer.Email,
		PhoneNumber: customer.PhoneNumber,
		Address:     customer.Address,
		Role:        "user",
	}

	if err := db.Create(&newCustomer).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: err.Error()})
	}

	response := models.CreateEmployeeResponse{
		Fullname:    newCustomer.Fullname,
		Username:    newCustomer.Username,
		Email:       newCustomer.Email,
		Password:    customer.Password,
		Role:        newCustomer.Role,
		PhoneNumber: newCustomer.PhoneNumber,
		Address:     newCustomer.Address,
	}

	return c.JSON(http.StatusOK, response)
}
