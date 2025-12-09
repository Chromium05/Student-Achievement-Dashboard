package middleware

import (
	"student-report/utils"
	"strings"

	"github.com/gofiber/fiber/v2"
)

// Middleware untuk memerlukan login
func AuthRequired() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Ambil token dari header Authorization
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(401).JSON(fiber.Map{
				"error": "Token akses diperlukan",
			})
		}

		// Extract token dari "Bearer TOKEN"
		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			return c.Status(401).JSON(fiber.Map{
				"error": "Format token tidak valid",
			})
		}

		// Validasi token
		claims, err := utils.ValidateToken(tokenParts[1])
		if err != nil {
			return c.Status(401).JSON(fiber.Map{
				"error": "Token tidak valid atau expired",
			})
		}

		// Simpan informasi user di context
		c.Locals("user_id", claims.UserID)
		c.Locals("username", claims.Username)
		c.Locals("role", claims.Role)
		c.Locals("permissions", claims.Permissions) // Store permissions di context

		return c.Next()
	}
}

// Middleware untuk memerlukan role admin
func AdminOnly() fiber.Handler {
	return func(c *fiber.Ctx) error {
		role := c.Locals("role").(string)
		if role != "admin" {
			return c.Status(403).JSON(fiber.Map{
				"error": "Akses ditolak. Hanya admin yang boleh mengakses.",
			})
		}
		return c.Next()
	}
}

// Middleware untuk memerlukan permission tertentu
func RequirePermission(permission string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		role := c.Locals("role").(string)
		
		// Admin memiliki akses ke semua permissions
		if role == "admin" {
			return c.Next()
		}
		
		// Check apakah user memiliki permission yang dibutuhkan
		permissions := c.Locals("permissions").([]string)
		hasPermission := false
		
		for _, p := range permissions {
			if p == permission {
				hasPermission = true
				break
			}
		}
		
		if !hasPermission {
			return c.Status(403).JSON(fiber.Map{
				"error": "Akses ditolak. Anda tidak memiliki permission: " + permission,
			})
		}
		
		return c.Next()
	}
}

// Middleware untuk memerlukan salah satu dari beberapa permission
func RequireAnyPermission(permissions ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		role := c.Locals("role").(string)
		
		// Admin memiliki akses ke semua permissions
		if role == "admin" {
			return c.Next()
		}
		
		userPermissions := c.Locals("permissions").([]string)
		hasPermission := false
		
		for _, requiredPerm := range permissions {
			for _, userPerm := range userPermissions {
				if userPerm == requiredPerm {
					hasPermission = true
					break
				}
			}
			if hasPermission {
				break
			}
		}
		
		if !hasPermission {
			return c.Status(403).JSON(fiber.Map{
				"error": "Akses ditolak. Anda tidak memiliki permission yang dibutuhkan.",
			})
		}
		
		return c.Next()
	}
}

// Middleware untuk memerlukan semua permission tertentu
func RequireAllPermissions(permissions ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		role := c.Locals("role").(string)
		
		// Admin memiliki akses ke semua permissions
		if role == "admin" {
			return c.Next()
		}
		
		userPermissions := c.Locals("permissions").([]string)
		
		for _, requiredPerm := range permissions {
			hasPermission := false
			for _, userPerm := range userPermissions {
				if userPerm == requiredPerm {
					hasPermission = true
					break
				}
			}
			if !hasPermission {
				return c.Status(403).JSON(fiber.Map{
					"error": "Akses ditolak. Anda tidak memiliki semua permission yang dibutuhkan.",
				})
			}
		}
		
		return c.Next()
	}
}