package ha

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/gorilla/websocket"
	"github.com/imilchev/hass-telegram-bot/pkg/bot/ha/model"
	"github.com/imilchev/hass-telegram-bot/pkg/config"
	"go.uber.org/zap"
)

type WsClient struct {
	url         *url.URL
	accessToken string
}

func NewWsClient(config config.HomeAssistantConfig) (*WsClient, error) {
	url, err := url.Parse(config.Url)
	if err != nil {
		zap.S().Errorf("Failed to parse HA websocket URL. %+v", err)
		return nil, err
	}
	return &WsClient{url: url, accessToken: config.AccessToken}, nil
}

func (ws *WsClient) Start() error {
	c, _, err := websocket.DefaultDialer.Dial(ws.url.String(), nil)
	if err != nil {
		zap.S().Fatal("Failed to open websocket connection to %s. %+v", ws.url.String(), err)
		return err
	}
	defer c.Close()

	errChan := make(chan error)
	ws.subscribeToEvents(c, errChan)
	return nil
}

func (ws *WsClient) subscribeToEvents(c *websocket.Conn, errChan chan error) {
	initMsg := model.NewAuthRequiredMessage()
	if err := readMessageOfType(c, initMsg); err != nil {
		zap.S().Error(err)
		errChan <- err
		return
	}

	// Authenticate
	zap.S().Infof("Connected to HA version %s. Authenticating...", initMsg.HomeAssistantVersion)
	if err := c.WriteJSON(model.NewAuthMessage(ws.accessToken)); err != nil {
		zap.S().Errorf("Failed to send authentication. %+v", err)
		errChan <- err
		return
	}

	// Verify authentication
	authOkMsg := &model.Message{Type: model.AuthOkMsgType}
	if err := readMessageOfType(c, authOkMsg); err != nil {
		zap.S().Error(err)
		errChan <- err
		return
	}
	zap.S().Infof("Successfully authenticated to HA.")

	// Subscribe to state_changed events
	if err := c.WriteJSON(model.NewSubscribeEventsMessage(1, model.StateChangedEventType)); err != nil {
		zap.S().Errorf("Failed to subscribe to events. %+v", err)
		errChan <- err
		return
	}

	// Verify subscription
	resMsg := model.NewResultMessage()
	if err := readMessageOfType(c, resMsg); err != nil {
		zap.S().Error(err)
		errChan <- err
		return
	}
	if !resMsg.Success {
		err := fmt.Errorf("Subscription to events failed. %+v", resMsg)
		errChan <- err
		return
	}

	for {
		eventMsg := model.NewEventMessage()
		if err := readMessageOfType(c, eventMsg); err != nil {
			zap.S().Error(err)
			errChan <- err
			return
		}
		zap.S().Debugf("recv: %+v", eventMsg) // TODO: make sure this is propagated to the necessary handler
	}
}

func readMessageOfType(c *websocket.Conn, msg model.HAMessage) error {
	_, msgRaw, err := c.ReadMessage()
	if err != nil {
		return fmt.Errorf("Error receiving message from HA. %+v", err)
	}

	if err := model.ExpectMessageType(msgRaw, msg.GetType()); err != nil {
		return err
	}
	return json.Unmarshal(msgRaw, msg)
}

// func Test() error {
// 	interrupt := make(chan os.Signal, 1)
// 	signal.Notify(interrupt, os.Interrupt)

// 	uStr := "ws://ha.home.io/api/websocket"
// 	u, err := url.Parse(uStr)
// 	if err != nil {
// 		zap.S().Errorf("Failed to parse HA websocket URL. %+v", err)
// 	}
// 	zap.S().Infof("Connecting to %s", u.String())

// 	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
// 	if err != nil {
// 		log.Fatal("dial:", err)
// 	}
// 	defer c.Close()

// 	done := make(chan struct{})

// 	go func() {
// 		defer close(done)
// 		_, msg, err := c.ReadMessage()
// 		if err != nil {
// 			zap.S().Errorf("Error receiving initial message from HA. %+v", err)
// 		}

// 		if err := model.ExpectMessageType(msg, model.AuthRequiredMsgType); err != nil {
// 			zap.S().Error(err)
// 			return
// 		}

// 		initMsg := &model.InitMessage{}
// 		if err := json.Unmarshal(msg, initMsg); err != nil {
// 			zap.S().Errorf("Failed to deserialize initial message %s. %+v", string(msg), err)
// 			return
// 		}

// 		zap.S().Infof("Connected to HA version %s. Authenticating...", initMsg.HomeAssistantVersion)
// 		if err := c.WriteJSON(model.NewAuthMessage(accessToken)); err != nil {
// 			zap.S().Errorf("Failed to send authentication. %+v", err)
// 			return
// 		}

// 		for {
// 			_, message, err := c.ReadMessage()
// 			if err != nil {
// 				zap.S().Debugf("read: %+v", err)
// 				return
// 			}
// 			zap.S().Debugf("recv: %s", message)
// 		}
// 	}()

// 	for {
// 		select {
// 		case <-done:
// 			return nil
// 		case <-interrupt:
// 			log.Println("interrupt")

// 			// Cleanly close the connection by sending a close message and then
// 			// waiting (with timeout) for the server to close the connection.
// 			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
// 			if err != nil {
// 				log.Println("write close:", err)
// 				return nil
// 			}
// 			select {
// 			case <-done:
// 			case <-time.After(time.Second):
// 			}
// 			return nil
// 		}
// 	}
// }
