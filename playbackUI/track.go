package playbackUI

import (
	tui "github.com/gizak/termui"
)

const (
	// progressTimeHeight  = 3
	// progressTimeWidth   = 3
	progressGuageWidth  = 12
	progressGuageHeight = 3
)

// Track represents track information that we want from the Spotify track object.
type Track struct {
	Name, Artists string
	DurationMs    float64
}

func getTrackInformationFromJSON(uiConfig *Config, context map[string]interface{}) Track {
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

	currentTrack := Track{
		Name:       trackName,
		Artists:    trackArtists,
		DurationMs: durationMs,
	}

	uiConfig.currentTrack = currentTrack
	return currentTrack
}

func updateTrackProgressGuage(uiConfig *Config, progress int) {
	newProgressGuage := createTrackProgressGuage(uiConfig, progress)

	tui.Body.Rows[1].Cols[0] = tui.NewCol(progressGuageWidth, 0, newProgressGuage)
	tui.Body.Align()
	tui.Render(tui.Body)
}

func createTrackProgressGuage(uiConfig *Config, progress int) *tui.Gauge {
	trackDurationMs := getTrackInformationFromJSON(uiConfig, uiConfig.context).DurationMs
	progressGuage := tui.NewGauge()
	progressGuage.Height = 3
	progressGuage.BarColor = themeProgressGuageColor
	progressGuage.BorderFg = themeBorderFg
	progressGuage.PercentColor = themePercentColor
	progressGuage.PercentColorHighlighted = themePercentColorHighlighted
	progressGuage.Percent = int((float64(progress*1000) / trackDurationMs) * 100)

	return progressGuage
}
