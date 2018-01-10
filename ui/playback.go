package ui

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	tui "github.com/gizak/termui"
)

const (
	currentlyPlayingContextURL    = "https://api.spotify.com/v1/me/player"
	currentlyPlayingContextMethod = "GET"
	playingText                   = "[ Playing ]"
	pausedText                    = "[ Paused ]"
	playChoiceNameText            = "[ 1 ] - Play"
	pauseChoiceNameText           = "[ 2 ] - Pause"
	previousChoiceNameText        = "[ 3 ] - Previous"
	nextChoiceNameText            = "[ 4 ] - Next"
	controlsBorderLabel           = "Controls"
	currentlyPlayingBorderLabel   = "Currently Playing"
)

// Playback is a component that contains all of the UI related to
// music playback, such as playing, pausing, current song, etc.
type Playback struct {
	view View
}

// Track represents track information that we want from the Spotify track object.
type Track struct {
	Name, Artists string
	DurationMs    float64
}

// Device represents device information that we want from the Sptify device object.
type Device struct {
	Name, DeviceType string
	ProgressMs       float64
	IsPlaying        bool
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

	// TODO: This line can throw a 'panic: interface conversion: interface {} is nil, not map[string]interface {}'
	// and needs to be fixed. I think this error arises when there are no tracks in the spotify player to begin with.
	trackInfo := getTrackInformationFromJSON(contextJSON)

	deviceInfo := getDeviceInformationFromJSON(contextJSON)
	uiConfig.SetCurrentlyPlayingContext(contextJSON)

	progressInSeconds := (uiConfig.timeElapsedFromTickerStart + int(deviceInfo.ProgressMs)) / 1000

	controls := createControls(uiConfig)
	currentlyPlayingUI := createCurrentlyPlayingUI(uiConfig, trackInfo, deviceInfo)
	trackProgressUI := createTrackProgressUI(uiConfig, progressInSeconds)

	if tui.Body != nil {
		ResetTerminal()
	} else {
		tui.Init()
	}

	tui.Body.AddRows(
		tui.NewRow(
			tui.NewCol(3, 0, controls),
			tui.NewCol(5, 0, currentlyPlayingUI),
		),
		tui.NewRow(
			tui.NewCol(3, 0, trackProgressUI),
		),
	)

	tui.Body.Align()
	tui.Render(tui.Body)
}

// https://beta.developer.spotify.com/documentation/web-api/reference/player/start-a-users-playback/
func playChoice() Choice {
	return Choice{
		Name:         playChoiceNameText,
		APIRoute:     "https://api.spotify.com/v1/me/player/play",
		APIMethod:    "PUT",
		ResponseType: "",
	}
}

// https://beta.developer.spotify.com/documentation/web-api/reference/player/pause-a-users-playback/
func pauseChoice() Choice {
	return Choice{
		Name:         pauseChoiceNameText,
		APIRoute:     "https://api.spotify.com/v1/me/player/pause",
		APIMethod:    "PUT",
		ResponseType: "",
	}
}

// https://developer.spotify.com/web-api/skip-users-playback-to-next-track/
func skipChoice() Choice {
	return Choice{
		Name:         previousChoiceNameText,
		APIRoute:     "https://api.spotify.com/v1/me/player/previous",
		APIMethod:    "POST",
		ResponseType: "",
	}
}

// https://developer.spotify.com/web-api/skip-users-playback-to-previous-track/
func backChoice() Choice {
	return Choice{
		Name:         nextChoiceNameText,
		APIRoute:     "https://api.spotify.com/v1/me/player/next",
		APIMethod:    "POST",
		ResponseType: "",
	}
}

func createControls(uiConfig *Config) *tui.List {
	controls := tui.NewList()
	controls.Border = true
	controls.BorderFg = tui.ColorMagenta
	controls.BorderLabel = controlsBorderLabel
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

func createCurrentlyPlayingUI(uiConfig *Config, trackInfo Track, deviceInfo Device) *tui.List {
	var playingState string

	// TODO: Abstract our ticker state from this method. Because updateCurrentlyPlayingUI
	// calls this method, some weird stuff starts to happen when we manipulate the play state.

	if deviceInfo.IsPlaying {
		if uiConfig.progressTicker != nil {
			uiConfig.progressTicker.Stop()
		}
		playingState = playingText
		uiConfig.progressTicker = time.NewTicker(time.Millisecond * 1000)
		uiConfig.timeElapsedFromTickerStart = 0
		go func() {
			for _ = range uiConfig.progressTicker.C {
				uiConfig.timeElapsedFromTickerStart += 1000

				if uiConfig.timeElapsedFromTickerStart+int(deviceInfo.ProgressMs) > int(trackInfo.DurationMs) {
					uiConfig.progressTicker.Stop()
					updateCurrentlyPlayingUI(uiConfig)
				}
				// TODO: Calculate progress into song based on needed variables. Use uiConfig
				// where necessary.
				progressInSeconds := (uiConfig.timeElapsedFromTickerStart + int(deviceInfo.ProgressMs)) / 1000
				updateTrackProgressUI(uiConfig, progressInSeconds)
			}
		}()
		// Skipping or going back a track always plays the track as well, so we will
		// only reach this else if a pause choice is chosen.
	} else {
		playingState = pausedText
		if uiConfig.progressTicker != nil {
			uiConfig.progressTicker.Stop()
		}
	}

	currentlyPlayingUI := tui.NewList()
	currentlyPlayingUI.Border = true
	currentlyPlayingUI.BorderLabel = currentlyPlayingBorderLabel
	currentlyPlayingUI.BorderFg = tui.ColorMagenta
	currentlyPlayingUI.Height = 10
	currentlyPlayingUI.Items = []string{
		NewLine,
		deviceInfo.DeviceType + " - " + deviceInfo.Name,
		NewLine + NewLine,
		trackInfo.Name,
		trackInfo.Artists,
		NewLine + NewLine,
		playingState,
	}
	currentlyPlayingUI.ItemFgColor = tui.ColorYellow

	return currentlyPlayingUI
}

func attachPlaybackComponentHandlers(uiConfig *Config) {
	playbackChoices := NewPlaybackComponent().view.Choices

	// Unfortunately, these have to be hardcoded. Handle() breaks when trying to
	// attach in a loop.
	tui.Handle("sys/kbd/1", func(e tui.Event) {
		req := playbackChoices[0].CreateAPIRequest(uiConfig.AccessToken)
		res := playbackChoices[0].SendAPIRequest(req)

		if res.StatusCode == 204 {
			// This is kind of hacky, but  wee need to wait here to give Spotify
			// playback information time to update.
			time.Sleep(250 * time.Millisecond)
			updateCurrentlyPlayingUI(uiConfig)
		}
	})

	tui.Handle("sys/kbd/2", func(e tui.Event) {
		req := playbackChoices[1].CreateAPIRequest(uiConfig.AccessToken)
		res := playbackChoices[1].SendAPIRequest(req)

		if res.StatusCode == 204 {
			// This is kind of hacky, but  wee need to wait here to give Spotify
			// playback information time to update.
			time.Sleep(250 * time.Millisecond)
			updateCurrentlyPlayingUI(uiConfig)
		}
	})

	tui.Handle("sys/kbd/3", func(e tui.Event) {
		req := playbackChoices[2].CreateAPIRequest(uiConfig.AccessToken)
		res := playbackChoices[2].SendAPIRequest(req)
		// Successful skips/backs return a 204 (no content)
		if res.StatusCode == 204 {
			// This is kind of hacky, but  wee need to wait here to give Spotify
			// playback information time to update.
			time.Sleep(250 * time.Millisecond)
			updateCurrentlyPlayingUI(uiConfig)
		}
	})

	tui.Handle("sys/kbd/4", func(e tui.Event) {
		req := playbackChoices[3].CreateAPIRequest(uiConfig.AccessToken)
		res := playbackChoices[3].SendAPIRequest(req)
		// Successful skips/backs return a 204 (no content)
		if res.StatusCode == 204 {
			// This is kind of hacky, but  wee need to wait here to give Spotify
			// playback information time to update.
			time.Sleep(250 * time.Millisecond)
			updateCurrentlyPlayingUI(uiConfig)
		}
	})
}

func getCurrentlyPlayingContext(uiConfig *Config) map[string]interface{} {
	var jsonMap map[string]interface{}
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

	return jsonMap
}

func getDeviceInformationFromJSON(context map[string]interface{}) Device {
	deviceName := context["device"].(map[string]interface{})["name"].(string)
	deviceType := context["device"].(map[string]interface{})["type"].(string)
	progressMs := context["progress_ms"].(float64)
	isPlaying := context["is_playing"].(bool)

	return Device{
		Name:       deviceName,
		DeviceType: deviceType,
		ProgressMs: progressMs,
		IsPlaying:  isPlaying,
	}
}

func getTrackInformationFromJSON(context map[string]interface{}) Track {
	trackArtists := ""
	trackName := context["item"].(map[string]interface{})["name"].(string)
	durationMs := context["item"].(map[string]interface{})["duration_ms"].(float64)
	artistJSONArray := context["item"].(map[string]interface{})["artists"].([]interface{})

	for i, artist := range artistJSONArray {
		trackArtists += artist.(map[string]interface{})["name"].(string)
		if i != len(artistJSONArray)-1 {
			trackArtists += ", "
		}
	}

	return Track{
		Name:       trackName,
		Artists:    trackArtists,
		DurationMs: durationMs,
	}
}

func updateCurrentlyPlayingUI(uiConfig *Config) {
	currentContext := getCurrentlyPlayingContext(uiConfig)
	currentTrack := getTrackInformationFromJSON(currentContext)
	deviceInfo := getDeviceInformationFromJSON(currentContext)
	newCurrentlyPlayingUI := createCurrentlyPlayingUI(uiConfig, currentTrack, deviceInfo)

	// Currently Playing box is row 1, column 2
	tui.Body.Rows[0].Cols[1] = tui.NewCol(5, 0, newCurrentlyPlayingUI)
	tui.Body.Align()
	tui.Render(tui.Body)
}

func createTrackProgressUI(uiConfig *Config, progress int) *tui.Par {
	timeString := fmt.Sprintf("%d:%.2d", progress/60, progress%60)
	progressUI := tui.NewPar(timeString)
	progressUI.Height = 3
	progressUI.Border = true
	progressUI.BorderFg = tui.ColorMagenta
	progressUI.BorderLabel = "Progress"

	return progressUI
}

func updateTrackProgressUI(uiConfig *Config, progress int) {
	newProgressUI := createTrackProgressUI(uiConfig, progress)

	tui.Body.Rows[1].Cols[0] = tui.NewCol(3, 0, newProgressUI)
	tui.Body.Align()
	tui.Render(tui.Body)
}
