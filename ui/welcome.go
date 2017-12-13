package ui

import tui "github.com/gizak/termui"

func welcomeComponent(uiConfig *Config) *tui.Par {
	welcomePar := tui.NewPar(welcomeText)
	welcomePar.Height = 10
	welcomePar.Border = false
	welcomePar.TextFgColor = tui.ColorGreen
	welcomePar.PaddingLeft = 2
	welcomePar.PaddingTop = 2

	return welcomePar
}
