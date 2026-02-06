package file

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/geelato/cli/pkg/crypto"
	"github.com/geelato/cli/pkg/logger"
)

type HTTPClient struct {
	baseURL string
	key     string
	client  *http.Client
}

type UploadBody struct {
	File     *os.File
	Filename string
}

func NewHTTPClient(baseURL, key string) (*HTTPClient, error) {
	if baseURL == "" {
		return nil, fmt.Errorf("API URL cannot be empty")
	}

	baseURL = strings.TrimSuffix(baseURL, "/")

	return &HTTPClient{
		baseURL: baseURL,
		key:     key,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}, nil
}

func (c *HTTPClient) request(method, path string, body interface{}) ([]byte, error) {
	url := fmt.Sprintf("%s%s", c.baseURL, path)

	var reqBody io.Reader
	if body != nil {
		jsonData, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		reqBody = bytes.NewBuffer(jsonData)
	}

	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	if c.key != "" {
		req.Header.Set("X-API-Key", c.key)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		respBody, _ := ioutil.ReadAll(resp.Body)
		return nil, fmt.Errorf("request failed with status %d: %s", resp.StatusCode, string(respBody))
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return respBody, nil
}

func (c *HTTPClient) Get(path string) ([]byte, error) {
	return c.request(http.MethodGet, path, nil)
}

func (c *HTTPClient) GetWithContext(ctx context.Context, path string) ([]byte, error) {
	url := fmt.Sprintf("%s%s", c.baseURL, path)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	if c.key != "" {
		req.Header.Set("X-API-Key", c.key)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		respBody, _ := ioutil.ReadAll(resp.Body)
		return nil, fmt.Errorf("request failed with status %d: %s", resp.StatusCode, string(respBody))
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return respBody, nil
}

func (c *HTTPClient) Post(path string, body interface{}) ([]byte, error) {
	return c.request(http.MethodPost, path, body)
}

func (c *HTTPClient) Put(path string, body interface{}) ([]byte, error) {
	return c.request(http.MethodPut, path, body)
}

func (c *HTTPClient) Delete(path string) ([]byte, error) {
	return c.request(http.MethodDelete, path, nil)
}

func (c *HTTPClient) Upload(ctx context.Context, path string, uploadBody *UploadBody) ([]byte, error) {
	url := fmt.Sprintf("%s%s", c.baseURL, path)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, uploadBody.File)
	if err != nil {
		return nil, err
	}

	if c.key != "" {
		req.Header.Set("X-API-Key", c.key)
	}

	fileInfo, _ := uploadBody.File.Stat()
	req.Header.Set("Content-Type", "application/zip")
	req.Header.Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", uploadBody.Filename))
	req.ContentLength = fileInfo.Size()

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		respBody, _ := ioutil.ReadAll(resp.Body)
		return nil, fmt.Errorf("upload failed with status %d: %s", resp.StatusCode, string(respBody))
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return respBody, nil
}

func (c *HTTPClient) Download(ctx context.Context, path string, writer io.Writer) ([]byte, error) {
	url := fmt.Sprintf("%s%s", c.baseURL, path)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	if c.key != "" {
		req.Header.Set("X-API-Key", c.key)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		respBody, _ := ioutil.ReadAll(resp.Body)
		return nil, fmt.Errorf("download failed with status %d: %s", resp.StatusCode, string(respBody))
	}

	_, err = io.Copy(writer, resp.Body)
	if err != nil {
		return nil, err
	}

	contentLength := resp.Header.Get("Content-Length")
	return []byte(contentLength), nil
}

func HashFile(path string) (string, error) {
	return crypto.HashFile(path, crypto.SHA256)
}

func WalkDirectory(root string, includePatterns []string, excludePatterns []string) ([]string, error) {
	var files []string

	excludeSet := make(map[string]bool)
	for _, pattern := range excludePatterns {
		excludeSet[pattern] = true
	}

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}

		if relPath == "." {
			return nil
		}

		dir := filepath.Dir(relPath)
		if excludeSet[dir] || excludeSet[filepath.Base(dir)] {
			if info.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		if info.IsDir() {
			return nil
		}

		if len(includePatterns) > 0 {
			ext := filepath.Ext(path)
			include := false
			for _, pattern := range includePatterns {
				if ext == pattern || strings.HasSuffix(path, pattern) {
					include = true
					break
				}
			}
			if !include {
				return nil
			}
		}

		files = append(files, relPath)
		return nil
	})

	return files, err
}

type FileInfo struct {
	Path     string
	Name     string
	Size     int64
	Mode     os.FileMode
	IsDir    bool
	Hash     string
	Modified time.Time
}

func GetFileInfo(path string) (*FileInfo, error) {
	info, err := os.Stat(path)
	if err != nil {
		return nil, err
	}

	hash, err := HashFile(path)
	if err != nil {
		logger.Warn("Failed to hash file: %s", err.Error())
		hash = ""
	}

	return &FileInfo{
		Path:     path,
		Name:     filepath.Base(path),
		Size:     info.Size(),
		Mode:     info.Mode(),
		IsDir:    info.IsDir(),
		Hash:     hash,
		Modified: info.ModTime(),
	}, nil
}
