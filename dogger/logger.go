package logger

import (
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var logFiles map[string]*os.File

func init() {
	logFiles = make(map[string]*os.File)
}

func CreateLogFileInDirectory(date string, directory string) (*os.File, error) {
	filename := filepath.Join(directory, "log_"+date+".txt")
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}
	logFiles[date] = file
	return file, nil
}
// func CreateLogFile(date string) (*os.File, error) {
// 	filename := "log_" + date + ".txt"
// 	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
// 	if err != nil {
// 		return nil, err
// 	}
// 	logFiles[date] = file
// 	return file, nil
// }

func WriteLogToFile(date string, message string,logDirectory string) {
	file, exists := logFiles[date]
	if !exists {
		file, err := CreateLogFileInDirectory(date,logDirectory)
		if err != nil {
			log.Println("Error creating log file:", err)
			return
		}
		file.WriteString("Log File Created: " + time.Now().Format("2006-01-02 15:04:05") + "\n")
	}
	file.WriteString(time.Now().Format("2001-01-02 15:04:05") + ": " + message + "\n")
}

func CloseLogFiles() {
	for _, file := range logFiles {
		file.Close()
	}
}

func CheckFilesWithDate(date string) {
	filename := "log_" + date + ".txt"
	if _, err := os.Stat(filename); err == nil {
		log.Println("File with date", date, "exists.")
	} else if os.IsNotExist(err) {
		log.Println("File with date", date, "does not exist.")
	} else {
		log.Println("Error checking file:", err)
	}
}

func FindLogFiles(targetDate string, logDirectory string) ([]string, error) {
	var existingFiles []string
	err := filepath.Walk(logDirectory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasPrefix(info.Name(), "log_") {
			fileDate := strings.TrimPrefix(info.Name(), "log_")
			fileDate = strings.TrimSuffix(fileDate, ".txt")
			if strings.Contains(fileDate, targetDate) {
				existingFiles = append(existingFiles, path)
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return existingFiles, nil
}
