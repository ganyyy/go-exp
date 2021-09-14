package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func main() {
	var e = echo.New()

	e.GET("/showInfo", func(context echo.Context) error {
		return context.String(http.StatusOK, "Hello world!")
	})

	e.GET("/showInfo", func(context echo.Context) error {
		return context.String(http.StatusOK, "1234")
	})

	e.GET("/showInfo", func(context echo.Context) error {
		return context.String(http.StatusOK, "gan")
	})

	e.Logger.Fatal(e.Start(":9900"))
}
