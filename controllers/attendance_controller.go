package controllers

import (
	"attendance/models"
	"attendance/utils"
	"fmt"
	"log"
	"net/http"
	"net/smtp"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
)

type AttendanceController struct{}

// ClockIn
// @Summary Clocks in an employee
// @Description Clocks in an employee and returns the clock-in time
// @Tags Attendance
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer {token}"
// @Success 200 {object} models.ClockResponse
// @Failure 401 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /attendance/clock-in/{id} [post]
func (ac *AttendanceController) ClockIn(c echo.Context) error {
	employeeID, _, err := utils.ExtractData(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, models.ErrorResponse{Error: err.Error()})
	}

	db, err := utils.Connect()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: err.Error()})
	}

	clockIn := models.ClockIn{EmployeeID: employeeID, ClockInTime: time.Now()}
	if err := db.Create(&clockIn).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: err.Error()})
	}

	// go ac.sendClockOutReminder(clockIn)

	return c.JSON(http.StatusOK, models.ClockResponse{
		ID:         clockIn.ID,
		EmployeeID: clockIn.EmployeeID,
		ClockType:  "clock_in",
		ClockTime:  time.Now(),
	})
}

// ClockOut
// @Summary Clocks out an employee
// @Description Clocks out an employee and returns the clock-out time and hours worked
// @Tags Attendance
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer {token}"
// @Success 200 {object} models.ClockResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 401 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /attendance/clock-out/{id} [post]
func (ac *AttendanceController) ClockOut(c echo.Context) error {
	employeeID, _, err := utils.ExtractData(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, models.ErrorResponse{Error: err.Error()})
	}

	db, err := utils.Connect()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: err.Error()})
	}

	var lastClockIn models.ClockIn
	if err := db.Where("employee_id = ?", employeeID).Last(&lastClockIn).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: err.Error()})
	}

	var existingClockOut models.ClockOut
	if err := db.Where("employee_id = ? AND id = ?", employeeID, lastClockIn.ID).First(&lastClockIn).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: err.Error()})
	}
	if existingClockOut.ID != 0 {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "You have already clocked out for this shift"})
	}

	clockOut := models.ClockOut{EmployeeID: employeeID, ClockInID: lastClockIn.ID, ClockOutTime: time.Now()}
	if err := db.Create(&clockOut).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: err.Error()})
	}

	// go ac.sendClockInReminder(clockOut.EmployeeID, clockOut.CreatedAt.AddDate(0, 0, 1))

	hoursWorked := clockOut.CreatedAt.Sub(lastClockIn.CreatedAt)
	hours := int(hoursWorked.Hours())
	minutes := int(hoursWorked.Minutes()) % 60
	workingHours := models.WorkingHours{EmployeeID: employeeID, HoursWorked: fmt.Sprintf("%d hour(s) %d minute(s)", hours, minutes)}
	if err := db.Create(&workingHours).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: err.Error()})
	}

	return c.JSON(http.StatusOK, models.ClockResponse{
		ID:         clockOut.ID,
		EmployeeID: clockOut.EmployeeID,
		ClockType:  "clock_out",
		ClockTime:  clockOut.ClockOutTime,
		Hours:      hours,
		Minutes:    minutes,
	})
}

// GetWorkHours godoc
// @Summary Get total work hours for an employee
// @Description Get the total number of hours an employee has worked based on their clock-in and clock-out entries
// @Tags Attendance
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer {token}"
// @ID get-work-hours
// @Produce json
// @Success 200 {object} models.ClockResponse
// @Failure 401 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /attendance/work-hours/{id} [get]
func (ac *AttendanceController) GetWorkHours(c echo.Context) error {
	// Get employee ID from JWT token
	employeeID, _, err := utils.ExtractData(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, models.ErrorResponse{Error: err.Error()})
	}

	// Find all clock-in and clock-out entries for the employee
	db, err := utils.Connect()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: err.Error()})
	}
	var clockIns []models.ClockIn
	if err := db.Where("employee_id = ?", employeeID).Find(&clockIns).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: err.Error()})
	}
	var clockOuts []models.ClockOut
	if err := db.Where("employee_id = ?", employeeID).Find(&clockOuts).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: err.Error()})
	}

	// Calculate the total number of hours worked
	var totalHours float64
	for _, clockIn := range clockIns {
		for _, clockOut := range clockOuts {
			if clockOut.ClockInID == clockIn.ID {
				totalHours += clockOut.CreatedAt.Sub(clockIn.CreatedAt).Hours()
				break
			}
		}
	}

	response := models.ClockResponse{
		EmployeeID: employeeID,
		ClockType:  "work_hours",
		// HoursWorked: totalHours,
	}
	return c.JSON(http.StatusOK, response)
}

func (ac *AttendanceController) sendClockInReminder(employeeID int, clockInTime time.Time) {
	// Find employee email address from the database
	db, err := utils.Connect()
	if err != nil {
		log.Println("Error connecting to database:", err)
		return
	}

	var employee models.Employee
	result := db.First(&employee, employeeID)
	if result.Error != nil {
		log.Println("Error finding employee:", result.Error)
		return
	}

	// Construct email message
	to := []string{employee.Email}
	subject := "Clock-in reminder for tomorrow"
	body := fmt.Sprintf("Hi %s,\n\nThis is a reminder that your clock-in time is tomorrow at %s.\n\nBest regards,\nThe Attendance App", employee.Fullname, clockInTime.Format("15:04:05"))

	// Set up SMTP client
	auth := smtp.PlainAuth("", "sender@example.com", "password", "smtp.example.com")
	msg := []byte("To: " + strings.Join(to, ",") + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"\r\n" + body + "\r\n")

	// Send email
	err = smtp.SendMail("smtp.example.com:587", auth, "sender@example.com", to, msg)
	if err != nil {
		log.Println("Error sending email:", err)
	}
}

func (ac *AttendanceController) sendClockOutReminder(clockIn models.ClockIn) {
	// Get employee email address
	db, err := utils.Connect()
	if err != nil {
		log.Println("Error connecting to database:", err.Error())
		return
	}

	var employee models.Employee
	result := db.Where("id = ?", clockIn.EmployeeID).First(&employee)
	if result.Error != nil {
		log.Println("Error retrieving employee:", result.Error.Error())
		return
	}

	// Construct email message
	to := employee.Email
	subject := "Reminder: Clock out time"
	body := fmt.Sprintf("Hello %s,\n\nThis is a reminder that your clock-out time is tomorrow at %s.\n\nBest regards,\nThe Attendance System", employee.Fullname, clockIn.CreatedAt.Add(time.Hour*8).Format("15:04:05"))

	// Send email using SMTP
	err = utils.SendEmail(to, subject, body)
	if err != nil {
		log.Println("Error sending email:", err.Error())
	}
}
