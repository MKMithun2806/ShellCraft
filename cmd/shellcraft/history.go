package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

type HistoryEntry struct {
	Timestamp string `json:"timestamp"`
	IP        string `json:"ip"`
	Port      int    `json:"port"`
	Type      string `json:"type"`
	Encoder   string `json:"encoder"`
	Payload   string `json:"payload"`
}

func getHistoryFile() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("cannot get home dir: %w", err)
	}
	dir := filepath.Join(home, ".shellcraft")
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", fmt.Errorf("cannot create config dir: %w", err)
	}
	return filepath.Join(dir, "history.json"), nil
}

func SaveToHistory(entry HistoryEntry) error {
	entry.Timestamp = time.Now().Format(time.RFC3339)
	
	history, err := LoadHistory()
	if err != nil {
		history = []HistoryEntry{}
	}
	history = append([]HistoryEntry{entry}, history...)
	
	if len(history) > 20 {
		history = history[:20]
	}

	data, err := json.MarshalIndent(history, "", "  ")
	if err != nil {
		return err
	}
	path, err := getHistoryFile()
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

func LoadHistory() ([]HistoryEntry, error) {
	path, err := getHistoryFile()
	if err != nil {
		return nil, err
	}
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return []HistoryEntry{}, nil
	}
	
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var history []HistoryEntry
	err = json.Unmarshal(data, &history)
	return history, err
}
