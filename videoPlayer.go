package rui

import (
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
	view.init(session)
	view.tag = "VideoPlayer"
	setInitParams(view, params)
	return view
}

func newVideoPlayer(session Session) View {
	return NewVideoPlayer(session, nil)
}

func (player *videoPlayerData) init(session Session) {
	player.mediaPlayerData.init(session)
	player.tag = "VideoPlayer"
}

func (player *videoPlayerData) String() string {
	return getViewString(player)
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
		player.session.removeProperty(player.htmlID(), "width")

	case VideoHeight:
		delete(player.properties, tag)
		player.session.removeProperty(player.htmlID(), "height")

	case Poster:
		delete(player.properties, tag)
		player.session.removeProperty(player.htmlID(), Poster)

	default:
		player.mediaPlayerData.remove(tag)
	}
}

func (player *videoPlayerData) Set(tag string, value any) bool {
	return player.set(strings.ToLower(tag), value)
}

func (player *videoPlayerData) set(tag string, value any) bool {
	if value == nil {
		player.remove(tag)
		return true
	}

	if player.mediaPlayerData.set(tag, value) {
		session := player.Session()
		updateSize := func(cssTag string) {
			if size, ok := floatTextProperty(player, tag, session, 0); ok {
				if size != "0" {
					session.updateProperty(player.htmlID(), cssTag, size)
				} else {
					session.removeProperty(player.htmlID(), cssTag)
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
				session.updateProperty(player.htmlID(), Poster, url)
			}
		}
		return true
	}

	return false
}

func (player *videoPlayerData) htmlProperties(self View, buffer *strings.Builder) {
	player.mediaPlayerData.htmlProperties(self, buffer)

	session := player.Session()

	if size, ok := floatTextProperty(player, VideoWidth, session, 0); ok && size != "0" {
		buffer.WriteString(` width="`)
		buffer.WriteString(size)
		buffer.WriteString(`"`)
	}

	if size, ok := floatTextProperty(player, VideoHeight, session, 0); ok && size != "0" {
		buffer.WriteString(` height="`)
		buffer.WriteString(size)
		buffer.WriteString(`"`)
	}

	if url, ok := stringProperty(player, Poster, session); ok && url != "" {
		buffer.WriteString(` poster="`)
		buffer.WriteString(url)
		buffer.WriteString(`"`)
	}
}
