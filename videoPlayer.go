package rui

import (
	"fmt"
	"strings"
)

const (
	// VideoWidth is the constant for the "video-width" property tag of VideoPlayer.
	// The "video-width" float property defines the width of the video's display area in pixels.
	VideoWidth = "video-width"
	// VideoHeight is the constant for the "video-height" property tag of VideoPlayer.
	// The "video-height" float property defines the height of the video's display area in pixels.
	VideoHeight = "video-height"
	// Poster is the constant for the "poster" property tag of VideoPlayer.
	// The "poster" property defines an URL for an image to be shown while the video is downloading.
	// If this attribute isn't specified, nothing is displayed until the first frame is available,
	// then the first frame is shown as the poster frame.
	Poster = "poster"
)

type VideoPlayer interface {
	MediaPlayer
}

type videoPlayerData struct {
	mediaPlayerData
}

// NewVideoPlayer create new MediaPlayer object and return it
func NewVideoPlayer(session Session, params Params) VideoPlayer {
	view := new(videoPlayerData)
	view.Init(session)
	view.tag = "VideoPlayer"
	setInitParams(view, params)
	return view
}

func newVideoPlayer(session Session) View {
	return NewVideoPlayer(session, nil)
}

func (player *videoPlayerData) Init(session Session) {
	player.mediaPlayerData.Init(session)
	player.tag = "VideoPlayer"
}

func (player *videoPlayerData) htmlTag() string {
	return "video"
}

func (player *videoPlayerData) Remove(tag string) {
	player.remove(strings.ToLower(tag))
}

func (player *videoPlayerData) remove(tag string) {
	switch tag {

	case VideoWidth:
		delete(player.properties, tag)
		removeProperty(player.htmlID(), "width", player.Session())

	case VideoHeight:
		delete(player.properties, tag)
		removeProperty(player.htmlID(), "height", player.Session())

	case Poster:
		delete(player.properties, tag)
		removeProperty(player.htmlID(), Poster, player.Session())

	default:
		player.mediaPlayerData.remove(tag)
	}
}

func (player *videoPlayerData) Set(tag string, value interface{}) bool {
	return player.set(strings.ToLower(tag), value)
}

func (player *videoPlayerData) set(tag string, value interface{}) bool {
	if value == nil {
		player.remove(tag)
		return true
	}

	if player.mediaPlayerData.set(tag, value) {
		session := player.Session()
		updateSize := func(cssTag string) {
			if size, ok := floatProperty(player, tag, session, 0); ok {
				if size > 0 {
					updateProperty(player.htmlID(), cssTag, fmt.Sprintf("%g", size), session)
				} else {
					removeProperty(player.htmlID(), cssTag, session)
				}
			}
		}

		switch tag {
		case VideoWidth:
			updateSize("width")

		case VideoHeight:
			updateSize("height")

		case Poster:
			if url, ok := stringProperty(player, Poster, session); ok {
				updateProperty(player.htmlID(), Poster, url, session)
			}
		}
		return true
	}

	return false
}

func (player *videoPlayerData) htmlProperties(self View, buffer *strings.Builder) {
	player.mediaPlayerData.htmlProperties(self, buffer)

	session := player.Session()

	if size, ok := floatProperty(player, VideoWidth, session, 0); ok && size > 0 {
		buffer.WriteString(fmt.Sprintf(` width="%g"`, size))
	}

	if size, ok := floatProperty(player, VideoHeight, session, 0); ok && size > 0 {
		buffer.WriteString(fmt.Sprintf(` height="%g"`, size))
	}

	if url, ok := stringProperty(player, Poster, session); ok && url != "" {
		buffer.WriteString(` poster="`)
		buffer.WriteString(url)
		buffer.WriteString(`"`)
	}
}
