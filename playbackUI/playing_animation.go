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

func createPlayingAnimationUI() *tui.BarChart {
	bars := tui.NewBarChart()
	bars.Data = rand.Perm(100)
	bars.DataLabels = barLabels
	bars.Height = playingAnimationUIHeight
	bars.BarColor = themeBarColor
	bars.BorderFg = themeBorderFg

	return bars
}

func updatePlayingAnimationUI() {
	newPlayingAnimationUI := createPlayingAnimationUI()

	tui.Body.Rows[3].Cols[0] = tui.NewCol(playingAnimationUIWidth, 0, newPlayingAnimationUI)
	tui.Body.Align()
	tui.Render(tui.Body)
}
