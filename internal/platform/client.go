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
	baseURL string
	apiKey  string
	client  *http.Client
	headers map[string]string
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
	AppID   string
	Version string
	Branch  string
	Message string
	Author  string
	Files   []FileEntry
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

func NewClientWithURL(apiURL string) *Client {
	timeout := 30
	if cfg := config.Get(); cfg != nil {
		timeout = cfg.API.Timeout
	}
	return &Client{
		baseURL: apiURL,
		apiKey:  "",
		client: &http.Client{
			Timeout: time.Duration(timeout) * time.Second,
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

type CloneRequest struct {
	AppCode string `json:"appCode"`
	Tenant  string `json:"tenant"`
	Version string `json:"version"`
}

type CloneResponse struct {
	Format string            `json:"format"`
	Data   CloneResponseData `json:"data"`
}

type CloneResponseData struct {
	AppID     string         `json:"appId"`
	AppCode   string         `json:"appCode"`
	AppName   string         `json:"name"`
	Version   string         `json:"version"`
	Tenant    string         `json:"tenant"`
	Entities  []EntityData   `json:"entities"`
	Pages     []PageData     `json:"pages"`
	APIs      []APIData      `json:"apis"`
	Workflows []WorkflowData `json:"workflows"`
}

type EntityData struct {
	EntityName  string                 `json:"entityName"`
	TableName   string                 `json:"tableName"`
	TableID     string                 `json:"tableId"`
	Meta        map[string]interface{} `json:"meta"`
	Define      map[string]interface{} `json:"define"`
	Columns     []ColumnData           `json:"columns"`
	Checks      []CheckData            `json:"checks"`
	ForeignKeys []ForeignKeyData       `json:"foreignKeys"`
	Views       []ViewData             `json:"views"`
}

type ColumnData struct {
	ID                     string `json:"id"`
	AppID                  string `json:"appId"`
	TableID                string `json:"tableId"`
	FieldName              string `json:"fieldName"`
	ColumnName             string `json:"columnName"`
	ColumnType             string `json:"columnType"`
	DataType               string `json:"dataType"`
	CharacterMaxinumLength int    `json:"characterMaxinumLength"`
	DatetimePrecision      string `json:"datetimePrecision"`
	OrdinalPosition        int    `json:"ordinalPosition"`
	IsNullable             bool   `json:"isNullable"`
	IsUnique               bool   `json:"isUnique"`
	ColumnDefault          string `json:"columnDefault"`
	Description            string `json:"description"`
	ColumnComment          string `json:"columnComment"`
	Title                  string `json:"title"`
	ColumnKey              string `json:"columnKey"`
	NumericPrecision       int    `json:"numericPrecision"`
	NumericScale           int    `json:"numericScale"`
	EnableStatus           int    `json:"enableStatus"`
	Synced                 int    `json:"synced"`
	DelStatus              int    `json:"delStatus"`
	SeqNo                  int    `json:"seqNo"`
	Drawed                 int    `json:"drawed"`
}

type CheckData struct {
	ID           string `json:"id"`
	AppID        string `json:"appId"`
	TableID      string `json:"tableId"`
	TableName    string `json:"tableName"`
	ColumnID     string `json:"columnId"`
	ColumnName   string `json:"columnName"`
	Code         string `json:"code"`
	Title        string `json:"title"`
	Type         string `json:"type"`
	CheckClause  string `json:"checkClause"`
	Description  string `json:"description"`
	EnableStatus int    `json:"enableStatus"`
	Synced       int    `json:"synced"`
	DelStatus    int    `json:"delStatus"`
	SeqNo        int    `json:"seqNo"`
}

type ForeignKeyData struct {
	ID              string `json:"id"`
	AppID           string `json:"appId"`
	MainTableID     string `json:"mainTableId"`
	MainTable       string `json:"mainTable"`
	MainTableCol    string `json:"mainTableCol"`
	ForeignTableID  string `json:"foreignTableId"`
	ForeignTable    string `json:"foreignTable"`
	ForeignTableCol string `json:"foreignTableCol"`
	DeleteAction    string `json:"deleteAction"`
	UpdateAction    string `json:"updateAction"`
	Description     string `json:"description"`
	EnableStatus    int    `json:"enableStatus"`
	DelStatus       int    `json:"delStatus"`
	SeqNo           int    `json:"seqNo"`
}

type ViewData struct {
	ID            string           `json:"id"`
	AppID         string           `json:"appId"`
	ConnectID     string           `json:"connectId"`
	ViewName      string           `json:"viewName"`
	ViewType      string           `json:"viewType"`
	EntityName    string           `json:"entityName"`
	TableName     string           `json:"tableName"`
	Description   string           `json:"description"`
	Title         string           `json:"title"`
	ViewConstruct string           `json:"viewConstruct"`
	ViewColumns   []ViewColumnData `json:"viewColumns"`
	EnableStatus  int              `json:"enableStatus"`
	Linked        int              `json:"linked"`
	DelStatus     int              `json:"delStatus"`
	SeqNo         int              `json:"seqNo"`
	OrderBy       string           `json:"orderBy"`
}

type ViewColumnData struct {
	ColumnName string `json:"columnName"`
	Alias      string `json:"alias"`
}

type PageData struct {
	ID             string `json:"id"`
	AppID          string `json:"appId"`
	PageName       string `json:"pageName"`
	Code           string `json:"code"`
	Description    string `json:"description"`
	Type           string `json:"type"`
	Title          string `json:"title"`
	SourceContent  string `json:"sourceContent"`
	ReleaseContent string `json:"releaseContent"`
	PreviewContent string `json:"previewContent"`
	ExtendID       string `json:"extendId"`
	CheckStatus    string `json:"checkStatus"`
	CheckUserID    string `json:"checkUserId"`
	CheckUserName  string `json:"checkUserName"`
	CheckAt        string `json:"checkAt"`
	Version        int    `json:"version"`
	SeqNo          int    `json:"seqNo"`
	CreatedAt      string `json:"createdAt"`
}

type APIData struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Module      string `json:"module"`
	Description string `json:"description"`
	Content     string `json:"content"`
	Method      string `json:"method"`
	Path        string `json:"path"`
}

type WorkflowData struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	FileName   string `json:"fileName"`
	Category   string `json:"category"`
	Definition string `json:"definition"`
}

func (c *Client) CloneApp(ctx context.Context, req *CloneRequest) (*CloneResponse, error) {
	body := map[string]interface{}{
		"appCode": req.AppCode,
		"tenant":  req.Tenant,
		"version": req.Version,
	}

	resp, err := c.Request(ctx, RequestOptions{
		Method: http.MethodPost,
		Path:   "/api/cli/app/clone",
		Body:   body,
	})
	if err != nil {
		return nil, err
	}

	var result CloneResponse
	if err := json.Unmarshal(resp.Body, &result); err != nil {
		return nil, fmt.Errorf("解析响应失败: %w", err)
	}

	return &result, nil
}

func (c *Client) DownloadAppPackageBytes(ctx context.Context, appID, version string) ([]byte, error) {
	path := fmt.Sprintf("/api/cli/app/download?appId=%s&version=%s", appID, version)
	resp, err := c.Request(ctx, RequestOptions{
		Method: http.MethodGet,
		Path:   path,
	})
	if err != nil {
		return nil, err
	}

	return resp.Body, nil
}
