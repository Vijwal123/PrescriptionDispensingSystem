package middleware

import (
	"prescription/utils"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
)
func VerifyToken(c *fiber.Ctx, allowedRoles ...string) error {
    authHeader := c.Get("Authorization")
    if authHeader == "" {
        return c.Status(404).JSON("token missing")
    }

    tokenStr := strings.TrimSpace(authHeader)
    if strings.HasPrefix(strings.ToLower(tokenStr), "bearer ") {
        tokenStr = strings.TrimSpace(tokenStr[7:])
    }

    if tokenStr == "" {
        return c.Status(fiber.StatusUnauthorized).JSON("invalid token format")
    }

    token, err := jwt.Parse(tokenStr, utils.ExtractSecertKey)
    if err != nil || !token.Valid {
        return c.Status(fiber.StatusUnauthorized).JSON("invalid token")
    }

    claims, ok := token.Claims.(jwt.MapClaims)
    if !ok {
        return c.Status(fiber.StatusUnauthorized).JSON("invalid token claims")
    }

    role, ok := claims["role"].(string)
    if !ok {
        return c.Status(fiber.StatusUnauthorized).JSON("invalid role")
    }

    allowed := false
    for _, r := range allowedRoles {
        if r == role {
            allowed = true
            break
        }
    }
    if !allowed {
        return c.Status(fiber.StatusUnauthorized).JSON("role not allowed")
    }

    c.Locals("user_id", claims["user_id"])
    c.Locals("role", role)
    return c.Next()
}

func AdminVerification(c *fiber.Ctx) error {
    return VerifyToken(c, "Admin")
}

func DoctorVerification(c *fiber.Ctx) error {
    return VerifyToken(c, "Doctor")
}

func PharmacistVerification(c *fiber.Ctx) error {
    return VerifyToken(c, "Pharmacist")
}
