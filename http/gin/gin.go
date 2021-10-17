package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	var r = gin.Default()

	r.GET("/ping", func(context *gin.Context) {
		context.String(http.StatusOK, "12344")
	})

	r.GET("/ping", func(context *gin.Context) {
		context.String(http.StatusOK, "12344")
	})

	_ = r.Run(":8899")
}
