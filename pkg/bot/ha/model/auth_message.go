package model

type AuthMessage struct {
	Message
	AccessToken string `json:"access_token"`
}

func NewAuthMessage(accessToken string) AuthMessage {
	return AuthMessage{
		Message:     Message{Type: AuthMsgType},
		AccessToken: accessToken,
	}
}
