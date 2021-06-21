package config

type Config struct {
	HomeAssistant HomeAssistantConfig `json:"homeAssistant"`
	Telegram      TelegramConfig      `json:"telegram"`
}

type HomeAssistantConfig struct {
	Url                  string   `json:"url"`
	AccessToken          string   `json:"accessToken"`
	TemperatureEntityIds []string `json:"temperatureEntityIds"`
	HumidityEntityIds    []string `json:"humidityEntityIds"`
}

type TelegramConfig struct {
	Token        string   `json:"token"`
	AllowedUsers []string `json:"allowedUsers"`
}
