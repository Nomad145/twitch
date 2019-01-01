package twitch

import (
	"github.com/grafov/m3u8"
	"os/exec"
	"time"
	"log"
)

type Player struct {
	Stream *StreamApi
}

var segments = make(chan string, 30)
var vlc = exec.Command("cvlc", "-")
var pipe, _ = vlc.StdinPipe()

func (player Player) StartStream(user string) {
	go player.processSegments(user)
	go player.downloadSegments()

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

			segments <- segment.URI
			segmentCache[segment.ProgramDateTime.String()] = segment.URI
		}

		time.Sleep(time.Second * 28)
		masterPlaylist = player.Stream.RefreshPlaylist(masterPlaylist)
	}
}

func (player Player) downloadSegments() {
	for segment := range segments {
		player.Stream.DownloadSegment(segment, pipe)
	}
}
