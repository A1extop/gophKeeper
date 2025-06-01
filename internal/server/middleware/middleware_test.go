// middleware_test.go
package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"unsafe"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"gophKeeper/internal/server/config"
)

// ---------------- Фейковая база данных для тестирования AuthorizeRoles ---------------- //

// fakeDB реализует интерфейс db.IDatabase, требуемый для MiddlewareService.
type fakeDB struct{}

func (f fakeDB) CheckMigrations(ctx context.Context, state string) error {
	return nil
}

// GetDB возвращает фейковый пул. Для этого мы создаём объект fakePool и небезопасно приводим его к типу *pgxpool.Pool.
func (f fakeDB) GetDB() *pgxpool.Pool {
	fp := &fakePool{}
	// NB! Используем unsafe для тестового преобразования. Это допустимо только в тестах!
	return (*pgxpool.Pool)(unsafe.Pointer(fp))
}

// fakePool — минимальная реализация пула, удовлетворяющая вызову QueryRow.
type fakePool struct{}

// QueryRow возвращает фейковую строку, которая всегда содержит user_type "admin".
func (p *fakePool) QueryRow(ctx context.Context, query string, args ...interface{}) pgx.Row {
	return fakeRow{userType: "admin"}
}

// fakeRow реализует интерфейс pgx.Row (метод Scan).
type fakeRow struct {
	userType string
}

func (r fakeRow) Scan(dest ...interface{}) error {
	if len(dest) > 0 {
		if ptr, ok := dest[0].(*string); ok {
			*ptr = r.userType
			return nil
		}
	}
	return nil
}

// TestCreateToken проверяет, что при создании токена в его claims присутствуют ожидаемые значения.
func TestCreateToken(t *testing.T) {
	db := fakeDB{}
	ms := NewMiddlewareService(&config.Config{}, db)
	tokenStr, err := ms.CreateToken(123, "testuser", "admin")
	if err != nil {
		t.Fatalf("CreateToken failed: %v", err)
	}
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return ms.GetJWTSecret(), nil
	})
	if err != nil || !token.Valid {
		t.Fatalf("Token is not valid: %v", err)
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		t.Fatal("Failed to get token claims")
	}
	if int(claims["user_id"].(float64)) != 123 {
		t.Errorf("Expected user_id 123, got %v", claims["user_id"])
	}
	if claims["user_name"].(string) != "testuser" {
		t.Errorf("Expected user_name 'testuser', got %v", claims["user_name"])
	}
	if claims["user_type"].(string) != "admin" {
		t.Errorf("Expected user_type 'admin', got %v", claims["user_type"])
	}
}

// TestMiddlewareJWT_ValidToken проверяет middleware MiddlewareJWT при наличии валидного токена.
func TestMiddlewareJWT_ValidToken(t *testing.T) {
	db := fakeDB{}
	ms := NewMiddlewareService(&config.Config{}, db)
	tokenStr, err := ms.CreateToken(456, "anotheruser", "user")
	if err != nil {
		t.Fatalf("CreateToken failed: %v", err)
	}
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(ms.MiddlewareJWT())
	r.GET("/test", func(c *gin.Context) {
		userId, exists := c.Get("userId")
		if !exists {
			c.String(http.StatusInternalServerError, "userId not set")
			return
		}
		c.String(http.StatusOK, "userId: %v", userId)
	})
	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", tokenStr)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("Expected 200 OK, got %d", w.Code)
	}
	if !strings.Contains(w.Body.String(), "456") {
		t.Errorf("Expected response to contain userId 456, got %s", w.Body.String())
	}
}

// TestMiddlewareJWT_MissingToken проверяет, что MiddlewareJWT возвращает 401, если токен отсутствует.
func TestMiddlewareJWT_MissingToken(t *testing.T) {
	db := fakeDB{}
	ms := NewMiddlewareService(&config.Config{}, db)
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(ms.MiddlewareJWT())
	r.GET("/test", func(c *gin.Context) {
		c.String(http.StatusOK, "OK")
	})
	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("Expected 401 Unauthorized for missing token, got %d", w.Code)
	}
}

// TestValidateUserId_WithToken проверяет ValidateUserId при наличии валидного токена.
func TestValidateUserId_WithToken(t *testing.T) {
	db := fakeDB{}
	ms := NewMiddlewareService(&config.Config{}, db)
	tokenStr, err := ms.CreateToken(789, "validuser", "user")
	if err != nil {
		t.Fatalf("CreateToken failed: %v", err)
	}
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(ms.ValidateUserId())
	r.GET("/validate", func(c *gin.Context) {
		userId, exists := c.Get("userId")
		if !exists {
			c.String(http.StatusInternalServerError, "userId not set")
			return
		}
		c.String(http.StatusOK, "userId: %v", userId)
	})
	req := httptest.NewRequest("GET", "/validate", nil)
	req.Header.Set("Authorization", tokenStr)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("Expected 200 OK, got %d", w.Code)
	}
	if !strings.Contains(w.Body.String(), "789") {
		t.Errorf("Expected response to contain userId 789, got %s", w.Body.String())
	}
}

// TestValidateUserId_NoToken проверяет ValidateUserId при отсутствии токена.
func TestValidateUserId_NoToken(t *testing.T) {
	db := fakeDB{}
	ms := NewMiddlewareService(&config.Config{}, db)
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(ms.ValidateUserId())
	r.GET("/validate", func(c *gin.Context) {
		c.String(http.StatusOK, "no token")
	})
	req := httptest.NewRequest("GET", "/validate", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("Expected 200 OK for no token, got %d", w.Code)
	}
	if !strings.Contains(w.Body.String(), "no token") {
		t.Errorf("Expected response 'no token', got %s", w.Body.String())
	}
}

// TestRateLimiter проверяет, что RateLimiter ограничивает количество запросов с одного IP.
func TestRateLimiter(t *testing.T) {

	limiters = sync.Map{}

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(RateLimiter())
	r.GET("/ratelimit", func(c *gin.Context) {
		c.String(http.StatusOK, "ok")
	})

	req := httptest.NewRequest("GET", "/ratelimit", nil)
	req.RemoteAddr = "192.168.1.100:1234"

	for i := 0; i < 5; i++ {
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		if w.Code != http.StatusOK {
			t.Fatalf("Expected 200 OK on request %d, got %d", i+1, w.Code)
		}
	}
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusTooManyRequests {
		t.Fatalf("Expected 429 Too Many Requests on 6th request, got %d", w.Code)
	}
}
