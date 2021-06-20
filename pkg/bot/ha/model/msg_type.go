package model

import (
	"encoding/json"

	"go.uber.org/zap"
)

type MsgType string

const (
	AuthRequiredMsgType    MsgType = "auth_required"
	AuthMsgType            MsgType = "auth"
	AuthOkMsgType          MsgType = "auth_ok"
	SubscribeEventsMsgType MsgType = "subscribe_events"
	ResultMsgType          MsgType = "result"
	EventMsgType           MsgType = "event"
)

func GetMessageType(rawMsg []byte) (MsgType, error) {
	msg := &Message{}
	if err := json.Unmarshal(rawMsg, msg); err != nil {
		zap.S().Errorf("Failed to deserialize message %s. %+v", string(rawMsg), err)
		return "", err
	}
	return msg.Type, nil
}
