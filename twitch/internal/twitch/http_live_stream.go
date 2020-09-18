package twitch

import (
	"github.com/grafov/m3u8"
	"io"
	"log"
)

type LiveStream interface {
	Start(user string, writer io.Writer)
}

type HttpLiveStream struct {
	provider MediaPlaylistProvider
}

var videoPlaylistSegments = make(chan string, 30)

func (stream HttpLiveStream) Start(user string, writer io.Writer) {
	go stream.parseMediaPlaylist(user)
	go stream.downloadVideoSegments(writer)
}

func (liveStream HttpLiveStream) parseMediaPlaylist(user string) {
	segmentCache := make(map[string]string)
	masterPlaylist := liveStream.provider.GetMasterPlaylist(user)

	if masterPlaylist.Playlist == nil {
		log.Fatal("Unable to fetch playlist")
	}

	stream := liveStream.provider.GetMediaPlaylist(masterPlaylist)

	for {
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

			videoPlaylistSegments <- segment.URI
			segmentCache[segment.ProgramDateTime.String()] = segment.URI
		}

		stream = liveStream.provider.RefreshPlaylist(stream)
	}
}

func (stream HttpLiveStream) downloadVideoSegments(writer io.Writer) {
	for segment := range videoPlaylistSegments {
		stream.provider.DownloadSegment(segment, writer)
	}
}
