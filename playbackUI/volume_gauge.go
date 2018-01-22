package playbackUI

import (
	tui "github.com/gizak/termui"
)

const (
	volumeGuageWidth  = 12
	volumeGuageHeight = 3
)

func createVolumeGauge(uiConfig *Config) *tui.Gauge {
	volumeGauge := tui.NewGauge()
	volumeGauge.Height = volumeGuageHeight
	volumeGauge.BarColor = themeVolumeGaugeColor
	volumeGauge.BorderFg = themeBorderFg
	volumeGauge.PercentColor = themePercentColor
	volumeGauge.PercentColorHighlighted = tui.ColorWhite
	volumeGauge.Percent = int(uiConfig.currentDevice.Volume)

	return volumeGauge
}
