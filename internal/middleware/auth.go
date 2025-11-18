package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/MarkAndrewKamau/Digital-Micro-Health-Assistant-Referral/internal/services"
	"github.com/MarkAndrewKamau/Digital-Micro-Health-Assistant-Referral/pkg/response"
)

func AuthMiddleware(authService *services.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Error(c, http.StatusUnauthorized, "MISSING_TOKEN", "Authorization header required")
			c.Abort()
			return
		}

		// Extract token from "Bearer <token>" format
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			response.Error(c, http.StatusUnauthorized, "INVALID_TOKEN_FORMAT", "Authorization header must be in format: Bearer <token>")
			c.Abort()
			return
		}

		sessionToken := parts[1]

		// Validate session
		user, err := authService.ValidateSession(c.Request.Context(), sessionToken)
		if err != nil {
			response.Error(c, http.StatusUnauthorized, "INVALID_SESSION", err.Error())
			c.Abort()
			return
		}

		// Store user in context
		c.Set("user", user)
		c.Set("user_id", user.ID)
		c.Set("user_role", user.Role)

		c.Next()
	}
}

// RoleMiddleware checks if user has required role
func RoleMiddleware(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, exists := c.Get("user_role")
		if !exists {
			response.Error(c, http.StatusUnauthorized, "UNAUTHORIZED", "User not authenticated")
			c.Abort()
			return
		}

		role := string(userRole.(string))
		for _, allowedRole := range allowedRoles {
			if role == allowedRole {
				c.Next()
				return
			}
		}

		response.Error(c, http.StatusForbidden, "FORBIDDEN", "Insufficient permissions")
		c.Abort()
	}
}