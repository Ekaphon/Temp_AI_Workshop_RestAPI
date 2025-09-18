package main

import (
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
)

// CreateCouponRequest defines the request body for creating a coupon
type CreateCouponRequest struct {
	Code      string     `json:"code"`
	Type      string     `json:"type"`
	Amount    float64    `json:"amount"`
	ExpiresAt *time.Time `json:"expires_at"`
	MaxUses   int        `json:"max_uses"`
}

// CreateCoupon creates a coupon
func CreateCoupon(c *fiber.Ctx) error {
	var req CreateCouponRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.ErrBadRequest
	}
	if req.Code == "" || (req.Type != "percent" && req.Type != "fixed") {
		return fiber.NewError(fiber.StatusBadRequest, "invalid coupon")
	}
	coupon := Coupon{
		Code: req.Code,
		Type: req.Type,
		Amount: req.Amount,
		ExpiresAt: req.ExpiresAt,
		MaxUses: req.MaxUses,
		UsedCount: 0,
		Active: true,
	}
	if err := DB.Create(&coupon).Error; err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "coupon already exists or invalid")
	}
	return c.Status(http.StatusCreated).JSON(coupon)
}

// ListCoupons returns all coupons
func ListCoupons(c *fiber.Ctx) error {
	var cs []Coupon
	DB.Find(&cs)
	return c.JSON(cs)
}

// GetCoupon returns a coupon by code
func GetCoupon(c *fiber.Ctx) error {
	code := c.Params("code")
	var cp Coupon
	if err := DB.First(&cp, "code = ?", code).Error; err != nil {
		return fiber.ErrNotFound
	}
	return c.JSON(cp)
}

// RedeemCoupon redeems a coupon by code
func RedeemCoupon(c *fiber.Ctx) error {
	code := c.Params("code")
	var cp Coupon
	if err := DB.First(&cp, "code = ?", code).Error; err != nil {
		return fiber.ErrNotFound
	}
	if !cp.Active {
		return fiber.NewError(fiber.StatusBadRequest, "coupon inactive")
	}
	if cp.ExpiresAt != nil {
		if time.Now().After(*cp.ExpiresAt) {
			return fiber.NewError(fiber.StatusBadRequest, "coupon expired")
		}
	}
	if cp.MaxUses > 0 && cp.UsedCount >= cp.MaxUses {
		return fiber.NewError(fiber.StatusBadRequest, "coupon fully used")
	}
	cp.UsedCount++
	if err := DB.Save(&cp).Error; err != nil {
		return fiber.ErrInternalServerError
	}
	return c.JSON(cp)
}
