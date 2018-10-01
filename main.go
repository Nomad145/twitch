package main

import (
	"./twitch"
	"flag"
	"fmt"
	"os"
	"text/tabwriter"
)

const CLIENT_ID = "g0m4aoe1qgv0lqais31yp27yzvw603"
const ACCESS_TOKEN = "iwrumiqe5ak9a89h586nd0xenro3b2"

var client = twitch.NewClient(CLIENT_ID, ACCESS_TOKEN)
var list = flag.Bool("list", false, "List live channels")
var tabWriter = tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)

func main() {
	flag.Parse()

	if *list {
		liveStreams := client.User.GetLiveStreams()

		for _, stream := range liveStreams {
			fmt.Fprintln(tabWriter, stream.Channel.Name + "\t" + stream.Channel.Game)
		}

		tabWriter.Flush()

		return
	}

	if len(flag.Args()) > 0 {
		user := client.User.GetUser(flag.Args()[0])

		client.Player.StartStream(user)
	}
}
