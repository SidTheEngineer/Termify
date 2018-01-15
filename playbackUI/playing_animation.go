package playbackUI

import (
	tui "github.com/gizak/termui"
)

const (
	playingAnimationUIHeight = 10
)

func createPlayingAnimationUI() *tui.BarChart {
	bars := tui.NewBarChart()
	bars.Data = []int{3, 2, 5, 3, 9, 5, 3, 2, 5, 8, 3, 2, 4, 3, 2, 4, 6, 4, 8, 5, 9, 4, 3, 6, 5, 3, 6}
	bars.Height = playingAnimationUIHeight
	bars.BarColor = tui.ColorGreen

	return bars
}
