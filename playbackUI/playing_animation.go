package playbackUI

import (
	"math/rand"
	"strconv"

	tui "github.com/gizak/termui"
)

const (
	playingAnimationUIHeight = 10
	playingAnimationUIWidth  = 12
)

// The termui bargraphs need bar labels for the bars to display, so we'll
// just use a bunch of empty strings to have them not show.
var barLabels = []string{
	"", "", "", "", "", "", "", "", "", "", "", "",
	"", "", "", "", "", "", "", "", "", "", "", "",
	"", "", "", "", "", "", "", "", "", "", "", "",
	"", "", "", "", "", "", "", "", "", "", "", "",
	"", "", "", "", "", "", "", "", "", "", "", "",
	"", "", "", "", "", "", "", "", "", "", "", "",
	"", "", "",
}

func createPlayingAnimationUI(uiConfig *Config) *tui.BarChart {
	bars := tui.NewBarChart()
	bars.Data = rand.Perm(100)
	bars.DataLabels = barLabels
	bars.Height = playingAnimationUIHeight

	// TODO: Change bar color based on BPM? (or some other track feature attribute?)
	// if uiConfig.currentTrack.BPM > 120 {
	// 	bars.BarColor = tui.ColorGreen
	// } else {
	// 	bars.BarColor = tui.ColorBlue
	// }

	bars.BarColor = themeBarColor
	bars.BorderFg = themeBorderFg
	bars.BorderLabel = strconv.Itoa(int(uiConfig.currentTrack.BPM)) + " BPM"

	return bars
}

func updatePlayingAnimationUI(uiConfig *Config) {
	newPlayingAnimationUI := createPlayingAnimationUI(uiConfig)

	tui.Body.Rows[3].Cols[0] = tui.NewCol(playingAnimationUIWidth, 0, newPlayingAnimationUI)
	tui.Body.Align()
	tui.Render(tui.Body)
}
