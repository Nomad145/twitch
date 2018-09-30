package main

import (
	"./twitch"
	"github.com/grafov/m3u8"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
)

const client_id = "g0m4aoe1qgv0lqais31yp27yzvw603"
const client_secret = "xhk8oay2vvhmpc63ingdsbavix2j7k"
const access_token = "twm4mh8qim2hawaf13jv8ia5i8eay5"

var client = twitch.NewClient(client_id)
var segments = make(chan string, 1024)
var vlc = exec.Command("cvlc", "-")
var pipe, _ = vlc.StdinPipe()

func main() {
	user := client.User.Get(string(os.Args[1]))
	go fetchSegments(user)

	downloadSegments()
}

func fetchSegments(user twitch.User) {
	for {
		masterPlaylist := client.Stream.GetStreamPlaylist(user)
		stream := client.Stream.GetBest(masterPlaylist)
		mediaPlaylist := stream.Playlist.(*m3u8.MediaPlaylist)

		for _, v := range mediaPlaylist.Segments {
			if v != nil {
				segments <- v.URI
			}
		}
	}
}

func downloadSegments() {
	vlc.Start()

	for v := range segments {
		data, _ := http.Get(v)
		body, _ := ioutil.ReadAll(data.Body)
		pipe.Write(body)
		data.Body.Close()
	}
}
