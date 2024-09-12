package rui

// AudioPlayer is a type of a [View] which can play audio files
type AudioPlayer interface {
	MediaPlayer
}

type audioPlayerData struct {
	mediaPlayerData
}

// NewAudioPlayer create new MediaPlayer object and return it
func NewAudioPlayer(session Session, params Params) AudioPlayer {
	view := new(audioPlayerData)
	view.init(session)
	view.tag = "AudioPlayer"
	setInitParams(view, params)
	return view
}

func newAudioPlayer(session Session) View {
	return NewAudioPlayer(session, nil)
}

func (player *audioPlayerData) init(session Session) {
	player.mediaPlayerData.init(session)
	player.tag = "AudioPlayer"
}

func (player *audioPlayerData) String() string {
	return getViewString(player, nil)
}

func (player *audioPlayerData) htmlTag() string {
	return "audio"
}
