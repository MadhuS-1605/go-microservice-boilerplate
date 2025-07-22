package middleware

import (
	"crypto/subtle"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"net/http"
	"os"
)

type Claims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

//func AuthMiddleware(jwtSecret string) gin.HandlerFunc {
//	return func(c *gin.Context) {
//		authHeader := c.GetHeader("Authorization")
//		if authHeader == "" {
//			response.Error(c, http.StatusUnauthorized, "Authorization header required", nil)
//			c.Abort()
//			return
//		}
//
//		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
//		if tokenString == authHeader {
//			response.Error(c, http.StatusUnauthorized, "Invalid authorization header format", nil)
//			c.Abort()
//			return
//		}
//
//		token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
//			return []byte(jwtSecret), nil
//		})
//
//		if err != nil || !token.Valid {
//			response.Error(c, http.StatusUnauthorized, "Invalid or expired token", nil)
//			c.Abort()
//			return
//		}
//
//		if claims, ok := token.Claims.(*Claims); ok {
//			c.Set("user_id", claims.UserID)
//			c.Set("email", claims.Email)
//		}
//
//		c.Next()
//	}
//}
//
//func OptionalAuth(jwtSecret string) gin.HandlerFunc {
//	return func(c *gin.Context) {
//		authHeader := c.GetHeader("Authorization")
//		if authHeader == "" {
//			c.Next()
//			return
//		}
//
//		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
//		if tokenString != authHeader {
//			token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
//				return []byte(jwtSecret), nil
//			})
//
//			if err == nil && token.Valid {
//				if claims, ok := token.Claims.(*Claims); ok {
//					c.Set("user_id", claims.UserID)
//					c.Set("email", claims.Email)
//				}
//			}
//		}
//
//		c.Next()
//	}
//}

// SwaggerAuth provides basic authentication for Swagger UI
func SwaggerAuth(username, password string) gin.HandlerFunc {
	return gin.BasicAuth(gin.Accounts{
		username: password,
	})
}

// SwaggerAuthAdvanced provides more flexible basic auth with custom validation
func SwaggerAuthAdvanced(validateFunc func(username, password string) bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		user, pass, hasAuth := c.Request.BasicAuth()

		if !hasAuth || !validateFunc(user, pass) {
			c.Header("WWW-Authenticate", "Basic realm=Restricted")
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		c.Next()
	}
}

// ValidateSwaggerCredentials provides custom validation function
func ValidateSwaggerCredentials(username, password string) bool {
	envUsername := os.Getenv("SWAGGER_USERNAME")
	envPassword := os.Getenv("SWAGGER_PASSWORD")

	// If environment variables are not set, disable authentication
	if envUsername == "" || envPassword == "" {
		return false
	}

	// Use constant time comparison to prevent timing attacks
	userMatch := subtle.ConstantTimeCompare([]byte(username), []byte(envUsername)) == 1
	passMatch := subtle.ConstantTimeCompare([]byte(password), []byte(envPassword)) == 1

	return userMatch && passMatch
}
