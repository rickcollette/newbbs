package main

import (
	"github.com/iammegalith/tvbbs/pkg/config"
	"github.com/iammegalith/tvbbs/pkg/database"
	"github.com/iammegalith/tvbbs/pkg/telnetbbs"
	"fmt"
	"os"
)

func main() {
	// Load the configuration
	cfg, err := config.LoadConfig("configs/bbs.ini")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
		os.Exit(1)
	}

	// Initialize the database 
	db, err := database.NewDatabase(cfg.Paths.DBPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error initializing database: %v\n", err)
		os.Exit(1)
	}
	defer db.Close()

	// Initialize and run the Telnet BBS server
	bbs := telnetbbs.NewTelnetBBS(":2323")
	if err := bbs.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error running Telnet BBS server: %v\n", err)
		os.Exit(1)
	}
}
