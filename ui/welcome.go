package ui

import tui "github.com/gizak/termui"

// Welcome defines the login component one sees upon starting Termify.
type Welcome struct {
	view View
}

// Render mounts/displays a welcome component in the terminal.
func (w Welcome) Render(uiConfig *Config) {
	welcome := createWelcomePar(uiConfig)

	tui.Body.AddRows(tui.NewRow(
		tui.NewCol(6, 0, welcome),
	))

	tui.Body.Align()
	tui.Render(tui.Body)
}

// NewWelcomeComponent returns the login component one sees upon starting Termify.
func NewWelcomeComponent() Welcome {
	return Welcome{
		view: View{
			Name: "welcome",
		},
	}
}

func createWelcomePar(uiConfig *Config) *tui.Par {
	welcomePar := tui.NewPar(welcomeText)
	welcomePar.Height = 10
	welcomePar.Border = false
	welcomePar.TextFgColor = tui.ColorGreen
	welcomePar.PaddingLeft = 2
	welcomePar.PaddingTop = 2

	return welcomePar
}
