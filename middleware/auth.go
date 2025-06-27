package middleware

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/Zauro25/Capstone-PerpusKominfosan/config"
	"github.com/Zauro25/Capstone-PerpusKominfosan/models"
)

var jwtSecret = []byte("perpustakaan-dpk-secret-key-2024")

type Claims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	UserType string `json:"user_type"`
	jwt.RegisteredClaims
}

func GenerateToken(userID uint, username, userType string) (string, int64, error) {
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		UserID:   userID,
		Username: username,
		UserType: userType,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtSecret)
	
	return tokenString, expirationTime.Unix(), err
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token tidak ditemukan"})
			c.Abort()
			return
		}

		if strings.HasPrefix(tokenString, "Bearer ") {
			tokenString = tokenString[7:]
		}

		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token tidak valid"})
			c.Abort()
			return
		}

		// Set user info in context
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("user_type", claims.UserType)

		c.Next()
	}
}

func RequireRole(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userType := c.GetString("user_type")
		
		allowed := false
		for _, role := range allowedRoles {
			if userType == role {
				allowed = true
				break
			}
		}

		if !allowed {
			c.JSON(http.StatusForbidden, gin.H{"error": "Akses ditolak"})
			c.Abort()
			return
		}

		c.Next()
	}
}

func AuditLogMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Record request details
		userID, exists := c.Get("user_id")
		if !exists {
			c.Next()
			return
		}

		userType := c.GetString("user_type")
		
		// Create audit log
		auditLog := models.AuditLog{
			UserType:  userType,
			UserID:    userID.(uint),
			Action:    c.Request.Method + " " + c.Request.URL.Path,
			IPAddress: c.ClientIP(),
			UserAgent: c.GetHeader("User-Agent"),
			Timestamp: time.Now(),
		}

		config.DB.Create(&auditLog)
		c.Next()
	}
}
