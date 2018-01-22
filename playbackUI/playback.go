package playbackUI

import (
	"github.com/SidTheEngineer/Termify/util"
	tui "github.com/gizak/termui"
)

const (
	updateUIWaitTime = 100
	playKey          = "sys/kbd/1"
	pauseKey         = "sys/kbd/2"
	prevKey          = "sys/kbd/3"
	nextKey          = "sys/kbd/4"
	volDownKey       = "sys/kbd/5"
	quitKey          = "sys/kbd/q"
)

// Playback is a component that contains all of the UI related to
// music playback, such as playing, pausing, current song, etc.
type Playback struct {
	view View
}

// NewPlaybackComponent returns a new component that contains
// all of the UI related to music playback, such as playing, pausing, current song, etc..
func NewPlaybackComponent(uiConfig *Config) Playback {
	return Playback{
		view: View{
			Name: "playback",
			Choices: []Choice{
				playChoice(),
				pauseChoice(),
				skipChoice(),
				backChoice(),
				volumeDownChoice(uiConfig),
			},
		},
	}
}

// Render mounts/displays a Playback component in the terminal.
func (p Playback) Render(uiConfig *Config) {

	tui.ResetHandlers()

	contextJSON := getCurrentlyPlayingContext(uiConfig)
	uiConfig.SetCurrentlyPlayingContext(contextJSON)
	trackInfo := getTrackInformationFromJSON(uiConfig, contextJSON)
	deviceInfo := getDeviceInformationFromJSON(uiConfig, contextJSON)

	progressInSeconds := (uiConfig.timeElapsedFromTickerStart + int(deviceInfo.ProgressMs)) / 1000

	controls := createControls(uiConfig)
	currentlyPlayingUI := createCurrentlyPlayingUI(uiConfig, trackInfo, deviceInfo)
	trackProgressGuage := createTrackProgressGauge(uiConfig, progressInSeconds)
	volumeGauge := createVolumeGauge(uiConfig)
	playingAnimationUI := createPlayingAnimationUI(progressInSeconds)

	if tui.Body != nil {
		util.ResetTerminal()
	} else {
		tui.Init()
	}

	tui.Body.AddRows(
		tui.NewRow(
			tui.NewCol(controlsWidth, 0, controls),
			tui.NewCol(currentlyPlayingWidth, 0, currentlyPlayingUI),
		),
		tui.NewRow(
			tui.NewCol(progressGuageWidth, 0, trackProgressGuage),
		),
		tui.NewRow(
			tui.NewCol(volumeGuageWidth, 0, volumeGauge),
		),
		tui.NewRow(
			tui.NewCol(playingAnimationUIWidth, 0, playingAnimationUI),
		),
	)

	tui.Body.Align()
	tui.Render(tui.Body)
}
