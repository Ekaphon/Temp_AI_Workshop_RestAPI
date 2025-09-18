package main

import "gorm.io/gorm"

// User model
type User struct {
	gorm.Model
	Email       string `gorm:"uniqueIndex;not null" json:"email"`
	Password    string `gorm:"not null" json:"-"`
	FirstName   string `gorm:"size:100" json:"first_name"`
	LastName    string `gorm:"size:100" json:"last_name"`
	Phone       string `gorm:"size:32" json:"phone"`
	MemberLevel string `gorm:"size:50" json:"member_level"`
	Points      int64  `json:"points"`
}
