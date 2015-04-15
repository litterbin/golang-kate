package main

import (
	"github.com/labstack/echo"
	mw "github.com/labstack/echo/middleware"
	"github.com/thoas/stats"
	"net/http"

	"fmt"
	"github.com/koding/kite"
)

func hello(c *echo.Context) {
	c.String(http.StatusOK, "hello world")
}

func main() {
	//kite
	k := kite.New("math", "1.0.0")
	k.Config.Port = 7000
	k.Config.DisableAuthentication = true

	k.HandleFunc("square", func(r *kite.Request) (interface{}, error) {
		a := r.Args.One().MustFloat64()
		r.Client.Go("kite.log", fmt.Sprintf("Message from %s: \"You have requested square of %.0f\"", r.LocalKite.Kite().Name, a))

		return a * a, nil
	})

	e := echo.New()

	// Use kite!
	e.Use(k)

	e.Use(mw.Logger)

	// https://github.com/thoas/stats
	s := stats.New()
	e.Use(s.Handler)

	// Route
	e.Get("/stats", func(c *echo.Context) {
		c.JSON(200, s.Data())
	})

	e.Get("/", func(c *echo.Context) {
		c.String(http.StatusOK, "HELLO")
	})

	e.Get("/hello", hello)

	//Start server
	//e.Run(":7000")

	k.HandleHTTP("/echo", e)
	k.Run()
}
