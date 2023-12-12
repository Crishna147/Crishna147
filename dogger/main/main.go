package main

import (
	//"fmt"
	"log"
	"os"
	"time"

	logger "github.com/EnsurityTechnologies/logger/dogger"
	)
	
func main() {
	//directoryname := "newdir"
	//err := os.Mkdir(directoryname, 0755)
	//if err != nil {
	//   fmt.Println(err) 
	//   return
	//}
	logDirectory := "D:/go/src/github.com/EnsurityTechnologies/logger/dogger/main/newdir"

	if err := os.MkdirAll(logDirectory, os.ModePerm); err != nil {
		log.Fatalf("Error creating log directory: %s", err)
	}
	logFilePath := "D:/go/src/github.com/EnsurityTechnologies/logger/dogger/main/newdir/logs.txt"
	f, err := os.OpenFile(logFilePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	//f, err := os.OpenFile("logs.txt", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Fatalf("Failed to open logfile: %v", err)
	}
	defer f.Close()

	log.SetOutput(f) // Set log output to the main log file
	targetDate := "2023-12-12"
	log.Println("Searching for log files with the date:", targetDate)
	logsForDate, err := logger.FindLogFiles(targetDate, logDirectory)
	if err != nil {
		log.Println("Error finding log files:", err)
		return
	}
	if len(logsForDate) > 0 {
		log.Println("Found  log files for date ", len(logsForDate), targetDate)
		for _, files := range logsForDate {
			log.Println(files)
		}
	} else {
		log.Println("No log files found for date ", targetDate)
	}
	 //newDate := "2023-12-05"
	// newLogFile, err := logger.CreateLogFileInDirectory(newDate, logDirectory)
	// if err != nil {
	// 	log.Println("Error creating log file for", newDate, ":", err)
	// 	return
	// }
	// defer newLogFile.Close()
	// log.Println("New log file created for date", newDate)

// 	newDate1 := "2023-12-06"
// 	newLogFile1, err := logger.CreateLogFile(newDate1)
// if err != nil {
// 	log.Println("Error creating log file for", newDate1, ":", err)
// 	return
// }
// defer newLogFile1.Close()
// log.Println("New log file created for date", newDate1)

	currentDate := time.Now().Format("2001-05-14")

	fp, err := logger.CreateLogFileInDirectory(currentDate,logDirectory)
	if err != nil {
		log.Fatal("Error creating log file:", err)
	}
	defer fp.Close()

	// Integration of logger functionalities
	l := log.New(fp, "LOG: ", log.Ldate|log.Ltime|log.Lshortfile)
	l.Println("This is a log message for today's date")
	//logger.CheckFilesWithDate("2023-12-05")
}


