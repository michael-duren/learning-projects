package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"time"
)

type event struct {
	Timestamp string `json:"timestamp"`
	Level     string `json:"level"`
	Path      string `json:"path"`
	Status    int    `json:"status"`
	Message   string `json:"message"`
	RequestID string `json:"request_id"`
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	if err := os.MkdirAll("/logs", 0755); err != nil {
		return err
	}
	file, err := os.OpenFile("/logs/app.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(io.MultiWriter(file, os.Stdout))
	scenarios := []event{
		{Level: "info", Path: "/orders", Status: 200, Message: "order created"},
		{Level: "debug", Path: "/health", Status: 200, Message: "health check"},
		{Level: "warn", Path: "/orders", Status: 429, Message: "rate limited"},
		{Level: "error", Path: "/checkout", Status: 500, Message: "payment failed"},
	}
	for {
		for _, entry := range scenarios {
			now := time.Now().UTC()
			entry.Timestamp = now.Format(time.RFC3339Nano)
			entry.RequestID = fmt.Sprintf("request-%d", now.UnixNano())
			if err := encoder.Encode(entry); err != nil {
				return err
			}
			time.Sleep(time.Second)
		}
	}
}
