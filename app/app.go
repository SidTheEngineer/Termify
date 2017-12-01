package app

// Config will hold all of the meta data about the current state of
// Termify, such as the currentView, history stack, etc.
type Config struct {
	CurrentView View
}

// SetCurrentView updates the current view of Termify. It does not necessarily
// mean that the view is being disiplayed, however.
func (v *Config) SetCurrentView(newView View) {
	v.CurrentView = newView
}
