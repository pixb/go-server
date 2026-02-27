package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestRateLimiter(t *testing.T) {
	// Create echo instance
	e := echo.New()

	// Create test request
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Create rate limiter with strict config
	config := RateLimiterConfig{
		Rate:  0.1, // 0.1 requests per second
		Burst: 0,   // 0 requests in a burst
		KeyFunc: func(c echo.Context) string {
			return "test-key" // Use a fixed key for testing
		},
	}
	rateLimiter := NewRateLimiter(config)

	// Create handler
	handler := rateLimiter(func(c echo.Context) error {
		return c.String(http.StatusOK, "ok")
	})

	// First request should pass
	err := handler(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "ok", rec.Body.String())

	// Second request should be rate limited
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)
	err = handler(c)
	assert.Error(t, err)
	assert.Equal(t, http.StatusTooManyRequests, rec.Code)
}

func TestCSRFTokenMiddleware(t *testing.T) {
	// Create echo instance
	e := echo.New()

	// Create test request
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Create CSRF token middleware
	csrfTokenMiddleware := CSRFTokenMiddleware(DefaultCSRFConfig())

	// Create handler
	handler := csrfTokenMiddleware(func(c echo.Context) error {
		return c.String(http.StatusOK, "ok")
	})

	// Test CSRF token generation
	err := handler(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "ok", rec.Body.String())

	// Check if CSRF token is set in cookie and header
	cookie := rec.Result().Cookies()
	assert.NotEmpty(t, cookie)
	header := rec.Result().Header.Get(CSRFTokenHeader)
	assert.NotEmpty(t, header)
}

func TestCSRF(t *testing.T) {
	// Create echo instance
	e := echo.New()

	// Generate CSRF token
	token, err := GenerateCSRFToken()
	assert.NoError(t, err)

	// Create test request with CSRF token
	req := httptest.NewRequest(http.MethodPost, "/", nil)
	req.AddCookie(&http.Cookie{
		Name:  CSRFTokenCookie,
		Value: token,
	})
	req.Header.Set(CSRFTokenHeader, token)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Create CSRF middleware
	csrfMiddleware := NewCSRF(DefaultCSRFConfig())

	// Create handler
	handler := csrfMiddleware(func(c echo.Context) error {
		return c.String(http.StatusOK, "ok")
	})

	// Test CSRF token validation
	err = handler(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "ok", rec.Body.String())
}

func TestCSRF_InvalidToken(t *testing.T) {
	// Create echo instance
	e := echo.New()

	// Generate valid CSRF token
	validToken, err := GenerateCSRFToken()
	assert.NoError(t, err)

	// Create test request with invalid CSRF token
	req := httptest.NewRequest(http.MethodPost, "/", nil)
	req.AddCookie(&http.Cookie{
		Name:  CSRFTokenCookie,
		Value: validToken,
	})

	// Set invalid CSRF token in header
	req.Header.Set(CSRFTokenHeader, "invalidtoken")

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Create CSRF middleware
	csrfMiddleware := NewCSRF(DefaultCSRFConfig())

	// Create handler
	handler := csrfMiddleware(func(c echo.Context) error {
		return c.String(http.StatusOK, "ok")
	})

	// Test CSRF token validation with invalid token
	err = handler(c)
	assert.Error(t, err)
	assert.Equal(t, http.StatusForbidden, rec.Code)
}

func TestUnifiedResponseMiddleware(t *testing.T) {
	// Create echo instance
	e := echo.New()

	// Test successful response
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Create unified response middleware
	unifiedResponseMiddleware := UnifiedResponseMiddleware()

	// Create handler that returns a success response
	handler := unifiedResponseMiddleware(func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"message": "success"})
	})

	// Test successful response
	err := handler(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	// Test error response
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)

	// Create handler that returns an error
	handler = unifiedResponseMiddleware(func(c echo.Context) error {
		return echo.NewHTTPError(http.StatusBadRequest, "bad request")
	})

	// Test error response
	err = handler(c)
	assert.Error(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}
