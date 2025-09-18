package main

import (
	"time"

	"gorm.io/gorm"
)

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

// Coupon model
type Coupon struct {
	gorm.Model
	Code      string    `gorm:"uniqueIndex;not null" json:"code"`
	Type      string    `gorm:"size:20;not null" json:"type"` // "percent" or "fixed"
	Amount    float64   `json:"amount"`                     // percent (0-100) or fixed amount
	ExpiresAt *time.Time `json:"expires_at" gorm:"index"`
	MaxUses   int       `json:"max_uses"`
	UsedCount int       `json:"used_count"`
	Active    bool      `json:"active"`
}
