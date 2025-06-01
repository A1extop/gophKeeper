package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/time/rate"
	"gophKeeper/internal/server/config"
	"gophKeeper/internal/server/db"
	"log/slog"
	"net/http"
	"sync"
	"time"
)

var jwtSecret = []byte("your-secret-key") // todo вынести

type IMiddlewareService interface {
	GetJWTSecret() []byte
	MiddlewareJWT() gin.HandlerFunc
	ValidateUserId() gin.HandlerFunc
	AuthorizeRoles(allowedRoles ...string) gin.HandlerFunc
	CreateToken(userId int, username, userType string) (string, error)
}

var (
	limiters        = sync.Map{}
	limiterTime     = time.Minute * 5
	cleanupInterval = time.Minute * 10
)

type limiterEntry struct {
	limiter    *rate.Limiter
	lastAccess time.Time
}

func getLimiter(ip string) *rate.Limiter {
	if entry, ok := limiters.Load(ip); ok {
		limiterEntry := entry.(*limiterEntry)
		limiterEntry.lastAccess = time.Now()
		return limiterEntry.limiter
	}

	limiter := rate.NewLimiter(5, 5)
	limiters.Store(ip, &limiterEntry{
		limiter:    limiter,
		lastAccess: time.Now(),
	})
	return limiter
}

func RateLimiter() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		limiter := getLimiter(ip)

		if !limiter.Allow() {
			c.JSON(http.StatusTooManyRequests, gin.H{"error": "too many requests"})
			c.Abort()
			return
		}

		c.Next()
	}
}

type MiddlewareService struct {
	config   *config.Config
	database db.IDatabase
}

func NewMiddlewareService(config *config.Config, database db.IDatabase) IMiddlewareService {
	return &MiddlewareService{config: config, database: database}
}

func (ms *MiddlewareService) CreateToken(userId int, username string, userType string) (string, error) {
	claims := jwt.MapClaims{
		"user_id":   userId,
		"user_name": username,
		"user_type": userType,
		"exp":       time.Now().Add(time.Hour * 72).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(ms.GetJWTSecret())
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func (ms *MiddlewareService) GetJWTSecret() []byte { // todo для конфига
	return jwtSecret
}

func (ms *MiddlewareService) MiddlewareJWT() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tokenString := ctx.GetHeader("Authorization")
		if tokenString == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "No token provided"})
			return
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return ms.GetJWTSecret(), nil
		})

		if err != nil || !token.Valid {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			return
		}

		userId := int(claims["user_id"].(float64))
		ctx.Set("userId", userId)

		userType, ok := claims["user_type"].(string)
		if ok {
			ctx.Set("userType", userType)
		}

		ctx.Next()
	}
}

// ValidateUserId todo избежать дублирования кода
func (ms *MiddlewareService) ValidateUserId() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tokenString := ctx.GetHeader("Authorization")
		if tokenString == "" {
			slog.Info("Not authorized user, next()")
			ctx.Next()
		} else {
			token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
				return ms.GetJWTSecret(), nil
			})

			if err != nil || !token.Valid {
				ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
				return
			}

			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
				return
			}

			userId := int(claims["user_id"].(float64))
			ctx.Set("userId", userId)

			ctx.Next()
		}
	}
}

func (ms *MiddlewareService) AuthorizeRoles(allowedRoles ...string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tokenString := ctx.GetHeader("Authorization")
		if tokenString == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "No token provided"})
			return
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return ms.GetJWTSecret(), nil
		})

		if err != nil || !token.Valid {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			return
		}

		userId := int(claims["user_id"].(float64))
		username := claims["user_name"].(string)

		var infoUser struct{ UserType string }
		query := `SELECT user_type FROM "users" WHERE user_id = $1 and username = $2`
		err = ms.database.GetDB().QueryRow(ctx, query, userId, username).Scan(&infoUser.UserType)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Access denied"})
			return
		}

		for _, role := range allowedRoles {
			if role == infoUser.UserType {
				ctx.Next()
				return
			}
		}

		ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Access denied"})
	}
}
