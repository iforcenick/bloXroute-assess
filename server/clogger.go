package main

import (
	"io"
	"log"
	"os"
)

var (
	InfoLogger    *log.Logger
	WarningLogger *log.Logger
	ErrorLogger   *log.Logger
)

/**
 * Initialize logger with the given log file path
 * If log_path is set to empty string, it uses stdout as a default log path
 */
func init_logger(log_path string) {
	var multi io.Writer
	if log_path == "" {
		multi = os.Stdout
	} else {
		file, err := os.OpenFile(log_path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
		if err != nil {
			log.Fatal(err)
		}
		multi = io.MultiWriter(file, os.Stdout)
	}

	InfoLogger = log.New(os.Stdout, "[INFO]: ", log.Ldate|log.Ltime|log.Lshortfile)
	WarningLogger = log.New(multi, "[WARNING]: ", log.Ldate|log.Ltime|log.Lshortfile)
	ErrorLogger = log.New(multi, "[ERROR]: ", log.Ldate|log.Ltime|log.Lshortfile)
}

// Log info-level message
func log_info(msg any) {
	InfoLogger.Print(msg)
}

// Log warning-level message
func log_warning(msg string) {
	WarningLogger.Print(msg)
}

// Log error-level message
func log_error(msg string) {
	ErrorLogger.Print(msg)
}
