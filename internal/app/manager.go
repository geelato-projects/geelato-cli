package app

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/geelato/cli/internal/config"
)

type Manager struct {
	appsDir string
}

func NewManager() *Manager {
	cfg := config.Get()
	var appsDir string
	if cfg != nil && cfg.Cache.Dir != "" {
		appsDir = filepath.Join(cfg.Cache.Dir, "apps")
	} else {
		appsDir = ".geelato/apps"
	}
	return &Manager{
		appsDir: appsDir,
	}
}

func (m *Manager) Exists(path string) bool {
	geelatoFile := filepath.Join(path, "geelato.json")
	_, err := os.Stat(geelatoFile)
	return err == nil
}

func (m *Manager) Load(path string) (*Application, error) {
	geelatoFile := filepath.Join(path, "geelato.json")
	if _, err := os.Stat(geelatoFile); os.IsNotExist(err) {
		return nil, nil
	}

	content, err := os.ReadFile(geelatoFile)
	if err != nil {
		return nil, err
	}

	var app Application
	if err := json.Unmarshal(content, &app); err != nil {
		return nil, err
	}

	app.Path = path
	return &app, nil
}

func (m *Manager) Save(app *Application) error {
	geelatoFile := filepath.Join(app.Path, "geelato.json")

	data, err := json.MarshalIndent(app, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(geelatoFile, data, 0644)
}

func (m *Manager) FindApps(searchDir string) ([]ApplicationInfo, error) {
	var apps []ApplicationInfo

	entries, err := os.ReadDir(searchDir)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		appDir := filepath.Join(searchDir, entry.Name())
		geelatoFile := filepath.Join(appDir, "geelato.json")

		if _, err := os.Stat(geelatoFile); err == nil {
			content, _ := os.ReadFile(geelatoFile)
			info := ApplicationInfo{
				Path: appDir,
			}
			parseAppInfo(string(content), &info)
			if stat, _ := os.Stat(geelatoFile); stat != nil {
				info.LastModTime = stat.ModTime()
			}
			apps = append(apps, info)
		}
	}

	return apps, nil
}

func parseAppInfo(content string, info *ApplicationInfo) {
	// Simplified JSON parsing for basic info
	// A full implementation would use json.Unmarshal
}
