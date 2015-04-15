package main

import (
	"fmt"
	"github.com/koding/kite"
	"net/http"
	"time"
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

	k.HandleHTTPFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		tm := time.Now().Format(time.RFC1123)
		w.Write([]byte("Hello! the time is: " + tm))
	})

	k.Run()
}
