package playbackUI

import (
	tui "github.com/gizak/termui"
)

const (
	updateUIWaitTime = 100
	playKey          = "sys/kbd/1"
	pauseKey         = "sys/kbd/2"
	prevKey          = "sys/kbd/3"
	nextKey          = "sys/kbd/4"
	quitKey          = "sys/kbd/q"
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

	contextJSON := getCurrentlyPlayingContext(uiConfig)
	uiConfig.SetCurrentlyPlayingContext(contextJSON)

	// TODO: This line can throw a 'panic: interface conversion: interface {} is nil, not map[string]interface {}'
	// and needs to be fixed. I think this error arises when there are no tracks in the spotify player to begin with.
	trackInfo := getTrackInformationFromJSON(contextJSON)
	deviceInfo := getDeviceInformationFromJSON(contextJSON)

	progressInSeconds := (uiConfig.timeElapsedFromTickerStart + int(deviceInfo.ProgressMs)) / 1000

	controls := createControls(uiConfig)
	currentlyPlayingUI := createCurrentlyPlayingUI(uiConfig, trackInfo, deviceInfo)
	trackProgressTime := createTrackProgressTime(uiConfig, progressInSeconds)
	trackProgressGuage := createTrackProgressGuage(uiConfig, progressInSeconds)

	if tui.Body != nil {
		ResetTerminal()
	} else {
		tui.Init()
	}

	tui.Body.AddRows(
		tui.NewRow(
			tui.NewCol(controlsWidth, 0, controls),
			tui.NewCol(currentlyPlayingWidth, 0, currentlyPlayingUI),
		),
		tui.NewRow(
			tui.NewCol(progressTimeWidth, 0, trackProgressTime),
			tui.NewCol(progressGuageWidth, 0, trackProgressGuage),
		),
	)

	tui.Body.Align()
	tui.Render(tui.Body)
}
