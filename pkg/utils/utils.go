package utils

import (
	"bufio"
	"database/sql"
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
func AnsiToggle(db *sql.DB, username string) string {
	err := users.ToggleAnsi(db, username)
	if err != nil {
		return "Error toggling ANSI mode."
	} else {
		return "ANSI mode toggled successfully."
	}
}
