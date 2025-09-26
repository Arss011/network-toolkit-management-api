package auth

import (
	"errors"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte("secrect-key-rahasia")

type JWTClaim struct {
	UserID   int    `json:"user_id"`
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

type AuthConfig struct {
	SecretKey     string
	TokenDuration time.Duration
}

type AuthService struct {
	config AuthConfig
}

func NewAuthService(config AuthConfig) *AuthService {
	if config.SecretKey == "" {
		config.SecretKey = "secrect-key-rahasia"
	}
	if config.TokenDuration == 0 {
		config.TokenDuration = 24 * time.Hour
	}

	return &AuthService{
		config: config,
	}
}

func (s *AuthService) GenerateToken(userID int, username, role string) (string, error) {
	claims := &JWTClaim{
		UserID:   userID,
		Username: username,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.config.TokenDuration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "toolkit-management",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.config.SecretKey))
}

func (s *AuthService) ValidateToken(tokenString string) (*JWTClaim, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaim{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.config.SecretKey), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*JWTClaim); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

func GetCurrentUser(c *gin.Context) (*JWTClaim, error) {
	user, exists := c.Get("user")
	if !exists {
		return nil, errors.New("user not found in context")
	}

	claims, ok := user.(*JWTClaim)
	if !ok {
		return nil, errors.New("invalid user claims type")
	}

	return claims, nil
}

func (s *AuthService) RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.JSON(401, gin.H{"success": false, "error": "Authorization header required"})
			c.Abort()
			return
		}

		if len(tokenString) > 7 && tokenString[:7] == "Bearer " {
			tokenString = tokenString[7:]
		}

		claims, err := s.ValidateToken(tokenString)
		if err != nil {
			c.JSON(401, gin.H{"success": false, "error": "Invalid token"})
			c.Abort()
			return
		}

		c.Set("user", claims)
		c.Next()
	}
}

// Check role
func (s *AuthService) RequireRole(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		claims, err := GetCurrentUser(c)
		if err != nil {
			c.JSON(401, gin.H{"success": false, "error": "Unauthorized"})
			c.Abort()
			return
		}

		for _, role := range roles {
			if claims.Role == role {
				c.Next()
				return
			}
		}

		c.JSON(403, gin.H{"success": false, "error": "Insufficient permissions"})
		c.Abort()
	}
}

// Check admin
func (s *AuthService) RequireAdmin() gin.HandlerFunc {
	return s.RequireRole("admin")
}
