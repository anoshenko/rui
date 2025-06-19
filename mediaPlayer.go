package rui

import (
	"fmt"
	"math"
	"reflect"
	"strconv"
	"strings"
)

// Constants which related to media player properties and events
const (
	// Controls is the constant for "controls" property tag.
	//
	// Used by AudioPlayer, VideoPlayer.
	//
	// Controls whether the browser need to provide controls to allow user to control audio playback, volume, seeking and
	// pause/resume playback. Default value is false.
	//
	// Supported types: bool, int, string.
	//
	// Values:
	//   - true, 1, "true", "yes", "on", "1" - The browser will offer controls to allow the user to control audio playback, volume, seeking and pause/resume playback.
	//   - false, 0, "false", "no", "off", "0" - No controls will be visible to the end user.
	Controls PropertyName = "controls"

	// Loop is the constant for "loop" property tag.
	//
	// Used by AudioPlayer, VideoPlayer.
	//
	// Controls whether the audio player will play media in a loop. Default value is false.
	//
	// Supported types: bool, int, string.
	//
	// Values:
	//   - true, 1, "true", "yes", "on", "1" - The audio player will automatically seek back to the start upon reaching the end of the audio.
	//   - false, 0, "false", "no", "off", "0" - Audio player will stop playing when the end of the media file has been reached.
	Loop PropertyName = "loop"

	// Muted is the constant for "muted" property tag.
	//
	// Used by AudioPlayer, VideoPlayer.
	//
	// Controls whether the audio will be initially silenced. Default value is false.
	//
	// Supported types: bool, int, string.
	//
	// Values:
	//   - true, 1, "true", "yes", "on", "1" - Audio will be muted.
	//   - false, 0, "false", "no", "off", "0" - Audio playing normally.
	Muted PropertyName = "muted"

	// Preload is the constant for "preload" property tag.
	//
	// Used by AudioPlayer, VideoPlayer.
	//
	// Property is intended to provide a hint to the browser about what the author thinks will lead to the best user
	// experience. Default value is different for each browser.
	//
	// Supported types: int, string.
	//
	// Values:
	//   - 0 (PreloadNone) or "none" - Media file must not be pre-loaded.
	//   - 1 (PreloadMetadata) or "metadata" - Only metadata is preloaded.
	//   - 2 (PreloadAuto) or "auto" - The entire media file can be downloaded even if the user doesn't have to use it.
	Preload PropertyName = "preload"

	// AbortEvent is the constant for "abort-event" property tag.
	//
	// Used by AudioPlayer, VideoPlayer.
	//
	// Fired when the resource was not fully loaded, but not as the result of an error.
	//
	// General listener format:
	//
	//  func(player rui.MediaPlayer)
	//
	// where:
	// player - Interface of a player which generated this event.
	//
	// Allowed listener formats:
	//
	//  func()
	AbortEvent PropertyName = "abort-event"

	// CanPlayEvent is the constant for "can-play-event" property tag.
	//
	// Used by AudioPlayer, VideoPlayer.
	//
	// Occur when the browser can play the media, but estimates that not enough data has been loaded to play the media up to
	// its end without having to stop for further buffering of content.
	//
	// General listener format:
	//
	//  func(player rui.MediaPlayer)
	//
	// where:
	// player - Interface of a player which generated this event.
	//
	// Allowed listener formats:
	//
	//  func()
	CanPlayEvent PropertyName = "can-play-event"

	// CanPlayThroughEvent is the constant for "can-play-through-event" property tag.
	//
	// Used by AudioPlayer, VideoPlayer.
	//
	// Occur when the browser estimates it can play the media up to its end without stopping for content buffering.
	//
	// General listener format:
	//
	//  func(player rui.MediaPlayer)
	//
	// where:
	// player - Interface of a player which generated this event.
	//
	// Allowed listener formats:
	//
	//  func()
	CanPlayThroughEvent PropertyName = "can-play-through-event"

	// CompleteEvent is the constant for "complete-event" property tag.
	//
	// Used by AudioPlayer, VideoPlayer.
	//
	// Occur when the rendering of an OfflineAudioContext has been terminated.
	//
	// General listener format:
	//
	//  func(player rui.MediaPlayer)
	//
	// where:
	// player - Interface of a player which generated this event.
	//
	// Allowed listener formats:
	//
	//  func()
	CompleteEvent PropertyName = "complete-event"

	// DurationChangedEvent is the constant for "duration-changed-event" property tag.
	//
	// Used by AudioPlayer, VideoPlayer.
	//
	// Occur when the duration attribute has been updated.
	//
	// General listener format:
	//
	//  func(player rui.MediaPlayer, duration float64).
	//
	// where:
	//   - player - Interface of a player which generated this event,
	//   - duration - Current duration.
	//
	// Allowed listener formats:
	//
	//  func(player rui.MediaPlayer),
	//  func(duration float64),
	//  func()
	DurationChangedEvent PropertyName = "duration-changed-event"

	// EmptiedEvent is the constant for "emptied-event" property tag.
	//
	// Used by AudioPlayer, VideoPlayer.
	//
	// Occur when the media has become empty; for example, this event is sent if the media has already been loaded(or
	// partially loaded), and the HTMLMediaElement.load method is called to reload it.
	//
	// General listener format:
	//
	//  func(player rui.MediaPlayer)
	//
	// where:
	// player - Interface of a player which generated this event.
	//
	// Allowed listener formats:
	//
	//  func()
	EmptiedEvent PropertyName = "emptied-event"

	// EndedEvent is the constant for "ended-event" property tag.
	//
	// Used by AudioPlayer, VideoPlayer.
	//
	// Occur when the playback has stopped because the end of the media was reached.
	//
	// General listener format:
	//
	//  func(player rui.MediaPlayer)
	//
	// where:
	// player - Interface of a player which generated this event.
	//
	// Allowed listener formats:
	//
	//  func()
	EndedEvent PropertyName = "ended-event"

	// LoadedDataEvent is the constant for "loaded-data-event" property tag.
	//
	// Used by AudioPlayer, VideoPlayer.
	//
	// Occur when the first frame of the media has finished loading.
	//
	// General listener format:
	//
	//  func(player rui.MediaPlayer)
	//
	// where:
	// player - Interface of a player which generated this event.
	//
	// Allowed listener formats:
	//
	//  func()
	LoadedDataEvent PropertyName = "loaded-data-event"

	// LoadedMetadataEvent is the constant for "loaded-metadata-event" property tag.
	//
	// Used by AudioPlayer, VideoPlayer.
	//
	// Occur when the metadata has been loaded.
	//
	// General listener format:
	//
	//  func(player rui.MediaPlayer)
	//
	// where:
	// player - Interface of a player which generated this event.
	//
	// Allowed listener formats:
	//
	//  func()
	LoadedMetadataEvent PropertyName = "loaded-metadata-event"

	// LoadStartEvent is the constant for "load-start-event" property tag.
	//
	// Used by AudioPlayer, VideoPlayer.
	//
	// Fired when the browser has started to load a resource.
	//
	// General listener format:
	//
	//  func(player rui.MediaPlayer)
	//
	// where:
	// player - Interface of a player which generated this event.
	//
	// Allowed listener formats:
	//
	//  func()
	LoadStartEvent PropertyName = "load-start-event"

	// PauseEvent is the constant for "pause-event" property tag.
	//
	// Used by AudioPlayer, VideoPlayer.
	//
	// Occur when the playback has been paused.
	//
	// General listener format:
	//
	//  func(player rui.MediaPlayer)
	//
	// where:
	// player - Interface of a player which generated this event.
	//
	// Allowed listener formats:
	//
	//  func()
	PauseEvent PropertyName = "pause-event"

	// PlayEvent is the constant for "play-event" property tag.
	//
	// Used by AudioPlayer, VideoPlayer.
	//
	// Occur when the playback has begun.
	//
	// General listener format:
	//
	//  func(player rui.MediaPlayer)
	//
	// where:
	// player - Interface of a player which generated this event.
	//
	// Allowed listener formats:
	//
	//  func()
	PlayEvent PropertyName = "play-event"

	// PlayingEvent is the constant for "playing-event" property tag.
	//
	// Used by AudioPlayer, VideoPlayer.
	//
	// Occur when the playback is ready to start after having been paused or delayed due to lack of data.
	//
	// General listener format:
	//
	//  func(player rui.MediaPlayer)
	//
	// where:
	// player - Interface of a player which generated this event.
	//
	// Allowed listener formats:
	//
	//  func()
	PlayingEvent PropertyName = "playing-event"

	// ProgressEvent is the constant for "progress-event" property tag.
	//
	// Used by AudioPlayer, VideoPlayer.
	//
	// Fired periodically as the browser loads a resource.
	//
	// General listener format:
	//
	//  func(player rui.MediaPlayer)
	//
	// where:
	// player - Interface of a player which generated this event.
	//
	// Allowed listener formats:
	//
	//  func()
	ProgressEvent PropertyName = "progress-event"

	// RateChangedEvent is the constant for "rate-changed-event" property tag.
	//
	// Used by AudioPlayer, VideoPlayer.
	//
	// Occur when the playback rate has changed.
	//
	// General listener format:
	//
	//  func(player rui.MediaPlayer, rate float64).
	//
	// where:
	//   - player - Interface of a player which generated this event,
	//   - rate - Playback rate.
	//
	// Allowed listener formats:
	//
	//  func(player rui.MediaPlayer),
	//  func(rate float64),
	//  func()
	RateChangedEvent PropertyName = "rate-changed-event"

	// SeekedEvent is the constant for "seeked-event" property tag.
	//
	// Used by AudioPlayer, VideoPlayer.
	//
	// Occur when a seek operation completed.
	//
	// General listener format:
	//
	//  func(player rui.MediaPlayer)
	//
	// where:
	// player - Interface of a player which generated this event.
	//
	// Allowed listener formats:
	//
	//  func()
	SeekedEvent PropertyName = "seeked-event"

	// SeekingEvent is the constant for "seeking-event" property tag.
	//
	// Used by AudioPlayer, VideoPlayer.
	//
	// Occur when a seek operation has began.
	//
	// General listener format:
	//
	//  func(player rui.MediaPlayer)
	//
	// where:
	// player - Interface of a player which generated this event.
	//
	// Allowed listener formats:
	//
	//  func()
	SeekingEvent PropertyName = "seeking-event"

	// StalledEvent is the constant for "stalled-event" property tag.
	//
	// Used by AudioPlayer, VideoPlayer.
	//
	// Occur when the user agent is trying to fetch media data, but data is unexpectedly not forthcoming.
	//
	// General listener format:
	//
	//  func(player rui.MediaPlayer)
	//
	// where:
	// player - Interface of a player which generated this event.
	//
	// Allowed listener formats:
	//
	//  func()
	StalledEvent PropertyName = "stalled-event"

	// SuspendEvent is the constant for "suspend-event" property tag.
	//
	// Used by AudioPlayer, VideoPlayer.
	//
	// Occur when the media data loading has been suspended.
	//
	// General listener format:
	//
	//  func(player rui.MediaPlayer)
	//
	// where:
	// player - Interface of a player which generated this event.
	//
	// Allowed listener formats:
	//
	//  func()
	SuspendEvent PropertyName = "suspend-event"

	// TimeUpdateEvent is the constant for "time-update-event" property tag.
	//
	// Used by AudioPlayer, VideoPlayer.
	//
	// Occur when the time indicated by the currentTime attribute has been updated.
	//
	// General listener format:
	//
	//  func(player rui.MediaPlayer, time float64).
	//
	// where:
	//   - player - Interface of a player which generated this event,
	//   - time - Current time.
	//
	// Allowed listener formats:
	//
	//  func(player rui.MediaPlayer),
	//  func(time float64),
	//  func()
	TimeUpdateEvent PropertyName = "time-update-event"

	// VolumeChangedEvent is the constant for "volume-changed-event" property tag.
	//
	// Used by AudioPlayer, VideoPlayer.
	//
	// Occur when the volume has changed.
	//
	// General listener format:
	//
	//  func(player rui.MediaPlayer, volume float64).
	//
	// where:
	//   - player - Interface of a player which generated this event,
	//   - volume - New volume level.
	//
	// Allowed listener formats:
	//
	//  func(player rui.MediaPlayer),
	//  func(volume float64),
	//  func()
	VolumeChangedEvent PropertyName = "volume-changed-event"

	// WaitingEvent is the constant for "waiting-event" property tag.
	//
	// Used by AudioPlayer, VideoPlayer.
	//
	// Occur when the playback has stopped because of a temporary lack of data.
	//
	// General listener format:
	//
	//  func(player rui.MediaPlayer)
	//
	// where:
	// player - Interface of a player which generated this event.
	//
	// Allowed listener formats:
	//
	//  func()
	WaitingEvent PropertyName = "waiting-event"

	// PlayerErrorEvent is the constant for "player-error-event" property tag.
	//
	// Used by AudioPlayer, VideoPlayer.
	//
	// Fired when the resource could not be loaded due to an error(for example, a network connectivity problem).
	//
	// General listener format:
	//
	//  func(player rui.MediaPlayer, code int, message string).
	//
	// where:
	//   - player - Interface of a player which generated this event,
	//   - code - Error code. See below,
	//   - message - Error message,
	//
	// Error codes:
	//   - 0 (PlayerErrorUnknown) - Unknown error,
	//   - 1 (PlayerErrorAborted) - Fetching the associated resource was interrupted by a user request,
	//   - 2 (PlayerErrorNetwork) - Some kind of network error has occurred that prevented the media from successfully ejecting, even though it was previously available,
	//   - 3 (PlayerErrorDecode) - Although the resource was previously identified as being used, an error occurred while trying to decode the media resource,
	//   - 4 (PlayerErrorSourceNotSupported) - The associated resource object or media provider was found to be invalid.
	//
	// Allowed listener formats:
	//
	//  func(code int, message string),
	//  func(player rui.MediaPlayer),
	//  func()
	PlayerErrorEvent PropertyName = "player-error-event"

	// PreloadNone - value of the view "preload" property: indicates that the audio/video should not be preloaded.
	PreloadNone = 0

	// PreloadMetadata - value of the view "preload" property: indicates that only audio/video metadata (e.g. length) is fetched.
	PreloadMetadata = 1

	// PreloadAuto - value of the view "preload" property: indicates that the whole audio file can be downloaded,
	// even if the user is not expected to use it.
	PreloadAuto = 2

	// PlayerErrorUnknown - MediaPlayer error code: An unknown error.
	PlayerErrorUnknown = 0

	// PlayerErrorAborted - MediaPlayer error code: The fetching of the associated resource was aborted by the user's request.
	PlayerErrorAborted = 1

	// PlayerErrorNetwork - MediaPlayer error code: Some kind of network error occurred which prevented the media
	// from being successfully fetched, despite having previously been available.
	PlayerErrorNetwork = 2

	// PlayerErrorDecode - MediaPlayer error code: Despite having previously been determined to be usable,
	// an error occurred while trying to decode the media resource, resulting in an error.
	PlayerErrorDecode = 3

	// PlayerErrorSourceNotSupported - MediaPlayer error code: The associated resource or media provider object has been found to be unsuitable.
	PlayerErrorSourceNotSupported = 4
)

// MediaPlayer is a common interface for media player views like [AudioPlayer] and [VideoPlayer].
type MediaPlayer interface {
	View

	// Play attempts to begin playback of the media.
	Play()

	// Pause will pause playback of the media, if the media is already in a paused state this method will have no effect.
	Pause()

	// SetCurrentTime sets the current playback time in seconds.
	SetCurrentTime(seconds float64)

	// CurrentTime returns the current playback time in seconds.
	CurrentTime() float64

	// Duration returns the value indicating the total duration of the media in seconds.
	// If no media data is available, the returned value is NaN.
	Duration() float64

	// SetPlaybackRate sets the rate at which the media is being played back. This is used to implement user controls
	// for fast forward, slow motion, and so forth. The normal playback rate is multiplied by this value to obtain
	// the current rate, so a value of 1.0 indicates normal speed.
	SetPlaybackRate(rate float64)

	// PlaybackRate returns the rate at which the media is being played back.
	PlaybackRate() float64

	// SetVolume sets the audio volume, from 0.0 (silent) to 1.0 (loudest).
	SetVolume(volume float64)

	// Volume returns the audio volume, from 0.0 (silent) to 1.0 (loudest).
	Volume() float64

	// IsEnded function tells whether the media element is ended.
	IsEnded() bool

	// IsPaused function tells whether the media element is paused.
	IsPaused() bool
}

type mediaPlayerData struct {
	viewData
}

// MediaSource represent one media file source
type MediaSource struct {
	// Url of the source
	Url string

	// MimeType of the source
	MimeType string
}

func (player *mediaPlayerData) init(session Session) {
	player.viewData.init(session)
	player.tag = "MediaPlayer"
	player.set = player.setFunc
	player.changed = player.propertyChanged
}

func (player *mediaPlayerData) Focusable() bool {
	return true
}

func (player *mediaPlayerData) setFunc(tag PropertyName, value any) []PropertyName {
	switch tag {

	case AbortEvent, CanPlayEvent, CanPlayThroughEvent, CompleteEvent, EmptiedEvent, LoadStartEvent,
		EndedEvent, LoadedDataEvent, LoadedMetadataEvent, PauseEvent, PlayEvent, PlayingEvent,
		ProgressEvent, SeekedEvent, SeekingEvent, StalledEvent, SuspendEvent, WaitingEvent:

		return setNoArgEventListener[MediaPlayer](player, tag, value)

	case DurationChangedEvent, RateChangedEvent, TimeUpdateEvent, VolumeChangedEvent:

		return setOneArgEventListener[MediaPlayer, float64](player, tag, value)

	case PlayerErrorEvent:
		if listeners, ok := valueToMediaPlayerErrorListeners(value); ok {
			return setArrayPropertyValue(player, tag, listeners)
		}
		notCompatibleType(tag, value)
		return nil

	case Source:
		return setMediaPlayerSource(player, value)
	}

	return player.viewData.setFunc(tag, value)
}

func setMediaPlayerSource(properties Properties, value any) []PropertyName {
	switch value := value.(type) {
	case string:
		src := MediaSource{Url: value, MimeType: ""}
		properties.setRaw(Source, []MediaSource{src})

	case MediaSource:
		properties.setRaw(Source, []MediaSource{value})

	case []MediaSource:
		properties.setRaw(Source, value)

	case DataObject:
		url, ok := value.PropertyValue("src")
		if !ok || url == "" {
			invalidPropertyValue(Source, value)
			return nil
		}

		mimeType, _ := value.PropertyValue("mime-type")
		src := MediaSource{Url: url, MimeType: mimeType}
		properties.setRaw(Source, []MediaSource{src})

	case []DataValue:
		src := []MediaSource{}
		for _, val := range value {
			if val.IsObject() {
				obj := val.Object()
				if url, ok := obj.PropertyValue("src"); ok && url != "" {
					mimeType, _ := obj.PropertyValue("mime-type")
					src = append(src, MediaSource{Url: url, MimeType: mimeType})
				} else {
					invalidPropertyValue(Source, value)
					return nil
				}
			} else {
				src = append(src, MediaSource{Url: val.Value(), MimeType: ""})
			}
		}

		if len(src) == 0 {
			invalidPropertyValue(Source, value)
			return nil
		}
		properties.setRaw(Source, src)

	default:
		notCompatibleType(Source, value)
		return nil
	}

	return []PropertyName{Source}
}

/*
	func valueToPlayerErrorListeners(value any) ([]func(MediaPlayer, int, string), bool) {
		if value == nil {
			return nil, true
		}

		switch value := value.(type) {
		case func(MediaPlayer, int, string):
			return []func(MediaPlayer, int, string){value}, true

		case func(int, string):
			fn := func(_ MediaPlayer, code int, message string) {
				value(code, message)
			}
			return []func(MediaPlayer, int, string){fn}, true

		case func(MediaPlayer):
			fn := func(player MediaPlayer, _ int, _ string) {
				value(player)
			}
			return []func(MediaPlayer, int, string){fn}, true

		case func():
			fn := func(MediaPlayer, int, string) {
				value()
			}
			return []func(MediaPlayer, int, string){fn}, true

		case []func(MediaPlayer, int, string):
			if len(value) == 0 {
				return nil, true
			}
			for _, fn := range value {
				if fn == nil {
					return nil, false
				}
			}
			return value, true

		case []func(int, string):
			count := len(value)
			if count == 0 {
				return nil, true
			}
			listeners := make([]func(MediaPlayer, int, string), count)
			for i, v := range value {
				if v == nil {
					return nil, false
				}
				listeners[i] = func(_ MediaPlayer, code int, message string) {
					v(code, message)
				}
			}
			return listeners, true

		case []func(MediaPlayer):
			count := len(value)
			if count == 0 {
				return nil, true
			}
			listeners := make([]func(MediaPlayer, int, string), count)
			for i, v := range value {
				if v == nil {
					return nil, false
				}
				listeners[i] = func(player MediaPlayer, _ int, _ string) {
					v(player)
				}
			}
			return listeners, true

		case []func():
			count := len(value)
			if count == 0 {
				return nil, true
			}
			listeners := make([]func(MediaPlayer, int, string), count)
			for i, v := range value {
				if v == nil {
					return nil, false
				}
				listeners[i] = func(MediaPlayer, int, string) {
					v()
				}
			}
			return listeners, true

		case []any:
			count := len(value)
			if count == 0 {
				return nil, true
			}
			listeners := make([]func(MediaPlayer, int, string), count)
			for i, v := range value {
				if v == nil {
					return nil, false
				}
				switch v := v.(type) {
				case func(MediaPlayer, int, string):
					listeners[i] = v

				case func(int, string):
					listeners[i] = func(_ MediaPlayer, code int, message string) {
						v(code, message)
					}

				case func(MediaPlayer):
					listeners[i] = func(player MediaPlayer, _ int, _ string) {
						v(player)
					}

				case func():
					listeners[i] = func(MediaPlayer, int, string) {
						v()
					}

				default:
					return nil, false
				}
			}
			return listeners, true
		}

		return nil, false
	}
*/
func mediaPlayerEvents() map[PropertyName]string {
	return map[PropertyName]string{
		AbortEvent:          "onabort",
		CanPlayEvent:        "oncanplay",
		CanPlayThroughEvent: "oncanplaythrough",
		CompleteEvent:       "oncomplete",
		EmptiedEvent:        "onemptied",
		EndedEvent:          "ended",
		LoadedDataEvent:     "onloadeddata",
		LoadedMetadataEvent: "onloadedmetadata",
		LoadStartEvent:      "onloadstart",
		PauseEvent:          "onpause",
		PlayEvent:           "onplay",
		PlayingEvent:        "onplaying",
		ProgressEvent:       "onprogress",
		SeekedEvent:         "onseeked",
		SeekingEvent:        "onseeking",
		StalledEvent:        "onstalled",
		SuspendEvent:        "onsuspend",
		WaitingEvent:        "onwaiting",
	}
}

func (player *mediaPlayerData) propertyChanged(tag PropertyName) {
	session := player.Session()

	switch tag {
	case Controls, Loop:
		value, _ := boolProperty(player, tag, session)
		if value {
			session.updateProperty(player.htmlID(), string(tag), value)
		} else {
			session.removeProperty(player.htmlID(), string(tag))
		}

	case Muted:
		value, _ := boolProperty(player, Muted, session)
		session.callFunc("setMediaMuted", player.htmlID(), value)

	case Preload:
		value, _ := enumProperty(player, Preload, session, 0)
		values := enumProperties[Preload].values
		session.updateProperty(player.htmlID(), string(Preload), values[value])

	case AbortEvent, CanPlayEvent, CanPlayThroughEvent, CompleteEvent, EmptiedEvent,
		EndedEvent, LoadedDataEvent, LoadedMetadataEvent, PauseEvent, PlayEvent, PlayingEvent, ProgressEvent,
		LoadStartEvent, SeekedEvent, SeekingEvent, StalledEvent, SuspendEvent, WaitingEvent:

		if cssTag, ok := mediaPlayerEvents()[tag]; ok {
			fn := ""
			if value := player.getRaw(tag); value != nil {
				if listeners, ok := value.([]func(MediaPlayer)); ok && len(listeners) > 0 {
					fn = fmt.Sprintf(`viewEvent(this, "%s")`, string(tag))
				}
			}
			session.updateProperty(player.htmlID(), cssTag, fn)
		}

	case TimeUpdateEvent:
		if value := player.getRaw(tag); value != nil {
			session.updateProperty(player.htmlID(), "ontimeupdate", "viewTimeUpdatedEvent(this)")
		} else {
			session.updateProperty(player.htmlID(), "ontimeupdate", "")
		}

	case VolumeChangedEvent:
		if value := player.getRaw(tag); value != nil {
			session.updateProperty(player.htmlID(), "onvolumechange", "viewVolumeChangedEvent(this)")
		} else {
			session.updateProperty(player.htmlID(), "onvolumechange", "")
		}

	case DurationChangedEvent:
		if value := player.getRaw(tag); value != nil {
			session.updateProperty(player.htmlID(), "ondurationchange", "viewDurationChangedEvent(this)")
		} else {
			session.updateProperty(player.htmlID(), "ondurationchange", "")
		}

	case RateChangedEvent:
		if value := player.getRaw(tag); value != nil {
			session.updateProperty(player.htmlID(), "onratechange", "viewRateChangedEvent(this)")
		} else {
			session.updateProperty(player.htmlID(), "onratechange", "")
		}

	case PlayerErrorEvent:
		if value := player.getRaw(tag); value != nil {
			session.updateProperty(player.htmlID(), "onerror", "viewErrorEvent(this)")
		} else {
			session.updateProperty(player.htmlID(), "onerror", "")
		}

	case Source:
		updateInnerHTML(player.htmlID(), session)

	default:
		player.viewData.propertyChanged(tag)
	}

}

func (player *mediaPlayerData) htmlSubviews(self View, buffer *strings.Builder) {
	if value := player.getRaw(Source); value != nil {
		if sources, ok := value.([]MediaSource); ok && len(sources) > 0 {
			session := player.session
			for _, src := range sources {
				if url, ok := session.resolveConstants(src.Url); ok && url != "" {
					buffer.WriteString(`<source src="`)
					buffer.WriteString(url)
					buffer.WriteRune('"')
					if mime, ok := session.resolveConstants(src.MimeType); ok && mime != "" {
						buffer.WriteString(` type="`)
						buffer.WriteString(mime)
						buffer.WriteRune('"')
					}
					buffer.WriteRune('>')
				}
			}
		}
	}
}

func (player *mediaPlayerData) htmlProperties(self View, buffer *strings.Builder) {
	player.viewData.htmlProperties(self, buffer)
	for _, tag := range []PropertyName{Controls, Loop, Muted, Preload} {
		if value, _ := boolProperty(player, tag, player.session); value {
			buffer.WriteRune(' ')
			buffer.WriteString(string(tag))
		}
	}

	if value, ok := enumProperty(player, Preload, player.session, 0); ok {
		values := enumProperties[Preload].values
		buffer.WriteString(` preload="`)
		buffer.WriteString(values[value])
		buffer.WriteRune('"')
	}

	for tag, cssTag := range mediaPlayerEvents() {
		if value := player.getRaw(tag); value != nil {
			if listeners, ok := value.([]func(MediaPlayer)); ok && len(listeners) > 0 {
				buffer.WriteString(` `)
				buffer.WriteString(cssTag)
				buffer.WriteString(`="playerEvent(this, '`)
				buffer.WriteString(string(tag))
				buffer.WriteString(`')"`)
			}
		}
	}

	if value := player.getRaw(TimeUpdateEvent); value != nil {
		buffer.WriteString(` ontimeupdate="playerTimeUpdatedEvent(this)"`)
	}

	if value := player.getRaw(VolumeChangedEvent); value != nil {
		buffer.WriteString(` onvolumechange="playerVolumeChangedEvent(this)"`)
	}

	if value := player.getRaw(DurationChangedEvent); value != nil {
		buffer.WriteString(` ondurationchange="playerDurationChangedEvent(this)"`)
	}

	if value := player.getRaw(RateChangedEvent); value != nil {
		buffer.WriteString(` onratechange="playerRateChangedEvent(this)"`)
	}

	if value := player.getRaw(PlayerErrorEvent); value != nil {
		buffer.WriteString(` onerror="playerErrorEvent(this)"`)
	}
}

func (player *mediaPlayerData) handleCommand(self View, command PropertyName, data DataObject) bool {
	switch command {
	case AbortEvent, CanPlayEvent, CanPlayThroughEvent, CompleteEvent, LoadStartEvent,
		EmptiedEvent, EndedEvent, LoadedDataEvent, LoadedMetadataEvent, PauseEvent, PlayEvent,
		PlayingEvent, ProgressEvent, SeekedEvent, SeekingEvent, StalledEvent, SuspendEvent,
		WaitingEvent:

		for _, listener := range getNoArgEventListeners[MediaPlayer](player, nil, command) {
			listener.Run(player)
		}

	case TimeUpdateEvent, DurationChangedEvent, RateChangedEvent, VolumeChangedEvent:
		time := dataFloatProperty(data, "value")
		for _, listener := range getOneArgEventListeners[MediaPlayer, float64](player, nil, command) {
			listener.Run(player, time)
		}

	case PlayerErrorEvent:
		if value := player.getRaw(command); value != nil {
			if listeners, ok := value.([]mediaPlayerErrorListener); ok {
				code, _ := dataIntProperty(data, "code")
				message, _ := data.PropertyValue("message")
				for _, listener := range listeners {
					listener.Run(player, code, message)
				}
			}
		}
	}

	return player.viewData.handleCommand(self, command, data)
}

func (player *mediaPlayerData) Play() {
	player.session.callFunc("mediaPlay", player.htmlID())
}

func (player *mediaPlayerData) Pause() {
	player.session.callFunc("mediaPause", player.htmlID())
}

func (player *mediaPlayerData) SetCurrentTime(seconds float64) {
	player.session.callFunc("mediaSetSetCurrentTime", player.htmlID(), seconds)
}

func (player *mediaPlayerData) getFloatPlayerProperty(tag string) (float64, bool) {
	value := player.session.htmlPropertyValue(player.htmlID(), tag)
	if value != "" {
		result, err := strconv.ParseFloat(value, 32)
		if err == nil {
			return result, true
		}
		ErrorLog(err.Error())
	}

	return 0, false
}

func (player *mediaPlayerData) CurrentTime() float64 {
	if result, ok := player.getFloatPlayerProperty("currentTime"); ok {
		return result
	}
	return 0
}

func (player *mediaPlayerData) Duration() float64 {
	if result, ok := player.getFloatPlayerProperty("duration"); ok {
		return result
	}
	return 0
}

func (player *mediaPlayerData) SetPlaybackRate(rate float64) {
	player.session.callFunc("mediaSetPlaybackRate", player.htmlID(), rate)
}

func (player *mediaPlayerData) PlaybackRate() float64 {
	if result, ok := player.getFloatPlayerProperty("playbackRate"); ok {
		return result
	}
	return 1
}

func (player *mediaPlayerData) SetVolume(volume float64) {
	if volume >= 0 && volume <= 1 {
		player.session.callFunc("mediaSetVolume", player.htmlID(), volume)
	}
}

func (player *mediaPlayerData) Volume() float64 {
	if result, ok := player.getFloatPlayerProperty("volume"); ok {
		return result
	}
	return 1
}

func (player *mediaPlayerData) getBoolPlayerProperty(tag string) (bool, bool) {
	switch value := player.session.htmlPropertyValue(player.htmlID(), tag); strings.ToLower(value) {
	case "0", "false", "off":
		return false, true

	case "1", "true", "on":
		return false, true
	}

	return false, false
}

func (player *mediaPlayerData) IsEnded() bool {
	if result, ok := player.getBoolPlayerProperty("ended"); ok {
		return result
	}
	return false
}

func (player *mediaPlayerData) IsPaused() bool {
	if result, ok := player.getBoolPlayerProperty("paused"); ok {
		return result
	}
	return false
}

// MediaPlayerPlay attempts to begin playback of the media.
func MediaPlayerPlay(view View, playerID string) {
	if playerID != "" {
		view = ViewByID(view, playerID)
	}

	if player, ok := view.(MediaPlayer); ok {
		player.Play()
	} else {
		ErrorLog(`The found View is not MediaPlayer`)
	}
}

// MediaPlayerPause will pause playback of the media, if the media is already in a paused state this method will have no effect.
func MediaPlayerPause(view View, playerID string) {
	if playerID != "" {
		view = ViewByID(view, playerID)
	}

	if player, ok := view.(MediaPlayer); ok {
		player.Pause()
	} else {
		ErrorLog(`The found View is not MediaPlayer`)
	}
}

// SetMediaPlayerCurrentTime sets the current playback time in seconds.
func SetMediaPlayerCurrentTime(view View, playerID string, seconds float64) {
	if playerID != "" {
		view = ViewByID(view, playerID)
	}

	if player, ok := view.(MediaPlayer); ok {
		player.SetCurrentTime(seconds)
	} else {
		ErrorLog(`The found View is not MediaPlayer`)
	}
}

// MediaPlayerCurrentTime returns the current playback time in seconds.
func MediaPlayerCurrentTime(view View, playerID string) float64 {
	if playerID != "" {
		view = ViewByID(view, playerID)
	}

	if player, ok := view.(MediaPlayer); ok {
		return player.CurrentTime()
	}

	ErrorLog(`The found View is not MediaPlayer`)
	return 0
}

// MediaPlayerDuration returns the value indicating the total duration of the media in seconds.
// If no media data is available, the returned value is NaN.
func MediaPlayerDuration(view View, playerID string) float64 {
	if playerID != "" {
		view = ViewByID(view, playerID)
	}

	if player, ok := view.(MediaPlayer); ok {
		return player.Duration()
	}

	ErrorLog(`The found View is not MediaPlayer`)
	return math.NaN()
}

// SetVolume sets the audio volume, from 0.0 (silent) to 1.0 (loudest).
func SetMediaPlayerVolume(view View, playerID string, volume float64) {
	if playerID != "" {
		view = ViewByID(view, playerID)
	}

	if player, ok := view.(MediaPlayer); ok {
		player.SetVolume(volume)
	} else {
		ErrorLog(`The found View is not MediaPlayer`)
	}
}

// Volume returns the audio volume, from 0.0 (silent) to 1.0 (loudest).
func MediaPlayerVolume(view View, playerID string) float64 {
	if playerID != "" {
		view = ViewByID(view, playerID)
	}

	if player, ok := view.(MediaPlayer); ok {
		return player.Volume()
	}

	ErrorLog(`The found View is not MediaPlayer`)
	return 1
}

// SetMediaPlayerPlaybackRate sets the rate at which the media is being played back. This is used to implement user controls
// for fast forward, slow motion, and so forth. The normal playback rate is multiplied by this value to obtain
// the current rate, so a value of 1.0 indicates normal speed.
func SetMediaPlayerPlaybackRate(view View, playerID string, rate float64) {
	if playerID != "" {
		view = ViewByID(view, playerID)
	}

	if player, ok := view.(MediaPlayer); ok {
		player.SetPlaybackRate(rate)
	} else {
		ErrorLog(`The found View is not MediaPlayer`)
	}
}

// MediaPlayerPlaybackRate returns the rate at which the media is being played back.
func MediaPlayerPlaybackRate(view View, playerID string) float64 {
	if playerID != "" {
		view = ViewByID(view, playerID)
	}

	if player, ok := view.(MediaPlayer); ok {
		return player.PlaybackRate()
	}

	ErrorLog(`The found View is not MediaPlayer`)
	return 1
}

// IsMediaPlayerEnded function tells whether the media element is ended.
func IsMediaPlayerEnded(view View, playerID string) bool {
	if playerID != "" {
		view = ViewByID(view, playerID)
	}

	if player, ok := view.(MediaPlayer); ok {
		return player.IsEnded()
	}

	ErrorLog(`The found View is not MediaPlayer`)
	return false
}

// IsMediaPlayerPaused function tells whether the media element is paused.
func IsMediaPlayerPaused(view View, playerID string) bool {
	if playerID != "" {
		view = ViewByID(view, playerID)
	}

	if player, ok := view.(MediaPlayer); ok {
		return player.IsPaused()
	}

	ErrorLog(`The found View is not MediaPlayer`)
	return false
}

type mediaPlayerErrorListener interface {
	Run(MediaPlayer, int, string)
	rawListener() any
}

type mediaPlayerErrorListener0 struct {
	fn func()
}

type mediaPlayerErrorListenerP struct {
	fn func(MediaPlayer)
}

type mediaPlayerErrorListenerI struct {
	fn func(int)
}

type mediaPlayerErrorListenerS struct {
	fn func(string)
}

type mediaPlayerErrorListenerPI struct {
	fn func(MediaPlayer, int)
}

type mediaPlayerErrorListenerPS struct {
	fn func(MediaPlayer, string)
}

type mediaPlayerErrorListenerIS struct {
	fn func(int, string)
}

type mediaPlayerErrorListenerPIS struct {
	fn func(MediaPlayer, int, string)
}

type mediaPlayerErrorListenerBinding struct {
	name string
}

func newMediaPlayerErrorListener0(fn func()) mediaPlayerErrorListener {
	obj := new(mediaPlayerErrorListener0)
	obj.fn = fn
	return obj
}

func (data *mediaPlayerErrorListener0) Run(_ MediaPlayer, _ int, _ string) {
	data.fn()
}

func (data *mediaPlayerErrorListener0) rawListener() any {
	return data.fn
}

func newMediaPlayerErrorListenerP(fn func(MediaPlayer)) mediaPlayerErrorListener {
	obj := new(mediaPlayerErrorListenerP)
	obj.fn = fn
	return obj
}

func (data *mediaPlayerErrorListenerP) Run(player MediaPlayer, _ int, _ string) {
	data.fn(player)
}

func (data *mediaPlayerErrorListenerP) rawListener() any {
	return data.fn
}

func newMediaPlayerErrorListenerI(fn func(int)) mediaPlayerErrorListener {
	obj := new(mediaPlayerErrorListenerI)
	obj.fn = fn
	return obj
}

func (data *mediaPlayerErrorListenerI) Run(_ MediaPlayer, code int, _ string) {
	data.fn(code)
}

func (data *mediaPlayerErrorListenerI) rawListener() any {
	return data.fn
}

func newMediaPlayerErrorListenerS(fn func(string)) mediaPlayerErrorListener {
	obj := new(mediaPlayerErrorListenerS)
	obj.fn = fn
	return obj
}

func (data *mediaPlayerErrorListenerS) Run(_ MediaPlayer, _ int, message string) {
	data.fn(message)
}

func (data *mediaPlayerErrorListenerS) rawListener() any {
	return data.fn
}

func newMediaPlayerErrorListenerPI(fn func(MediaPlayer, int)) mediaPlayerErrorListener {
	obj := new(mediaPlayerErrorListenerPI)
	obj.fn = fn
	return obj
}

func (data *mediaPlayerErrorListenerPI) Run(player MediaPlayer, code int, _ string) {
	data.fn(player, code)
}

func (data *mediaPlayerErrorListenerPI) rawListener() any {
	return data.fn
}

func newMediaPlayerErrorListenerPS(fn func(MediaPlayer, string)) mediaPlayerErrorListener {
	obj := new(mediaPlayerErrorListenerPS)
	obj.fn = fn
	return obj
}

func (data *mediaPlayerErrorListenerPS) Run(player MediaPlayer, _ int, message string) {
	data.fn(player, message)
}

func (data *mediaPlayerErrorListenerPS) rawListener() any {
	return data.fn
}

func newMediaPlayerErrorListenerIS(fn func(int, string)) mediaPlayerErrorListener {
	obj := new(mediaPlayerErrorListenerIS)
	obj.fn = fn
	return obj
}

func (data *mediaPlayerErrorListenerIS) Run(_ MediaPlayer, code int, message string) {
	data.fn(code, message)
}

func (data *mediaPlayerErrorListenerIS) rawListener() any {
	return data.fn
}

func newMediaPlayerErrorListenerPIS(fn func(MediaPlayer, int, string)) mediaPlayerErrorListener {
	obj := new(mediaPlayerErrorListenerPIS)
	obj.fn = fn
	return obj
}

func (data *mediaPlayerErrorListenerPIS) Run(player MediaPlayer, code int, message string) {
	data.fn(player, code, message)
}

func (data *mediaPlayerErrorListenerPIS) rawListener() any {
	return data.fn
}

func newMediaPlayerErrorListenerBinding(name string) mediaPlayerErrorListener {
	obj := new(mediaPlayerErrorListenerBinding)
	obj.name = name
	return obj
}

func (data *mediaPlayerErrorListenerBinding) Run(player MediaPlayer, code int, message string) {
	bind := player.binding()
	if bind == nil {
		ErrorLogF(`There is no a binding object for call "%s"`, data.name)
		return
	}

	val := reflect.ValueOf(bind)
	method := val.MethodByName(data.name)
	if !method.IsValid() {
		ErrorLogF(`The "%s" method is not valid`, data.name)
		return
	}

	methodType := method.Type()
	var args []reflect.Value = nil

	switch methodType.NumIn() {
	case 0:
		args = []reflect.Value{}

	case 1:
		switch methodType.In(0) {
		case reflect.TypeOf(player):
			args = []reflect.Value{reflect.ValueOf(player)}
		case reflect.TypeOf(code):
			args = []reflect.Value{reflect.ValueOf(code)}
		case reflect.TypeOf(message):
			args = []reflect.Value{reflect.ValueOf(message)}
		}

	case 2:
		in0 := methodType.In(0)
		in1 := methodType.In(1)
		if in0 == reflect.TypeOf(player) {
			if in1 == reflect.TypeOf(code) {
				args = []reflect.Value{reflect.ValueOf(player), reflect.ValueOf(code)}
			} else if in1 == reflect.TypeOf(message) {
				args = []reflect.Value{reflect.ValueOf(player), reflect.ValueOf(message)}
			}
		} else if in0 == reflect.TypeOf(code) && in1 == reflect.TypeOf(message) {
			args = []reflect.Value{reflect.ValueOf(code), reflect.ValueOf(message)}
		}

	case 3:
		if methodType.In(0) == reflect.TypeOf(player) &&
			methodType.In(1) == reflect.TypeOf(code) &&
			methodType.In(2) == reflect.TypeOf(message) {
			args = []reflect.Value{
				reflect.ValueOf(player),
				reflect.ValueOf(code),
				reflect.ValueOf(message),
			}
		}
	}

	if args != nil {
		method.Call(args)
	} else {
		ErrorLogF(`Unsupported prototype of "%s" method`, data.name)
	}
}

func (data *mediaPlayerErrorListenerBinding) rawListener() any {
	return data.name
}

func valueToMediaPlayerErrorListeners(value any) ([]mediaPlayerErrorListener, bool) {
	if value == nil {
		return nil, true
	}

	switch value := value.(type) {
	case []mediaPlayerErrorListener:
		return value, true

	case mediaPlayerErrorListener:
		return []mediaPlayerErrorListener{value}, true

	case string:
		return []mediaPlayerErrorListener{newMediaPlayerErrorListenerBinding(value)}, true

	case func():
		return []mediaPlayerErrorListener{newMediaPlayerErrorListener0(value)}, true

	case func(MediaPlayer):
		return []mediaPlayerErrorListener{newMediaPlayerErrorListenerP(value)}, true

	case func(int):
		return []mediaPlayerErrorListener{newMediaPlayerErrorListenerI(value)}, true

	case func(string):
		return []mediaPlayerErrorListener{newMediaPlayerErrorListenerS(value)}, true

	case func(MediaPlayer, int):
		return []mediaPlayerErrorListener{newMediaPlayerErrorListenerPI(value)}, true

	case func(MediaPlayer, string):
		return []mediaPlayerErrorListener{newMediaPlayerErrorListenerPS(value)}, true

	case func(int, string):
		return []mediaPlayerErrorListener{newMediaPlayerErrorListenerIS(value)}, true

	case func(MediaPlayer, int, string):
		return []mediaPlayerErrorListener{newMediaPlayerErrorListenerPIS(value)}, true

	case []func():
		result := make([]mediaPlayerErrorListener, 0, len(value))
		for _, fn := range value {
			if fn != nil {
				result = append(result, newMediaPlayerErrorListener0(fn))
			}
		}
		return result, len(result) > 0

	case []func(MediaPlayer):
		result := make([]mediaPlayerErrorListener, 0, len(value))
		for _, fn := range value {
			if fn != nil {
				result = append(result, newMediaPlayerErrorListenerP(fn))
			}
		}
		return result, len(result) > 0

	case []func(int):
		result := make([]mediaPlayerErrorListener, 0, len(value))
		for _, fn := range value {
			if fn != nil {
				result = append(result, newMediaPlayerErrorListenerI(fn))
			}
		}
		return result, len(result) > 0

	case []func(string):
		result := make([]mediaPlayerErrorListener, 0, len(value))
		for _, fn := range value {
			if fn != nil {
				result = append(result, newMediaPlayerErrorListenerS(fn))
			}
		}
		return result, len(result) > 0

	case []func(MediaPlayer, int):
		result := make([]mediaPlayerErrorListener, 0, len(value))
		for _, fn := range value {
			if fn != nil {
				result = append(result, newMediaPlayerErrorListenerPI(fn))
			}
		}
		return result, len(result) > 0

	case []func(MediaPlayer, string):
		result := make([]mediaPlayerErrorListener, 0, len(value))
		for _, fn := range value {
			if fn != nil {
				result = append(result, newMediaPlayerErrorListenerPS(fn))
			}
		}
		return result, len(result) > 0

	case []func(int, string):
		result := make([]mediaPlayerErrorListener, 0, len(value))
		for _, fn := range value {
			if fn != nil {
				result = append(result, newMediaPlayerErrorListenerIS(fn))
			}
		}
		return result, len(result) > 0

	case []func(MediaPlayer, int, string):
		result := make([]mediaPlayerErrorListener, 0, len(value))
		for _, fn := range value {
			if fn != nil {
				result = append(result, newMediaPlayerErrorListenerPIS(fn))
			}
		}
		return result, len(result) > 0

	case []any:
		result := make([]mediaPlayerErrorListener, 0, len(value))
		for _, v := range value {
			if v != nil {
				switch v := v.(type) {
				case func():
					result = append(result, newMediaPlayerErrorListener0(v))

				case func(MediaPlayer):
					result = append(result, newMediaPlayerErrorListenerP(v))

				case func(int):
					result = append(result, newMediaPlayerErrorListenerI(v))

				case func(string):
					result = append(result, newMediaPlayerErrorListenerS(v))

				case func(MediaPlayer, int):
					result = append(result, newMediaPlayerErrorListenerPI(v))

				case func(MediaPlayer, string):
					result = append(result, newMediaPlayerErrorListenerPS(v))

				case func(int, string):
					result = append(result, newMediaPlayerErrorListenerIS(v))

				case func(MediaPlayer, int, string):
					result = append(result, newMediaPlayerErrorListenerPIS(v))

				case string:
					result = append(result, newMediaPlayerErrorListenerBinding(v))

				default:
					return nil, false
				}
			}
		}
		return result, len(result) > 0
	}

	return nil, false
}

func getMediaPlayerErrorListenerBinding(listeners []mediaPlayerErrorListener) string {
	for _, listener := range listeners {
		raw := listener.rawListener()
		if text, ok := raw.(string); ok && text != "" {
			return text
		}
	}
	return ""
}
