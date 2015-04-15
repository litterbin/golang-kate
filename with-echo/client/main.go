package main

import (
	"fmt"
	"github.com/koding/kite"
)

func main() {
	k := kite.New("second", "1.0.0")

	client := k.NewClient("http://localhost:7000/kite")
	client.Dial()

	pong, _ := client.Tell("kite.ping")
	fmt.Println(pong.MustString())

	response, _ := client.Tell("square", 4)
	fmt.Println(response.MustFloat64())
}
