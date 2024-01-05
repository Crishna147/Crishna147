package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"
)

func main() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("Error getting user's home directory: %s", err)
	}

	logDirectory := filepath.Join(homeDir, "logger", "dogger", "main", "newdir")
	log.Printf("Log directory path: %s\n", logDirectory) // Print log directory path

	if _, err := os.Stat(logDirectory); os.IsNotExist(err) {
		
		err := os.MkdirAll(logDirectory, os.ModePerm)
		if err != nil {
			log.Fatalf("Error creating log directory: %s", err)
		}
	}

	createLogDirectories(logDirectory) // Create day, week, and month directories

	// Function to create a log file with a unique identifier for the current date
	createLogFile := func() *os.File {
		now := time.Now()
		uniqueIdentifier := fmt.Sprintf("%02d%02d%02d", now.Hour(), now.Minute(), now.Second())
		logFilePath := filepath.Join(logDirectory, "day", now.Format("2006-01-02"), "_logs_"+uniqueIdentifier+".txt")
		log.Printf("Log file path: %s\n", logFilePath) // Print log file path

		f, err := os.OpenFile(logFilePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
		if err != nil {
			log.Fatalf("Failed to open logfile: %v", err)
		}
		return f
	}

	logFile := createLogFile()
	defer logFile.Close()

	log.SetOutput(logFile)

	// Function to check and rotate log file at 12 am every day
	rotateLogFile := func() {
		for {
			now := time.Now()
			nextMidnight := time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 0, 0, now.Location())

			timeUntilNextMidnight := nextMidnight.Sub(now)
			time.Sleep(timeUntilNextMidnight)

			logFile.Close()
			createLogDirectories(logDirectory) // Recreate day, week, and month directories
			logFile = createLogFile()
			log.SetOutput(logFile)
		}
	}

	go rotateLogFile()

	// Prevent the main goroutine from exiting
	select {}
}

func createLogDirectories(logDirectory string) {
	now := time.Now()
	createDirIfNotExist(filepath.Join(logDirectory, "day", now.Format("2006-01-02")))
	createDirIfNotExist(filepath.Join(logDirectory, "week", now.Format("2006-W01")))
	createDirIfNotExist(filepath.Join(logDirectory, "month", now.Format("2006-01")))
}

func createDirIfNotExist(dirPath string) {
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		err := os.MkdirAll(dirPath, os.ModePerm)
		if err != nil {
			log.Fatalf("Error creating directory: %s", err)
		}
	}
}
