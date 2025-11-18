package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/MarkAndrewKamau/Digital-Micro-Health-Assistant-Referral/internal/services"
	"github.com/MarkAndrewKamau/Digital-Micro-Health-Assistant-Referral/pkg/response"
)

type AuthHandler struct {
	authService *services.AuthService
}

func NewAuthHandler(authService *services.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

type LoginRequest struct {
	Phone string `json:"phone" binding:"required"`
}

type LoginResponse struct {
	SessionToken string      `json:"session_token"`
	User         interface{} `json:"user"`
}

// Login handles POST /v1/auth/login
func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body")
		return
	}

	// Login with phone number
	user, err := h.authService.Login(c.Request.Context(), req.Phone)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "LOGIN_FAILED", "Failed to login")
		return
	}

	// Create session
	userAgent := c.Request.UserAgent()
	ipAddress := c.ClientIP()
	session, err := h.authService.CreateSession(c.Request.Context(), user.ID, userAgent, ipAddress)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "SESSION_CREATE_FAILED", "Failed to create session")
		return
	}

	response.Success(c, http.StatusOK, LoginResponse{
		SessionToken: session.SessionToken,
		User:         user,
	})
}

// Logout handles POST /v1/auth/logout
func (h *AuthHandler) Logout(c *gin.Context) {
	sessionToken := c.GetHeader("Authorization")
	if sessionToken == "" {
		response.Error(c, http.StatusBadRequest, "MISSING_TOKEN", "Authorization header required")
		return
	}

	// Remove "Bearer " prefix if present
	if len(sessionToken) > 7 && sessionToken[:7] == "Bearer " {
		sessionToken = sessionToken[7:]
	}

	if err := h.authService.DeleteSession(c.Request.Context(), sessionToken); err != nil {
		response.Error(c, http.StatusInternalServerError, "LOGOUT_FAILED", "Failed to logout")
		return
	}

	response.Success(c, http.StatusOK, gin.H{"message": "Logged out successfully"})
}

// Me handles GET /v1/auth/me
func (h *AuthHandler) Me(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "UNAUTHORIZED", "User not authenticated")
		return
	}

	response.Success(c, http.StatusOK, user)
}