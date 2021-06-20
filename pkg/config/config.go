package config

type Config struct {
	HomeAssistant HomeAssistantConfig `json:"homeAssistant"`
}

type HomeAssistantConfig struct {
	Url         string `json:"url"`
	AccessToken string `json:"accessToken"`
}
