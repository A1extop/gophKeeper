package v1

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gophKeeper/internal/server/config"
	"gophKeeper/internal/server/domain"
	"gophKeeper/internal/server/middleware"
	"gophKeeper/internal/server/services/auth/models"
	"gophKeeper/internal/server/services/auth/usecase"

	"net/http"
	"net/http/httptest"
	"testing"
)

// Тест на авторизацию
func TestAuthHandler_Login(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockAuthUsecase := usecase.NewAuthUsecaseMock().(*usecase.AuthUsecaseMock)
	mockMiddlewareService := middleware.NewMock().(*middleware.MockMiddlewareService)
	mockAuthUsecase.On("CheckUser", mock.Anything, &models.AuthUser{Username: "test", Password: "password"}).
		Return(&models.InfoUser{
			UserId:   1,
			Username: "test",
			UserType: "admin",
		}, nil)

	handler := AuthHandler{
		config:  &config.Config{},
		service: mockAuthUsecase,
		mware:   mockMiddlewareService,
	}

	r := gin.Default()
	engine := r.Group("/v1")
	NewAuthHandler(handler.config, engine, mockAuthUsecase, mockMiddlewareService)

	t.Run("Success", func(t *testing.T) {
		mockAuthUsecase.On("CheckUser", mock.Anything, &models.AuthUser{Username: "test", Password: "password"}).
			Return(&models.InfoUser{
				UserId:   1,
				Username: "test",
				UserType: "admin",
			}, nil)

		mockMiddlewareService.On("CreateToken", 1, "test", "admin").Return("mock_token", nil)

		payload := &models.AuthUser{
			Username: "test",
			Password: "password",
		}
		body, _ := json.Marshal(payload)

		req, _ := http.NewRequest("POST", "/v1/auth/login", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "mock_token")

		mockAuthUsecase.AssertExpectations(t)
		mockMiddlewareService.AssertExpectations(t)
	})

	t.Run("InvalidInput", func(t *testing.T) {
		body := []byte(`{invalid_json}`)

		req, _ := http.NewRequest("POST", "/v1/auth/login", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), domain.ErrInvalidInput.Error())
	})

	t.Run("InvalidCredentials", func(t *testing.T) {
		mockAuthUsecase.On("CheckUser", mock.Anything, &models.AuthUser{Username: "test", Password: "wrongpassword"}).
			Return(nil, domain.ErrInvalidCredentials)

		payload := &models.AuthUser{
			Username: "test",
			Password: "wrongpassword",
		}
		body, _ := json.Marshal(payload)

		req, _ := http.NewRequest("POST", "/v1/auth/login", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.Contains(t, w.Body.String(), domain.ErrInvalidCredentials)
	})

	t.Run("InternalError", func(t *testing.T) {
		mockAuthUsecase.ExpectedCalls = nil
		mockMiddlewareService.ExpectedCalls = nil

		mockAuthUsecase.On("CheckUser", mock.Anything, mock.MatchedBy(func(user *models.AuthUser) bool {
			return user.Username == "test" && user.Password == "password"
		})).Return(&models.InfoUser{UserId: 1, Username: "test", UserType: "admin"}, nil)

		mockMiddlewareService.On("CreateToken", 1, "test", "admin").Return("", domain.ErrTokenCreation)

		payload := &models.AuthUser{Username: "test", Password: "password"}
		body, _ := json.Marshal(payload)

		req, _ := http.NewRequest("POST", "/v1/auth/login", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
	})
}
