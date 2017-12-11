package ui

import tui "github.com/gizak/termui"

// NewPlaybackView returns a View corresponding to music playback
// (play, pause, currently playing, etc.)
func NewPlaybackView() View {
	return View{
		Name: "playback",
		Choices: []Choice{
			playChoice(),
			pauseChoice(),
		},
	}
}

// https://beta.developer.spotify.com/documentation/web-api/reference/player/start-a-users-playback/
func playChoice() Choice {
	return Choice{
		Name:         "Play",
		APIRoute:     "https://api.spotify.com/v1/me/player/play",
		APIMethod:    "PUT",
		ResponseType: "",
	}
}

// https://beta.developer.spotify.com/documentation/web-api/reference/player/pause-a-users-playback/
func pauseChoice() Choice {
	return Choice{
		Name:         "Pause",
		APIRoute:     "https://api.spotify.com/v1/me/player/pause",
		APIMethod:    "PUT",
		ResponseType: "",
	}
}

func playbackComponent() *tui.List {
	choiceList := tui.NewList()
	choiceList.Border = true
	choiceList.BorderFg = tui.ColorGreen
	choiceList.Height = 50
	choiceList.Items = []string{
		playChoice().Name,
		pauseChoice().Name,
	}

	return choiceList
}
