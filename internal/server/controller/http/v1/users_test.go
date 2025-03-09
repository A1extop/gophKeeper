package v1

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gophKeeper/internal/server/config"
	"gophKeeper/internal/server/middleware"
	"gophKeeper/internal/server/services/users/models"
	"gophKeeper/internal/server/services/users/usecase"
	"gophKeeper/util"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRegisterUser(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockUsecase := new(usecase.UserUsecaseMock)
	mockMiddleware := new(middleware.MockMiddlewareService)

	router := gin.Default()
	apiGroup := router.Group("/api")
	NewUserHandler(&config.Config{}, apiGroup, mockUsecase, mockMiddleware)

	t.Run("should register user successfully", func(t *testing.T) {
		user := models.User{
			Username: "testuser",
			Password: "securepassword",
			UserType: util.Attendee,
		}
		mockUsecase.On("CreateUser", mock.Anything, mock.AnythingOfType("*models.User")).Return(nil)

		body, _ := json.Marshal(user)
		req, _ := http.NewRequest(http.MethodPost, "/api/users/", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
		mockUsecase.AssertExpectations(t)
	})

	t.Run("should return error on invalid JSON", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodPost, "/api/users/", bytes.NewBuffer([]byte("{invalid json")))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

}
