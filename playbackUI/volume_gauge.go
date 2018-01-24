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
	volumeGauge.BorderLabel = "Volume"
	volumeGauge.Percent = int(uiConfig.currentDevice.Volume)

	return volumeGauge
}

func updateVolumeGauge(uiConfig *Config, incrementAmount int) {
	newGuage := createVolumeGauge(uiConfig)
	newGuage.Percent = newGuage.Percent + incrementAmount

	tui.Body.Rows[2].Cols[0] = tui.NewCol(volumeGuageWidth, 0, newGuage)
}
