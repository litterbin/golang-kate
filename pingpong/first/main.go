package main

import (
	"fmt"
	"github.com/koding/kite"
)

var (
	Clients = []*kite.Client{}
)

func main() {
	k := kite.New("first", "1.0.0")
	k.Config.Port = 6000
	k.Config.DisableAuthentication = true

	// Add OnConnect handler method
	k.OnConnect(func(c *kite.Client) {
		Clients = append(Clients, c)
		fmt.Println(Clients)
	})

	// Add pre handler method
	k.PreHandleFunc(func(r *kite.Request) (interface{}, error) {
		fmt.Println("\nThis pre handler is executed before the method is executed")

		for _, c := range Clients {
			pong, _ := c.Tell("kite.ping")
			fmt.Println("c:", c.Name, "pong:", pong.MustString())
		}

		// let us return an hello to base square method!
		return "hello from pre handler!", nil
	})

	// Add post handler method
	//TODO: I don't know how do use PostHandleFunc.
	/*
		k.PostHandleFunc(func(r *kite.Request) (interface{}, error) {
			fmt.Println("This post handler is executed after the method is executed")

			// pass the response from the previous square method back to the
			// client, this is imporant if you use post handler
			return r.Response, nil
		})
	*/

	k.HandleFunc("square", func(r *kite.Request) (interface{}, error) {
		a := r.Args.One().MustFloat64()
		r.Client.Go("kite.log", fmt.Sprintf("Message from %s: \"You have requested square of %.0f\"", r.LocalKite.Kite().Name, a))

		return a * a, nil
	})

	k.Run()
}
