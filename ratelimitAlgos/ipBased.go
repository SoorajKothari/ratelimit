package main

import (
	"net/http"
	"sync"
	"time"
	"fmt"
	"math/rand"
	"github.com/gin-gonic/gin"
)

const (
	maxRequests = 2
	perMinutePeriod = 1 * time.Minute
)

var (
	ipRequestsCounts = make(map[string]int)
	mutex = &sync.Mutex{}
)

func rateLimiter(context *gin.Context) {
	ip := context.ClientIP()
	fmt.Println(ip)
	mutex.Lock()
	defer mutex.Unlock()
	count := ipRequestsCounts[ip]
	if count >= maxRequests{
		context.AbortWithStatus(http.StatusTooManyRequests)
		return
	}

	ipRequestsCounts[ip] = count + 1
	fmt.Print(ipRequestsCounts[ip])
	time.AfterFunc(perMinutePeriod, func() {
		mutex.Lock()
		defer mutex.Unlock()

		ipRequestsCounts[ip] = ipRequestsCounts[ip] - 1
	})

	context.Next()
}

func getVisits(context *gin.Context) {
	// some handler
	context.IndentedJSON(http.StatusOK, rand.Intn(1000))
}

func main() {
	router := gin.Default()
	router.Use(rateLimiter)
	router.GET("/getVisits", getVisits)
	router.Run("localhost:8080")
}
