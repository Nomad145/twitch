package main

import (
	"./twitch"
	"flag"
	"fmt"
	"os"
	"text/tabwriter"
)

const CLIENT_ID = "g0m4aoe1qgv0lqais31yp27yzvw603"
const CLIENT_SECRET = "rftkcblmnloq5q8ndm4i92kpbt6o8t"

var client = twitch.NewClient()
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
