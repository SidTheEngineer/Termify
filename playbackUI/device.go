package playbackUI

// Device represents device information that we want from the Sptify device object.
type Device struct {
	Name, DeviceType string
	ProgressMs       float64
	IsPlaying        bool
	Volume           float64
}

func getDeviceInformationFromJSON(uiConfig *Config, context map[string]interface{}) Device {
	// TODO: This needs to be fixed, it will break the UI when switching songs/updating.
	// if context["device"] == nil {
	// 	return Device{
	// 		Name:       "No device currently active",
	// 		DeviceType: "N/A",
	// 		ProgressMs: 0.0,
	// 		IsPlaying:  false,
	// 		Volume:     0.0,
	// 	}
	// }
	deviceName := context["device"].(map[string]interface{})["name"].(string)
	deviceType := context["device"].(map[string]interface{})["type"].(string)
	progressMs := context["progress_ms"].(float64)
	isPlaying := context["is_playing"].(bool)

	deviceVolume := 0.0 // default when nil from Spotify
	spotifyDeviceVolume := context["device"].(map[string]interface{})["volume_percent"]
	if spotifyDeviceVolume != nil {
		deviceVolume = spotifyDeviceVolume.(float64)
	}

	currentDevice := Device{
		Name:       deviceName,
		DeviceType: deviceType,
		ProgressMs: progressMs,
		IsPlaying:  isPlaying,
		Volume:     deviceVolume,
	}

	uiConfig.currentDevice = currentDevice
	return currentDevice
}
