package main

import (
	"./twitch"
	"os"
)

const client_id = "g0m4aoe1qgv0lqais31yp27yzvw603"

var client = twitch.NewClient(client_id)

func main() {
	user := client.User.Get(string(os.Args[1]))

	client.Player.StartStream(user)
}
