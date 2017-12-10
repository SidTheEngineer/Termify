package api

import "Termify/ui"

// PlayChoice returns a Choice corresponding to the Spotfify "play" endpoint.
// https://beta.developer.spotify.com/documentation/web-api/reference/player/start-a-users-playback/
func PlayChoice() ui.Choice {
	return ui.Choice{
		Name:         "Play",
		APIRoute:     "https://api.spotify.com/v1/me/player/play",
		APIMethod:    "PUT",
		ResponseType: "",
	}
}

// PauseChoice returns a Choice corresponding to the Spotfify "pause" endpoint.
// https://beta.developer.spotify.com/documentation/web-api/reference/player/pause-a-users-playback/
func PauseChoice() ui.Choice {
	return ui.Choice{
		Name:         "Pause",
		APIRoute:     "https://api.spotify.com/v1/me/player/pause",
		APIMethod:    "PUT",
		ResponseType: "",
	}
}
