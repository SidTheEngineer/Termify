package playbackUI

import (
	"fmt"

	tui "github.com/gizak/termui"
)

const (
	progressTimeHeight  = 3
	progressTimeWidth   = 3
	progressGuageWidth  = 9
	progressGuageHeight = 3
)

// Track represents track information that we want from the Spotify track object.
type Track struct {
	Name, Artists string
	DurationMs    float64
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

func createTrackProgressTime(uiConfig *Config, progress int) *tui.Par {
	trackDurationMs := getTrackInformationFromJSON(uiConfig.context).DurationMs
	trackDurationSecs := int(trackDurationMs / 1000)
	trackDurationMins := trackDurationSecs / 60
	trackDurationRemaining := trackDurationSecs % 60

	timeString := fmt.Sprintf("%d:%.2d/%d:%.2d", progress/60, progress%60, trackDurationMins, trackDurationRemaining)
	progressTime := tui.NewPar(timeString)
	progressTime.Height = progressTimeHeight
	progressTime.Border = true
	progressTime.BorderFg = themeBorderFg
	progressTime.TextFgColor = themeTextFgColor
	progressTime.BorderLabel = "Progress"

	return progressTime
}

func updatetrackProgressTime(uiConfig *Config, progress int) {
	newProgressTime := createTrackProgressTime(uiConfig, progress)
	newProgressGuage := createTrackProgressGuage(uiConfig, progress)

	tui.Body.Rows[1].Cols[0] = tui.NewCol(progressTimeWidth, 0, newProgressTime)
	tui.Body.Rows[1].Cols[1] = tui.NewCol(progressGuageWidth, 0, newProgressGuage)
	tui.Body.Align()
	tui.Render(tui.Body)
}

func createTrackProgressGuage(uiConfig *Config, progress int) *tui.Gauge {
	trackDurationMs := getTrackInformationFromJSON(uiConfig.context).DurationMs
	progressGuage := tui.NewGauge()
	progressGuage.Height = 3
	progressGuage.BarColor = themeProgressGuageColor
	progressGuage.BorderFg = themeBorderFg
	progressGuage.PercentColor = themePercentColor
	progressGuage.PercentColorHighlighted = themePercentColorHighlighted
	progressGuage.Percent = int((float64(progress*1000) / trackDurationMs) * 100)

	return progressGuage
}
