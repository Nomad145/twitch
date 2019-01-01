package main

import (
	"./twitch"
	"flag"
	"fmt"
	"os"
	"text/tabwriter"
)

var client = twitch.NewClient()
var list = flag.Bool("list", false, "List live channels")
var tabWriter = tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)

func main() {
	flag.Parse()

	if *list {
		liveStreams := client.User.GetLiveStreams()

		for _, stream := range liveStreams {
			fmt.Fprintln(tabWriter, stream.Channel.Name + "\t" + stream.Channel.Game + "\t" + stream.Channel.Status)
		}

		tabWriter.Flush()

		return
	}

	if len(flag.Args()) > 0 {
		client.Player.StartStream(flag.Args()[0])
	}
}
