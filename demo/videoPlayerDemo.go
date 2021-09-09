package main

import (
	"fmt"

	"github.com/anoshenko/rui"
)

const videoPlayerDemoText = `
GridLayout {
	cell-height = "auto, 1fr, auto, auto", width = 100%, height = 100%,
	content = [
		ListLayout {
			row = 0, orientation = start-to-end, padding = 4px,
			content = [
				Checkbox { 
					id = showVideoPlayerControls, content = "Controls"
				},
				Checkbox { 
					id = showVideoPlayerLoop, content = "Loop"
				},
				Checkbox { 
					id = showVideoPlayerMuted, content = "Muted"
				},
			],
		},
		GridLayout {
			row = 1, id = videoPlayerContainer,
			content = VideoPlayer {
				id = videoPlayer, src = "https://alxanosoft.com/testVideo.mp4", video-width = 320,
			},
		},
		ListLayout {
			row = 2, orientation = start-to-end, vertical-align = top, padding = 8px,
			content = [
				NumberPicker {
					id = videoPlayerSlider, width = 200px, type = slider
				}
				Button {
					id = videoPlayerPlay, content = "Play", margin-left = 16px
				}
			]
		},
		Resizable {
			row = 3, side = top, background-color = lightgrey, height = 200px,
			content = EditView {
				id = videoPlayerEventsLog, type = multiline, read-only = true, wrap = true,
			}
		},
	]
}`

var videoPlayerPause = true

func createVideoPlayerDemo(session rui.Session) rui.View {
	view := rui.CreateViewFromText(session, videoPlayerDemoText)
	if view == nil {
		return nil
	}

	createListener := func(event string) func() {
		return func() {
			rui.AppendEditText(view, "videoPlayerEventsLog", event+"\n")
			rui.ScrollViewToEnd(view, "videoPlayerEventsLog")
		}
	}
	createListener2 := func(event string) func(value float64) {
		return func(value float64) {
			rui.AppendEditText(view, "videoPlayerEventsLog", fmt.Sprintf("%s: %g\n", event, value))
			rui.ScrollViewToEnd(view, "videoPlayerEventsLog")
		}
	}

	rui.Set(view, "videoPlayerContainer", rui.ResizeEvent, func(frame rui.Frame) {
		rui.Set(view, "videoPlayer", rui.VideoWidth, frame.Width)
		rui.Set(view, "videoPlayer", rui.VideoHeight, frame.Height)
	})

	rui.Set(view, "showVideoPlayerControls", rui.CheckboxChangedEvent, func(state bool) {
		rui.Set(view, "videoPlayer", rui.Controls, state)
	})

	rui.Set(view, "showVideoPlayerLoop", rui.CheckboxChangedEvent, func(state bool) {
		rui.Set(view, "videoPlayer", rui.Loop, state)
	})

	rui.Set(view, "showVideoPlayerMuted", rui.CheckboxChangedEvent, func(state bool) {
		rui.Set(view, "videoPlayer", rui.Muted, state)
	})

	for _, event := range []string{rui.AbortEvent, rui.CanPlayEvent, rui.CanPlayThroughEvent,
		rui.CompleteEvent, rui.EmptiedEvent, rui.EndedEvent, rui.LoadStartEvent,
		rui.LoadedMetadataEvent, rui.PlayingEvent, rui.SeekedEvent, rui.SeekingEvent,
		rui.StalledEvent, rui.SuspendEvent, rui.WaitingEvent} {

		rui.Set(view, "videoPlayer", event, createListener(event))
	}

	for _, event := range []string{rui.DurationChangedEvent, rui.RateChangedEvent, rui.VolumeChangedEvent} {

		rui.Set(view, "videoPlayer", event, createListener2(event))
	}

	rui.Set(view, "videoPlayer", rui.PlayEvent, func() {
		rui.AppendEditText(view, "videoPlayerEventsLog", "play-event\n")
		rui.ScrollViewToEnd(view, "videoPlayerEventsLog")
		rui.Set(view, "videoPlayerPlay", rui.Content, "Pause")
		videoPlayerPause = false
	})

	rui.Set(view, "videoPlayer", rui.PauseEvent, func() {
		rui.AppendEditText(view, "videoPlayerEventsLog", "pause-event\n")
		rui.ScrollViewToEnd(view, "videoPlayerEventsLog")
		rui.Set(view, "videoPlayerPlay", rui.Content, "Play")
		videoPlayerPause = true
	})

	rui.Set(view, "videoPlayer", rui.LoadedDataEvent, func() {
		rui.AppendEditText(view, "videoPlayerEventsLog", "loaded-data-event\n")
		rui.ScrollViewToEnd(view, "videoPlayerEventsLog")
		rui.Set(view, "videoPlayerSlider", rui.Max, rui.MediaPlayerDuration(view, "videoPlayer"))
	})

	rui.Set(view, "videoPlayer", rui.TimeUpdateEvent, func(time float64) {
		rui.AppendEditText(view, "videoPlayerEventsLog", fmt.Sprintf("time-update-event %gs\n", time))
		rui.ScrollViewToEnd(view, "videoPlayerEventsLog")
		rui.Set(view, "videoPlayerSlider", rui.Value, time)
	})

	rui.Set(view, "vodeoPlayer", rui.PlayerErrorEvent, func(code int, message string) {
		rui.AppendEditText(view, "vodeoPlayerEventsLog", fmt.Sprintf("player-error-event: code = %d, message = '%s'\n", code, message))
		rui.ScrollViewToEnd(view, "vodeoPlayerEventsLog")
	})

	rui.Set(view, "videoPlayerPlay", rui.ClickEvent, func() {
		if videoPlayerPause {
			rui.MediaPlayerPlay(view, "videoPlayer")
		} else {
			rui.MediaPlayerPause(view, "videoPlayer")
		}
	})

	rui.Set(view, "videoPlayerSlider", rui.NumberChangedEvent, func(value float64) {
		if videoPlayerPause {
			rui.SetMediaPlayerCurrentTime(view, "videoPlayer", value)
		}
	})

	return view
}

//MAH00054.MP4
