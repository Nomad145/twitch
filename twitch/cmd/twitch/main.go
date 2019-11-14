package main

import (
	"flag"
	"fmt"
	"github.com/michaeljoelphillips/twitch/internal/twitch"
	"os"
	"text/tabwriter"
)

var client = twitch.NewClient()
var list = flag.Bool("list", false, "List live channels")
var tabWriter = tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)

func main() {
	flag.Parse()

	if *list {
		liveStreams := client.ListStreams()

		for _, stream := range liveStreams {
			fmt.Fprintln(tabWriter, stream.Channel.Name+"\t"+stream.Channel.Game+"\t"+stream.Channel.Status)
		}

		tabWriter.Flush()

		return
	}

	if len(flag.Args()) > 0 {
		client.Play(flag.Args()[0])
	}
}
