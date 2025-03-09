package v1

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gophKeeper/internal/server/config"
	"gophKeeper/internal/server/middleware"
	"gophKeeper/internal/server/services/lockbox/models"
	"gophKeeper/internal/server/services/lockbox/usecase"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCreateLockBox(t *testing.T) {
	mockService := usecase.NewLockBoxUsecaseMock()
	mockConfig := &config.Config{}

	handler := LockBoxHandler{
		config:         mockConfig,
		lockBoxService: mockService,
	}

	router := gin.Default()
	router.POST("/lock_boxes/create", handler.createLockBox)

	t.Run("should create lockbox successfully", func(t *testing.T) {

		mockService.On("CreateLock", mock.Anything, mock.Anything).Return(123, nil)

		lockBoxData := models.Data{Name: "testBox", UserID: 1}
		body, _ := json.Marshal(lockBoxData)
		req, _ := http.NewRequest(http.MethodPost, "/lock_boxes/create", bytes.NewReader(body))
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)

		mockService.AssertExpectations(t)
	})

}

func TestDeleteLockBox(t *testing.T) {
	mockService := usecase.NewLockBoxUsecaseMock()
	mockMiddlewareService := middleware.NewMock().(*middleware.MockMiddlewareService)
	mockConfig := &config.Config{}

	handler := LockBoxHandler{
		config:         mockConfig,
		lockBoxService: mockService,
		mware:          mockMiddlewareService,
	}

	router := gin.Default()

	router.Use(func(ctx *gin.Context) {
		ctx.Set("userId", 1)
		ctx.Next()
	})

	router.DELETE("/lock_boxes/:name", handler.deleteLockBox)

	t.Run("should delete lockbox successfully", func(t *testing.T) {
		mockService.On("DeleteLock", mock.Anything, "testBox", 1).Return(nil)

		req, _ := http.NewRequest(http.MethodDelete, "/lock_boxes/testBox", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNoContent, w.Code)
		mockService.AssertExpectations(t)
	})

}

func TestGetLockBox(t *testing.T) {
	mockService := usecase.NewLockBoxUsecaseMock()
	mockMiddlewareService := middleware.NewMock().(*middleware.MockMiddlewareService)
	mockConfig := &config.Config{}

	handler := LockBoxHandler{
		config:         mockConfig,
		lockBoxService: mockService,
		mware:          mockMiddlewareService,
	}

	router := gin.Default()
	router.GET("/lock_boxes/:name", handler.getLockBox)

	t.Run("should get lockbox successfully", func(t *testing.T) {
		mockService.On("GetLockByName", mock.Anything, "testBox", 1).Return(&models.Data{Name: "testBox", UserID: 1}, nil)

		req, _ := http.NewRequest(http.MethodGet, "/lock_boxes/testBox", nil)
		req = req.WithContext(context.WithValue(req.Context(), "userId", 1))
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		mockService.AssertExpectations(t)
	})

}
