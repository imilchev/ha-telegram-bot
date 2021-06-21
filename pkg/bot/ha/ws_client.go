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
	TempChan  chan model.EventMessage
	HumidChan chan model.EventMessage
	url       *url.URL
	config    config.HomeAssistantConfig
	ws        *websocket.Conn
}

func NewWsClient(config config.HomeAssistantConfig) (*WsClient, error) {
	url, err := url.Parse(config.Url)
	if err != nil {
		zap.S().Errorf("Failed to parse HA websocket URL. %+v", err)
		return nil, err
	}

	c, _, err := websocket.DefaultDialer.Dial(url.String(), nil)
	if err != nil {
		zap.S().Fatal("Failed to open websocket connection to %s. %+v", url.String(), err)
		return nil, err
	}

	return &WsClient{
		TempChan:  make(chan model.EventMessage, 100),
		HumidChan: make(chan model.EventMessage, 100),
		url:       url,
		config:    config,
		ws:        c,
	}, nil
}

func (ws *WsClient) Start() error {
	errChan := make(chan error)
	return ws.subscribeToEvents(errChan)
}

func (ws *WsClient) Stop() error {
	return ws.ws.Close()
}

func (ws *WsClient) subscribeToEvents(errChan chan error) error {
	initMsg := model.NewAuthRequiredMessage()
	if err := ws.readMessageOfType(initMsg); err != nil {
		zap.S().Error(err)
		errChan <- err
		return err
	}

	// Authenticate
	zap.S().Infof("Connected to HA version %s. Authenticating...", initMsg.HomeAssistantVersion)
	if err := ws.ws.WriteJSON(model.NewAuthMessage(ws.config.AccessToken)); err != nil {
		zap.S().Errorf("Failed to send authentication. %+v", err)
		errChan <- err
		return err
	}

	// Verify authentication
	authOkMsg := &model.Message{Type: model.AuthOkMsgType}
	if err := ws.readMessageOfType(authOkMsg); err != nil {
		zap.S().Error(err)
		errChan <- err
		return err
	}
	zap.S().Infof("Successfully authenticated to HA.")

	// Subscribe to state_changed events
	if err := ws.ws.WriteJSON(model.NewSubscribeEventsMessage(1, model.StateChangedEventType)); err != nil {
		zap.S().Errorf("Failed to subscribe to HA events. %+v", err)
		errChan <- err
		return err
	}

	// Verify subscription
	resMsg := model.NewResultMessage()
	if err := ws.readMessageOfType(resMsg); err != nil {
		zap.S().Error(err)
		errChan <- err
		return err
	}
	if !resMsg.Success {
		err := fmt.Errorf("Subscription to events failed. %+v", resMsg)
		errChan <- err
		return err
	}
	zap.S().Info("Successfully subscribed to HA events.")

	for {
		eventMsg := model.NewEventMessage()
		if err := ws.readMessageOfType(eventMsg); err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				zap.S().Error(err)
				errChan <- err
				return err
			}
			zap.S().Info("HA client shut down.")
			close(ws.TempChan)
			close(ws.HumidChan)
			return nil
		}
		ws.filterMessage(*eventMsg)
	}
}

func (ws *WsClient) readMessageOfType(msg model.HAMessage) error {
	_, msgRaw, err := ws.ws.ReadMessage()
	if err != nil {
		return fmt.Errorf("Error receiving message from HA. %+v", err)
	}

	if err := model.ExpectMessageType(msgRaw, msg.GetType()); err != nil {
		return err
	}
	return json.Unmarshal(msgRaw, msg)
}

func (ws *WsClient) filterMessage(msg model.EventMessage) {
	for _, t := range ws.config.TemperatureEntityIds {
		if msg.Event.Data.EntityId == t {
			ws.TempChan <- msg
			zap.S().Debugf("Message written to temperature chan %+v.", msg)
			return
		}
	}

	for _, h := range ws.config.HumidityEntityIds {
		if msg.Event.Data.EntityId == h {
			ws.HumidChan <- msg
			zap.S().Debugf("Message written to humidity chan %+v.", msg)
			return
		}
	}
	zap.S().Debugf("Message for entity %q has been ignored.", msg.Event.Data.EntityId)
}
