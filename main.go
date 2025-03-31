package main

import (
	"log"

	"github.com/lyricalsoul/crystalline/seedlink"
)

func test() {
	client := seedlink.SeedLinkConnection{
		URL: "seisrequest.iag.usp.br:18000",
	}

	if err := client.Connect(); err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer client.Disconnect()

	client.SetupConnection()
	// start receiving messages on the channel in the main thread
	for {
		select {
		case msg, ok := <-client.Messages:
			if !ok {
				log.Println("Channel closed")
				return
			}
			switch msg.(type) {
			case seedlink.HelloMessage:
				helloMsg := msg.(seedlink.HelloMessage)
				log.Printf("Connected to server %s. Ran by %s", helloMsg.ClientName, helloMsg.Institution)
			case seedlink.OKMessage:
				log.Println("Received OK message")
			case seedlink.ErrorMessage:
				errorMsg := msg.(seedlink.ErrorMessage)
				log.Printf("Received an error with type %s: %s", errorMsg.Type, errorMsg.Message)
			default:
				log.Printf("Received unknown message type: %v", msg)
			}
		}
	}
}
