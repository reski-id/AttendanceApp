package models

import "time"

type ClockIn struct {
	ID          uint      `gorm:"primary_key" json:"id"`
	EmployeeID  int       `gorm:"not null" json:"employee_id"`
	ClockInTime time.Time `gorm:"not null" json:"clock_in_time"`
	CreatedAt   time.Time `gorm:"not null" json:"created_at"`
}

type ClockOut struct {
	ID           uint      `gorm:"primary_key" json:"id"`
	EmployeeID   int       `gorm:"not null" json:"employee_id"`
	ClockOutTime time.Time `gorm:"not null" json:"clock_out_time"`
	ClockInID    uint
	CreatedAt    time.Time `gorm:"not null" json:"created_at"`
}

type WorkingHours struct {
	ID          uint      `json:"id" gorm:"primary_key"`
	EmployeeID  int       `json:"employee_id" gorm:"not null"`
	HoursWorked string    `json:"hours_worked" gorm:"not null"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type ClockResponse struct {
	ID         uint      `json:"id"`
	EmployeeID int       `json:"employee_id"`
	ClockType  string    `json:"clock_type"`
	ClockTime  time.Time `json:"clock_time"`
	Hours      int       `json:"hours"`
	Minutes    int       `json:"minutes"`
}
