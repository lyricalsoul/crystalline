package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/lyricalsoul/crystalline/seedlink"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run cli.go server:port")
		return
	}

	server := os.Args[1]
	splitted := strings.Split(server, ":")
	if len(splitted) != 2 {
		fmt.Println("Usage: go run cli.go server:port")
		return
	}

	fmt.Printf("Connecting to server %s on port %s...\n", splitted[0], splitted[1])
	client := seedlink.SeedLinkConnection{
		URL:                   server,
		PrintIncomingMessages: true,
	}

	if err := client.Connect(); err != nil {
		fmt.Println("Error connecting to server:", err)
		return
	}

	defer client.Disconnect()
	fmt.Println("Crystalline Interactive Shell. Type 'exit' to quit.")

	fmt.Print("> ")
	scanner := bufio.NewScanner(os.Stdin)
	for {
		if !scanner.Scan() {
			break
		}
		input := strings.ToUpper(scanner.Text())

		if strings.TrimSpace(input) == "EXIT" {
			fmt.Println("Exiting...")
			client.SendRaw("BYE")
			os.Exit(0)
		}

		if err := client.SendRaw(strings.TrimSpace(input)); err != nil {
			fmt.Println("Error:", err)
		}
	}
}
