package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/labstack/echo/v4"
	"golang.org/x/time/rate"
)

// RateLimiterConfig defines the configuration for the rate limiter middleware
type RateLimiterConfig struct {
	// Rate is the number of requests allowed per second
	Rate float64
	// Burst is the maximum number of requests allowed in a burst
	Burst int
	// KeyFunc is a function that returns a unique key for each client
	KeyFunc func(c echo.Context) string
	// ErrorHandler is a function that handles rate limit errors
	ErrorHandler func(c echo.Context) error
}

// DefaultRateLimiterConfig returns the default configuration for the rate limiter middleware
func DefaultRateLimiterConfig() RateLimiterConfig {
	return RateLimiterConfig{
		Rate:  10, // 10 requests per second
		Burst: 20, // 20 requests in a burst
		KeyFunc: func(c echo.Context) string {
			return c.RealIP() // Use client IP as the key
		},
		ErrorHandler: func(c echo.Context) error {
			return c.JSON(http.StatusTooManyRequests, NewErrorResponse(429, "rate limit exceeded"))
		},
	}
}

// RateLimiter is a middleware that limits the rate of requests
type RateLimiter struct {
	config       RateLimiterConfig
	rateLimiters map[string]*rate.Limiter
	mu           sync.Mutex
}

// NewRateLimiter creates a new rate limiter middleware
func NewRateLimiter(config RateLimiterConfig) echo.MiddlewareFunc {
	// Use default config if not provided
	if config.Rate < 0 {
		config.Rate = DefaultRateLimiterConfig().Rate
	}
	if config.Burst < 0 {
		config.Burst = DefaultRateLimiterConfig().Burst
	}
	if config.KeyFunc == nil {
		config.KeyFunc = DefaultRateLimiterConfig().KeyFunc
	}
	if config.ErrorHandler == nil {
		config.ErrorHandler = DefaultRateLimiterConfig().ErrorHandler
	}

	rl := &RateLimiter{
		config:       config,
		rateLimiters: make(map[string]*rate.Limiter),
	}

	// Start a background goroutine to clean up expired rate limiters
	go func() {
		for {
			time.Sleep(10 * time.Minute)
			rl.cleanup()
		}
	}()

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			key := rl.config.KeyFunc(c)
			limiter := rl.getLimiter(key)

			if !limiter.Allow() {
				return rl.config.ErrorHandler(c)
			}

			return next(c)
		}
	}
}

// getLimiter returns a rate limiter for the given key, creating a new one if needed
func (rl *RateLimiter) getLimiter(key string) *rate.Limiter {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	limiter, exists := rl.rateLimiters[key]
	if !exists {
		limiter = rate.NewLimiter(rate.Limit(rl.config.Rate), rl.config.Burst)
		rl.rateLimiters[key] = limiter
	}

	return limiter
}

// cleanup removes expired rate limiters
func (rl *RateLimiter) cleanup() {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	// For simplicity, we'll just clear all rate limiters for now
	// In a production environment, you might want to implement a more sophisticated cleanup strategy
	rl.rateLimiters = make(map[string]*rate.Limiter)
}
