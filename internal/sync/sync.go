package sync

import (
	"archive/zip"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/geelato/cli/internal/file"
	"github.com/geelato/cli/pkg/logger"
)

type SyncService struct {
	cwd    string
	url    string
	key    string
	client *file.HTTPClient
}

type DiffResult struct {
	Added    []string `json:"added"`
	Modified []string `json:"modified"`
	Deleted  []string `json:"deleted"`
}

type FileRecord struct {
	Path     string `json:"path"`
	Hash     string `json:"hash"`
	Type     string `json:"type"`
}

func NewSyncService(cwd, url, key string) (*SyncService, error) {
	client, err := file.NewHTTPClient(url, key)
	if err != nil {
		return nil, err
	}

	return &SyncService{
		cwd:    cwd,
		url:    url,
		key:    key,
		client: client,
	}, nil
}

func (s *SyncService) Push(ctx context.Context, message string) error {
	logger.Info("Creating package...")

	zipPath, err := s.createPackage()
	if err != nil {
		return fmt.Errorf("failed to create package: %w", err)
	}
	defer os.Remove(zipPath)

	logger.Info("Uploading package...")
	err = s.uploadPackage(ctx, zipPath, message)
	if err != nil {
		return fmt.Errorf("failed to upload package: %w", err)
	}

	return nil
}

func (s *SyncService) Pull(ctx context.Context) error {
	logger.Info("Downloading package...")

	zipPath, err := s.downloadPackage(ctx)
	if err != nil {
		return fmt.Errorf("failed to download package: %w", err)
	}
	defer os.Remove(zipPath)

	logger.Info("Extracting package...")
	err = s.extractPackage(zipPath)
	if err != nil {
		return fmt.Errorf("failed to extract package: %w", err)
	}

	return nil
}

func (s *SyncService) GetDiff() (*DiffResult, error) {
	result := &DiffResult{
		Added:    []string{},
		Modified: []string{},
		Deleted:  []string{},
	}

	localFiles, err := s.scanLocalFiles()
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Get("/api/v1/apps/files")
	if err != nil {
		return nil, err
	}

	var remoteFiles []FileRecord
	if err := json.Unmarshal(resp, &remoteFiles); err != nil {
		return nil, err
	}

	localMap := make(map[string]string)
	for _, f := range localFiles {
		localMap[f.Path] = f.Hash
	}

	remoteMap := make(map[string]string)
	for _, f := range remoteFiles {
		remoteMap[f.Path] = f.Hash
	}

	for _, f := range localFiles {
		if _, exists := remoteMap[f.Path]; !exists {
			result.Added = append(result.Added, f.Path)
		} else if localMap[f.Path] != remoteMap[f.Path] {
			result.Modified = append(result.Modified, f.Path)
		}
	}

	for _, f := range remoteFiles {
		if _, exists := localMap[f.Path]; !exists {
			result.Deleted = append(result.Deleted, f.Path)
		}
	}

	return result, nil
}

func (s *SyncService) createPackage() (string, error) {
	tmpDir, err := ioutil.TempDir("", "geelato-push")
	if err != nil {
		return "", err
	}

	files, err := s.scanLocalFiles()
	if err != nil {
		return "", err
	}

	manifest := map[string]interface{}{
		"createdAt": time.Now().Format(time.RFC3339),
		"files":     files,
	}

	manifestData, err := json.MarshalIndent(manifest, "", "  ")
	if err != nil {
		return "", err
	}

	manifestPath := filepath.Join(tmpDir, "manifest.json")
	if err := ioutil.WriteFile(manifestPath, manifestData, 0644); err != nil {
		return "", err
	}

	baseName := filepath.Base(s.cwd)
	zipPath := filepath.Join(tmpDir, baseName+".zip")

	if err := s.zipDirectory(s.cwd, zipPath, files); err != nil {
		return "", err
	}

	return zipPath, nil
}

func (s *SyncService) scanLocalFiles() ([]FileRecord, error) {
	var files []FileRecord

	ignoreDirs := map[string]bool{
		".git":     true,
		"node_modules": true,
		".idea":    true,
		".vscode":  true,
		"vendor":   true,
		"__pycache__": true,
	}

	err := filepath.Walk(s.cwd, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(s.cwd, path)
		if err != nil {
			return err
		}

		if relPath == "." || relPath == "" {
			return nil
		}

		dir := filepath.Dir(relPath)
		if ignoreDirs[dir] || ignoreDirs[filepath.Base(dir)] {
			if info.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		if info.IsDir() {
			return nil
		}

		hash, err := file.HashFile(path)
		if err != nil {
			return err
		}

		files = append(files, FileRecord{
			Path: relPath,
			Hash: hash,
			Type: getFileType(relPath),
		})

		return nil
	})

	return files, err
}

func (s *SyncService) zipDirectory(source, target string, files []FileRecord) error {
	zipFile, err := os.Create(target)
	if err != nil {
		return err
	}
	defer zipFile.Close()

	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	fileSet := make(map[string]bool)
	for _, f := range files {
		fileSet[f.Path] = true
	}

	err = filepath.Walk(source, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(source, path)
		if err != nil {
			return err
		}

		if relPath == "." || relPath == "" {
			return nil
		}

		if !fileSet[relPath] {
			if info.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}

		header.Name = relPath
		header.Method = zip.Deflate

		writer, err := zipWriter.CreateHeader(header)
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		f, err := os.Open(path)
		if err != nil {
			return err
		}
		defer f.Close()

		_, err = io.Copy(writer, f)
		return err
	})

	return err
}

func (s *SyncService) uploadPackage(ctx context.Context, zipPath, message string) error {
	zipFile, err := os.Open(zipPath)
	if err != nil {
		return err
	}
	defer zipFile.Close()

	body := &file.UploadBody{
		File:     zipFile,
		Filename: filepath.Base(zipPath),
	}

	resp, err := s.client.Upload(ctx, "/api/v1/apps/upload", body)
	if err != nil {
		return err
	}

	logger.Debug(string(resp))
	return nil
}

func (s *SyncService) downloadPackage(ctx context.Context) (string, error) {
	resp, err := s.client.GetWithContext(ctx, "/api/v1/apps/current")
	if err != nil {
		return "", err
	}

	var metadata struct {
		ID   string `json:"id"`
		Path string `json:"path"`
	}
	if err := json.Unmarshal(resp, &metadata); err != nil {
		return "", err
	}

	tmpDir, err := ioutil.TempDir("", "geelato-pull")
	if err != nil {
		return "", err
	}

	zipPath := filepath.Join(tmpDir, "app.zip")
	zipFile, err := os.Create(zipPath)
	if err != nil {
		return "", err
	}
	defer zipFile.Close()

	_, err = s.client.Download(ctx, "/api/v1/apps/download", zipFile)
	if err != nil {
		return "", err
	}

	return zipPath, nil
}

func (s *SyncService) extractPackage(zipPath string) error {
	reader, err := zip.OpenReader(zipPath)
	if err != nil {
		return err
	}
	defer reader.Close()

	for _, zipEntry := range reader.File {
		path := filepath.Join(s.cwd, zipEntry.Name)

		if zipEntry.FileInfo().IsDir() {
			os.MkdirAll(path, 0755)
			continue
		}

		os.MkdirAll(filepath.Dir(path), 0755)

		rc, err := zipEntry.Open()
		if err != nil {
			return err
		}

		data, err := ioutil.ReadAll(rc)
		rc.Close()
		if err != nil {
			return err
		}

		if err := ioutil.WriteFile(path, data, 0644); err != nil {
			return err
		}
	}

	return nil
}

func getFileType(path string) string {
	ext := filepath.Ext(path)
	switch ext {
	case ".json":
		return "model"
	case ".js":
		return "api"
	case ".xml", ".bpmn":
		return "workflow"
	default:
		return "other"
	}
}

type UploadResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	ID      string `json:"id"`
}

var (
	ErrNotConfigured = errors.New("API URL not configured")
	ErrUnauthorized  = errors.New("Unauthorized")
)
