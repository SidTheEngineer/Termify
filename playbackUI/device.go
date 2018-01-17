package playbackUI

// Device represents device information that we want from the Sptify device object.
type Device struct {
	Name, DeviceType string
	ProgressMs       float64
	IsPlaying        bool
}

func getDeviceInformationFromJSON(uiConfig *Config, context map[string]interface{}) Device {
	deviceName := context["device"].(map[string]interface{})["name"].(string)
	deviceType := context["device"].(map[string]interface{})["type"].(string)
	progressMs := context["progress_ms"].(float64)
	isPlaying := context["is_playing"].(bool)

	currentDevice := Device{
		Name:       deviceName,
		DeviceType: deviceType,
		ProgressMs: progressMs,
		IsPlaying:  isPlaying,
	}

	uiConfig.currentDevice = currentDevice
	return currentDevice
}
