package utils

import (
    "fmt"
    "os"
    "path/filepath"
    "time"
)

type Logger struct {
    logDir string
}

func NewLogger() *Logger {
    logDir := "logs"
    if err := os.MkdirAll(logDir, 0755); err != nil {
        panic(err)
    }
    return &Logger{logDir: logDir}
}

func (l *Logger) getLogFile() (*os.File, error) {
    today := time.Now().Format("20060102")
    logPath := filepath.Join(l.logDir, fmt.Sprintf("nalaai-%s.log", today))
    
    return os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
}

func (l *Logger) Log(entryType, message string) {
    file, err := l.getLogFile()
    if err != nil {
        fmt.Printf("Error opening log file: %v\n", err)
        return
    }
    defer file.Close()
    
    timestamp := time.Now().Format("20060102 15:04:05")
    logEntry := fmt.Sprintf("%s : [%s] %s\n", timestamp, entryType, message)
    
    file.WriteString(logEntry)
}

func (l *Logger) LogRequestResponse(request, response string) {
    file, err := l.getLogFile()
    if err != nil {
        fmt.Printf("Error opening log file: %v\n", err)
        return
    }
    defer file.Close()
    
    timestamp := time.Now().Format("20060102 15:04:05")
    
    file.WriteString(fmt.Sprintf("%s : [REQUEST] %s\n", timestamp, request))
    file.WriteString(fmt.Sprintf("%s : [RESPONSE] %s\n", timestamp, response))
    file.WriteString(fmt.Sprintf("%s : [SEPARATOR] --------------------\n", timestamp))
}