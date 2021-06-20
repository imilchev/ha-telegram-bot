package model

import "fmt"

func ExpectMessageType(rawMsg []byte, expectedType MsgType) error {
	msgType, err := GetMessageType(rawMsg)
	if err != nil {
		return err
	}

	if msgType != expectedType {
		return fmt.Errorf("Expected message type %q, received %q instead.", expectedType, msgType)
	}
	return nil
}
