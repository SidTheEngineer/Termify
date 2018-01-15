package playbackUI

import (
	"math/rand"

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

func createPlayingAnimationUI(currentTime int) *tui.BarChart {
	bars := tui.NewBarChart()
	bars.Data = rand.Perm(100)
	bars.DataLabels = barLabels
	bars.Height = playingAnimationUIHeight
	bars.BarColor = tui.ColorGreen
	bars.BorderFg = tui.ColorMagenta

	return bars
}

func updatePlayingAnimationUI(currentTime int) {
	newPlayingAnimationUI := createPlayingAnimationUI(currentTime)

	tui.Body.Rows[2].Cols[0] = tui.NewCol(playingAnimationUIWidth, 0, newPlayingAnimationUI)
	tui.Body.Align()
	tui.Render(tui.Body)
}
