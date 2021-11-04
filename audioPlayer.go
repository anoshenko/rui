package rui

type AudioPlayer interface {
	MediaPlayer
}

type audioPlayerData struct {
	mediaPlayerData
}

// NewAudioPlayer create new MediaPlayer object and return it
func NewAudioPlayer(session Session, params Params) AudioPlayer {
	view := new(audioPlayerData)
	view.Init(session)
	view.tag = "AudioPlayer"
	setInitParams(view, params)
	return view
}

func newAudioPlayer(session Session) View {
	return NewAudioPlayer(session, nil)
}

func (player *audioPlayerData) Init(session Session) {
	player.mediaPlayerData.Init(session)
	player.tag = "AudioPlayer"
}

func (player *audioPlayerData) htmlTag() string {
	return "audio"
}
