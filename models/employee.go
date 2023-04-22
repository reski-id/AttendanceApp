package models

import (
	"time"
)

type Employee struct {
	Model
	Username    string `json:"username" form:"username" binding:"required"`
	Fullname    string `json:"fullname" form:"fullname" binding:"required"`
	Password    string `json:"password" form:"password" binding:"required"`
	Email       string `json:"email" form:"email" binding:"required" gorm:"unique"`
	Role        string `json:"role" form:"role"`
	PhoneNumber string `json:"phoneNumber" form:"phoneNumber"`
	Address     string `json:"address" form:"address"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type LoginData struct {
	ID       uint   `gorm:"primaryKey"`
	Username string `json:"username" form:"username" validate:"required"`
	Password string `json:"password" form:"password" validate:"required"`
}

type Model struct {
	ID        uint `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type TokenResponse struct {
	Token    string `json:"token"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Role     string `json:"role"`
}

type CreateEmployeeResponse struct {
	Fullname    string `json:"fullname"`
	Username    string `json:"username"`
	Password    string `json:"password"`
	Email       string `json:"email"`
	Role        string `json:"role"`
	PhoneNumber string `json:"phoneNumber"`
	Address     string `json:"address"`
}
