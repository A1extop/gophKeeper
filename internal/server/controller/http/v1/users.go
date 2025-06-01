package v1

import (
	"github.com/gin-gonic/gin"
	"gophKeeper/internal/server/config"
	"gophKeeper/internal/server/domain"
	"gophKeeper/internal/server/middleware"
	"gophKeeper/internal/server/services/users/models"
	"gophKeeper/internal/server/services/users/usecase"
	"net/http"
	"strconv"
)

type UserHandler struct {
	config      *config.Config
	userService usecase.IUserUsecase
	mware       middleware.IMiddlewareService
}

func NewUserHandler(config *config.Config, router *gin.RouterGroup, userService usecase.IUserUsecase, mware middleware.IMiddlewareService) {
	userHandler := UserHandler{
		config:      config,
		userService: userService,
		mware:       mware,
	}

	userRouter := router.Group("/users")
	{
		//userRouter.GET("", userHandler.GetUsers)

		userRouter.POST("/", middleware.RateLimiter(), userHandler.RegisterUser)
		//userRouter.PUT("/:id", userHandler.UpdateUser)
		//userRouter.DELETE("/:id", userHandler.DeleteUser)

	}
}

func (uh *UserHandler) RegisterUser(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := uh.userService.CreateUser(c, &user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusCreated)
}
func (uh *UserHandler) UpdateUser(c *gin.Context) {
	userId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": domain.ErrInvalidUserID.Error()})
		return
	}

	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user.UserId = userId

	if err := uh.userService.UpdateUser(&user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}

func (uh *UserHandler) DeleteUser(c *gin.Context) {
	userId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": domain.ErrInvalidUserID.Error()})
		return
	}

	if err := uh.userService.DeleteUser(userId); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
