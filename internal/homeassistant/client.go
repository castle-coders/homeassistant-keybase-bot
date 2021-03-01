package homeassistant

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"

	status "github.com/castle-coders/homeassistant-keybase-bot/internal/common"
	"github.com/castle-coders/homeassistant-keybase-bot/internal/keybase"

	"net/http"

	"github.com/gorilla/websocket"
)

const (
	authRequired = 0
	authOK       = 1
	subscribed   = 2

	// RX
	typeAuthRequired = "auth_required"
	typeAuthOK       = "auth_ok"
	typeAuthInvalid  = "auth_invalid"
	typeResult       = "result"
	typeEvent        = "event"

	//TX
	typeAuth            = "auth"
	typeSubscribeEvents = "subscribe_events"

	// EVENT TYPES
	typeEventNotifyKeybase = "notify_keybase"
)

// HAClientConfig contains client configuration
type HAClientConfig struct {
	WebsocketEndpoint string
	AccessToken       string
}

// NewClient creates a new client impl
func NewClient(haClientConfig HAClientConfig, statusChan chan int, kbClient *keybase.Client) *Client {
	c := Client{
		msgCtr:         -1,
		conn:           nil,
		haClientConfig: haClientConfig,
		statusChan:     statusChan,
		kbClient:       kbClient,
	}
	return &c
}

// Client is an implementation of homeassistant Client
type Client struct {
	msgCtr         int
	conn           *websocket.Conn
	haClientConfig HAClientConfig
	statusChan     chan int
	kbClient       *keybase.Client
}

func (c Client) fail(msg string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, msg+"\n", args...)
	c.statusChan <- status.Error
	close(c.statusChan)
}

func (c Client) handleNextMessage() (message []byte, messageType string, err error) {
	if c.conn == nil {
		return nil, "", errors.New("not ready")
	}
	msgType, data, err := c.conn.ReadMessage()
	if err != nil {
		return nil, "", err
	}
	if msgType != websocket.TextMessage {
		return nil, "", errors.New("unexpected message type in websocket response")
	}
	var envelope HAEnvelope
	err = json.Unmarshal(data, &envelope)
	if err != nil {
		return nil, "", err
	}
	return data, envelope.Type, nil
}

func (c Client) auth() error {
	return c.conn.WriteJSON(AuthMessage{typeAuth, c.haClientConfig.AccessToken})
}

func (c Client) subscribe() error {
	c.msgCtr++
	return c.conn.WriteJSON(SubscribeMessage{c.msgCtr, typeSubscribeEvents, typeEventNotifyKeybase})
}

func (c Client) checkSub(message []byte) error {
	var resultMessage ResultMessage
	err := json.Unmarshal(message, &resultMessage)
	if err != nil {
		return err
	}
	if !resultMessage.Success {
		return fmt.Errorf("subscription FAILED: %s (code %s)", resultMessage.Error.Message, resultMessage.Error.Code)
	}
	return nil
}

func (c Client) handleEvent(incommingMessage []byte) error {
	var eventMessage EventMessage
	err := json.Unmarshal(incommingMessage, &eventMessage)
	if err != nil {
		return err
	}

	messageToSend := *&eventMessage.Event.Data.Message
	if eventMessage.Event.Data.User != nil {
		user := *eventMessage.Event.Data.User
		err = c.kbClient.SendMessageToUser(user, messageToSend)
	}
	if err == nil && eventMessage.Event.Data.Team != nil {
		team := *eventMessage.Event.Data.Team
		channel := eventMessage.Event.Data.Channel
		err = c.kbClient.SendMessageToTeam(team, channel, messageToSend)
	}
	return err
}

// Start to run homeassistant client
func (c *Client) Start() {
	log.Println("starting up ha client")
	headers := http.Header{}
	conn, _, err := websocket.DefaultDialer.Dial(c.haClientConfig.WebsocketEndpoint, headers)
	if err != nil {
		c.fail("Error connecting to API: %s", err.Error())
		return
	}
	c.conn = conn

	log.Println("connected")

	for {
		message, messageType, err := c.handleNextMessage()
		if err != nil {
			c.fail("Error rcv: %s", err)
		}
		switch messageType {
		case typeAuthRequired:
			log.Printf("got auth required message")
			err = c.auth()
		case typeAuthOK:
			log.Printf("got auth OK message")
			c.msgCtr = 0
			// auth successful, next subcribe to events
			err = c.subscribe()
		case typeAuthInvalid:
			err = errors.New("invalid auth")
		case typeResult:
			err = c.checkSub(message)
		case typeEvent:
			err = c.handleEvent(message)
		default:
			err = errors.New("unexpected ha message type")
		}
		if err != nil {
			c.fail("Error sending: %s", err)
			return
		}
	}

}
