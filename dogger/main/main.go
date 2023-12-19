package main

import (
	"log"
	"os"
	"path/filepath"
	"time"

	logger "github.com/EnsurityTechnologies/logger/dogger"
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

	logFilePath := filepath.Join(logDirectory, "logs.txt")
	log.Printf("Log file path: %s\n", logFilePath) // Print log file path

	f, err := os.OpenFile(logFilePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Fatalf("Failed to open logfile: %v", err)
	}
	defer f.Close()

	log.SetOutput(f)

	// Check if log file for the current date exists
	currentDateLog := filepath.Join(logDirectory, time.Now().Format("2006-01-02")+"_logs.txt")
	if _, err := os.Stat(currentDateLog); err == nil {
		log.Printf("Log file for today's date '%s' exists.\n", time.Now().Format("2006-01-02"))
	} else if os.IsNotExist(err) {
		log.Printf("Log file for today's date '%s' does not exist.\n", time.Now().Format("2006-01-02"))
	} else {
		log.Println("Error checking current date log file:", err)
	}

	// Check log files in the current week
	startOfWeek := time.Now().AddDate(0, 0, -int(time.Now().Weekday()))
	for i := 0; i < 7; i++ {
		date := startOfWeek.AddDate(0, 0, i)
		logFileForDate := filepath.Join(logDirectory, date.Format("2006-01-02")+"_logs.txt")
		if _, err := os.Stat(logFileForDate); err == nil {
			log.Printf("Log file for date '%s' exists: %s\n", date.Format("2006-01-02"), logFileForDate)
		} else if os.IsNotExist(err) {
			log.Printf("Log file for date '%s' does not exist.\n", date.Format("2006-01-02"))
		} else {
			log.Printf("Error checking log file for date '%s': %s\n", date.Format("2006-01-02"), err)
		}
	}

	targetDate := "2023-12-13"
	log.Println("Searching for log files with the date:", targetDate)
	logsForDate, err := logger.FindLogFiles(targetDate, logDirectory)
	if err != nil {
		log.Println("Error finding log files:", err)
		return
	}
	if len(logsForDate) > 0 {
		log.Println("Found log files for date", len(logsForDate), targetDate)
		for _, files := range logsForDate {
			log.Println(files)
		}
	} else {
		log.Println("No log files found for date", targetDate)
	}

	currentDate := time.Now().Format("2006-01-02")
	fp, err := logger.CreateLogFileInDirectory(currentDate, logDirectory)
	if err != nil {
		
		log.Fatal("Error creating log file:", err)
	}
	defer fp.Close()

	// Integration of logger functionalities
	l := log.New(fp, "LOG: ", log.Ldate|log.Ltime|log.Lshortfile)
	l.Println("This is a log message for today's date")
}




























// package main

// import (
// 	//"fmt"
// 	"log"
// 	"os"
// 	"time"

// 	logger "github.com/EnsurityTechnologies/logger/dogger"
// 	)

// func main() {
// 	logDirectory := "path/logger/dogger/main/newdir"
// 	if err := os.MkdirAll(logDirectory, os.ModePerm); err != nil {
// 		directoryname := ""
// 	    err := os.Mkdir(directoryname, 0755)
// 		log.Fatalf("Error creating log directory: %s", err)
// 	}
// 	logFilePath := "path/logger/dogger/main/newdir/logs.txt"
// 	f, err := os.OpenFile(logFilePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
// 	//f, err := os.OpenFile("logs.txt", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
// 	if err != nil {
// 		log.Fatalf("Failed to open logfile: %v", err)
// 	}
// 	defer f.Close()

// 	log.SetOutput(f) // Set log output to the main log file
// 	targetDate := "2023-12-13"
// 	log.Println("Searching for log files with the date:", targetDate)
// 	logsForDate, err := logger.FindLogFiles(targetDate, logDirectory)
// 	if err != nil {
// 		log.Println("Error finding log files:", err)
// 		return
// 	}
// 	if len(logsForDate) > 0 {
// 		log.Println("Found  log files for date ", len(logsForDate), targetDate)
// 		for _, files := range logsForDate {
// 			log.Println(files)
// 		}
// 	} else {
// 		log.Println("No log files found for date ", targetDate)
// 	}
// 	//  newDate := "2023-12-05"
// 	// newLogFile, err := logger.CreateLogFileInDirectory(newDate, logDirectory)
// 	// if err != nil {
// 	// 	log.Println("Error creating log file for", newDate, ":", err)
// 	// 	return
// 	// }
// 	// defer newLogFile.Close()
// 	// log.Println("New log file created for date", newDate)

// // 	newDate1 := "2023-12-06"
// // 	newLogFile1, err := logger.CreateLogFile(newDate1)
// // if err != nil {
// // 	log.Println("Error creating log file for", newDate1, ":", err)
// // 	return
// // }
// // defer newLogFile1.Close()
// // log.Println("New log file created for date", newDate1)

// 	currentDate := time.Now().Format("2001-05-14")

// 	fp, err := logger.CreateLogFileInDirectory(currentDate,logDirectory)
// 	if err != nil {
// 		log.Fatal("Error creating log file:", err)
// 	}
// 	defer fp.Close()

// 	// Integration of logger functionalities
// 	l := log.New(fp, "LOG: ", log.Ldate|log.Ltime|log.Lshortfile)
// 	l.Println("This is a log message for today's date")
// 	//logger.CheckFilesWithDate("2023-12-05")
// }
