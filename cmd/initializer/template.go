package initializer

import (
	"bytes"
	"embed"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"time"

	"github.com/geelato/cli/pkg/logger"
)

//go:embed templates/* templates/**/* templates/**/**/* templates/**/**/**/* templates/**/**/**/**/*
var templatesFS embed.FS

// TemplateData holds the data for template rendering
type TemplateData struct {
	AppName   string
	CreatedAt string
	Repo      string
}

// TemplateManager manages embedded templates
type TemplateManager struct {
	fs embed.FS
}

// NewTemplateManager creates a new template manager
func NewTemplateManager() *TemplateManager {
	return &TemplateManager{
		fs: templatesFS,
	}
}

// RenderTemplate renders a template with the given data
func (tm *TemplateManager) RenderTemplate(templatePath string, data TemplateData) (string, error) {
	content, err := tm.fs.ReadFile(templatePath)
	if err != nil {
		return "", fmt.Errorf("failed to read template %s: %w", templatePath, err)
	}

	tmpl, err := template.New("template").Parse(string(content))
	if err != nil {
		return "", fmt.Errorf("failed to parse template %s: %w", templatePath, err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("failed to execute template %s: %w", templatePath, err)
	}

	return buf.String(), nil
}

// CopyStaticFile copies a static file from embedded FS to destination
func (tm *TemplateManager) CopyStaticFile(srcPath, destPath string) error {
	content, err := tm.fs.ReadFile(srcPath)
	if err != nil {
		return fmt.Errorf("failed to read file %s: %w", srcPath, err)
	}

	if err := os.MkdirAll(filepath.Dir(destPath), 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	if err := os.WriteFile(destPath, content, 0644); err != nil {
		return fmt.Errorf("failed to write file %s: %w", destPath, err)
	}

	return nil
}

// WalkTemplates walks through the templates directory and executes a callback for each file
func (tm *TemplateManager) WalkTemplates(root string, callback func(path string, isDir bool) error) error {
	return fs.WalkDir(tm.fs, root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		return callback(path, d.IsDir())
	})
}

// InitializeApp initializes a new application using embedded templates
func InitializeApp(appName string, repo string) (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("failed to get current directory: %w", err)
	}

	appDir := filepath.Join(cwd, appName)

	if _, err := os.Stat(appDir); err == nil {
		return "", fmt.Errorf("directory '%s' already exists", appName)
	}

	if err := os.MkdirAll(appDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create app directory: %w", err)
	}

	tm := NewTemplateManager()
	data := TemplateData{
		AppName:   appName,
		CreatedAt: time.Now().Format(time.RFC3339),
		Repo:      repo,
	}

	// Create main app files
	appTemplates := map[string]string{
		"templates/app/geelato.json.tmpl": "geelato.json",
		"templates/app/README.md.tmpl":    "README.md",
		"templates/app/.gitignore.tmpl":   ".gitignore",
	}

	for templatePath, destName := range appTemplates {
		content, err := tm.RenderTemplate(templatePath, data)
		if err != nil {
			os.RemoveAll(appDir)
			return "", fmt.Errorf("failed to render template %s: %w", templatePath, err)
		}

		destPath := filepath.Join(appDir, destName)
		if err := os.WriteFile(destPath, []byte(content), 0644); err != nil {
			os.RemoveAll(appDir)
			return "", fmt.Errorf("failed to create file %s: %w", destName, err)
		}
	}

	// Create directories
	dirs := []string{
		"doc",
		"example/user-management-app/meta/User",
		"example/user-management-app/meta/Department",
		"example/user-management-app/api/user",
		"example/user-management-app/page",
		"example/user-management-app/workflow",
		"meta",
		"api",
		"page",
		"workflow",
	}

	for _, dir := range dirs {
		fullPath := filepath.Join(appDir, dir)
		if err := os.MkdirAll(fullPath, 0755); err != nil {
			os.RemoveAll(appDir)
			return "", fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}

	// Download doc files from GitHub
	docFiles := map[string]string{
		"doc/实体模型定义规范.md":     "geelato-cli/templates/实体模型定义规范.md",
		"doc/API脚本编写指引和规范.md": "geelato-cli/templates/API脚本编写指引和规范.md",
		"doc/工作流流程编排指引和规范.md": "geelato-cli/templates/工作流流程编排指引和规范.md",
		"doc/MCP配置指引.md":      "geelato-cli/templates/MCP配置指引.md",
	}

	for destPath, srcPath := range docFiles {
		filePath := filepath.Join(appDir, destPath)
		if err := downloadFromGitHub(filePath, srcPath); err != nil {
			logger.Warn("Failed to download %s: %v", filepath.Base(destPath), err)
			if err := createDownloadNotice(filePath); err != nil {
				os.RemoveAll(appDir)
				return "", fmt.Errorf("failed to create doc file %s: %w", destPath, err)
			}
		} else {
			logger.Info("Downloaded %s", filepath.Base(destPath))
		}
	}

	// Create example app files
	exampleTemplates := map[string]string{
		"templates/example/user-management-app/.gitignore.tmpl":   "example/user-management-app/.gitignore",
		"templates/example/user-management-app/README.md.tmpl":    "example/user-management-app/README.md",
		"templates/example/user-management-app/geelato.json.tmpl": "example/user-management-app/geelato.json",
	}

	for templatePath, destPath := range exampleTemplates {
		content, err := tm.RenderTemplate(templatePath, data)
		if err != nil {
			os.RemoveAll(appDir)
			return "", fmt.Errorf("failed to render template %s: %w", templatePath, err)
		}

		fullDestPath := filepath.Join(appDir, destPath)
		if err := os.WriteFile(fullDestPath, []byte(content), 0644); err != nil {
			os.RemoveAll(appDir)
			return "", fmt.Errorf("failed to create example file %s: %w", destPath, err)
		}
	}

	// Copy static API files
	apiFiles := map[string]string{
		"templates/example/user-management-app/api/user/getList.api.js.tmpl":   "example/user-management-app/api/user/getList.api.js",
		"templates/example/user-management-app/api/user/getDetail.api.js.tmpl": "example/user-management-app/api/user/getDetail.api.js",
		"templates/example/user-management-app/api/user/saveUser.api.js.tmpl":  "example/user-management-app/api/user/saveUser.api.js",
	}

	for srcPath, destPath := range apiFiles {
		fullDestPath := filepath.Join(appDir, destPath)
		if err := tm.CopyStaticFile(srcPath, fullDestPath); err != nil {
			os.RemoveAll(appDir)
			return "", fmt.Errorf("failed to copy API file %s: %w", destPath, err)
		}
	}

	// Copy static meta files
	metaFiles := map[string]string{
		"templates/example/user-management-app/meta/User/User.define.json.tmpl":                  "example/user-management-app/meta/User/User.define.json",
		"templates/example/user-management-app/meta/User/User.columns.json.tmpl":                 "example/user-management-app/meta/User/User.columns.json",
		"templates/example/user-management-app/meta/User/User.check.json.tmpl":                   "example/user-management-app/meta/User/User.check.json",
		"templates/example/user-management-app/meta/User/User.fk.json.tmpl":                      "example/user-management-app/meta/User/User.fk.json",
		"templates/example/user-management-app/meta/User/User.default.view.sql.tmpl":             "example/user-management-app/meta/User/User.default.view.sql",
		"templates/example/user-management-app/meta/Department/Department.define.json.tmpl":      "example/user-management-app/meta/Department/Department.define.json",
		"templates/example/user-management-app/meta/Department/Department.columns.json.tmpl":     "example/user-management-app/meta/Department/Department.columns.json",
		"templates/example/user-management-app/meta/Department/Department.check.json.tmpl":       "example/user-management-app/meta/Department/Department.check.json",
		"templates/example/user-management-app/meta/Department/Department.fk.json.tmpl":          "example/user-management-app/meta/Department/Department.fk.json",
		"templates/example/user-management-app/meta/Department/Department.default.view.sql.tmpl": "example/user-management-app/meta/Department/Department.default.view.sql",
	}

	for srcPath, destPath := range metaFiles {
		fullDestPath := filepath.Join(appDir, destPath)
		if err := tm.CopyStaticFile(srcPath, fullDestPath); err != nil {
			os.RemoveAll(appDir)
			return "", fmt.Errorf("failed to copy meta file %s: %w", destPath, err)
		}
	}

	// Create .gitkeep files
	gitkeepContent := "# Keep this directory\n"
	gitkeepFiles := []string{
		"api/.gitkeep",
		"meta/.gitkeep",
		"page/.gitkeep",
		"workflow/.gitkeep",
	}

	for _, path := range gitkeepFiles {
		filePath := filepath.Join(appDir, path)
		if err := os.WriteFile(filePath, []byte(gitkeepContent), 0644); err != nil {
			os.RemoveAll(appDir)
			return "", fmt.Errorf("failed to create .gitkeep file %s: %w", path, err)
		}
	}

	return appDir, nil
}

const githubRawURL = "https://raw.githubusercontent.com/geelato/geelato-doc/main"

func downloadFromGitHub(destPath, srcPath string) error {
	url := fmt.Sprintf("%s/%s", githubRawURL, srcPath)

	logger.Info("Downloading %s...", filepath.Base(destPath))

	client := &http.Client{Timeout: 30 * time.Second}
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to download: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("HTTP %d", resp.StatusCode)
	}

	if err := os.MkdirAll(filepath.Dir(destPath), 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	file, err := os.Create(destPath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}

// ModelTemplateData holds data for model template rendering
type ModelTemplateData struct {
	ModelName       string
	EntityName      string
	EntityNameLower string
	EntityNameUpper string
	TableName       string
	TableID         string
	CreatedAt       string
	AppID           string
}

// APITemplateData holds data for API template rendering
type APITemplateData struct {
	APIName string
	APIPath string
}

// WorkflowTemplateData holds data for workflow template rendering
type WorkflowTemplateData struct {
	WorkflowName string
	WorkflowDesc string
	CreatedAt    string
	UpdatedAt    string
}

// ColumnTemplateData holds data for column template rendering
type ColumnTemplateData struct {
	EntityName        string
	EntityNameLower   string
	FieldName         string
	FieldNameLower    string
	AppID             string
	TableID           string
	TableName         string
	ColumnName        string
	ColumnType        string
	DataType          string
	Length            int
	DateTimePrecision string
	OrdinalPosition   int
	Comment           string
}

// RenderColumnTemplate renders a column template with the given data
func (tm *TemplateManager) RenderColumnTemplate(templatePath string, data ColumnTemplateData) (string, error) {
	content, err := tm.fs.ReadFile(templatePath)
	if err != nil {
		return "", fmt.Errorf("failed to read template %s: %w", templatePath, err)
	}

	tmpl, err := template.New("column").Parse(string(content))
	if err != nil {
		return "", fmt.Errorf("failed to parse template %s: %w", templatePath, err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("failed to execute template %s: %w", templatePath, err)
	}

	return buf.String(), nil
}

// FkTemplateData holds data for FK template rendering
type FkTemplateData struct {
	EntityName         string
	EntityNameLower    string
	FieldName          string
	AppID              string
	TableID            string
	TableName          string
	ForeignTableSchema string
	ForeignTableID     string
	ForeignTable       string
	Description        string
	SeqNo              int
}

// CheckTemplateData holds data for check template rendering
type CheckTemplateData struct {
	EntityName      string
	EntityNameLower string
	CheckID         string
	Title           string
	Code            string
	TableID         string
	TableName       string
	Type            string
	CheckClause     string
	Description     string
	AppID           string
}

// ViewTemplateData holds data for view template rendering
type ViewTemplateData struct {
	EntityName      string
	EntityNameLower string
	ViewName        string
	ViewNameLower   string
	AppID           string
	Description     string
	Title           string
	TableName       string
	SelectColumns   string
	OrderBy         string
	SeqNo           int
}

// RenderFkTemplate renders an FK template with the given data
func (tm *TemplateManager) RenderFkTemplate(templatePath string, data FkTemplateData) (string, error) {
	content, err := tm.fs.ReadFile(templatePath)
	if err != nil {
		return "", fmt.Errorf("failed to read template %s: %w", templatePath, err)
	}

	tmpl, err := template.New("fk").Parse(string(content))
	if err != nil {
		return "", fmt.Errorf("failed to parse template %s: %w", templatePath, err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("failed to execute template %s: %w", templatePath, err)
	}

	return buf.String(), nil
}

// RenderCheckTemplate renders a check template with the given data
func (tm *TemplateManager) RenderCheckTemplate(templatePath string, data CheckTemplateData) (string, error) {
	content, err := tm.fs.ReadFile(templatePath)
	if err != nil {
		return "", fmt.Errorf("failed to read template %s: %w", templatePath, err)
	}

	tmpl, err := template.New("check").Parse(string(content))
	if err != nil {
		return "", fmt.Errorf("failed to parse template %s: %w", templatePath, err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("failed to execute template %s: %w", templatePath, err)
	}

	return buf.String(), nil
}

// RenderViewTemplate renders a view template with the given data
func (tm *TemplateManager) RenderViewTemplate(templatePath string, data ViewTemplateData) (string, error) {
	content, err := tm.fs.ReadFile(templatePath)
	if err != nil {
		return "", fmt.Errorf("failed to read template %s: %w", templatePath, err)
	}

	tmpl, err := template.New("view").Parse(string(content))
	if err != nil {
		return "", fmt.Errorf("failed to parse template %s: %w", templatePath, err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("failed to execute template %s: %w", templatePath, err)
	}

	return buf.String(), nil
}

// RenderPageTemplate renders a page template with the given data
func (tm *TemplateManager) RenderPageTemplate(templatePath string, data PageTemplateData) (string, error) {
	content, err := tm.fs.ReadFile(templatePath)
	if err != nil {
		return "", fmt.Errorf("failed to read template %s: %w", templatePath, err)
	}

	tmpl, err := template.New("page").Parse(string(content))
	if err != nil {
		return "", fmt.Errorf("failed to parse template %s: %w", templatePath, err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("failed to execute template %s: %w", templatePath, err)
	}

	return buf.String(), nil
}

// CreateModelFiles creates model definition files using templates
func CreateModelFiles(entityDir, entityName, tableName, tableID, createdAt, appID string) error {
	tm := NewTemplateManager()
	data := ModelTemplateData{
		ModelName:       entityName,
		EntityName:      entityName,
		EntityNameLower: strings.ToLower(entityName),
		TableName:       tableName,
		TableID:         tableID,
		CreatedAt:       createdAt,
		AppID:           appID,
	}

	data.EntityNameUpper = strings.ToUpper(entityName)

	files := map[string]string{
		"templates/meta/define.json.tmpl":  entityName + ".define.json",
		"templates/meta/columns.json.tmpl": entityName + ".columns.json",
		"templates/meta/check.json.tmpl":   entityName + ".check.json",
		"templates/meta/fk.json.tmpl":      entityName + ".fk.json",
		"templates/meta/view.sql.tmpl":     entityName + ".default.view.sql",
	}

	for templatePath, destName := range files {
		content, err := tm.RenderModelTemplate(templatePath, data)
		if err != nil {
			return fmt.Errorf("failed to render template %s: %w", templatePath, err)
		}

		destPath := filepath.Join(entityDir, destName)
		if err := os.WriteFile(destPath, []byte(content), 0644); err != nil {
			return fmt.Errorf("failed to create file %s: %w", destName, err)
		}
	}

	return nil
}

// PageTemplateData holds data for page template rendering
type PageTemplateData struct {
	PageName       string
	PageID         string
	AppID          string
	CreatedAt      string
	Description    string
	Type           string
	SourceContent  string
	ReleaseContent string
	PreviewContent string
}

// CreatePageFiles creates page definition files using templates
func CreatePageFiles(pageDir, pageName, pageID, createdAt, appID, description, pageType string) error {
	tm := NewTemplateManager()
	data := PageTemplateData{
		PageName:       pageName,
		PageID:         pageID,
		AppID:          appID,
		CreatedAt:      createdAt,
		Description:    description,
		Type:           pageType,
		SourceContent:  "",
		ReleaseContent: "",
		PreviewContent: "",
	}

	files := map[string]string{
		"templates/page/page.define.json.tmpl":  pageName + ".define.json",
		"templates/page/page.source.json.tmpl":  pageName + ".source.json",
		"templates/page/page.release.json.tmpl": pageName + ".release.json",
		"templates/page/page.preview.json.tmpl": pageName + ".preview.json",
	}

	for templatePath, destName := range files {
		content, err := tm.RenderPageTemplate(templatePath, data)
		if err != nil {
			return fmt.Errorf("failed to render template %s: %w", templatePath, err)
		}

		destPath := filepath.Join(pageDir, destName)
		if err := os.WriteFile(destPath, []byte(content), 0644); err != nil {
			return fmt.Errorf("failed to create file %s: %w", destName, err)
		}
	}

	return nil
}

// CreateAPIFile creates an API file using templates
func CreateAPIFile(filePath, apiName, apiType string) error {
	tm := NewTemplateManager()
	data := APITemplateData{
		APIName: apiName,
		APIPath: strings.ToLower(apiName),
	}

	var templatePath string
	switch strings.ToLower(apiType) {
	case "python", "py":
		templatePath = "templates/api/api.py.tmpl"
	case "go":
		templatePath = "templates/api/api.go.tmpl"
	default:
		templatePath = "templates/api/api.js.tmpl"
	}

	content, err := tm.RenderAPITemplate(templatePath, data)
	if err != nil {
		return fmt.Errorf("failed to render API template: %w", err)
	}

	if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
		return fmt.Errorf("failed to create API file: %w", err)
	}

	return nil
}

// CreateWorkflowFile creates a workflow file using templates
func CreateWorkflowFile(filePath, workflowName, workflowDesc, createdAt, updatedAt string) error {
	tm := NewTemplateManager()
	data := WorkflowTemplateData{
		WorkflowName: workflowName,
		WorkflowDesc: workflowDesc,
		CreatedAt:    createdAt,
		UpdatedAt:    updatedAt,
	}

	content, err := tm.RenderWorkflowTemplate("templates/workflow/workflow.json.tmpl", data)
	if err != nil {
		return fmt.Errorf("failed to render workflow template: %w", err)
	}

	if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
		return fmt.Errorf("failed to create workflow file: %w", err)
	}

	return nil
}

// RenderModelTemplate renders a model template with the given data
func (tm *TemplateManager) RenderModelTemplate(templatePath string, data ModelTemplateData) (string, error) {
	content, err := tm.fs.ReadFile(templatePath)
	if err != nil {
		return "", fmt.Errorf("failed to read template %s: %w", templatePath, err)
	}

	tmpl, err := template.New("model").Parse(string(content))
	if err != nil {
		return "", fmt.Errorf("failed to parse template %s: %w", templatePath, err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("failed to execute template %s: %w", templatePath, err)
	}

	return buf.String(), nil
}

// RenderAPITemplate renders an API template with the given data
func (tm *TemplateManager) RenderAPITemplate(templatePath string, data APITemplateData) (string, error) {
	content, err := tm.fs.ReadFile(templatePath)
	if err != nil {
		return "", fmt.Errorf("failed to read template %s: %w", templatePath, err)
	}

	tmpl, err := template.New("api").Parse(string(content))
	if err != nil {
		return "", fmt.Errorf("failed to parse template %s: %w", templatePath, err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("failed to execute template %s: %w", templatePath, err)
	}

	return buf.String(), nil
}

// RenderWorkflowTemplate renders a workflow template with the given data
func (tm *TemplateManager) RenderWorkflowTemplate(templatePath string, data WorkflowTemplateData) (string, error) {
	content, err := tm.fs.ReadFile(templatePath)
	if err != nil {
		return "", fmt.Errorf("failed to read template %s: %w", templatePath, err)
	}

	tmpl, err := template.New("workflow").Parse(string(content))
	if err != nil {
		return "", fmt.Errorf("failed to parse template %s: %w", templatePath, err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("failed to execute template %s: %w", templatePath, err)
	}

	return buf.String(), nil
}

func createDownloadNotice(destPath string) error {
	docDir := filepath.Dir(destPath)
	if err := os.MkdirAll(docDir, 0755); err != nil {
		return err
	}

	baseName := filepath.Base(destPath)
	githubURL := fmt.Sprintf("https://github.com/geelato/geelato-doc/raw/main/geelato-cli/templates/%s", baseName)

	content := "# " + strings.TrimSuffix(baseName, ".md") + "\n\n" +
		"**无法从 GitHub 下载此文档**\n\n" +
		"## 手动下载指引\n\n" +
		"请从以下地址手动下载文档：\n\n" +
		"**GitHub 原始文件：**\n" + githubURL + "\n\n" +
		"## 操作步骤\n\n" +
		"1. 复制上面的链接地址\n" +
		"2. 在浏览器中打开\n" +
		"3. 右键点击 \"Raw\" 按钮，选择 \"链接另存为\"\n" +
		"4. 保存到当前目录的 doc 文件夹下\n\n" +
		"## 快速命令\n\n" +
		"```bash\n" +
		"# 进入应用目录\n" +
		"cd <your-app-name>\n\n" +
		"# 使用 curl 下载（Linux/Mac）\n" +
		"curl -o \"doc/" + baseName + "\" \"" + githubURL + "\"\n\n" +
		"# 使用 PowerShell 下载（Windows）\n" +
		"Invoke-WebRequest -Uri \"" + githubURL + "\" -OutFile \"doc/" + baseName + "\"\n" +
		"```\n\n" +
		"## GitHub 页面\n\n" +
		"也可以直接在 GitHub 网页上查看文档内容：\n\n" +
		"https://github.com/geelato/geelato-doc/blob/main/geelato-cli/templates/" + baseName + "\n\n" +
		"---\n" +
		"*Generated by Geelato CLI*\n"

	return os.WriteFile(destPath, []byte(content), 0644)
}
