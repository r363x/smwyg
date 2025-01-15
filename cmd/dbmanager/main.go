package main

import (
	"fmt"
	"log"
	"os"

	"github.com/r363x/dbmanager/internal/config"
	"github.com/r363x/dbmanager/internal/db"
	"github.com/r363x/dbmanager/internal/tui"
	tea "github.com/charmbracelet/bubbletea"
)

func setupLogging(cfg config.LoggingConfig) *os.File {
    f, err := tea.LogToFile(cfg.LogFile, cfg.LogLevel)
    if err != nil {
        fmt.Println("fatal:", err)
        os.Exit(1)
    }
    return f
}

func main() {
    var dbManager db.Manager

	// Load configuration
	cfg, err := config.Load()

	if err != nil {
		log.Printf("Failed to load configuration: %v. Using interactive menu.", err)
        logFile := setupLogging(config.LoggingConfig{
            LogFile: "debug.log",
            LogLevel: "debug",
        })
        defer logFile.Close()

    } else {
        // Initialize logging
        logFile := setupLogging(cfg.LoggingConfig)
        defer logFile.Close()

        // Initialize database manager
        dbManager, err = db.NewManager(cfg.DatabaseConfig)
        if err != nil {
            log.Fatalf("Failed to initialize database manager: %v", err)
        }

    }

    // Initialize TUI
    ui, err := tui.New(dbManager)
    if err != nil {
        log.Fatalf("Failed to initialize TUI: %v", err)
    }

	// Run the application
	if _, err := ui.Run();
    err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
