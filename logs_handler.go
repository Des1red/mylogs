package mylogs

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

var (
	logDir      = "logs/logfiles/" // Default log directory
	logFileLock sync.Mutex
)

// SetLogDirectory allows customizing the log directory
func SetLogDirectory(dir string) {
	logDir = dir
}

// LogMessage logs a message to a file specific to the given layer
func LogMessage(message, layer string) {
	// Ensure the logs directory exists
	if err := ensureDirectoryExists(logDir); err != nil {
		fmt.Println("Error ensuring logs directory exists:", err)
		return
	}

	// Construct the full path to the log file
	filename := filepath.Join(logDir, layer+".log")

	// Append a timestamp to the log message
	timestampedMessage := fmt.Sprintf("[%s] %s", time.Now().Format("2006-01-02 15:04:05"), message)

	// Write the message to the log file
	if err := writeToFile(timestampedMessage, filename); err != nil {
		fmt.Println("Error writing to log file:", err)
	}
}

// ensureDirectoryExists creates the log directory if it doesn't exist
func ensureDirectoryExists(dir string) error {
	_, err := os.Stat(dir)
	if os.IsNotExist(err) {
		if err = os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory: %v", err)
		}
	} else if err != nil {
		return fmt.Errorf("failed to stat directory: %v", err)
	}
	return nil
}

// writeToFile safely appends a message to the specified log file
func writeToFile(message, filename string) error {
	logFileLock.Lock()
	defer logFileLock.Unlock()

	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open log file: %v", err)
	}
	defer file.Close()

	if _, err := file.WriteString(message + "\n"); err != nil {
		return fmt.Errorf("failed to write to log file: %v", err)
	}
	return nil
}