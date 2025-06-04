package middleware

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	entity "github.com/Schieck/packs-calculator/internal/domain/entity/auth"
	"github.com/Schieck/packs-calculator/internal/domain/errs"
)

type mockValidateTokenUseCase struct {
	jwtKey string
}

func (m *mockValidateTokenUseCase) Execute(token string) (*entity.JWTClaims, error) {
	jwtToken, err := jwt.ParseWithClaims(token, &entity.JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(m.jwtKey), nil
	})

	if err != nil {
		return nil, errs.ErrInvalidToken
	}

	if !jwtToken.Valid {
		return nil, errs.ErrInvalidToken
	}

	claims, ok := jwtToken.Claims.(*entity.JWTClaims)
	if !ok {
		return nil, errs.ErrInvalidToken
	}

	return claims, nil
}

type MiddlewareTestSuite struct {
	suite.Suite
	router               *gin.Engine
	logger               *slog.Logger
	jwtKey               string
	validateTokenUseCase entity.ValidateTokenUseCase
}

func (suite *MiddlewareTestSuite) SetupTest() {
	gin.SetMode(gin.TestMode)
	suite.router = gin.New()

	var buf bytes.Buffer
	suite.logger = slog.New(slog.NewJSONHandler(&buf, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))

	suite.jwtKey = "test-secret-key"
	suite.validateTokenUseCase = &mockValidateTokenUseCase{jwtKey: suite.jwtKey}
}

func (suite *MiddlewareTestSuite) TestCORSMiddleware() {
	suite.router.Use(CORS())
	suite.router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "test"})
	})

	req := httptest.NewRequest("OPTIONS", "/test", nil)
	req.Header.Set("Origin", "http://localhost:3000")
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusNoContent, w.Code)
	assert.Equal(suite.T(), "*", w.Header().Get("Access-Control-Allow-Origin"))
	assert.Contains(suite.T(), w.Header().Get("Access-Control-Allow-Methods"), "GET")
	assert.Contains(suite.T(), w.Header().Get("Access-Control-Allow-Headers"), "Authorization")

	req = httptest.NewRequest("GET", "/test", nil)
	w = httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)
	assert.Equal(suite.T(), "*", w.Header().Get("Access-Control-Allow-Origin"))
}

func (suite *MiddlewareTestSuite) TestLoggerMiddleware() {
	var logOutput bytes.Buffer
	testLogger := slog.New(slog.NewJSONHandler(&logOutput, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))

	suite.router.Use(Logger(testLogger))
	suite.router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "test"})
	})

	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("User-Agent", "test-agent")
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)

	logContent := logOutput.String()
	assert.Contains(suite.T(), logContent, "HTTP Request")
	assert.Contains(suite.T(), logContent, "GET")
	assert.Contains(suite.T(), logContent, "/test")
	assert.Contains(suite.T(), logContent, "test-agent")
}

func (suite *MiddlewareTestSuite) TestRecoveryMiddleware() {
	var logOutput bytes.Buffer
	testLogger := slog.New(slog.NewJSONHandler(&logOutput, &slog.HandlerOptions{
		Level: slog.LevelError,
	}))

	suite.router.Use(Recovery(testLogger))
	suite.router.GET("/panic", func(c *gin.Context) {
		panic("test panic")
	})

	req := httptest.NewRequest("GET", "/panic", nil)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusInternalServerError, w.Code)

	var errorResponse errs.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &errorResponse)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "Internal server error", errorResponse.Error)
}

func (suite *MiddlewareTestSuite) TestJWTMiddleware() {
	suite.router.Use(JWT(suite.validateTokenUseCase))
	suite.router.GET("/protected", func(c *gin.Context) {
		authenticated := IsAuthenticated(c)
		subject, _ := GetSubject(c)
		c.JSON(http.StatusOK, gin.H{
			"authenticated": authenticated,
			"subject":       subject,
		})
	})

	req := httptest.NewRequest("GET", "/protected", nil)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusUnauthorized, w.Code)
	var errorResponse errs.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &errorResponse)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "Authorization header required", errorResponse.Error)

	req = httptest.NewRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", "InvalidToken")
	w = httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusUnauthorized, w.Code)
	err = json.Unmarshal(w.Body.Bytes(), &errorResponse)
	assert.NoError(suite.T(), err)
	assert.Contains(suite.T(), errorResponse.Error, "Invalid authorization header format")

	req = httptest.NewRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", "Bearer invalid-token")
	w = httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusUnauthorized, w.Code)

	token := suite.createValidJWT()
	req = httptest.NewRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w = httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)
	var successResponse map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &successResponse)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), true, successResponse["authenticated"])
	assert.Equal(suite.T(), "test-subject", successResponse["subject"])
}

func (suite *MiddlewareTestSuite) TestJWTMiddlewareExpiredToken() {
	suite.router.Use(JWT(suite.validateTokenUseCase))
	suite.router.GET("/protected", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	token := suite.createExpiredJWT()
	req := httptest.NewRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusUnauthorized, w.Code)
	var errorResponse errs.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &errorResponse)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "Invalid or expired token", errorResponse.Error)
}

func (suite *MiddlewareTestSuite) TestIsAuthenticated() {
	suite.router.Use(JWT(suite.validateTokenUseCase))
	suite.router.GET("/test", func(c *gin.Context) {
		authenticated := IsAuthenticated(c)
		c.JSON(http.StatusOK, gin.H{"authenticated": authenticated})
	})

	token := suite.createValidJWT()
	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)
	var successResponse map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &successResponse)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), true, successResponse["authenticated"])
}

func (suite *MiddlewareTestSuite) TestGetSubject() {
	suite.router.Use(JWT(suite.validateTokenUseCase))
	suite.router.GET("/test", func(c *gin.Context) {
		subject, exists := GetSubject(c)
		if exists {
			c.JSON(http.StatusOK, gin.H{"subject": subject})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": "no subject"})
		}
	})

	token := suite.createValidJWT()
	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)
	var successResponse map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &successResponse)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "test-subject", successResponse["subject"])
}

func (suite *MiddlewareTestSuite) createValidJWT() string {
	claims := &entity.JWTClaims{
		Subject: "test-subject",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(suite.jwtKey))
	assert.NoError(suite.T(), err)
	return tokenString
}

func (suite *MiddlewareTestSuite) createExpiredJWT() string {
	claims := &entity.JWTClaims{
		Subject: "test-subject",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(-time.Hour)), // Expired 1 hour ago
			IssuedAt:  jwt.NewNumericDate(time.Now().Add(-2 * time.Hour)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(suite.jwtKey))
	assert.NoError(suite.T(), err)
	return tokenString
}

func (suite *MiddlewareTestSuite) TestJWTMiddlewareInvalidSigningMethod() {
	suite.router.Use(JWT(suite.validateTokenUseCase))
	suite.router.GET("/protected", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	// Create an invalid token that will fail validation
	invalidToken := "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOnRydWV9.EkN-DOsnsuRjRO6BxXemmJDm3HbxrbRzXglbN2S4sOkopdU4IsDxTI8jO19W_A4K8ZPJijNLis4EZsHeY559a4DFOd50_OqgHs3PpPIXJI5wfLtY9TUCzl4TgQ1F2z0x9c4AUEw_jyD3FzS_sNJZF"

	req := httptest.NewRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+invalidToken)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusUnauthorized, w.Code)
	var errorResponse errs.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &errorResponse)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "Invalid or expired token", errorResponse.Error)
}

func (suite *MiddlewareTestSuite) TestMiddlewareChaining() {
	var logOutput bytes.Buffer
	testLogger := slog.New(slog.NewJSONHandler(&logOutput, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))

	suite.router.Use(Logger(testLogger))
	suite.router.Use(Recovery(testLogger))
	suite.router.Use(CORS())
	suite.router.Use(JWT(suite.validateTokenUseCase))

	suite.router.GET("/test", func(c *gin.Context) {
		authenticated := IsAuthenticated(c)
		c.JSON(http.StatusOK, gin.H{"authenticated": authenticated})
	})

	token := suite.createValidJWT()
	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Origin", "http://localhost:3000")
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)
	assert.Equal(suite.T(), "*", w.Header().Get("Access-Control-Allow-Origin"))

	var successResponse map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &successResponse)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), true, successResponse["authenticated"])

	logContent := logOutput.String()
	assert.Contains(suite.T(), logContent, "HTTP Request")
}

func TestMiddlewareSuite(t *testing.T) {
	suite.Run(t, new(MiddlewareTestSuite))
}
