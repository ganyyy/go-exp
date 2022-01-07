package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func TMMiddleWare(c *gin.Context) {
	var begin = time.Now()
	c.Next()
	log.Printf("[INF] time %v", time.Now().Sub(begin))
}

func main() {
	var r = gin.Default()
	r.Use(TMMiddleWare)
	r.GET("/ping", func(context *gin.Context) {
		context.String(http.StatusOK, "12344")
	})

	type Data struct {
		Name    string `form:"name"`
		Age     int    `form:"age"`
		Address string `form:"address"`
	}

	r.POST("/postform", func(context *gin.Context) {
		var d Data
		var err = context.Bind(&d)
		log.Println(err)
		var bs, _ = json.Marshal(d)
		log.Printf("%+v", string(bs))
		context.JSON(http.StatusOK, d)
	})

	r.POST("/postform2", func(context *gin.Context) {
		var d Data
		_ = context.BindJSON(&d)
		log.Printf("%+v", d)
		context.JSON(http.StatusOK, d)
	})

	r.POST("/voucher_post_callback")

	_ = r.Run(":8899")
}
