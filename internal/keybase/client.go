package keybase

import (
	"flag"
	"fmt"

	"github.com/keybase/go-keybase-chat-bot/kbchat"
)

// NewClient creates a new client impl
func NewClient() (*Client, error) {
	c := Client{
		kbc: nil,
	}

	var kbLoc string
	var err error

	flag.StringVar(&kbLoc, "keybase", "keybase", "the location of the Keybase app")
	flag.Parse()

	if c.kbc, err = kbchat.Start(kbchat.RunOptions{KeybaseLocation: kbLoc}); err != nil {
		return nil, err
	}

	return &c, nil
}

// Client is an implementation of keybase Client
type Client struct {
	kbc *kbchat.API
}

// SendMessageToUser to send a message to a user
func (c Client) SendMessageToUser(user string, message string) error {
	tlfName := fmt.Sprintf("%s,%s", c.kbc.GetUsername(), user)
	_, err := c.kbc.SendMessageByTlfName(tlfName, message)
	return err
}

// SendMessageToTeam to send a message to a team
func (c Client) SendMessageToTeam(team string, channel *string, message string) error {
	_, err := c.kbc.SendMessageByTeamName(team, channel, message)
	return err
}

// Stop to stop the keybase client
func (c Client) Stop() {
	fmt.Println("gracefully shutting down keybase client")
	c.kbc.Shutdown()
}
