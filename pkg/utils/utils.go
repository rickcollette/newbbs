package utils

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"github.com/iammegalith/tvbbs/pkg/users"
)

// ShowText reads a file specified by the path and sends its content to the connection
func ShowText(conn net.Conn, path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		fmt.Fprintln(conn, scanner.Text())
	}

	return scanner.Err()
}

// AnsiToggle toggles the ANSI mode for the user and sends a confirmation message to the connection
func AnsiToggle(conn net.Conn, username string) {
	err := users.ToggleAnsi(username)
	if err != nil {
		fmt.Fprintln(conn, "Error toggling ANSI mode.")
	} else {
		fmt.Fprintln(conn, "ANSI mode toggled successfully.")
	}
}