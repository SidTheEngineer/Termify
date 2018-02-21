package playbackUI

// Device represents device information that we want from the Sptify device object.
type Device struct {
	Name, DeviceType, ID string
	ProgressMs           float64
	IsPlaying            bool
	Volume               float64
}

func getDeviceInformationFromJSON(uiConfig *Config, context map[string]interface{}) Device {
	if context["device"] == nil {
		return Device{
			Name:       "No device currently active",
			ID:         "N/A",
			DeviceType: "N/A",
			ProgressMs: 0.0,
			IsPlaying:  false,
			Volume:     100.0,
		}
	}
	deviceName := context["device"].(map[string]interface{})["name"].(string)
	deviceType := context["device"].(map[string]interface{})["type"].(string)
	progressMs := context["progress_ms"].(float64)
	isPlaying := context["is_playing"].(bool)

	deviceVolume := 100.0 // default when nil from Spotify
	deviceID := "N/A"     // default when nil form Spotify

	spotifyDeviceVolume := context["device"].(map[string]interface{})["volume_percent"]
	if spotifyDeviceVolume != nil {
		deviceVolume = spotifyDeviceVolume.(float64)
	}

	spotifyDeviceID := context["device"].(map[string]interface{})["id"]
	if spotifyDeviceID != nil {
		deviceID = spotifyDeviceID.(string)
	}

	currentDevice := Device{
		Name:       deviceName,
		DeviceType: deviceType,
		ID:         deviceID,
		ProgressMs: progressMs,
		IsPlaying:  isPlaying,
		Volume:     deviceVolume,
	}

	uiConfig.currentDevice = currentDevice
	return currentDevice
}
