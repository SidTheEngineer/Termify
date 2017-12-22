package ui

// CurrentlyPlaying is a component that encompasses all of the UI
// pertaining to a currently playing track and the current device
// that music is being played on.
type CurrentlyPlaying struct {
	view View
}

// NewCurrentlyPlayingComponent returns a new component that contains UI
// pertaining to the currently playing track as well as the current device
// that music is being played on.
func NewCurrentlyPlayingComponent() CurrentlyPlaying {
	return CurrentlyPlaying{
		view: View{
			Name:    "currently playing",
			Choices: []Choice{},
		},
	}
}
