package middleware

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

const (
	// CSRFTokenHeader is the header name for CSRF token
	CSRFTokenHeader = "X-CSRF-Token"
	// CSRFTokenCookie is the cookie name for CSRF token
	CSRFTokenCookie = "csrf_token"
	// CSRFTokenLength is the length of the CSRF token
	CSRFTokenLength = 32
)

// CSRFConfig defines the configuration for the CSRF middleware
type CSRFConfig struct {
	// TokenLookup is a function that returns the CSRF token from the request
	TokenLookup func(c echo.Context) string
	// CookieName is the name of the cookie that stores the CSRF token
	CookieName string
	// CookiePath is the path of the cookie that stores the CSRF token
	CookiePath string
	// CookieDomain is the domain of the cookie that stores the CSRF token
	CookieDomain string
	// CookieSecure is whether the cookie that stores the CSRF token is secure
	CookieSecure bool
	// CookieHTTPOnly is whether the cookie that stores the CSRF token is HTTP only
	CookieHTTPOnly bool
	// CookieSameSite is the SameSite attribute of the cookie that stores the CSRF token
	CookieSameSite http.SameSite
	// TokenExpiration is the expiration time of the CSRF token
	TokenExpiration time.Duration
	// ErrorHandler is a function that handles CSRF token errors
	ErrorHandler func(c echo.Context) error
}

// DefaultCSRFConfig returns the default configuration for the CSRF middleware
func DefaultCSRFConfig() CSRFConfig {
	return CSRFConfig{
		TokenLookup: func(c echo.Context) string {
			return c.Request().Header.Get(CSRFTokenHeader)
		},
		CookieName:      CSRFTokenCookie,
		CookiePath:      "/",
		CookieDomain:    "",
		CookieSecure:    false,
		CookieHTTPOnly:  true,
		CookieSameSite:  http.SameSiteLaxMode,
		TokenExpiration: 24 * time.Hour,
		ErrorHandler: func(c echo.Context) error {
			return c.JSON(http.StatusForbidden, NewErrorResponse(403, "CSRF token mismatch"))
		},
	}
}

// CSRF is a middleware that protects against CSRF attacks
type CSRF struct {
	config CSRFConfig
}

// NewCSRF creates a new CSRF middleware
func NewCSRF(config CSRFConfig) echo.MiddlewareFunc {
	// Use default config if not provided
	if config.TokenLookup == nil {
		config.TokenLookup = DefaultCSRFConfig().TokenLookup
	}
	if config.CookieName == "" {
		config.CookieName = DefaultCSRFConfig().CookieName
	}
	if config.CookiePath == "" {
		config.CookiePath = DefaultCSRFConfig().CookiePath
	}
	if config.TokenExpiration <= 0 {
		config.TokenExpiration = DefaultCSRFConfig().TokenExpiration
	}
	if config.ErrorHandler == nil {
		config.ErrorHandler = DefaultCSRFConfig().ErrorHandler
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Skip CSRF protection for GET, HEAD, OPTIONS requests
			if c.Request().Method == http.MethodGet || c.Request().Method == http.MethodHead || c.Request().Method == http.MethodOptions {
				return next(c)
			}

			// Get CSRF token from cookie
			cookie, err := c.Cookie(config.CookieName)
			if err != nil {
				return config.ErrorHandler(c)
			}

			// Get CSRF token from request
			token := config.TokenLookup(c)
			if token == "" {
				return config.ErrorHandler(c)
			}

			// Compare tokens
			if token != cookie.Value {
				return config.ErrorHandler(c)
			}

			// Call the handler
			return next(c)
		}
	}
}

// GenerateCSRFToken generates a new CSRF token
func GenerateCSRFToken() (string, error) {
	b := make([]byte, CSRFTokenLength)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

// SetCSRFToken sets the CSRF token in the cookie
func SetCSRFToken(c echo.Context, token string, config CSRFConfig) {
	cookie := new(http.Cookie)
	cookie.Name = config.CookieName
	cookie.Value = token
	cookie.Path = config.CookiePath
	cookie.Domain = config.CookieDomain
	cookie.Secure = config.CookieSecure
	cookie.HttpOnly = config.CookieHTTPOnly
	cookie.SameSite = config.CookieSameSite
	cookie.Expires = time.Now().Add(config.TokenExpiration)

	c.SetCookie(cookie)
	c.Response().Header().Set(CSRFTokenHeader, token)
}

// CSRFTokenMiddleware is a middleware that generates and sets CSRF tokens
func CSRFTokenMiddleware(config CSRFConfig) echo.MiddlewareFunc {
	// Use default config if not provided
	if config.CookieName == "" {
		config.CookieName = DefaultCSRFConfig().CookieName
	}
	if config.CookiePath == "" {
		config.CookiePath = DefaultCSRFConfig().CookiePath
	}
	if config.TokenExpiration <= 0 {
		config.TokenExpiration = DefaultCSRFConfig().TokenExpiration
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Check if CSRF token exists in cookie
			_, err := c.Cookie(config.CookieName)
			if err != nil {
				// Generate new CSRF token
				token, err := GenerateCSRFToken()
				if err != nil {
					return err
				}

				// Set CSRF token in cookie
				SetCSRFToken(c, token, config)
			}

			// Call the handler
			return next(c)
		}
	}
}
