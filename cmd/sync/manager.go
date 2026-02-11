package sync

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/geelato/cli/internal/file"
	"github.com/geelato/cli/internal/platform"
	"github.com/geelato/cli/pkg/logger"
)

type Change struct {
	Type        string
	Path        string
	Content     []byte
	LocalHash   string
	RemoteHash  string
	Description string
}

type ChangeType string

const (
	ChangeTypeAdded    ChangeType = "added"
	ChangeTypeModified ChangeType = "modified"
	ChangeTypeDeleted  ChangeType = "deleted"
)

type Manager struct {
	syncDir        string
	stateFile      string
	platformClient *platform.Client
}

func NewManager() *Manager {
	return &Manager{
		syncDir:        ".geelato",
		stateFile:      filepath.Join(".geelato", "sync-state.json"),
		platformClient: platform.NewClient(),
	}
}

func NewManagerWithAPI(apiURL string) *Manager {
	return &Manager{
		syncDir:        ".geelato",
		stateFile:      filepath.Join(".geelato", "sync-state.json"),
		platformClient: platform.NewClientWithURL(apiURL),
	}
}

func (m *Manager) DetectChanges(rootDir string) ([]Change, error) {
	var changes []Change

	syncState, _ := m.loadSyncState()

	patterns := []string{
		"meta/**/*.json",
		"meta/**/*.sql",
		"api/**/*.js",
		"api/**/*.py",
		"api/**/*.go",
		"page/**/*.json",
		"workflow/**/*.json",
	}

	for _, pattern := range patterns {
		absPattern := filepath.Join(rootDir, pattern)
		files, _ := filepath.Glob(absPattern)
		for _, filePath := range files {
			relPath, _ := filepath.Rel(rootDir, filePath)
			hash := m.computeHash(filePath)

			lastSync, exists := syncState.Files[relPath]
			if !exists {
				changes = append(changes, Change{
					Type:        string(ChangeTypeAdded),
					Path:        relPath,
					Content:     m.readContent(filePath),
					LocalHash:   hash,
					Description: "新增文件",
				})
			} else if lastSync != hash {
				changes = append(changes, Change{
					Type:        string(ChangeTypeModified),
					Path:        relPath,
					Content:     m.readContent(filePath),
					LocalHash:   hash,
					RemoteHash:  lastSync,
					Description: "修改文件",
				})
			}
		}
	}

	return changes, nil
}

func (m *Manager) Push(changes []Change, message string) (*platform.SyncStatus, error) {
	if len(changes) == 0 {
		return nil, fmt.Errorf("没有需要推送的变更")
	}

	files := make([]platform.FileEntry, 0, len(changes))
	for _, change := range changes {
		files = append(files, platform.FileEntry{
			Path:    change.Path,
			Content: change.Content,
			Hash:    change.LocalHash,
		})
	}

	version := time.Now().Format("20060102150405")

	req := &platform.UploadRequest{
		AppID:   m.getAppID(),
		Version: version,
		Branch:  "main",
		Message: message,
		Author:  "Developer",
		Files:   files,
	}

	newVersion, err := m.platformClient.UploadAppPackage(nil, req)
	if err != nil {
		return nil, fmt.Errorf("上传失败: %w", err)
	}

	m.updateSyncState(changes, newVersion)

	return &platform.SyncStatus{
		AppID:   m.getAppID(),
		Version: newVersion,
		Status:  "success",
	}, nil
}

func (m *Manager) Pull(version string) (*platform.SyncStatus, error) {
	return nil, fmt.Errorf("拉取功能待实现")
}

func (m *Manager) GetStatus() (*SyncStatusView, error) {
	syncState, _ := m.loadSyncState()

	changes, _ := m.DetectChanges(".")

	var ahead, behind int
	for _, change := range changes {
		if change.Type == string(ChangeTypeAdded) || change.Type == string(ChangeTypeModified) {
			ahead++
		}
	}

	return &SyncStatusView{
		LocalVersion:  syncState.Version,
		RemoteVersion: "latest",
		AheadBy:       ahead,
		BehindBy:      behind,
		HasConflict:   false,
		LastSyncTime:  syncState.LastSyncAt,
	}, nil
}

func (m *Manager) DetectConflicts() ([]Conflict, error) {
	return nil, fmt.Errorf("冲突检测待实现")
}

func (m *Manager) ResolveWithLocal(conflict Conflict) error {
	return nil
}

func (m *Manager) ResolveWithRemote(conflict Conflict) error {
	return nil
}

func (m *Manager) computeHash(filePath string) string {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return ""
	}
	hash := sha256.Sum256(data)
	return hex.EncodeToString(hash[:])
}

func (m *Manager) readContent(filePath string) []byte {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil
	}
	return data
}

func (m *Manager) loadSyncState() (*SyncState, error) {
	if !file.Exists(m.stateFile) {
		return &SyncState{
			Files: make(map[string]string),
		}, nil
	}

	data, err := os.ReadFile(m.stateFile)
	if err != nil {
		return nil, err
	}

	var state SyncState
	if err := json.Unmarshal(data, &state); err != nil {
		return nil, err
	}

	if state.Files == nil {
		state.Files = make(map[string]string)
	}

	return &state, nil
}

func (m *Manager) updateSyncState(changes []Change, version string) {
	state, _ := m.loadSyncState()

	for _, change := range changes {
		if change.Type == string(ChangeTypeAdded) || change.Type == string(ChangeTypeModified) {
			state.Files[change.Path] = change.LocalHash
		} else if change.Type == string(ChangeTypeDeleted) {
			delete(state.Files, change.Path)
		}
	}

	state.Version = version
	state.LastSyncAt = time.Now().Format(time.RFC3339)

	data, _ := json.MarshalIndent(state, "", "  ")
	os.WriteFile(m.stateFile, data, 0644)
}

func (m *Manager) getAppID() string {
	data, err := os.ReadFile("geelato.json")
	if err != nil {
		return ""
	}
	var config struct {
		Meta struct {
			AppID string `json:"appId"`
		} `json:"meta"`
	}
	json.Unmarshal(data, &config)
	return config.Meta.AppID
}

type SyncState struct {
	Version    string            `json:"version"`
	LastSyncAt string            `json:"lastSyncAt"`
	Files      map[string]string `json:"files"`
}

type SyncStatusView struct {
	LocalVersion  string
	RemoteVersion string
	AheadBy       int
	BehindBy      int
	HasConflict   bool
	LastSyncTime  string
}

type Conflict struct {
	Path       string
	LocalHash  string
	RemoteHash string
}
