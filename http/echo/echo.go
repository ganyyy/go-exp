package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

func main() {
	var e = echo.New()

	e.Logger.SetLevel(log.DEBUG)

	e.GET("/showInfo", func(context echo.Context) error {
		return context.String(http.StatusOK, "Hello world!")
	})

	e.GET("/showInfo", func(context echo.Context) error {
		return context.String(http.StatusOK, "1234")
	})

	e.GET("/showInfo", func(context echo.Context) error {
		context.Logger().Debugf("this is a request. Params:%v, Values:%v", context.ParamNames(), context.ParamValues())
		return context.String(http.StatusOK, "gan")
	})

	e.Logger.Fatal(e.Start(":8899"))
}
