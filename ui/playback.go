package ui

import (
	tui "github.com/gizak/termui"
)

// Playback is a component that contains all of the UI related to
// music playback, such as playing, pausing, current song, etc.
type Playback struct {
	view View
}

// NewPlaybackComponent returns a new component that contains
// all of the UI related to music playback, such as playing, pausing, current song, etc..
func NewPlaybackComponent() Playback {
	return Playback{
		view: View{
			Name: "playback",
			Choices: []Choice{
				playChoice(),
				pauseChoice(),
				skipChoice(),
				backChoice(),
			},
		},
	}
}

// Render mounts/displays a Playback component in the terminal.
func (p Playback) Render(uiConfig *Config) {

	tui.ResetHandlers()
	controls := createControls(uiConfig)

	if tui.Body != nil {
		ResetRows()
	} else {
		tui.Init()
	}

	tui.Body.AddRows(tui.NewRow(
		tui.NewCol(2, 0, controls),
	))

	tui.Body.Align()
	tui.Render(tui.Body)

}

// https://beta.developer.spotify.com/documentation/web-api/reference/player/start-a-users-playback/
func playChoice() Choice {
	return Choice{
		Name:         "[ 1 ] - Play",
		APIRoute:     "https://api.spotify.com/v1/me/player/play",
		APIMethod:    "PUT",
		ResponseType: "",
	}
}

// https://beta.developer.spotify.com/documentation/web-api/reference/player/pause-a-users-playback/
func pauseChoice() Choice {
	return Choice{
		Name:         "[ 2 ] - Pause",
		APIRoute:     "https://api.spotify.com/v1/me/player/pause",
		APIMethod:    "PUT",
		ResponseType: "",
	}
}

// https://developer.spotify.com/web-api/skip-users-playback-to-next-track/
func skipChoice() Choice {
	return Choice{
		Name:         "[ 3 ] - Next",
		APIRoute:     "https://api.spotify.com/v1/me/player/next",
		APIMethod:    "POST",
		ResponseType: "",
	}
}

// https://developer.spotify.com/web-api/skip-users-playback-to-previous-track/
func backChoice() Choice {
	return Choice{
		Name:         "[ 4 ] - Previous",
		APIRoute:     "https://api.spotify.com/v1/me/player/previous",
		APIMethod:    "POST",
		ResponseType: "",
	}
}

func createControls(uiConfig *Config) *tui.List {
	controls := tui.NewList()
	controls.Border = true
	controls.BorderFg = tui.ColorGreen
	controls.BorderLabel = "Controls"
	controls.Height = 10
	controls.ItemFgColor = tui.ColorYellow
	controls.Items = []string{
		NewLine,
		ExitText,
		NewLine,
		playChoice().Name,
		pauseChoice().Name,
		skipChoice().Name,
		backChoice().Name,
	}

	tui.Handle("/sys/kbd/q", func(tui.Event) {
		tui.StopLoop()
	})

	attachPlaybackComponentHandlers(uiConfig)

	return controls
}

func attachPlaybackComponentHandlers(uiConfig *Config) {
	playbackChoices := NewPlaybackComponent().view.Choices
	tui.Handle("/sys/kbd/q", func(e tui.Event) {
		tui.StopLoop()
	})

	// Unfortunately, these have to be hardcoded. Handle() breaks when trying to
	// attach in a loop.
	tui.Handle("sys/kbd/1", func(e tui.Event) {
		req := playbackChoices[0].CreateAPIRequest(uiConfig.AccessToken)
		playbackChoices[0].SendAPIRequest(req)
	})

	tui.Handle("sys/kbd/2", func(e tui.Event) {
		req := playbackChoices[1].CreateAPIRequest(uiConfig.AccessToken)
		playbackChoices[1].SendAPIRequest(req)
	})

	tui.Handle("sys/kbd/3", func(e tui.Event) {
		req := playbackChoices[2].CreateAPIRequest(uiConfig.AccessToken)
		playbackChoices[2].SendAPIRequest(req)
	})

	tui.Handle("sys/kbd/4", func(e tui.Event) {
		req := playbackChoices[3].CreateAPIRequest(uiConfig.AccessToken)
		playbackChoices[3].SendAPIRequest(req)
	})
}
