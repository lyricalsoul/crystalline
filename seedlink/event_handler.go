package seedlink

import (
	"bufio"
	"fmt"
)

// Launches the parser goroutine.
func (s *SeedLinkConnection) StartEventHandler() {
	// Start a goroutine to handle events
	go func() {
		reader := bufio.NewReader(s.Conn)

		for {
			select {
			case <-s.Stop:
				return
			default:
			}

			// We get the last message type and switch on it
			if s.LastMessageType == nil {
				continue
			}

			switch string(*s.LastMessageType) {
			case string(HELLO_COMMAND):
				// Handle HELLO message
				helloMsg, err := NewHelloMessage(*reader)
				if err != nil {
					fmt.Println("Error reading hello message:", err)
					continue
				}

				s.Messages <- helloMsg
			default:
				// Must be an OK or an ERROR!
				resultMsg, err := NewResultMessage(*reader)
				if err != nil {
					fmt.Println("Error reading result message:", err)
					continue
				}

				if resultMsg == nil {
					fmt.Println("Received nil message")
					continue
				}

				s.Messages <- resultMsg
			}
			// We set the last message type to nil, since we must have consumed it by now
			s.LastMessageType = nil
		}
	}()
}

// A handler that just prints the messages to stdout
func (s *SeedLinkConnection) StartStdoutHandler() {
	go func() {
		reader := bufio.NewReader(s.Conn)
		for {
			select {
			case <-s.Stop:
				return
			default:
			}

			test, err := reader.ReadString('\n')
			if err != nil {
				fmt.Println("Error reading message:", err)
				return
			}

			fmt.Println(test)
			fmt.Print("> ")
		}
	}()
}
