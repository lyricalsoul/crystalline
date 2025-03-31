package seedlink

import (
	"bufio"
	"fmt"
	"strings"
)

// The commands we can send to the server.
type SeedLinkCommand string

const (
	// HELLO command.
	HELLO_COMMAND SeedLinkCommand = "HELLO"
	// USERAGENT command.
	USERAGENT_COMMAND SeedLinkCommand = "USERAGENT"
)

// What the server responds to the HELLO message.
type HelloMessage struct {
	// The name of the server. Something like "SeedLink v3.3 (2024.020)"
	ClientName string
	// The institution running the server. Something like "Centro de Sismologia da USP"
	Institution string
}

// The OK message signals the last operation was successful.
type OKMessage struct{}

// The error types that an ERROR message can contain.
type ErrorType string

const (
	// Command wasn't recognized or not supported.
	UNSUPPORTED_ERROR ErrorType = "UNSUPPORTED"
	// Unexpected command.
	UNEXPECTED_ERROR ErrorType = "UNEXPECTED"
	// Client isn't authorized to send the command. Perhaps log in?
	UNAUTHORIZED_ERROR ErrorType = "UNAUTHORIZED"
	// The client exceeded a limit - e.g., too many STATION or SELECT commands.
	LIMIT_ERROR ErrorType = "LIMIT"
	// Incorrect arguments were sent to the command.
	ARGUMENTS_ERROR ErrorType = "ARGUMENTS"
	// Authentication failed. Maybe wrong user/pass, token...
	AUTH_ERROR ErrorType = "AUTH"
	// An internal error occurred on the server.
	INTERNAL_ERROR ErrorType = "INTERNAL"
)

// The ERROR message signals the last operation failed. It contains the error type and, optionally, a message.
type ErrorMessage struct {
	// The type of error.
	Type ErrorType
	// The message of the error.
	Message string
}

func NewHelloMessage(reader bufio.Reader) (HelloMessage, error) {
	// receive the messages from the server
	client, err := reader.ReadString('\n')
	institution, err2 := reader.ReadString('\n')
	if err != nil || err2 != nil {
		if err != nil {
			return HelloMessage{}, err
		}

		return HelloMessage{}, err2
	}

	msg := HelloMessage{
		// trim the new line characters
		ClientName:  strings.TrimSpace(client),
		Institution: strings.TrimSpace(institution),
	}

	return msg, nil
}

func NewResultMessage(reader bufio.Reader) (any, error) {
	// receive the messages from the server
	result, err := reader.ReadString('\n')
	if err != nil {
		return nil, err
	}

	// if the message starts with ok, we return an OK message
	if strings.HasPrefix(result, "OK") {
		return OKMessage{}, nil
	}
	// we split the message by spaces and remove the first two elements - the ERROR part and the type.
	// the rest is the message
	splitted := strings.Split(result, " ")
	if len(splitted) < 2 {
		return nil, nil
	}

	// we remove the first two elements
	first_part := splitted[0]
	if first_part != "ERROR" {
		fmt.Printf("Expected ERROR, got %s", first_part)
		fmt.Println("Full message:", result)
		return nil, nil
	}

	// The type.
	error_type := splitted[1]
	// Map the type to the ErrorType enum
	var error_type_enum ErrorType
	switch error_type {
	case "UNSUPPORTED":
		error_type_enum = UNSUPPORTED_ERROR
	case "UNEXPECTED":
		error_type_enum = UNEXPECTED_ERROR
	case "UNAUTHORIZED":
		error_type_enum = UNAUTHORIZED_ERROR
	case "LIMIT":
		error_type_enum = LIMIT_ERROR
	case "ARGUMENTS":
		error_type_enum = ARGUMENTS_ERROR
	case "AUTH":
		error_type_enum = AUTH_ERROR
	case "INTERNAL":
		error_type_enum = INTERNAL_ERROR
	default:
		fmt.Printf("Unknown error type: %s", error_type)
		return nil, nil
	}

	// The message.
	error_message := strings.Join(splitted[2:], " ")

	return ErrorMessage{
		Type:    error_type_enum,
		Message: error_message,
	}, nil
}
