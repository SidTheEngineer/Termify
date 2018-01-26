package playbackUI

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	tui "github.com/gizak/termui"
)

const (
	progressGuageWidth  = 12
	progressGuageHeight = 3
	trackFeaturesURL    = "https://api.spotify.com/v1/audio-features/"
	trackFeaturesMethod = "GET"
)

// Track represents track information that we want from the Spotify track object.
type Track struct {
	Name, Artists string
	DurationMs    float64
	ID            string
	BPM           float64
}

// TODO: This still has that bug where no track is in the player. "interface{} is type
// nil and not map[string]interface{}
func getTrackInformationFromJSON(uiConfig *Config, context map[string]interface{}) Track {
	var jsonMap map[string]interface{}
	// This state can be reached when there is no context information returned from Spotify's end.
	if context["item"] == nil {
		return Track{
			Name:       "NO CURRENT TRACK INFORMATION, START SPOTIFY.",
			Artists:    "",
			DurationMs: 0.0,
		}
	}

	trackArtists := ""
	trackName := context["item"].(map[string]interface{})["name"].(string)
	durationMs := context["item"].(map[string]interface{})["duration_ms"].(float64)
	trackID := context["item"].(map[string]interface{})["id"].(string)
	artistJSONArray := context["item"].(map[string]interface{})["artists"].([]interface{})

	for i, artist := range artistJSONArray {
		trackArtists += artist.(map[string]interface{})["name"].(string)
		if i != len(artistJSONArray)-1 {
			trackArtists += ", "
		}
	}

	client := http.Client{}
	req, _ := http.NewRequest(trackFeaturesMethod, trackFeaturesURL+trackID, nil)

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", uiConfig.AccessToken.Token))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Accept", "application/json")

	resp, _ := client.Do(req)

	bytes, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal(bytes, &jsonMap)

	BPM := jsonMap["tempo"].(float64)

	currentTrack := Track{
		Name:       trackName,
		Artists:    trackArtists,
		DurationMs: durationMs,
		ID:         trackID,
		BPM:        BPM,
	}

	uiConfig.currentTrack = currentTrack
	return currentTrack
}

func createTrackProgressGauge(uiConfig *Config, progress int) *tui.Gauge {
	trackDurationMs := getTrackInformationFromJSON(uiConfig, uiConfig.context).DurationMs
	progressGuage := tui.NewGauge()
	progressGuage.Height = 3
	progressGuage.BarColor = themeProgressGuageColor
	progressGuage.BorderFg = themeBorderFg
	progressGuage.PercentColor = themePercentColor
	progressGuage.PercentColorHighlighted = themePercentColorHighlighted
	progressGuage.BorderLabel = "Progress"
	progressGuage.Percent = int((float64(progress*1000) / trackDurationMs) * 100)
	progressGuage.Label = createTrackProgressTimeString(uiConfig, progress)

	return progressGuage
}

func createTrackProgressTimeString(uiConfig *Config, progress int) string {
	trackDurationMs := getTrackInformationFromJSON(uiConfig, uiConfig.context).DurationMs
	trackDurationSecs := int(trackDurationMs / 1000)
	trackDurationMins := trackDurationSecs / 60
	trackDurationRemaining := trackDurationSecs % 60

	timeString := fmt.Sprintf("%d:%.2d/%d:%.2d", progress/60, progress%60, trackDurationMins, trackDurationRemaining)

	return timeString
}

func updateTrackProgressGauge(uiConfig *Config, progress int) {
	newProgressGuage := createTrackProgressGauge(uiConfig, progress)

	tui.Body.Rows[1].Cols[0] = tui.NewCol(progressGuageWidth, 0, newProgressGuage)
	tui.Body.Align()
	tui.Render(tui.Body)
}
