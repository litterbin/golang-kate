package main

import (
	"github.com/gin-gonic/gin"
	"net/http"

	"fmt"
	"github.com/koding/kite"
)

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

	r := gin.Default()
	r.GET("/hello", func(c *gin.Context) {
		c.String(http.StatusOK, "hello world")
	})

	k.HandleHTTP("/hello", r)
	k.Run()
}
