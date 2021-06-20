package model

type ResultMessage struct {
	Message
	Id      int  `json:"id"`
	Success bool `json:"success"`
}

func NewResultMessage() *ResultMessage {
	return &ResultMessage{Message: Message{Type: ResultMsgType}}
}
