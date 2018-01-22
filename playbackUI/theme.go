package playbackUI

import tui "github.com/gizak/termui"

const (
	themeBorderFg                = tui.ColorMagenta
	themeTextFgColor             = tui.ColorYellow
	themeBarColor                = tui.ColorGreen
	themeProgressGuageColor      = tui.ColorYellow
	themePercentColor            = tui.ColorYellow
	themePercentColorHighlighted = tui.ColorMagenta
	themeVolumeGaugeColor        = tui.ColorMagenta

	// By default, the border label Fg is green, so they aren't being set throughout
	// the app. Change this and set the values to it to make a change.
	themeBorderLabelFg = tui.ColorGreen
)
