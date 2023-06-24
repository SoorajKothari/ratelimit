package main

import (
	"math/rand"
	"net/http"
	"github.com/gin-gonic/gin"
)

func getVisits(context *gin.Context) {
	// some handler
	context.IndentedJSON(http.StatusOK, rand.Intn(1000))
}


func main() {
	router := gin.Default()
	router.GET("/getVisits", getVisits)
	router.Run("localhost:8080")
}
