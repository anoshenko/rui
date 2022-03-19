package rui

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

const (
	// Controls is the constant for the "autoplay" controls tag.
	// If the "controls" bool property is "true", the browser will offer controls to allow the user
	// to control audio/video playback, including volume, seeking, and pause/resume playback.
	// Its default value is false.
	Controls = "controls"
	// Loop is the constant for the "loop" property tag.
	// If the "loop" bool property is "true", the audio/video player will automatically seek back
	// to the start upon reaching the end of the audio/video.
	// Its default value is false.
	Loop = "loop"
	// Muted is the constant for the "muted" property tag.
	// The "muted" bool property indicates whether the audio/video will be initially silenced.
	// Its default value is false.
	Muted = "muted"
	// Preload is the constant for the "preload" property tag.
	// The "preload" int property is intended to provide a hint to the browser about what
	// the author thinks will lead to the best user experience. It may have one of the following values:
	// PreloadNone (0), PreloadMetadata (1), and PreloadAuto (2)
	// The default value is different for each browser.
	Preload = "preload"

	// AbortEvent is the constant for the "abort-event" property tag.
	// The "abort-event" event fired when the resource was not fully loaded, but not as the result of an error.
	AbortEvent = "abort-event"
	// CanPlayEvent is the constant for the "can-play-event" property tag.
	// The "can-play-event" event occurs when the browser can play the media, but estimates that not enough data has been
	// loaded to play the media up to its end without having to stop for further buffering of content.
	CanPlayEvent = "can-play-event"
	// CanPlayThroughEvent is the constant for the "can-play-through-event" property tag.
	// The "can-play-through-event" event occurs when the browser estimates it can play the media up
	// to its end without stopping for content buffering.
	CanPlayThroughEvent = "can-play-through-event"
	// CompleteEvent is the constant for the "complete-event" property tag.
	// The "complete-event" event occurs when the rendering of an OfflineAudioContext is terminated.
	CompleteEvent = "complete-event"
	// DurationChangedEvent is the constant for the "duration-changed-event" property tag.
	// The "duration-changed-event" event occurs when the duration attribute has been updated.
	DurationChangedEvent = "duration-changed-event"
	// EmptiedEvent is the constant for the "emptied-event" property tag.
	// The "emptied-event" event occurs when the media has become empty; for example, this event is sent if the media has already been loaded
	// (or partially loaded), and the HTMLMediaElement.load method is called to reload it.
	EmptiedEvent = "emptied-event"
	// EndedEvent is the constant for the "ended-event" property tag.
	// The "ended-event" event occurs when the playback has stopped because the end of the media was reached.
	EndedEvent = "ended-event"
	// LoadedDataEvent is the constant for the "loaded-data-event" property tag.
	// The "loaded-data-event" event occurs when the first frame of the media has finished loading.
	LoadedDataEvent = "loaded-data-event"
	// LoadedMetadataEvent is the constant for the "loaded-metadata-event" property tag.
	// The "loaded-metadata-event" event occurs when the metadata has been loaded.
	LoadedMetadataEvent = "loaded-metadata-event"
	// LoadStartEvent is the constant for the "load-start-event" property tag.
	// The "load-start-event" event is fired when the browser has started to load a resource.
	LoadStartEvent = "load-start-event"
	// PauseEvent is the constant for the "pause-event" property tag.
	// The "pause-event" event occurs when the playback has been paused.
	PauseEvent = "pause-event"
	// PlayEvent is the constant for the "play-event" property tag.
	// The "play-event" event occurs when the playback has begun.
	PlayEvent = "play-event"
	// PlayingEvent is the constant for the "playing-event" property tag.
	// The "playing-event" event occurs when the playback is ready to start after having been paused or delayed due to lack of data.
	PlayingEvent = "playing-event"
	// ProgressEvent is the constant for the "progress-event" property tag.
	// The "progress-event" event is fired periodically as the browser loads a resource.
	ProgressEvent = "progress-event"
	// RateChangeEvent is the constant for the "rate-change-event" property tag.
	// The "rate-change-event" event occurs when the playback rate has changed.
	RateChangedEvent = "rate-changed-event"
	// SeekedEvent is the constant for the "seeked-event" property tag.
	// The "seeked-event" event occurs when a seek operation completed.
	SeekedEvent = "seeked-event"
	// SeekingEvent is the constant for the "seeking-event" property tag.
	// The "seeking-event" event occurs when a seek operation began.
	SeekingEvent = "seeking-event"
	// StalledEvent is the constant for the "stalled-event" property tag.
	// The "stalled-event" event occurs when the user agent is trying to fetch media data, but data is unexpectedly not forthcoming.
	StalledEvent = "stalled-event"
	// SuspendEvent is the constant for the "suspend-event" property tag.
	// The "suspend-event" event occurs when the media data loading has been suspended.
	SuspendEvent = "suspend-event"
	// TimeUpdateEvent is the constant for the "time-update-event" property tag.
	// The "time-update-event" event occurs when the time indicated by the currentTime attribute has been updated.
	TimeUpdateEvent = "time-update-event"
	// VolumeChangedEvent is the constant for the "volume-change-event" property tag.
	// The "volume-change-event" event occurs when the volume has changed.
	VolumeChangedEvent = "volume-changed-event"
	// WaitingEvent is the constant for the "waiting-event" property tag.
	// The "waiting-event" event occurs when the playback has stopped because of a temporary lack of data
	WaitingEvent = "waiting-event"
	// PlayerErrorEvent is the constant for the "player-error-event" property tag.
	// The "player-error-event" event is fired when the resource could not be loaded due to an error
	// (for example, a network connectivity problem).
	PlayerErrorEvent = "player-error-event"

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

type MediaSource struct {
	Url      string
	MimeType string
}

func (player *mediaPlayerData) Init(session Session) {
	player.viewData.Init(session)
	player.tag = "MediaPlayer"
}

func (player *mediaPlayerData) Focusable() bool {
	return true
}

func (player *mediaPlayerData) Remove(tag string) {
	player.remove(strings.ToLower(tag))
}

func (player *mediaPlayerData) remove(tag string) {
	player.viewData.remove(tag)
	player.propertyChanged(tag)
}

func (player *mediaPlayerData) Set(tag string, value interface{}) bool {
	return player.set(strings.ToLower(tag), value)
}

func (player *mediaPlayerData) set(tag string, value interface{}) bool {
	if value == nil {
		player.remove(tag)
		return true
	}

	switch tag {
	case Controls, Loop, Muted, Preload:
		if player.viewData.set(tag, value) {
			player.propertyChanged(tag)
			return true
		}

	case AbortEvent, CanPlayEvent, CanPlayThroughEvent, CompleteEvent, EmptiedEvent, LoadStartEvent,
		EndedEvent, LoadedDataEvent, LoadedMetadataEvent, PauseEvent, PlayEvent, PlayingEvent,
		ProgressEvent, SeekedEvent, SeekingEvent, StalledEvent, SuspendEvent, WaitingEvent:
		if listeners, ok := valueToPlayerListeners(value); ok {
			if listeners == nil {
				delete(player.properties, tag)
			} else {
				player.properties[tag] = listeners
			}
			player.propertyChanged(tag)
			player.propertyChangedEvent(tag)
			return true
		}
		notCompatibleType(tag, value)

	case DurationChangedEvent, RateChangedEvent, TimeUpdateEvent, VolumeChangedEvent:
		if listeners, ok := valueToPlayerTimeListeners(value); ok {
			if listeners == nil {
				delete(player.properties, tag)
			} else {
				player.properties[tag] = listeners
			}
			player.propertyChanged(tag)
			player.propertyChangedEvent(tag)
			return true
		}
		notCompatibleType(tag, value)

	case PlayerErrorEvent:
		if listeners, ok := valueToPlayerErrorListeners(value); ok {
			if listeners == nil {
				delete(player.properties, tag)
			} else {
				player.properties[tag] = listeners
			}
			player.propertyChanged(tag)
			player.propertyChangedEvent(tag)
			return true
		}
		notCompatibleType(tag, value)

	case Source:
		if player.setSource(value) {
			player.propertyChanged(tag)
			player.propertyChangedEvent(tag)
			return true
		}

	default:
		return player.viewData.set(tag, value)
	}

	return false
}

func (player *mediaPlayerData) setSource(value interface{}) bool {
	switch value := value.(type) {
	case string:
		src := MediaSource{Url: value, MimeType: ""}
		player.properties[Source] = []MediaSource{src}

	case MediaSource:
		player.properties[Source] = []MediaSource{value}

	case []MediaSource:
		player.properties[Source] = value

	case DataObject:
		url, ok := value.PropertyValue("src")
		if !ok || url == "" {
			invalidPropertyValue(Source, value)
			return false
		}

		mimeType, _ := value.PropertyValue("mime-type")
		src := MediaSource{Url: url, MimeType: mimeType}
		player.properties[Source] = []MediaSource{src}

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
					return false
				}
			} else {
				src = append(src, MediaSource{Url: val.Value(), MimeType: ""})
			}
		}

		if len(src) == 0 {
			invalidPropertyValue(Source, value)
			return false
		}
		player.properties[Source] = src

	default:
		notCompatibleType(Source, value)
		return false
	}

	return true
}

func valueToPlayerListeners(value interface{}) ([]func(MediaPlayer), bool) {
	if value == nil {
		return nil, true
	}

	switch value := value.(type) {
	case func(MediaPlayer):
		return []func(MediaPlayer){value}, true

	case func():
		fn := func(MediaPlayer) {
			value()
		}
		return []func(MediaPlayer){fn}, true

	case []func(MediaPlayer):
		if len(value) == 0 {
			return nil, true
		}
		for _, fn := range value {
			if fn == nil {
				return nil, false
			}
		}
		return value, true

	case []func():
		count := len(value)
		if count == 0 {
			return nil, true
		}
		listeners := make([]func(MediaPlayer), count)
		for i, v := range value {
			if v == nil {
				return nil, false
			}
			listeners[i] = func(MediaPlayer) {
				v()
			}
		}
		return listeners, true

	case []interface{}:
		count := len(value)
		if count == 0 {
			return nil, true
		}
		listeners := make([]func(MediaPlayer), count)
		for i, v := range value {
			if v == nil {
				return nil, false
			}
			switch v := v.(type) {
			case func(MediaPlayer):
				listeners[i] = v

			case func():
				listeners[i] = func(MediaPlayer) {
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

func valueToPlayerTimeListeners(value interface{}) ([]func(MediaPlayer, float64), bool) {
	if value == nil {
		return nil, true
	}

	switch value := value.(type) {
	case func(MediaPlayer, float64):
		return []func(MediaPlayer, float64){value}, true

	case func(float64):
		fn := func(player MediaPlayer, time float64) {
			value(time)
		}
		return []func(MediaPlayer, float64){fn}, true

	case func(MediaPlayer):
		fn := func(player MediaPlayer, time float64) {
			value(player)
		}
		return []func(MediaPlayer, float64){fn}, true

	case func():
		fn := func(player MediaPlayer, time float64) {
			value()
		}
		return []func(MediaPlayer, float64){fn}, true

	case []func(MediaPlayer, float64):
		if len(value) == 0 {
			return nil, true
		}
		for _, fn := range value {
			if fn == nil {
				return nil, false
			}
		}
		return value, true

	case []func(float64):
		count := len(value)
		if count == 0 {
			return nil, true
		}
		listeners := make([]func(MediaPlayer, float64), count)
		for i, v := range value {
			if v == nil {
				return nil, false
			}
			listeners[i] = func(player MediaPlayer, time float64) {
				v(time)
			}
		}
		return listeners, true

	case []func(MediaPlayer):
		count := len(value)
		if count == 0 {
			return nil, true
		}
		listeners := make([]func(MediaPlayer, float64), count)
		for i, v := range value {
			if v == nil {
				return nil, false
			}
			listeners[i] = func(player MediaPlayer, time float64) {
				v(player)
			}
		}
		return listeners, true

	case []func():
		count := len(value)
		if count == 0 {
			return nil, true
		}
		listeners := make([]func(MediaPlayer, float64), count)
		for i, v := range value {
			if v == nil {
				return nil, false
			}
			listeners[i] = func(player MediaPlayer, time float64) {
				v()
			}
		}
		return listeners, true

	case []interface{}:
		count := len(value)
		if count == 0 {
			return nil, true
		}
		listeners := make([]func(MediaPlayer, float64), count)
		for i, v := range value {
			if v == nil {
				return nil, false
			}
			switch v := v.(type) {
			case func(MediaPlayer, float64):
				listeners[i] = v

			case func(float64):
				listeners[i] = func(player MediaPlayer, time float64) {
					v(time)
				}

			case func(MediaPlayer):
				listeners[i] = func(player MediaPlayer, time float64) {
					v(player)
				}

			case func():
				listeners[i] = func(player MediaPlayer, time float64) {
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

func valueToPlayerErrorListeners(value interface{}) ([]func(MediaPlayer, int, string), bool) {
	if value == nil {
		return nil, true
	}

	switch value := value.(type) {
	case func(MediaPlayer, int, string):
		return []func(MediaPlayer, int, string){value}, true

	case func(int, string):
		fn := func(player MediaPlayer, code int, message string) {
			value(code, message)
		}
		return []func(MediaPlayer, int, string){fn}, true

	case func(MediaPlayer):
		fn := func(player MediaPlayer, code int, message string) {
			value(player)
		}
		return []func(MediaPlayer, int, string){fn}, true

	case func():
		fn := func(player MediaPlayer, code int, message string) {
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
			listeners[i] = func(player MediaPlayer, code int, message string) {
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
			listeners[i] = func(player MediaPlayer, code int, message string) {
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
			listeners[i] = func(player MediaPlayer, code int, message string) {
				v()
			}
		}
		return listeners, true

	case []interface{}:
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
				listeners[i] = func(player MediaPlayer, code int, message string) {
					v(code, message)
				}

			case func(MediaPlayer):
				listeners[i] = func(player MediaPlayer, code int, message string) {
					v(player)
				}

			case func():
				listeners[i] = func(player MediaPlayer, code int, message string) {
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

func playerEvents() []struct{ tag, cssTag string } {
	return []struct{ tag, cssTag string }{
		{AbortEvent, "onabort"},
		{CanPlayEvent, "oncanplay"},
		{CanPlayThroughEvent, "oncanplaythrough"},
		{CompleteEvent, "oncomplete"},
		{EmptiedEvent, "onemptied"},
		{EndedEvent, "ended"},
		{LoadedDataEvent, "onloadeddata"},
		{LoadedMetadataEvent, "onloadedmetadata"},
		{LoadStartEvent, "onloadstart"},
		{PauseEvent, "onpause"},
		{PlayEvent, "onplay"},
		{PlayingEvent, "onplaying"},
		{ProgressEvent, "onprogress"},
		{SeekedEvent, "onseeked"},
		{SeekingEvent, "onseeking"},
		{StalledEvent, "onstalled"},
		{SuspendEvent, "onsuspend"},
		{WaitingEvent, "onwaiting"},
	}
}

func (player *mediaPlayerData) propertyChanged(tag string) {
	if player.created {
		switch tag {
		case Controls, Loop:
			value, _ := boolProperty(player, tag, player.Session())
			if value {
				updateBoolProperty(player.htmlID(), tag, value, player.Session())
			} else {
				removeProperty(player.htmlID(), tag, player.Session())
			}

		case Muted:
			value, _ := boolProperty(player, tag, player.Session())
			if value {
				player.Session().runScript("setMediaMuted('" + player.htmlID() + "', true)")
			} else {
				player.Session().runScript("setMediaMuted('" + player.htmlID() + "', false)")
			}

		case Preload:
			value, _ := enumProperty(player, tag, player.Session(), 0)
			values := enumProperties[Preload].values
			updateProperty(player.htmlID(), tag, values[value], player.Session())

		case AbortEvent, CanPlayEvent, CanPlayThroughEvent, CompleteEvent, EmptiedEvent,
			EndedEvent, LoadedDataEvent, LoadedMetadataEvent, PauseEvent, PlayEvent, PlayingEvent, ProgressEvent,
			LoadStartEvent, SeekedEvent, SeekingEvent, StalledEvent, SuspendEvent, WaitingEvent:

			for _, event := range playerEvents() {
				if event.tag == tag {
					if value := player.getRaw(event.tag); value != nil {
						switch value := value.(type) {
						case []func(MediaPlayer):
							if len(value) > 0 {
								fn := fmt.Sprintf(`playerEvent(this, "%s")`, event.tag)
								updateProperty(player.htmlID(), event.cssTag, fn, player.Session())
								return
							}
						}
					}
					updateProperty(player.htmlID(), tag, "", player.Session())
					break
				}

			}
		case TimeUpdateEvent:
			if value := player.getRaw(tag); value != nil {
				updateProperty(player.htmlID(), "ontimeupdate", "playerTimeUpdatedEvent(this)", player.Session())
			} else {
				updateProperty(player.htmlID(), "ontimeupdate", "", player.Session())
			}

		case VolumeChangedEvent:
			if value := player.getRaw(tag); value != nil {
				updateProperty(player.htmlID(), "onvolumechange", "playerVolumeChangedEvent(this)", player.Session())
			} else {
				updateProperty(player.htmlID(), "onvolumechange", "", player.Session())
			}

		case DurationChangedEvent:
			if value := player.getRaw(tag); value != nil {
				updateProperty(player.htmlID(), "ondurationchange", "playerDurationChangedEvent(this)", player.Session())
			} else {
				updateProperty(player.htmlID(), "ondurationchange", "", player.Session())
			}

		case RateChangedEvent:
			if value := player.getRaw(tag); value != nil {
				updateProperty(player.htmlID(), "onratechange", "playerRateChangedEvent(this)", player.Session())
			} else {
				updateProperty(player.htmlID(), "onratechange", "", player.Session())
			}

		case PlayerErrorEvent:
			if value := player.getRaw(tag); value != nil {
				updateProperty(player.htmlID(), "onerror", "playerErrorEvent(this)", player.Session())
			} else {
				updateProperty(player.htmlID(), "onerror", "", player.Session())
			}

		case Source:
			updateInnerHTML(player.htmlID(), player.Session())
		}
	}
}

func (player *mediaPlayerData) htmlSubviews(self View, buffer *strings.Builder) {
	if value := player.getRaw(Source); value != nil {
		if sources, ok := value.([]MediaSource); ok && len(sources) > 0 {
			session := player.Session()
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
	for _, tag := range []string{Controls, Loop, Muted, Preload} {
		if value, _ := boolProperty(player, tag, player.Session()); value {
			buffer.WriteRune(' ')
			buffer.WriteString(tag)
		}
	}

	if value, ok := enumProperty(player, Preload, player.Session(), 0); ok {
		values := enumProperties[Preload].values
		buffer.WriteString(` preload="`)
		buffer.WriteString(values[value])
		buffer.WriteRune('"')
	}

	for _, event := range playerEvents() {
		if value := player.getRaw(event.tag); value != nil {
			switch value := value.(type) {
			case []func(MediaPlayer):
				if len(value) > 0 {
					buffer.WriteString(` `)
					buffer.WriteString(event.cssTag)
					buffer.WriteString(`="playerEvent(this, \'`)
					buffer.WriteString(event.tag)
					buffer.WriteString(`\')"`)
				}
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

func (player *mediaPlayerData) handleCommand(self View, command string, data DataObject) bool {
	switch command {
	case AbortEvent, CanPlayEvent, CanPlayThroughEvent, CompleteEvent, LoadStartEvent,
		EmptiedEvent, EndedEvent, LoadedDataEvent, LoadedMetadataEvent, PauseEvent, PlayEvent,
		PlayingEvent, ProgressEvent, SeekedEvent, SeekingEvent, StalledEvent, SuspendEvent,
		WaitingEvent:

		if value := player.getRaw(command); value != nil {
			if listeners, ok := value.([]func(MediaPlayer)); ok {
				for _, listener := range listeners {
					listener(player)
				}
			}
		}

	case TimeUpdateEvent, DurationChangedEvent, RateChangedEvent, VolumeChangedEvent:
		if value := player.getRaw(command); value != nil {
			if listeners, ok := value.([]func(MediaPlayer, float64)); ok {
				time := dataFloatProperty(data, "value")
				for _, listener := range listeners {
					listener(player, time)
				}
			}
		}

	case PlayerErrorEvent:
		if value := player.getRaw(command); value != nil {
			if listeners, ok := value.([]func(MediaPlayer, int, string)); ok {
				code, _ := dataIntProperty(data, "code")
				message, _ := data.PropertyValue("message")
				for _, listener := range listeners {
					listener(player, code, message)
				}
			}
		}
	}

	return player.viewData.handleCommand(self, command, data)
}

func (player *mediaPlayerData) Play() {
	player.session.runScript(fmt.Sprintf(`mediaPlay('%v');`, player.htmlID()))
}

func (player *mediaPlayerData) Pause() {
	player.session.runScript(fmt.Sprintf(`mediaPause('%v');`, player.htmlID()))
}

func (player *mediaPlayerData) SetCurrentTime(seconds float64) {
	player.session.runScript(fmt.Sprintf(`mediaSetSetCurrentTime('%v', %v);`, player.htmlID(), seconds))
}

func (player *mediaPlayerData) getFloatPlayerProperty(tag string) (float64, bool) {

	script := allocStringBuilder()
	defer freeStringBuilder(script)

	script.WriteString(`const element = document.getElementById('`)
	script.WriteString(player.htmlID())
	script.WriteString(`');
if (element && element.`)
	script.WriteString(tag)
	script.WriteString(`) {
	sendMessage('answer{answerID=' + answerID + ',`)
	script.WriteString(tag)
	script.WriteString(`=' + element.`)
	script.WriteString(tag)
	script.WriteString(` + '}');
} else {
	sendMessage('answer{answerID=' + answerID + ',`)
	script.WriteString(tag)
	script.WriteString(`=0}');
}`)

	result := player.Session().runGetterScript(script.String())
	switch result.Tag() {
	case "answer":
		if value, ok := result.PropertyValue(tag); ok {
			w, err := strconv.ParseFloat(value, 32)
			if err == nil {
				return w, true
			}
			ErrorLog(err.Error())
		}

	case "error":
		if text, ok := result.PropertyValue("errorText"); ok {
			ErrorLog(text)
		} else {
			ErrorLog("error")
		}

	default:
		ErrorLog("Unknown answer: " + result.Tag())
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
	player.session.runScript(fmt.Sprintf(`mediaSetPlaybackRate('%v', %v);`, player.htmlID(), rate))
}

func (player *mediaPlayerData) PlaybackRate() float64 {
	if result, ok := player.getFloatPlayerProperty("playbackRate"); ok {
		return result
	}
	return 1
}

func (player *mediaPlayerData) SetVolume(volume float64) {
	if volume >= 0 && volume <= 1 {
		player.session.runScript(fmt.Sprintf(`mediaSetVolume('%v', %v);`, player.htmlID(), volume))
	}
}

func (player *mediaPlayerData) Volume() float64 {
	if result, ok := player.getFloatPlayerProperty("volume"); ok {
		return result
	}
	return 1
}

func (player *mediaPlayerData) getBoolPlayerProperty(tag string) (bool, bool) {

	script := allocStringBuilder()
	defer freeStringBuilder(script)

	script.WriteString(`const element = document.getElementById('`)
	script.WriteString(player.htmlID())
	script.WriteString(`');
if (element && element.`)
	script.WriteString(tag)
	script.WriteString(`) {
	sendMessage('answer{answerID=' + answerID + ',`)
	script.WriteString(tag)
	script.WriteString(`=1}')
} else {
	sendMessage('answer{answerID=' + answerID + ',`)
	script.WriteString(tag)
	script.WriteString(`=0}')
}`)

	result := player.Session().runGetterScript(script.String())
	switch result.Tag() {
	case "answer":
		if value, ok := result.PropertyValue(tag); ok {
			if value == "1" {
				return true, true
			}
			return false, true
		}

	case "error":
		if text, ok := result.PropertyValue("errorText"); ok {
			ErrorLog(text)
		} else {
			ErrorLog("error")
		}

	default:
		ErrorLog("Unknown answer: " + result.Tag())
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
