package v1

import (
	"github.com/gin-gonic/gin"
	"gophKeeper/internal/server/config"
	"gophKeeper/internal/server/middleware"
	"gophKeeper/internal/server/services/auth/models"
	"gophKeeper/internal/server/services/auth/usecase"
	"net/http"
)

type AuthHandler struct {
	config  *config.Config
	service usecase.IAuthUsecase
	mware   middleware.IMiddlewareService
}

func NewAuthHandler(config *config.Config, engine *gin.RouterGroup, service usecase.IAuthUsecase, mware middleware.IMiddlewareService) {
	handler := AuthHandler{
		config:  config,
		service: service,
		mware:   mware,
	}

	router := engine.Group("/auth")
	{
		router.POST("/login", handler.login)
	}
}

func (h *AuthHandler) login(c *gin.Context) {
	var user models.AuthUser
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	userInfo, err := h.service.CheckUser(c, &user)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	token, err := h.mware.CreateToken(userInfo.UserId, userInfo.Username, userInfo.UserType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}
