package cmd

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/geelato/cli/pkg/logger"
	"github.com/spf13/cobra"
)

var (
	cloneOutput      string
	cloneVersion     string
	cloneSkipExtract bool
)

var cloneCmd = &cobra.Command{
	Use:   "clone <url>",
	Short: "clone(从服务器克隆应用)",
	Long: `从 Geelato 服务器克隆应用，包括模型定义、API 接口、工作流和页面配置。

URL 格式：
  http://{host}:{port}/{tenant}/{app-code}

克隆过程：
1. 解析 URL 提取 tenant 和 appCode
2. 连接 Geelato 服务器
3. 下载应用数据
4. 渲染模型文件（define.json, columns.json, check.json, fk.json, view.sql）
5. 渲染页面文件
6. 保存到本地目录

示例：
  geelato clone http://localhost:8080/default/myapp       # 克隆应用到 myapp 目录
  geelato clone http://localhost:8080/mytenant/myapp      # 指定租户
  geelato clone http://localhost:8080/default/myapp -o ./projects  # 指定输出目录`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		repoURL := args[0]
		return runClone(repoURL, cloneOutput, cloneVersion, cloneSkipExtract)
	},
}

func init() {
	cloneCmd.Flags().StringVarP(&cloneOutput, "output", "o", "", "Output directory (default: app code)")
	cloneCmd.Flags().StringVar(&cloneVersion, "version", "latest", "App version")
	cloneCmd.Flags().BoolVar(&cloneSkipExtract, "skip-extract", false, "Skip extracting zip file")
}

type CloneResult struct {
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
	ID             string `json:"id"`
	AppID          string `json:"appId"`
	Code           string `json:"code"`
	Name           string `json:"name"`
	Module         string `json:"module"`
	GroupName      string `json:"groupName"`
	Description    string `json:"description"`
	SourceContent  string `json:"sourceContent"`
	ReleaseContent string `json:"releaseContent"`
	Method         string `json:"method"`
	Path           string `json:"path"`
	ResponseType   string `json:"responseType"`
	ResponseFormat string `json:"responseFormat"`
	Version        int    `json:"version"`
	EnableStatus   int    `json:"enableStatus"`
	Anonymous      int    `json:"anonymous"`
	Paging         int    `json:"paging"`
}

type WorkflowData struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	FileName   string `json:"fileName"`
	Category   string `json:"category"`
	Definition string `json:"definition"`
}

func runClone(repoURL, outputDir, version string, skipExtract bool) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	tenant, appCode, apiURL, err := parseCloneURL(repoURL)
	if err != nil {
		return fmt.Errorf("invalid URL format: %w", err)
	}

	logger.Infof("Parsing URL: tenant=%s, appCode=%s, apiURL=%s", tenant, appCode, apiURL)

	if outputDir == "" {
		outputDir = appCode
	}

	logger.Infof("Cloning app '%s' to '%s'...", appCode, outputDir)

	body := map[string]interface{}{
		"appCode": appCode,
		"tenant":  tenant,
		"version": version,
	}
	bodyBytes, _ := json.Marshal(body)

	cloneURL := fmt.Sprintf("%s/api/cli/app/clone", apiURL)
	logger.Infof("Requesting: %s", cloneURL)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, cloneURL, bytes.NewReader(bodyBytes))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 5 * time.Minute}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("clone failed with status %d: %s", resp.StatusCode, string(respBody))
	}

	var result struct {
		Code int             `json:"code"`
		Data json.RawMessage `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
	}

	logger.Infof("Response code: %d, data length: %d", result.Code, len(result.Data))

	// 解析外层结构，检查是否有 format 字段
	var outerData struct {
		Format string          `json:"format"`
		Data   json.RawMessage `json:"data"`
	}
	if err := json.Unmarshal(result.Data, &outerData); err != nil {
		return fmt.Errorf("failed to parse outer data: %w", err)
	}

	// 如果有 format 字段，说明数据被包装了一层
	var appDataBytes []byte
	if outerData.Format != "" && outerData.Data != nil {
		logger.Infof("Found nested format: %s", outerData.Format)
		appDataBytes = outerData.Data
	} else {
		appDataBytes = result.Data
	}

	var appData CloneResponseData
	if err := json.Unmarshal(appDataBytes, &appData); err != nil {
		return fmt.Errorf("failed to parse app data: %w", err)
	}

	logger.Infof("Parsed: entities=%d, pages=%d, apis=%d, workflows=%d",
		len(appData.Entities), len(appData.Pages), len(appData.APIs), len(appData.Workflows))

	manager := NewCloneManager(appCode, repoURL)
	if err := manager.RenderAndSave(ctx, &appData, outputDir); err != nil {
		return fmt.Errorf("failed to render and save: %w", err)
	}

	logger.Success("Clone completed successfully!")
	logger.Infof("App cloned to: %s", outputDir)
	return nil
}

func parseCloneURL(repoURL string) (tenant, appCode, apiURL string, err error) {
	repoURL = strings.TrimSpace(repoURL)

	if !strings.HasPrefix(repoURL, "http://") && !strings.HasPrefix(repoURL, "https://") {
		repoURL = "http://" + repoURL
	}

	u, err := url.Parse(repoURL)
	if err != nil {
		return "", "", "", fmt.Errorf("failed to parse URL: %w", err)
	}

	path := strings.TrimPrefix(u.Path, "/")
	path = strings.TrimSuffix(path, "/")

	parts := strings.Split(path, "/")
	if len(parts) < 2 {
		return "", "", "", fmt.Errorf("URL path should contain tenant and app code (e.g., /tenant/app-code)")
	}

	tenant = parts[0]
	appCode = parts[1]

	if tenant == "" || appCode == "" {
		return "", "", "", fmt.Errorf("tenant and app code cannot be empty")
	}

	apiURL = u.Scheme + "://" + u.Host

	return tenant, appCode, apiURL, nil
}

type CloneManager struct {
	AppCode string
	RepoURL string
}

func NewCloneManager(appCode string, repoURL string) *CloneManager {
	return &CloneManager{AppCode: appCode, RepoURL: repoURL}
}

func (m *CloneManager) RenderAndSave(ctx context.Context, data *CloneResponseData, outputDir string) error {
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	if err := m.renderAppConfig(outputDir, data); err != nil {
		return fmt.Errorf("failed to render app config: %w", err)
	}

	logger.Infof("Processing entities: %d", len(data.Entities))
	for _, entity := range data.Entities {
		entityDir := filepath.Join(outputDir, "meta", entity.EntityName)
		if err := os.MkdirAll(entityDir, 0755); err != nil {
			return fmt.Errorf("failed to create entity directory: %w", err)
		}

		if err := m.renderEntity(entity, entityDir); err != nil {
			return fmt.Errorf("failed to render entity %s: %w", entity.EntityName, err)
		}
	}

	logger.Infof("Processing pages: %d", len(data.Pages))
	// 确保 page 目录存在，即使数据为空
	pageDir := filepath.Join(outputDir, "page")
	if err := os.MkdirAll(pageDir, 0755); err != nil {
		return fmt.Errorf("failed to create page directory: %w", err)
	}

	// 如果 page 数据为空，创建 .gitkeep 占位文件
	if len(data.Pages) == 0 {
		gitkeepPath := filepath.Join(pageDir, ".gitkeep")
		if _, err := os.Stat(gitkeepPath); os.IsNotExist(err) {
			os.WriteFile(gitkeepPath, []byte("# Keep this directory\n"), 0644)
		}
	}

	for _, page := range data.Pages {
		// 使用 code 作为目录名，如果为空则使用 id
		pageDirName := page.Code
		if pageDirName == "" {
			pageDirName = page.ID
		}
		pageSubDir := filepath.Join(outputDir, "page", pageDirName)
		if err := os.MkdirAll(pageSubDir, 0755); err != nil {
			logger.Errorf("Failed to create page dir %s: %v", pageSubDir, err)
			continue
		}

		if err := m.renderPage(page, pageSubDir); err != nil {
			logger.Errorf("Failed to render page %s: %v", page.Code, err)
			continue
		}
	}

	logger.Infof("Processing APIs: %d", len(data.APIs))
	// 确保 api 目录存在，即使数据为空
	apiDir := filepath.Join(outputDir, "api")
	if err := os.MkdirAll(apiDir, 0755); err != nil {
		return fmt.Errorf("failed to create api directory: %w", err)
	}

	// 如果 api 数据为空，创建 .gitkeep 占位文件
	if len(data.APIs) == 0 {
		gitkeepPath := filepath.Join(apiDir, ".gitkeep")
		if _, err := os.Stat(gitkeepPath); os.IsNotExist(err) {
			os.WriteFile(gitkeepPath, []byte("# Keep this directory\n"), 0644)
		}
	}

	for _, api := range data.APIs {
		// 使用 code 作为目录名，如果为空则使用 id
		apiDirName := api.Code
		if apiDirName == "" {
			apiDirName = api.ID
		}
		apiSubDir := filepath.Join(outputDir, "api", apiDirName)
		if err := os.MkdirAll(apiSubDir, 0755); err != nil {
			logger.Errorf("Failed to create api dir %s: %v", apiSubDir, err)
			continue
		}

		if err := m.renderAPI(api, apiSubDir); err != nil {
			logger.Errorf("Failed to render api %s: %v", api.Code, err)
			continue
		}
	}

	logger.Infof("Processing workflows: %d", len(data.Workflows))
	wfDir := filepath.Join(outputDir, "workflow")
	if err := os.MkdirAll(wfDir, 0755); err != nil {
		return fmt.Errorf("failed to create workflow directory: %w", err)
	}

	for _, wf := range data.Workflows {
		if err := m.renderWorkflow(wf, wfDir); err != nil {
			return fmt.Errorf("failed to render workflow %s: %w", wf.Name, err)
		}
	}

	return nil
}

func (m *CloneManager) renderAppConfig(outputDir string, data *CloneResponseData) error {
	config := map[string]interface{}{
		"meta": map[string]interface{}{
			"appId":   data.AppID,
			"appCode": data.AppCode,
			"name":    data.AppName,
			"version": data.Version,
			"tenant":  data.Tenant,
		},
		"config": map[string]interface{}{
			"api": map[string]interface{}{
				"url":     "",
				"timeout": 30,
			},
			"repo": map[string]interface{}{
				"url": m.RepoURL,
			},
			"sync": map[string]interface{}{
				"autoPush": false,
				"autoPull": false,
			},
		},
	}

	content, _ := json.MarshalIndent(config, "", "  ")
	return os.WriteFile(filepath.Join(outputDir, "geelato.json"), content, 0644)
}

func (m *CloneManager) renderEntity(entity EntityData, entityDir string) error {
	content, _ := json.MarshalIndent(entity.Define, "", "  ")
	if err := os.WriteFile(filepath.Join(entityDir, entity.EntityName+".define.json"), content, 0644); err != nil {
		return err
	}

	if len(entity.Columns) > 0 {
		content, _ := json.MarshalIndent(map[string]interface{}{
			"meta":    entity.Meta,
			"columns": entity.Columns,
		}, "", "  ")
		if err := os.WriteFile(filepath.Join(entityDir, entity.EntityName+".columns.json"), content, 0644); err != nil {
			return err
		}
	}

	if len(entity.Checks) > 0 {
		content, _ := json.MarshalIndent(map[string]interface{}{
			"meta":   entity.Meta,
			"checks": entity.Checks,
		}, "", "  ")
		if err := os.WriteFile(filepath.Join(entityDir, entity.EntityName+".check.json"), content, 0644); err != nil {
			return err
		}
	}

	if len(entity.ForeignKeys) > 0 {
		content, _ := json.MarshalIndent(map[string]interface{}{
			"meta":        entity.Meta,
			"foreignKeys": entity.ForeignKeys,
		}, "", "  ")
		if err := os.WriteFile(filepath.Join(entityDir, entity.EntityName+".fk.json"), content, 0644); err != nil {
			return err
		}
	}

	for _, view := range entity.Views {
		viewContent := m.renderViewSQL(view)
		viewFileName := fmt.Sprintf("%s.%s.view.sql", entity.EntityName, view.ViewName)
		if err := os.WriteFile(filepath.Join(entityDir, viewFileName), []byte(viewContent), 0644); err != nil {
			return err
		}
	}

	return nil
}

func (m *CloneManager) renderViewSQL(view ViewData) string {
	var sb strings.Builder
	sb.WriteString("-- @meta\n")
	sb.WriteString(fmt.Sprintf("-- @id view_%s_%s\n", strings.ToLower(m.AppCode), view.ViewName))
	sb.WriteString(fmt.Sprintf("-- @appId %s\n", view.AppID))
	sb.WriteString(fmt.Sprintf("-- @connectId %s\n", view.ConnectID))
	sb.WriteString(fmt.Sprintf("-- @description %s\n", view.Description))
	sb.WriteString(fmt.Sprintf("-- @title %s\n", view.Title))
	sb.WriteString(fmt.Sprintf("-- @viewName %s.%s\n", view.EntityName, view.ViewName))
	sb.WriteString(fmt.Sprintf("-- @viewType %s\n", view.ViewType))
	sb.WriteString(fmt.Sprintf("-- @entityName %s\n", view.EntityName))
	sb.WriteString(fmt.Sprintf("-- @enableStatus %d\n", view.EnableStatus))
	sb.WriteString(fmt.Sprintf("-- @linked %d\n", view.Linked))
	sb.WriteString(fmt.Sprintf("-- @delStatus %d\n", view.DelStatus))
	sb.WriteString(fmt.Sprintf("-- @seqNo %d\n", view.SeqNo))
	sb.WriteString("\n")

	if len(view.ViewColumns) > 0 {
		var cols []string
		for _, col := range view.ViewColumns {
			cols = append(cols, fmt.Sprintf("t.%s", col.ColumnName))
		}
		sb.WriteString(fmt.Sprintf("SELECT\n  %s\n", strings.Join(cols, ",\n  ")))
	} else {
		sb.WriteString("SELECT\n  t.id\n")
	}
	sb.WriteString(fmt.Sprintf("FROM %s t\n", view.TableName))
	sb.WriteString("WHERE t.del_status = 0\n")
	if view.OrderBy != "" {
		sb.WriteString(fmt.Sprintf("ORDER BY %s\n", view.OrderBy))
	}

	return sb.String()
}

func (m *CloneManager) renderPage(page PageData, pageDir string) error {
	// 使用 code 作为文件名前缀，如果为空则使用 id
	filePrefix := page.Code
	if filePrefix == "" {
		filePrefix = page.ID
	}

	define := map[string]interface{}{
		"meta": map[string]interface{}{
			"version":   "1.0.0",
			"createdAt": page.CreatedAt,
		},
		"page": map[string]interface{}{
			"id":            page.ID,
			"appId":         page.AppID,
			"extendId":      page.ExtendID,
			"code":          page.Code,
			"description":   page.Description,
			"type":          page.Type,
			"title":         page.Title,
			"checkStatus":   page.CheckStatus,
			"checkUserId":   page.CheckUserID,
			"checkUserName": page.CheckUserName,
			"checkAt":       page.CheckAt,
			"version":       page.Version,
			"seqNo":         page.SeqNo,
		},
	}

	content, _ := json.MarshalIndent(define, "", "  ")
	if err := os.WriteFile(filepath.Join(pageDir, filePrefix+".define.json"), content, 0644); err != nil {
		return err
	}

	if page.SourceContent != "" {
		if err := os.WriteFile(filepath.Join(pageDir, filePrefix+".source.json"), []byte(page.SourceContent), 0644); err != nil {
			return err
		}
	}

	if page.ReleaseContent != "" {
		if err := os.WriteFile(filepath.Join(pageDir, filePrefix+".release.json"), []byte(page.ReleaseContent), 0644); err != nil {
			return err
		}
	}

	if page.PreviewContent != "" {
		if err := os.WriteFile(filepath.Join(pageDir, filePrefix+".preview.json"), []byte(page.PreviewContent), 0644); err != nil {
			return err
		}
	}

	return nil
}

func (m *CloneManager) renderAPI(api APIData, apiDir string) error {
	filePrefix := api.Code
	if filePrefix == "" {
		filePrefix = api.ID
	}

	define := map[string]interface{}{
		"meta": map[string]interface{}{
			"version":   "1.0.0",
			"updatedAt": time.Now().Format(time.RFC3339),
		},
		"api": map[string]interface{}{
			"id":             api.ID,
			"appId":          api.AppID,
			"code":           api.Code,
			"name":           api.Name,
			"module":         api.Module,
			"groupName":      api.GroupName,
			"description":    api.Description,
			"method":         api.Method,
			"path":           api.Path,
			"responseType":   api.ResponseType,
			"responseFormat": api.ResponseFormat,
			"version":        api.Version,
			"enableStatus":   api.EnableStatus,
			"anonymous":      api.Anonymous,
			"paging":         api.Paging,
		},
	}

	content, _ := json.MarshalIndent(define, "", "  ")
	if err := os.WriteFile(filepath.Join(apiDir, filePrefix+".define.json"), content, 0644); err != nil {
		return err
	}

	if api.ReleaseContent != "" {
		jsContent := fmt.Sprintf(`/**
 * @api
 * @name %s
 * @path %s
 * @method %s
 * @description %s
 * @version %d
 */

%s
`, api.Name, api.Path, api.Method, api.Description, api.Version, api.ReleaseContent)
		if err := os.WriteFile(filepath.Join(apiDir, filePrefix+".api.js"), []byte(jsContent), 0644); err != nil {
			return err
		}
	}

	return nil
}

func (m *CloneManager) renderWorkflow(wf WorkflowData, wfDir string) error {
	content, _ := json.MarshalIndent(wf.Definition, "", "  ")
	return os.WriteFile(filepath.Join(wfDir, wf.FileName), content, 0644)
}
