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

	"github.com/xitongsys/parquet-go/writer"
)

// ensureDirectory ensures the directory exists, creating it if necessary.
func ensureDirectory(dir string) error {
	return os.MkdirAll(dir, os.ModePerm)
}

// Funcția principală
func main() {
	// Deschide fișierul config.json
	file, configErr := os.Open("./config.json")

	if configErr != nil {
		log.Fatalf("Failed to open config.json: %v", configErr)
	}

	defer file.Close()

	// Decodează JSON-ul
	var system utils.System

	decoder := json.NewDecoder(file)

	if decodeError := decoder.Decode(&system); decodeError != nil {
		log.Fatalf("Failed to decode JSON: %v", decodeError)
	}

	// Open a Parquet file named with the current date\
	fileName := fmt.Sprintf("logs/%s", time.Now().Format("2006-01-02")+".parquet")

	fileNameError := ensureDirectory("logs")

	if fileNameError != nil {
		log.Fatalf("Failed to create logs directory: %v", fileNameError)
	}

	fw, localParquetError := utils.NewLocalFileWriter(fileName)
	if localParquetError != nil {
		log.Println("Can't create local file", localParquetError)
		return
	}

	pw, parquetWritterError := writer.NewParquetWriter(fw, new(utils.LogEntry), 4)

	if parquetWritterError != nil {
		log.Println("Can't create parquet writer", parquetWritterError)
		return
	}

	// pw.RowGroupSize = 128 * 1024 * 1024 //128M
	// pw.RowGroupSize = 1 * 1024 //1k
	// pw.PageSize = 1 * 1024     //1K
	// pw.CompressionType = parquet.CompressionCodec_SNAPPY

	// Ensure resources are closed properly
	var cleanupOnce sync.Once
	cleanup := func() {
		fmt.Println("Cleaning up resources...")
		if err := pw.WriteStop(); err != nil {
			log.Printf("Error stopping Parquet writer: %v", err)
		}
		if err := fw.Close(); err != nil {
			log.Printf("Error closing file writer: %v", err)
		}
	}

	defer cleanupOnce.Do(cleanup)

	// Signal channel to handle interrupts
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	// Create a ticker for logging every second
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	fmt.Println("Logged data to Parquet file")

	logs := computing.ComputeSystem(system)

	fmt.Println(logs)

	// go func() {
	// 	// Log system data every second
	// 	for range ticker.C {
	// 		// Collect logs from the system
	// 		// Rulează calculul pe baza configurării încărcate
	// 		logs := computing.ComputeSystem(system)

	// 		// Write logs to the Parquet file
	// 		for _, logEntry := range logs {
	// 			if err := pw.Write(logEntry); err != nil {
	// 				log.Printf("Failed to write log entry: %v", err)
	// 			}
	// 			// Force flush to immediately write the log to disk
	// 			if err := pw.Flush(true); err != nil {
	// 				log.Printf("Failed to flush parquet writer: %v", err)
	// 			}
	// 		}

	// 		fmt.Println("Logged data to Parquet file")
	// 	}
	// }()

	// // Wait for termination signal
	// sig := <-signals
	// log.Printf("Received signal: %v. Shutting down gracefully...", sig)
	// cleanupOnce.Do(cleanup)

	// zones(system)
}
