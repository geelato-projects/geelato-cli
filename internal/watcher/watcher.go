package watcher

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/geelato/cli/internal/file"
	"github.com/geelato/cli/pkg/logger"
)

type Watcher struct {
	rootDir  string
	interval time.Duration
	lastHash map[string]string
	client   *file.HTTPClient
}

type WatchEvent struct {
	Type   string
	Path   string
	RelPath string
}

func NewWatcher(rootDir string) (*Watcher, error) {
	_, err := os.Stat(rootDir)
	if err != nil {
		return nil, err
	}

	return &Watcher{
		rootDir:  rootDir,
		interval: 2 * time.Second,
		lastHash: make(map[string]string),
	}, nil
}

func (w *Watcher) SetClient(client *file.HTTPClient) {
	w.client = client
}

func (w *Watcher) Start(ctx context.Context) error {
	if err := w.scanInitial(); err != nil {
		return err
	}

	ticker := time.NewTicker(w.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
			if err := w.checkChanges(ctx); err != nil {
				logger.Warn("Error checking changes: %s", err.Error())
			}
		}
	}
}

func (w *Watcher) scanInitial() error {
	files, err := file.WalkDirectory(
		w.rootDir,
		[]string{".json", ".js", ".xml", ".bpmn"},
		[]string{".git", "node_modules", ".idea", ".vscode"},
	)
	if err != nil {
		return err
	}

	for _, relPath := range files {
		fullPath := filepath.Join(w.rootDir, relPath)
		hash, err := file.HashFile(fullPath)
		if err != nil {
			logger.Warn("Failed to hash file %s: %s", relPath, err.Error())
			continue
		}
		w.lastHash[relPath] = hash
	}

	logger.Infof("Watching %d files", len(files))
	return nil
}

func (w *Watcher) checkChanges(ctx context.Context) error {
	files, err := file.WalkDirectory(
		w.rootDir,
		[]string{".json", ".js", ".xml", ".bpmn"},
		[]string{".git", "node_modules", ".idea", ".vscode"},
	)
	if err != nil {
		return err
	}

	currentHash := make(map[string]string)
	for _, relPath := range files {
		fullPath := filepath.Join(w.rootDir, relPath)
		hash, err := file.HashFile(fullPath)
		if err != nil {
			logger.Warn("Failed to hash file %s: %s", relPath, err.Error())
			continue
		}
		currentHash[relPath] = hash
	}

	for relPath, hash := range currentHash {
		if _, exists := w.lastHash[relPath]; !exists {
			w.onEvent(WatchEvent{Type: "created", Path: filepath.Join(w.rootDir, relPath), RelPath: relPath})
		} else if w.lastHash[relPath] != hash {
			w.onEvent(WatchEvent{Type: "modified", Path: filepath.Join(w.rootDir, relPath), RelPath: relPath})
		}
	}

	for relPath := range w.lastHash {
		if _, exists := currentHash[relPath]; !exists {
			w.onEvent(WatchEvent{Type: "deleted", Path: relPath, RelPath: relPath})
		}
	}

	w.lastHash = currentHash
	return nil
}

func (w *Watcher) onEvent(event WatchEvent) {
	switch event.Type {
	case "created":
		logger.Success("[CREATED] %s", event.RelPath)
	case "modified":
		logger.Info("[MODIFIED] %s", event.RelPath)
	case "deleted":
		logger.Warn("[DELETED] %s", event.RelPath)
	}

	if w.client != nil {
		w.syncEvent(event)
	}
}

func (w *Watcher) syncEvent(event WatchEvent) {
	logger.Debug("Syncing event to cloud: %s", event.Type)

	eventData := map[string]string{
		"type":   event.Type,
		"path":   event.RelPath,
		"root":   w.rootDir,
	}

	_, err := w.client.Post("/api/v1/sync/event", eventData)
	if err != nil {
		logger.Warn("Failed to sync event: %s", err.Error())
	}
}

func isIgnoredDir(dirName string) bool {
	ignored := []string{"node_modules", ".git", "vendor", "__pycache__", ".idea", ".vscode"}
	for _, ig := range ignored {
		if dirName == ig {
			return true
		}
	}
	return false
}

func isIgnoredFile(fileName string) bool {
	ignoredPrefixes := []string{".", "_"}
	for _, prefix := range ignoredPrefixes {
		if strings.HasPrefix(fileName, prefix) {
			return true
		}
	}
	return false
}
