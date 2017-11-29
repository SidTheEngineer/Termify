package app

import (
	"fmt"

	"github.com/CrowdSurge/banner"
	"github.com/fatih/color"
)

const browse = "browse"

// ViewTypes holds the types of views and their choices that users
// can choose from.
var ViewTypes = make(map[string]*View)

// Config will hold all of the meta data about the current state of
// Termify, such as the currentView, history stack, etc.
type Config struct {
	CurrentView *View
}

// View is a struct that contains a view's information and behaviors, such
// as the type of view it is, and it being displayed in the console.
type View struct {
	name    string
	choices []string
}

// Print prints a View's informatio to the console.
func (v View) Print() {
	color.Cyan(banner.PrintS(v.name))
	fmt.Printf("\n\n")
	for i, choice := range v.choices {
		color.HiGreen(fmt.Sprintf("%d. %s\n", i, choice))
	}
	fmt.Printf("\n\n")
}

// Init initializes the app configuration information, such as the types of
// views that can be used and their corresponding choices.
func Init() {
	ViewTypes[browse] = &View{
		browse,
		[]string{
			"Devices",
			"Categories",
			"Featured Playlists"}}
}

// NewView is a View constructor that returns a pointer to a View
func NewView(name string) *View {
	return &View{name, ViewTypes[name].choices}
}
