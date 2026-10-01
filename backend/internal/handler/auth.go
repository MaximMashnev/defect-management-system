package handler

import (
	"defect-management-system/backend/internal/model"
	"defect-management-system/backend/internal/repository"
	"defect-management-system/backend/pkg/hash"
	"defect-management-system/backend/pkg/jwt"
	"defect-management-system/backend/pkg/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	userRepo  *repository.UserRepository
	jwtSecret string
}

func NewAuthHandler(userRepo *repository.UserRepository, jwtSecret string) *AuthHandler {
	return &AuthHandler{userRepo: userRepo, jwtSecret: jwtSecret}
}

type LoginRequest struct {
	Login    string `json:"login" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Data  model.AuthUser `json:"data"`
	Token string         `json:"token"`
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid input. Both login and password are required.", err)
		return
	}

	user, err := h.userRepo.GetByLogin(req.Login)
	if err != nil || user == nil {
		response.Unauthorized(c, "Invalid credentials")
		return
	}

	if !hash.CheckPassword(req.Password, user.PasswordHash) {
		response.Unauthorized(c, "Invalid credentials")
		return
	}

	token, err := jwt.GenerateToken(user.ID, user.Login, h.jwtSecret)
	if err != nil {
		response.Internal(c, "Failed to generate token", err)
		return
	}

	c.JSON(http.StatusOK, LoginResponse{
		Data: model.AuthUser{
			ID:    user.ID,
			Login: user.Login,
		},
		Token: token,
	})
}

func (h *AuthHandler) AuthMe(c *gin.Context) {
	userIDVal, exists := c.Get("userID")
	if !exists {
		response.Unauthorized(c, "Unauthorized")
		return
	}

	userID := userIDVal.(int)
	user, err := h.userRepo.GetByID(userID)
	if err != nil || user == nil {
		response.NotFound(c, "User not found")
		return
	}

	c.JSON(http.StatusOK, model.AuthUser{
		ID:    user.ID,
		Login: user.Login,
	})
}
