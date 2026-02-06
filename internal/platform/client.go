package platform

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/geelato/cli/internal/config"
	"github.com/geelato/cli/pkg/crypto"
)

type Client struct {
	baseURL   string
	apiKey    string
	client    *http.Client
	headers   map[string]string
}

type RequestOptions struct {
	Method      string
	Path        string
	Body        interface{}
	QueryParams map[string]string
	Headers     map[string]string
	Timeout     time.Duration
}

type Response struct {
	StatusCode int
	Headers    http.Header
	Body       []byte
	Data       interface{}
}

type UploadRequest struct {
	AppID     string
	Version   string
	Branch    string
	Message   string
	Author    string
	Files     []FileEntry
}

type FileEntry struct {
	Path    string
	Content []byte
	Hash    string
}

type SyncStatus struct {
	AppID     string    `json:"appId"`
	Version   string    `json:"version"`
	Branch    string    `json:"branch"`
	Status    string    `json:"status"`
	Message   string    `json:"message"`
	Author    string    `json:"author"`
	CreatedAt time.Time `json:"createdAt"`
}

type ConflictInfo struct {
	Path       string `json:"path"`
	LocalHash  string `json:"localHash"`
	RemoteHash string `json:"remoteHash"`
}

type ConflictCheckResult struct {
	HasConflict bool           `json:"hasConflict"`
	Conflicts   []ConflictInfo `json:"conflicts"`
}

func NewClient() *Client {
	cfg := config.Get()
	return NewClientWithConfig(cfg)
}

func NewClientWithConfig(cfg *config.Config) *Client {
	return &Client{
		baseURL: cfg.API.URL,
		apiKey:  cfg.API.Key,
		client: &http.Client{
			Timeout: time.Duration(cfg.API.Timeout) * time.Second,
		},
		headers: make(map[string]string),
	}
}

func (c *Client) SetHeader(key, value string) {
	c.headers[key] = value
}

func (c *Client) SetAuthToken(token string) {
	c.apiKey = token
	c.headers["Authorization"] = "Bearer " + token
}

func (c *Client) Request(ctx context.Context, opts RequestOptions) (*Response, error) {
	if opts.Timeout == 0 {
		opts.Timeout = 30 * time.Second
	}

	url := c.baseURL + opts.Path

	var bodyReader io.Reader
	var contentType string

	if opts.Body != nil {
		data, err := json.Marshal(opts.Body)
		if err != nil {
			return nil, fmt.Errorf("序列化请求体失败: %w", err)
		}
		bodyReader = bytes.NewReader(data)
		contentType = "application/json"
	}

	req, err := http.NewRequest(opts.Method, url, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %w", err)
	}

	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", contentType)

	for k, v := range c.headers {
		req.Header.Set(k, v)
	}
	for k, v := range opts.Headers {
		req.Header.Set(k, v)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("请求失败: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应体失败: %w", err)
	}

	response := &Response{
		StatusCode: resp.StatusCode,
		Headers:    resp.Header.Clone(),
		Body:       body,
	}

	if resp.StatusCode >= 400 {
		return response, c.handleError(resp.StatusCode, body)
	}

	return response, nil
}

func (c *Client) handleError(statusCode int, body []byte) error {
	var errResp map[string]interface{}
	json.Unmarshal(body, &errResp)

	message := "未知错误"
	if msg, ok := errResp["message"].(string); ok {
		message = msg
	}

	switch statusCode {
	case 401:
		return fmt.Errorf("认证失败: %s", message)
	case 403:
		return fmt.Errorf("权限不足: %s", message)
	case 404:
		return fmt.Errorf("资源不存在: %s", message)
	case 409:
		return fmt.Errorf("同步冲突: %s", message)
	default:
		return fmt.Errorf("服务器错误 %d: %s", statusCode, message)
	}
}

func (c *Client) Get(ctx context.Context, path string) (*Response, error) {
	return c.Request(ctx, RequestOptions{
		Method: http.MethodGet,
		Path:   path,
	})
}

func (c *Client) Post(ctx context.Context, path string, body interface{}) (*Response, error) {
	return c.Request(ctx, RequestOptions{
		Method: http.MethodPost,
		Path:   path,
		Body:   body,
	})
}

func (c *Client) Put(ctx context.Context, path string, body interface{}) (*Response, error) {
	return c.Request(ctx, RequestOptions{
		Method: http.MethodPut,
		Path:   path,
		Body:   body,
	})
}

func (c *Client) Delete(ctx context.Context, path string) (*Response, error) {
	return c.Request(ctx, RequestOptions{
		Method: http.MethodDelete,
		Path:   path,
	})
}

func (c *Client) UploadAppPackage(ctx context.Context, req *UploadRequest) (string, error) {
	files := make([]FileEntry, 0, len(req.Files))
	for _, f := range req.Files {
		content := f.Content
		if content == nil && f.Path != "" {
			data, err := os.ReadFile(f.Path)
			if err != nil {
				return "", fmt.Errorf("读取文件失败: %w", err)
			}
			content = data
		}

		hash := f.Hash
		if hash == "" {
			hash = crypto.SHA256String(content)
		}

		files = append(files, FileEntry{
			Path:    f.Path,
			Content: content,
			Hash:    hash,
		})
	}

	body := map[string]interface{}{
		"appId":   req.AppID,
		"version": req.Version,
		"branch":  req.Branch,
		"message": req.Message,
		"author":  req.Author,
		"files":   files,
	}

	resp, err := c.Request(ctx, RequestOptions{
		Method: http.MethodPost,
		Path:   "/api/cli/app/upload",
		Body:   body,
	})
	if err != nil {
		return "", err
	}

	var result struct {
		Version string `json:"version"`
	}
	if err := json.Unmarshal(resp.Body, &result); err != nil {
		return "", fmt.Errorf("解析响应失败: %w", err)
	}

	return result.Version, nil
}

func (c *Client) DownloadAppPackage(ctx context.Context, appID, version, outputDir string) error {
	path := fmt.Sprintf("/api/cli/app/download?appId=%s&version=%s", appID, version)
	resp, err := c.Request(ctx, RequestOptions{
		Method: http.MethodGet,
		Path:   path,
	})
	if err != nil {
		return err
	}

	return os.WriteFile(filepath.Join(outputDir, "app.zip"), resp.Body, 0644)
}

func (c *Client) CheckConflict(ctx context.Context, appID, version string, files []FileEntry) (*ConflictCheckResult, error) {
	body := map[string]interface{}{
		"appId":   appID,
		"version": version,
		"files":   files,
	}

	resp, err := c.Request(ctx, RequestOptions{
		Method: http.MethodPost,
		Path:   "/api/cli/app/check-conflict",
		Body:   body,
	})
	if err != nil {
		return nil, err
	}

	var result ConflictCheckResult
	if err := json.Unmarshal(resp.Body, &result); err != nil {
		return nil, fmt.Errorf("解析响应失败: %w", err)
	}

	return &result, nil
}

func (c *Client) GetSyncStatus(ctx context.Context, appID string) (*SyncStatus, error) {
	path := fmt.Sprintf("/api/cli/app/status?appId=%s", appID)
	resp, err := c.Request(ctx, RequestOptions{
		Method: http.MethodGet,
		Path:   path,
	})
	if err != nil {
		return nil, err
	}

	var status SyncStatus
	if err := json.Unmarshal(resp.Body, &status); err != nil {
		return nil, fmt.Errorf("解析响应失败: %w", err)
	}

	return &status, nil
}

func (c *Client) Ping(ctx context.Context) error {
	_, err := c.Request(ctx, RequestOptions{
		Method:  http.MethodGet,
		Path:    "/health",
		Timeout: 10 * time.Second,
	})
	return err
}
