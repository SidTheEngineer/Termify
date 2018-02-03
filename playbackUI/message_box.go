package playbackUI

import tui "github.com/gizak/termui"

const (
	messageBoxHeight = 2
)

func createMessageBox(uiConfig *Config, message string) *tui.Par {
	return tui.NewPar(message)
}
