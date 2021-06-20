package model

type AuthRequiredMessage struct {
	Message
	HomeAssistantVersion string `json:"ha_version"`
}

func NewAuthRequiredMessage() *AuthRequiredMessage {
	return &AuthRequiredMessage{Message: Message{Type: AuthRequiredMsgType}}
}
