package logger

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

type Logger struct {
	fileLogger *log.Logger
	file       *os.File
}

type LogEntry struct {
	Timestamp string      `json:"timestamp"`
	Message   string      `json:"message"`
	Packet    interface{} `json:"packet,omitempty"`
	RawBytes  string      `json:"raw_bytes,omitempty"`
}

func NewLogger(logPath string) (*Logger, error) {
	file, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}

	return &Logger{
		fileLogger: log.New(file, "", log.LstdFlags),
		file:       file,
	}, nil
}

func (l *Logger) LogPacket(message string, packet interface{}, rawBytes []byte) {
	entry := LogEntry{
		Timestamp: l.fileLogger.Prefix(),
		Message:   message,
		Packet:    packet,
		RawBytes:  fmt.Sprintf("%x", rawBytes),
	}

	jsonData, err := json.Marshal(entry)
	if err != nil {
		l.fileLogger.Printf("Failed to marshal JSON: %v", err)
		return
	}
	l.fileLogger.Println(string(jsonData))

	prettyPrint(entry)
}

func (l *Logger) Log(message string) {
	entry := LogEntry{
		Timestamp: l.fileLogger.Prefix(),
		Message:   message,
	}
	jsonData, _ := json.Marshal(entry)
	l.fileLogger.Println(string(jsonData))
	prettyPrint(entry)
}

func prettyPrint(entry LogEntry) {
	fmt.Printf("\033[1;34m%s\033[0m\n", entry.Timestamp)
	fmt.Printf("  Message: %s\n", entry.Message)
	if entry.Packet != nil {
		fmt.Printf("  Packet: \033[1;32m%+v\033[0m\n", entry.Packet)
	}
	if entry.RawBytes != "" {
		fmt.Printf("  Raw Bytes: \033[1;90m%s\033[0m\n", entry.RawBytes)
	}
	fmt.Println("---")
}

func (l *Logger) Close() {
	l.file.Close()
}
