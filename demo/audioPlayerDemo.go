package main

import (
	"fmt"

	"github.com/anoshenko/rui"
)

const audioPlayerDemoText = `
GridLayout {
	cell-height = "auto, auto, 1fr, auto", width = 100%, height = 100%,
	content = [
		ListLayout {
			row = 0, orientation = start-to-end, padding = 4px,
			content = [
				Checkbox { 
					id = showAudioPlayerControls, content = "Controls"
				},
				Checkbox { 
					id = showAudioPlayerLoop, content = "Loop"
				},
				Checkbox { 
					id = showAudioPlayerMuted, content = "Muted"
				},
			],
		},
		AudioPlayer {
			row = 1, id = audioPlayer, src = "https://alxanosoft.com/jazzy-loop-1.mp3", 
		},
		ListLayout {
			row = 2, orientation = start-to-end, vertical-align = top, padding = 8px,
			content = [
				NumberPicker {
					id = audioPlayerSlider, width = 200px, type = slider
				}
				Button {
					id = audioPlayerPlay, content = "Play", margin-left = 16px
				}
			]
		},
		Resizable {
			row = 3, side = top, background-color = lightgrey, height = 200px,
			content = EditView {
				id = audioPlayerEventsLog, type = multiline, read-only = true, wrap = true,
			}
		},
	]
}`

var audioPlayerPause = true

func createAudioPlayerDemo(session rui.Session) rui.View {
	view := rui.CreateViewFromText(session, audioPlayerDemoText)
	if view == nil {
		return nil
	}

	createListener := func(event string) func() {
		return func() {
			rui.AppendEditText(view, "audioPlayerEventsLog", event+"\n")
			rui.ScrollViewToEnd(view, "audioPlayerEventsLog")
		}
	}
	createListener2 := func(event string) func(value float64) {
		return func(value float64) {
			rui.AppendEditText(view, "audioPlayerEventsLog", fmt.Sprintf("%s: %g\n", event, value))
			rui.ScrollViewToEnd(view, "audioPlayerEventsLog")
		}
	}

	rui.Set(view, "showAudioPlayerControls", rui.CheckboxChangedEvent, func(state bool) {
		rui.Set(view, "audioPlayer", rui.Controls, state)
	})

	rui.Set(view, "showAudioPlayerLoop", rui.CheckboxChangedEvent, func(state bool) {
		rui.Set(view, "audioPlayer", rui.Loop, state)
	})

	rui.Set(view, "showAudioPlayerMuted", rui.CheckboxChangedEvent, func(state bool) {
		rui.Set(view, "audioPlayer", rui.Muted, state)
	})

	for _, event := range []string{rui.AbortEvent, rui.CanPlayEvent, rui.CanPlayThroughEvent,
		rui.CompleteEvent, rui.EmptiedEvent, rui.EndedEvent, rui.LoadStartEvent,
		rui.LoadedMetadataEvent, rui.PlayingEvent, rui.SeekedEvent, rui.SeekingEvent,
		rui.StalledEvent, rui.SuspendEvent, rui.WaitingEvent} {

		rui.Set(view, "audioPlayer", event, createListener(event))
	}

	for _, event := range []string{rui.DurationChangedEvent, rui.RateChangedEvent, rui.VolumeChangedEvent} {

		rui.Set(view, "audioPlayer", event, createListener2(event))
	}

	rui.Set(view, "audioPlayer", rui.PlayEvent, func() {
		rui.AppendEditText(view, "audioPlayerEventsLog", "play-event\n")
		rui.ScrollViewToEnd(view, "audioPlayerEventsLog")
		rui.Set(view, "audioPlayerPlay", rui.Content, "Pause")
		audioPlayerPause = false
	})

	rui.Set(view, "audioPlayer", rui.PauseEvent, func() {
		rui.AppendEditText(view, "audioPlayerEventsLog", "pause-event\n")
		rui.ScrollViewToEnd(view, "audioPlayerEventsLog")
		rui.Set(view, "audioPlayerPlay", rui.Content, "Play")
		audioPlayerPause = true
	})

	rui.Set(view, "audioPlayer", rui.LoadedDataEvent, func() {
		rui.AppendEditText(view, "audioPlayerEventsLog", "loaded-data-event\n")
		rui.ScrollViewToEnd(view, "audioPlayerEventsLog")
		rui.Set(view, "audioPlayerSlider", rui.Max, rui.MediaPlayerDuration(view, "audioPlayer"))
	})

	rui.Set(view, "audioPlayer", rui.TimeUpdateEvent, func(time float64) {
		rui.AppendEditText(view, "audioPlayerEventsLog", fmt.Sprintf("time-update-event %gs\n", time))
		rui.ScrollViewToEnd(view, "audioPlayerEventsLog")
		rui.Set(view, "audioPlayerSlider", rui.Value, time)
	})

	rui.Set(view, "audioPlayer", rui.PlayerErrorEvent, func(code int, message string) {
		rui.AppendEditText(view, "audioPlayerEventsLog", fmt.Sprintf("player-error-event: code = %d, message = '%s'\n", code, message))
		rui.ScrollViewToEnd(view, "audioPlayerEventsLog")
	})

	rui.Set(view, "audioPlayerPlay", rui.ClickEvent, func() {
		if audioPlayerPause {
			rui.MediaPlayerPlay(view, "audioPlayer")
		} else {
			rui.MediaPlayerPause(view, "audioPlayer")
		}
	})

	rui.Set(view, "audioPlayerSlider", rui.NumberChangedEvent, func(value float64) {
		if audioPlayerPause {
			rui.SetMediaPlayerCurrentTime(view, "audioPlayer", value)
		}
	})

	return view
}
