package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
)

// Mock for IMiddlewareService
type MockMiddlewareService struct {
	mock.Mock
}

func (m *MockMiddlewareService) GetJWTSecret() []byte {
	args := m.Called()
	if secret, ok := args.Get(0).([]byte); ok {
		return secret
	}
	return nil
}

func (m *MockMiddlewareService) MiddlewareJWT() gin.HandlerFunc {
	args := m.Called()
	if handler, ok := args.Get(0).(gin.HandlerFunc); ok {
		return handler
	}
	return nil
}

func (m *MockMiddlewareService) ValidateUserId() gin.HandlerFunc {
	args := m.Called()
	if handler, ok := args.Get(0).(gin.HandlerFunc); ok {
		return handler
	}
	return nil
}

func (m *MockMiddlewareService) AuthorizeRoles(allowedRoles ...string) gin.HandlerFunc {
	args := m.Called(allowedRoles)
	if handler, ok := args.Get(0).(gin.HandlerFunc); ok {
		return handler
	}
	return nil
}

func (m *MockMiddlewareService) CreateToken(userId int, username, userType string) (string, error) {
	args := m.Called(userId, username, userType)
	return args.String(0), args.Error(1)
}

func NewMock() IMiddlewareService {
	return &MockMiddlewareService{}
}
