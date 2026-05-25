package main

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type Template struct {
	Name    string `json:"name"`
	IP      string `json:"ip"`
	Port    int    `json:"port"`
	Type    string `json:"type"`
	Encoder string `json:"encoder"`
}

func getTemplateDir() string {
	home, _ := os.UserHomeDir()
	dir := filepath.Join(home, ".shellcraft", "templates")
	os.MkdirAll(dir, 0755)
	return dir
}

func SaveTemplate(t Template) error {
	path := filepath.Join(getTemplateDir(), t.Name+".json")
	data, err := json.MarshalIndent(t, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

func ListTemplates() ([]Template, error) {
	dir := getTemplateDir()
	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	var templates []Template
	for _, f := range files {
		if filepath.Ext(f.Name()) == ".json" {
			data, _ := os.ReadFile(filepath.Join(dir, f.Name()))
			var t Template
			if err := json.Unmarshal(data, &t); err == nil {
				templates = append(templates, t)
			}
		}
	}
	return templates, nil
}
