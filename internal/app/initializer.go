package app

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

type Initializer struct {
	templatesDir string
}

func NewInitializer() *Initializer {
	return &Initializer{
		templatesDir: "templates",
	}
}

func (i *Initializer) Initialize(app *Application) error {
	if err := i.createDirectoryStructure(app); err != nil {
		return fmt.Errorf("创建目录结构失败: %w", err)
	}

	if err := i.createConfigFile(app); err != nil {
		return fmt.Errorf("创建配置文件失败: %w", err)
	}

	if err := i.createGitKeepFiles(); err != nil {
		return fmt.Errorf("创建占位文件失败: %w", err)
	}

	return nil
}

func (i *Initializer) createDirectoryStructure(app *Application) error {
	dirs := []string{
		"meta",
		"api",
		"page",
		"workflow",
		".geelato",
	}

	for _, dir := range dirs {
		fullPath := filepath.Join(app.Path, dir)
		if err := os.MkdirAll(fullPath, 0755); err != nil {
			return err
		}
	}

	return nil
}

func (i *Initializer) createConfigFile(app *Application) error {
	now := time.Now().Format(time.RFC3339)

	configData := map[string]interface{}{
		"meta": map[string]interface{}{
			"version":    "1.0.0",
			"appId":      fmt.Sprintf("app_%s", app.Name),
			"name":       app.Name,
			"description": app.Description,
			"createdAt":  now,
		},
		"config": map[string]interface{}{
			"api": map[string]interface{}{
				"url":     "",
				"timeout": 30,
			},
			"sync": map[string]interface{}{
				"autoPush": false,
				"autoPull": false,
			},
		},
	}

	data, err := json.MarshalIndent(configData, "", "  ")
	if err != nil {
		return err
	}

	configPath := filepath.Join(app.Path, "geelato.json")
	return os.WriteFile(configPath, data, 0644)
}

func (i *Initializer) createGitKeepFiles() error {
	gitkeepFiles := []string{
		"api/.gitkeep",
		"meta/.gitkeep",
		"page/.gitkeep",
		"workflow/.gitkeep",
	}

	content := "# Keep this directory\n"

	for _, path := range gitkeepFiles {
		if err := os.WriteFile(path, []byte(content), 0644); err != nil {
			return err
		}
	}

	return nil
}

func (i *Initializer) InitializeWithExample(app *Application) error {
	if err := i.Initialize(app); err != nil {
		return err
	}

	exampleDirs := []string{
		"example/user-management-app",
		"example/user-management-app/meta/User",
		"example/user-management-app/meta/Department",
		"example/user-management-app/api/user",
		"example/user-management-app/page/user",
		"example/user-management-app/workflow",
	}

	for _, dir := range exampleDirs {
		fullPath := filepath.Join(app.Path, dir)
		if err := os.MkdirAll(fullPath, 0755); err != nil {
			return err
		}
	}

	return nil
}
