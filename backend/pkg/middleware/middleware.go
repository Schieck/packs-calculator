package middleware

import (
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	entity "github.com/Schieck/packs-calculator/internal/domain/entity/auth"
	"github.com/Schieck/packs-calculator/internal/domain/errs"
)

func Logger(logger *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		c.Next()

		latency := time.Since(start)
		clientIP := c.ClientIP()
		method := c.Request.Method
		statusCode := c.Writer.Status()

		if raw != "" {
			path = path + "?" + raw
		}

		logger.Info("HTTP Request",
			slog.String("method", method),
			slog.String("path", path),
			slog.Int("status", statusCode),
			slog.String("ip", clientIP),
			slog.Duration("latency", latency),
			slog.String("user_agent", c.Request.UserAgent()),
		)
	}
}

func Recovery(logger *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				logger.Error("Panic recovered",
					slog.String("error", err.(string)),
					slog.String("path", c.Request.URL.Path),
					slog.String("method", c.Request.Method),
				)

				c.JSON(http.StatusInternalServerError, errs.ErrorResponse{
					Error: "Internal server error",
				})
				c.Abort()
			}
		}()
		c.Next()
	}
}

func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		c.Header("Access-Control-Expose-Headers", "Content-Length")
		c.Header("Access-Control-Allow-Credentials", "true")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

func JWT(validateTokenUseCase entity.ValidateTokenUseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, errs.ErrorResponse{
				Error: "Authorization header required",
			})
			c.Abort()
			return
		}

		tokenString := ""
		if strings.HasPrefix(authHeader, "Bearer ") {
			tokenString = authHeader[7:]
		} else {
			c.JSON(http.StatusUnauthorized, errs.ErrorResponse{
				Error: "Invalid authorization header format. Use 'Bearer <token>'",
			})
			c.Abort()
			return
		}

		claims, err := validateTokenUseCase.Execute(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, errs.ErrorResponse{
				Error: "Invalid or expired token",
			})
			c.Abort()
			return
		}

		c.Set("authenticated", true)
		c.Set("subject", claims.Subject)
		c.Next()
	}
}

func IsAuthenticated(c *gin.Context) bool {
	authenticated, exists := c.Get("authenticated")
	if !exists {
		return false
	}
	auth, ok := authenticated.(bool)
	return ok && auth
}

func GetSubject(c *gin.Context) (string, bool) {
	subject, exists := c.Get("subject")
	if !exists {
		return "", false
	}
	sub, ok := subject.(string)
	return sub, ok
}
