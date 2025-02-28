package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Email         string    `gorm:"unique;not null" json:"email"`
	Password      string    `json:"password,omitempty"`
	FirstName     string    `json:"first_name"`
	LastName      string    `json:"last_name"`
	LastLogin     time.Time `json:"last_login"`
	IsActive      bool      `gorm:"default:true" json:"is_active"`
	EmailVerified bool      `gorm:"default:false" json:"email_verified"`
	// add phone number field and country code field as optional fields
	PhoneNumber string `json:"phone_number,omitempty"`
	CountryCode string `json:"country_code,omitempty"`
}
