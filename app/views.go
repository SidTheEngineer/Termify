package app

import (
	"fmt"

	"github.com/CrowdSurge/banner"
	"github.com/fatih/color"
)

const browse = "browse"

// View is a struct that contains a view's information and behaviors, such
// as the type of view it is, and it being displayed in the console.
type View struct {
	name    string
	choices []Choice
}

// Print prints a View's informatio to the console.
func (v View) Print() {
	color.Cyan(banner.PrintS(v.name))
	fmt.Printf("\n\n")
	for i, choice := range v.choices {
		color.HiGreen(fmt.Sprintf("%d. %s\n", i, choice.name))
	}
	fmt.Printf("\n\n")
}

// NewBrowseView is a View constructor that returns a browse menu view (the main menu)
func NewBrowseView() View {
	return View{
		browse,
		[]Choice{
			Devices,
			Categories,
			FeaturedPlaylists,
		},
	}
}
