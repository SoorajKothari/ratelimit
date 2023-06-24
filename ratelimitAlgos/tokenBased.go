package main

import (
	"math/rand"
	"net/http"
	"github.com/gin-gonic/gin"
	"sync"
	"time"
)

type RateLimiter struct {
	tokens int
	maxTokens int
	refillRate int
	mutex  sync.Mutex
}

func NewRateLimiter(maxTokens, refillRate int) *RateLimiter {
	return &RateLimiter{
		tokens:     maxTokens,
		maxTokens:  maxTokens,
		refillRate: refillRate,
	}
}
func (rl *RateLimiter) RefillTokens() {
	for {
		time.Sleep(time.Second)
		rl.mutex.Lock()
		if rl.tokens < rl.maxTokens {
			rl.tokens++
		}
		rl.mutex.Unlock()
	}
}

func (rl *RateLimiter) AcquireToken() bool {
	rl.mutex.Lock()
	defer rl.mutex.Unlock()

	if rl.tokens > 0 {
		rl.tokens--
		return true
	}
	return false
}

func getVisits(context *gin.Context, rateLimiter *RateLimiter) {
	if rateLimiter.AcquireToken() {
		context.IndentedJSON(http.StatusOK, rand.Intn(1000))
	} else {
		context.AbortWithStatus(http.StatusTooManyRequests)
	}
}

func main() {
	// Create a rate limiter with a maximum of 1 token and refill rate of 2 tokens per second
	rateLimiter := NewRateLimiter(1, 2)

	// Start the token refill routine
	go rateLimiter.RefillTokens()

	router := gin.Default()
	router.GET("/getVisits", func(context *gin.Context) {
		getVisits(context, rateLimiter)
	})
	router.Run("localhost:8080")
}

