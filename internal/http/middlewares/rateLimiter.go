package middlewares

import (
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

type RateLimiter struct {
	visitors map[string]*rate.Limiter
	mu       sync.Mutex
	rate     rate.Limit
	burst    int
}

// NewRateLimiter initializes the rate limiter with a rate and burst limit
func NewRateLimiter(r rate.Limit, b int) *RateLimiter {
	return &RateLimiter{
		visitors: make(map[string]*rate.Limiter),
		rate:     r,
		burst:    b,
	}
}

// getVisitor retrieves or creates a rate limiter for a visitor based on device ID
func (rl *RateLimiter) getVisitor(deviceID string) *rate.Limiter {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	limiter, exists := rl.visitors[deviceID]
	if !exists {
		limiter = rate.NewLimiter(rl.rate, rl.burst)
		rl.visitors[deviceID] = limiter
	}

	return limiter
}

// cleanupVisitors removes visitors who haven't been seen for some time
func (rl *RateLimiter) cleanupVisitors() {
	for {
		time.Sleep(time.Minute)
		rl.mu.Lock()

		for deviceID, limiter := range rl.visitors {
			// If the visitor's limiter hasn't been used recently, remove it.
			if limiter.AllowN(time.Now(), rl.burst) {
				delete(rl.visitors, deviceID)
			}
		}

		rl.mu.Unlock()
	}
}

// Limit is the middleware function that enforces rate limiting
func (rl *RateLimiter) Limit() gin.HandlerFunc {
	go rl.cleanupVisitors() // Start the cleanup routine

	return func(c *gin.Context) {
		clientIP := c.ClientIP()
		if clientIP == "" {
			c.AbortWithStatusJSON(400, gin.H{"error": "Could not determine Client IP!"})
			return
		}

		limiter := rl.getVisitor(clientIP)

		if !limiter.Allow() {
			c.AbortWithStatus(429) // HTTP 429 Too Many Requests
			return
		}

		c.Next() // Allow the request to proceed
	}
}
