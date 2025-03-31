package seedlink

import (
	"fmt"
	"net"
	"strings"
)

// SeedLinkConnection represents a connection to a SeedLink server.
type SeedLinkConnection struct {
	// The connection to the server
	Conn net.Conn
	// The URL of the server
	URL string
	// Whether it's running.
	IsConnected bool
	// A channel to receive messages from the server.
	Messages chan interface{}
	// A channel to inform the listener when to stop.
	Stop chan bool
	// The last message type we sent. Used so we can properly generate the response payloads.
	// Once a message is fully received and parsed, this becomes nil, signifying that we are ready to send new messages in a safe way.
	LastMessageType *SeedLinkCommand
	// Whether we should print incoming messages to stdout.
	PrintIncomingMessages bool
}

// Connect to the server. It just opens a connection to the server.
func (s *SeedLinkConnection) Connect() error {
	conn, err := net.Dial("tcp", s.URL)
	if err != nil {
		return err
	}

	s.Conn = conn
	s.IsConnected = true
	s.Messages = make(chan any, 1)
	s.Stop = make(chan bool)

	s.StartListening()
	return nil
}

// Disconnect from the server. It just closes the connection to the server.
func (s *SeedLinkConnection) Disconnect() error {
	if s.Conn != nil {
		s.IsConnected = false
		if s.Stop != nil {
			s.Stop <- true
			close(s.Stop)
		}

		return s.Conn.Close()
	}

	s.IsConnected = false
	return nil
}

func (s *SeedLinkConnection) SendMessageWithData(command SeedLinkCommand, data string) error {
	if !s.IsConnected {
		return fmt.Errorf("not connected to server")
	}

	commandWCRLF := string(command) + " " + data

	if !strings.HasSuffix(commandWCRLF, "\r\n") {
		commandWCRLF = commandWCRLF + "\r\n"
	}

	_, err := fmt.Fprintf(s.Conn, commandWCRLF)

	s.LastMessageType = &command
	return err
}

func (s *SeedLinkConnection) SendMessage(command SeedLinkCommand) error {
	if !s.IsConnected {
		return fmt.Errorf("not connected to server")
	}

	commandWCRLF := string(command)

	if !strings.HasSuffix(commandWCRLF, "\r\n") {
		commandWCRLF = commandWCRLF + "\r\n"
	}

	_, err := fmt.Fprintf(s.Conn, commandWCRLF)

	s.LastMessageType = &command
	return err
}

// Send raw data to the server.
func (s *SeedLinkConnection) SendRaw(data string) error {
	if !s.IsConnected {
		return fmt.Errorf("not connected to server")
	}

	if !strings.HasSuffix(data, "\r\n") {
		data = data + "\r\n"
	}

	_, err := fmt.Fprintf(s.Conn, data)

	return err
}

// Starts the listening process. We run it on a separate goroutine.
func (s *SeedLinkConnection) StartListening() error {
	// Check if we are connected
	if !s.IsConnected {
		return fmt.Errorf("not connected to server")
	}

	if s.PrintIncomingMessages {
		s.StartStdoutHandler()
	} else {
		s.StartEventHandler()
	}
	return nil
}

// A simple command to set up the connection.
// It automatically sends the HELLO and the USERAGENT commands.
func (s *SeedLinkConnection) SetupConnection() error {
	if err := s.SendMessage(HELLO_COMMAND); err != nil {
		return fmt.Errorf("failed to send HELLO command: %w", err)
	}

	//if err := s.SendMessageWithData(USERAGENT_COMMAND, utilities.AssembleCrystallineUserAgent()); err != nil {
	//	return fmt.Errorf("failed to send USERAGENT command: %w", err)
	//}

	return nil
}
