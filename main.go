package main

import (
	"fmt"
	"./twitch"
)

const client_id = "g0m4aoe1qgv0lqais31yp27yzvw603"
const client_secret = "xhk8oay2vvhmpc63ingdsbavix2j7k"
const access_token = "twm4mh8qim2hawaf13jv8ia5i8eay5"

var client = twitch.NewClient(client_id)

func main() {
	user := client.User.Get("specsnstats")
	playlist := client.Stream.GetStreamPlaylist(user)

	fmt.Println(playlist)
}
