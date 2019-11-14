package twitch

import (
	"github.com/grafov/m3u8"
	"io/ioutil"
	"log"
	"os/exec"
)

type Player struct {
	Stream *StreamApi
}

var segments = make(chan string, 30)
var adSegments = make(chan string, 30)
var vlc = exec.Command("cvlc", "-")
var pipe, _ = vlc.StdinPipe()

func (player Player) StartStream(user string) {
	go player.processSegments(user)
	go player.downloadSegments()
	go player.downloadAdSegments()

	vlc.Run()
}

func (player Player) processSegments(user string) {
	segmentCache := make(map[string]string)
	masterPlaylist := player.Stream.GetMasterPlaylist(user)

	for {
		if masterPlaylist.Playlist == nil {
			masterPlaylist = player.Stream.GetMasterPlaylist(user)

			log.Println("Stream playlist was empty")

			continue
		}

		stream := player.Stream.GetMediaPlaylist(masterPlaylist)

		mediaPlaylist := stream.Playlist.(*m3u8.MediaPlaylist)

		for _, segment := range mediaPlaylist.Segments {
			if segment == nil {
				continue
			}

			if _, hit := segmentCache[segment.ProgramDateTime.String()]; hit {
				continue
			}

			if segment.Title != "live" {

				continue
			}

			segments <- segment.URI
			segmentCache[segment.ProgramDateTime.String()] = segment.URI
		}

		masterPlaylist = player.Stream.RefreshPlaylist(masterPlaylist)
	}
}

func (player Player) downloadSegments() {
	for segment := range segments {
		player.Stream.DownloadSegment(segment, pipe)
	}
}

func (player Player) downloadAdSegments() {
	for segment := range adSegments {
		player.Stream.DownloadSegment(segment, ioutil.Discard)
	}
}
