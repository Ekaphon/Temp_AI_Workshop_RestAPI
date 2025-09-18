package main

import (
	"fmt"
	"net/mail"
	"os"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var jwtSecret = []byte("replace-with-secure-secret")

func init() {
	if s := os.Getenv("JWT_SECRET"); s != "" {
		jwtSecret = []byte(s)
	}
}

func hashPassword(password string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(b), err
}

func checkPassword(hash, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

func generateJWT(userID uint) (string, error) {
	claims := jwt.MapClaims{
		"sub": userID,
		"exp": time.Now().Add(24 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

// RequireAuth is a middleware that validates JWT and sets the user in locals
func RequireAuth(c *fiber.Ctx) error {
	auth := c.Get("Authorization")
	if auth == "" {
		return fiber.ErrUnauthorized
	}
	parts := strings.SplitN(auth, " ", 2)
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		return fiber.ErrUnauthorized
	}
	tokenString := parts[1]
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return jwtSecret, nil
	})
	if err != nil || !token.Valid {
		return fiber.ErrUnauthorized
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return fiber.ErrUnauthorized
	}
	sub, ok := claims["sub"]
	if !ok {
		return fiber.ErrUnauthorized
	}
	// numbers come back as float64
	var uid uint
	switch v := sub.(type) {
	case float64:
		uid = uint(v)
	case int:
		uid = uint(v)
	case int64:
		uid = uint(v)
	default:
		return fiber.ErrUnauthorized
	}
	var user User
	if err := DB.First(&user, uid).Error; err != nil {
		return fiber.ErrUnauthorized
	}
	c.Locals("user", &user)
	return c.Next()
}

func Register(c *fiber.Ctx) error {
	var body struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := c.BodyParser(&body); err != nil {
		return fiber.ErrBadRequest
	}
	body.Email = strings.TrimSpace(body.Email)
	if body.Email == "" || body.Password == "" {
		return fiber.ErrBadRequest
	}
	// validate email
	if _, err := mail.ParseAddress(body.Email); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid email")
	}
	if len(body.Password) < 6 {
		return fiber.NewError(fiber.StatusBadRequest, "password must be at least 6 characters")
	}
	// hash
	h, err := hashPassword(body.Password)
	if err != nil {
		return fiber.ErrInternalServerError
	}
	user := User{Email: body.Email, Password: h}
	if err := DB.Create(&user).Error; err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "email already exists")
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"email": user.Email, "id": user.ID})
}

func Login(c *fiber.Ctx) error {
	var body struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := c.BodyParser(&body); err != nil {
		return fiber.ErrBadRequest
	}
	if body.Email == "" || body.Password == "" {
		return fiber.ErrBadRequest
	}
	var user User
	if err := DB.First(&user, "email = ?", body.Email).Error; err != nil {
		return fiber.ErrUnauthorized
	}
	if err := checkPassword(user.Password, body.Password); err != nil {
		return fiber.ErrUnauthorized
	}
	tok, err := generateJWT(user.ID)
	if err != nil {
		return fiber.ErrInternalServerError
	}
	return c.JSON(fiber.Map{"token": tok})
}

func UpdateProfile(c *fiber.Ctx) error {
	var body struct {
		FirstName   string `json:"first_name"`
		LastName    string `json:"last_name"`
		Phone       string `json:"phone"`
		MemberLevel string `json:"member_level"`
		Points      *int64 `json:"points"`
	}
	if err := c.BodyParser(&body); err != nil {
		return fiber.ErrBadRequest
	}
	// simple validation
	if body.FirstName == "" && body.LastName == "" && body.Phone == "" && body.MemberLevel == "" && body.Points == nil {
		return fiber.NewError(fiber.StatusBadRequest, "no fields to update")
	}
	u := c.Locals("user").(*User)
	if body.FirstName != "" {
		u.FirstName = body.FirstName
	}
	if body.LastName != "" {
		u.LastName = body.LastName
	}
	if body.Phone != "" {
		u.Phone = body.Phone
	}
	if body.MemberLevel != "" {
		u.MemberLevel = body.MemberLevel
	}
	if body.Points != nil {
		u.Points = *body.Points
	}
	if err := DB.Save(u).Error; err != nil {
		return fiber.ErrInternalServerError
	}
	return c.JSON(fiber.Map{"email": u.Email, "id": u.ID, "first_name": u.FirstName, "last_name": u.LastName, "phone": u.Phone, "member_level": u.MemberLevel, "points": u.Points})
}
