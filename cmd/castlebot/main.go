package main

import (
	"log"
	"os"

	"github.com/castle-coders/homeassistant-keybase-bot/internal/homeassistant"
	"github.com/castle-coders/homeassistant-keybase-bot/internal/keybase"
)

func main() {
	haStatusChan := make(chan int)

	kbClient, err := keybase.NewClient()
	if err != nil {
		log.Fatal(err)
		return
	}

	haClientConfig := homeassistant.HAClientConfig{
		WebsocketEndpoint: os.Getenv("HA_WS_ENDPOINT"),
		AccessToken:       os.Getenv("HA_ACCESS_TOKEN"),
	}

	hac := homeassistant.NewClient(haClientConfig, haStatusChan, kbClient)
	hac.Start()
}
