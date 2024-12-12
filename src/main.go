package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"contor-system/src/computing"
	"contor-system/src/utils"
)

// ensureDirectory ensures the directory exists, creating it if necessary.
func ensureDirectory(dir string) error {
	return os.MkdirAll(dir, os.ModePerm)
}

// Function to ensure we read config changes.
func loadConfig(filePath string) (utils.System, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return utils.System{}, fmt.Errorf("failed to open config.json: %v", err)
	}
	defer file.Close()

	var system utils.System
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&system); err != nil {
		return utils.System{}, fmt.Errorf("failed to decode config.json: %v", err)
	}

	return system, nil
}

func writeToFile(logs []utils.LogEntry, filePath string) error {
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()

	for _, logEntry := range logs {
		line := fmt.Sprintf("%s | %s | %s\n", logEntry.Timestamp, logEntry.ComponentID, logEntry.Message)
		if _, err := file.WriteString(line); err != nil {
			return fmt.Errorf("failed to write to file: %v", err)
		}
	}

	return nil
}

// Funcția principală
func main() {
	configPath := "./config.json"

	// Initial config load
	config, err := loadConfig(configPath)
	if err != nil {
		log.Fatalf("Failed to load initial config: %v", err)
	}

	lastConfigJSON, _ := os.ReadFile(configPath)

	// Ensure graceful shutdown on interrupt
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
	defer close(signals)

	cleanupManager := sync.Once{}
	cleanup := func() {
		log.Println("Cleaning up resources...")
	}
	defer cleanupManager.Do(cleanup)

	// Infinite loop to monitor logs and handle writes
	// ticker := time.NewTicker(1 * time.Minute)
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	// Generate the initial log file path
	logFilePath := fmt.Sprintf("logs/%s.txt", time.Now().Format("2006-01-02_150405"))

	for {
		select {
		case <-signals:
			log.Println("Received shutdown signal.")
			cleanupManager.Do(cleanup)
			log.Println("Shutdown complete.")
			return
		case <-ticker.C:
			// Reload configuration every minute and check for changes
			currentConfig, err := loadConfig(configPath)
			if err != nil {
				log.Printf("Failed to reload config: %v", err)
				continue
			}

			currentConfigJSON, _ := os.ReadFile(configPath)
			if string(currentConfigJSON) != string(lastConfigJSON) {
				log.Println("Configuration has changed. New configuration loaded.")
				lastConfigJSON = currentConfigJSON
				config = currentConfig

				// Generate a new log file path when the config changes
				logFilePath = fmt.Sprintf("logs/%s.txt", time.Now().Format("2006-01-02_150405"))
				log.Printf("Switched to new log file: %s", logFilePath)
			}

			// Simulate log calculation
			logEntries := computing.ComputeSystem(config)

			// Ensure the log directory exists
			if err := os.MkdirAll("logs", os.ModePerm); err != nil {
				log.Printf("Error creating log directory: %v", err)
				continue
			}

			// Write logs to the new log file
			if err := writeToFile(logEntries, logFilePath); err != nil {
				log.Printf("Error writing to file: %v", err)
			} else {
				log.Printf("Logged %d entries to %s", len(logEntries), logFilePath)
			}
		}
	}
}
