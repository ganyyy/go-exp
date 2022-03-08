package src

import (
	"log"
	"net/http"
	"sync/atomic"

	"github.com/gin-gonic/gin"
)

var globalContent atomic.Value

func init() {
	globalContent.Store("default content")
}

func RunGinServer() {
	var engine = gin.New()

	engine.LoadHTMLGlob("templates/*.tmpl")

	engine.GET("/", func(context *gin.Context) {
		context.HTML(http.StatusOK, "index.tmpl", map[string]interface{}{
			"content": globalContent.Load().(string),
		})
	})

	engine.GET("/content", func(context *gin.Context) {
		context.HTML(http.StatusOK, "content.tmpl", map[string]interface{}{
			"content": globalContent.Load().(string),
		})
	})

	engine.POST("/content", func(context *gin.Context) {
		var content PostContent
		var err = context.Bind(&content)
		if err != nil {
			log.Printf("error:%v", err)
		}
		log.Printf("[POST] receive post request:[%v]", content.Content)
		globalContent.Store(content.Content)
		context.Redirect(http.StatusFound, "/")
	})

	_ = engine.Run("0.0.0.0:8899")
}
