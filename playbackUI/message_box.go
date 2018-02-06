package playbackUI

import tui "github.com/gizak/termui"

const (
	messageBoxWidth       = 12
	messageBoxHeight      = 3
	messageBoxBorderLabel = "Notifications"
)

func createMessageBox(uiConfig *Config, message string) *tui.Par {
	messageBox := tui.NewPar(message)
	messageBox.Height = messageBoxHeight
	messageBox.Border = true
	messageBox.BorderFg = themeBorderFg
	messageBox.BorderLabel = messageBoxBorderLabel
	messageBox.TextFgColor = themeTextFgColor

	return messageBox
}

func updateMessageBox(uiConfig *Config, newMessage string) {
	newMessageBox := createMessageBox(uiConfig, newMessage)

	tui.Body.Rows[4].Cols[0] = tui.NewCol(messageBoxWidth, 0, newMessageBox)
	tui.Body.Align()
	tui.Render(tui.Body)
}
