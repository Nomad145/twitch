package twitch

import (
	"github.com/grafov/m3u8"
	"os/exec"
	"time"
)

type Player struct {
	Stream *StreamApi
}

var segments = make(chan string, 30)
var vlc = exec.Command("cvlc", "-")
var pipe, _ = vlc.StdinPipe()

func (player Player) StartStream(user User) {
	go player.fetchSegments(user)

	player.downloadSegments()
}

func (player Player) fetchSegments(user User) {
	segmentCache := make(map[string]string)
	masterPlaylist := player.Stream.GetMasterPlaylist(user)

	for {
		stream := player.Stream.GetMediaPlaylist(masterPlaylist)
		mediaPlaylist := stream.Playlist.(*m3u8.MediaPlaylist)

		for _, v := range mediaPlaylist.Segments {
			if v == nil {
				continue
			}

			if _, hit := segmentCache[v.ProgramDateTime.String()]; hit {
				continue
			}

			segments <- v.URI
			segmentCache[v.ProgramDateTime.String()] = v.URI
		}

		time.Sleep(time.Second * 28)
		masterPlaylist = player.Stream.RefreshPlaylist(masterPlaylist)
	}
}

func (player Player) downloadSegments() {
	vlc.Start()

	for segment := range segments {
		player.Stream.DownloadSegment(segment, pipe)
	}
}
