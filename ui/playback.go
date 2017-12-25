package ui

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	tui "github.com/gizak/termui"
)

const (
	currentlyPlayingContextURL    = "https://api.spotify.com/v1/me/player"
	currentlyPlayingContextMethod = "GET"
)

// Playback is a component that contains all of the UI related to
// music playback, such as playing, pausing, current song, etc.
type Playback struct {
	view View
}

// Track represents track information that we want from the Spotify track object.
type Track struct {
	Name, Artists string
}

// Device represents device information that we want from the Sptify device object.
type Device struct {
	Name, DeviceType string
}

// NewPlaybackComponent returns a new component that contains
// all of the UI related to music playback, such as playing, pausing, current song, etc..
func NewPlaybackComponent() Playback {
	return Playback{
		view: View{
			Name: "playback",
			Choices: []Choice{
				playChoice(),
				pauseChoice(),
				skipChoice(),
				backChoice(),
			},
		},
	}
}

// Render mounts/displays a Playback component in the terminal.
func (p Playback) Render(uiConfig *Config) {

	tui.ResetHandlers()

	contextJSON := getCurrentlyPlayingContext(uiConfig)
	trackInfo := getTrackInformationFromJSON(contextJSON)
	deviceInfo := getDeviceInformationFromJSON(contextJSON)
	uiConfig.SetCurrentlyPlayingContext(contextJSON)

	controls := createControls(uiConfig)
	currentlyPlayingUI := createCurrentlyPlayingUI(uiConfig, trackInfo, deviceInfo)

	if tui.Body != nil {
		ResetTerminal()
	} else {
		tui.Init()
	}

	tui.Body.AddRows(tui.NewRow(
		tui.NewCol(2, 0, controls),
		tui.NewCol(4, 0, currentlyPlayingUI),
	))

	tui.Body.Align()
	tui.Render(tui.Body)
}

// https://beta.developer.spotify.com/documentation/web-api/reference/player/start-a-users-playback/
func playChoice() Choice {
	return Choice{
		Name:         "[ 1 ] - Play",
		APIRoute:     "https://api.spotify.com/v1/me/player/play",
		APIMethod:    "PUT",
		ResponseType: "",
	}
}

// https://beta.developer.spotify.com/documentation/web-api/reference/player/pause-a-users-playback/
func pauseChoice() Choice {
	return Choice{
		Name:         "[ 2 ] - Pause",
		APIRoute:     "https://api.spotify.com/v1/me/player/pause",
		APIMethod:    "PUT",
		ResponseType: "",
	}
}

// https://developer.spotify.com/web-api/skip-users-playback-to-next-track/
func skipChoice() Choice {
	return Choice{
		Name:         "[ 3 ] - Previous",
		APIRoute:     "https://api.spotify.com/v1/me/player/previous",
		APIMethod:    "POST",
		ResponseType: "",
	}
}

// https://developer.spotify.com/web-api/skip-users-playback-to-previous-track/
func backChoice() Choice {
	return Choice{
		Name:         "[ 4 ] - Next",
		APIRoute:     "https://api.spotify.com/v1/me/player/next",
		APIMethod:    "POST",
		ResponseType: "",
	}
}

func createControls(uiConfig *Config) *tui.List {
	controls := tui.NewList()
	controls.Border = true
	controls.BorderFg = tui.ColorMagenta
	controls.BorderLabel = "Controls"
	controls.Height = 10
	controls.ItemFgColor = tui.ColorYellow
	controls.Items = []string{
		NewLine,
		ExitText,
		NewLine,
		playChoice().Name,
		pauseChoice().Name,
		skipChoice().Name,
		backChoice().Name,
	}

	tui.Handle("/sys/kbd/q", func(tui.Event) {
		tui.StopLoop()
	})

	attachPlaybackComponentHandlers(uiConfig)

	return controls
}

func createCurrentlyPlayingUI(uiConfig *Config, trackInfo Track, deviceInfo Device) *tui.Par {
	currentlyPlayingUI := tui.NewPar(
		fmt.Sprintf(
			"\n%s - %s\n\n%s\n%s",
			deviceInfo.DeviceType,
			deviceInfo.Name,
			trackInfo.Name,
			trackInfo.Artists,
		),
	)
	currentlyPlayingUI.BorderLabel = "Currently Playing"
	currentlyPlayingUI.BorderFg = tui.ColorMagenta
	currentlyPlayingUI.TextFgColor = tui.ColorYellow
	currentlyPlayingUI.Height = 10

	return currentlyPlayingUI
}

func attachPlaybackComponentHandlers(uiConfig *Config) {
	playbackChoices := NewPlaybackComponent().view.Choices

	// Unfortunately, these have to be hardcoded. Handle() breaks when trying to
	// attach in a loop.
	tui.Handle("sys/kbd/1", func(e tui.Event) {
		req := playbackChoices[0].CreateAPIRequest(uiConfig.AccessToken)
		playbackChoices[0].SendAPIRequest(req)
	})

	tui.Handle("sys/kbd/2", func(e tui.Event) {
		req := playbackChoices[1].CreateAPIRequest(uiConfig.AccessToken)
		playbackChoices[1].SendAPIRequest(req)
	})

	tui.Handle("sys/kbd/3", func(e tui.Event) {
		req := playbackChoices[2].CreateAPIRequest(uiConfig.AccessToken)
		playbackChoices[2].SendAPIRequest(req)
		updateCurrentlyPlayingUI(uiConfig)
	})

	tui.Handle("sys/kbd/4", func(e tui.Event) {
		req := playbackChoices[3].CreateAPIRequest(uiConfig.AccessToken)
		playbackChoices[3].SendAPIRequest(req)
		updateCurrentlyPlayingUI(uiConfig)
	})
}

func getCurrentlyPlayingContext(uiConfig *Config) map[string]interface{} {
	var jsonMap interface{}
	client := http.Client{}
	req, _ := http.NewRequest(currentlyPlayingContextMethod, currentlyPlayingContextURL, nil)

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", uiConfig.AccessToken.Token))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Accept", "application/json")

	resp, err := client.Do(req)

	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	bytes, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal(bytes, &jsonMap)

	return jsonMap.(map[string]interface{})
}

func getDeviceInformationFromJSON(context map[string]interface{}) Device {
	deviceName := context["device"].(map[string]interface{})["name"].(string)
	deviceType := context["device"].(map[string]interface{})["type"].(string)
	return Device{
		Name:       deviceName,
		DeviceType: deviceType,
	}
}

func getTrackInformationFromJSON(context map[string]interface{}) Track {
	trackArtists := ""
	trackName := context["item"].(map[string]interface{})["name"].(string)
	artistJSONArray := context["item"].(map[string]interface{})["artists"].([]interface{})

	for i, artist := range artistJSONArray {
		trackArtists += artist.(map[string]interface{})["name"].(string)
		if i != len(artistJSONArray)-1 {
			trackArtists += ", "
		}
	}

	return Track{
		Name:    trackName,
		Artists: trackArtists,
	}
}

func updateCurrentlyPlayingUI(uiConfig *Config) {
	currentContext := getCurrentlyPlayingContext(uiConfig)
	currentTrack := getTrackInformationFromJSON(currentContext)
	deviceInfo := getDeviceInformationFromJSON(currentContext)
	newCurrentlyPlayingUI := createCurrentlyPlayingUI(uiConfig, currentTrack, deviceInfo)

	tui.Body.Rows[0].Cols[1] = tui.NewCol(4, 0, newCurrentlyPlayingUI)
	tui.Body.Align()
	tui.Render(tui.Body)
}
