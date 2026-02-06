package cmd

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/geelato/cli/pkg/logger"
	"github.com/geelato/cli/pkg/prompt"
	"github.com/spf13/cobra"
)

const (
	githubRawURL = "https://raw.githubusercontent.com/geelato/geelato-doc/main"
)

func initCmdFn() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init <app-name>",
		Short: "Initialize a new application",
		Long: `Initialize a new Geelato application

Examples:
  geelato init my-app`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var appName string
			if len(args) > 0 {
				appName = args[0]
			} else {
				input, err := prompt.Input("Enter app name (app-name):")
				if err != nil {
					return fmt.Errorf("failed to input app name: %w", err)
				}
				appName = strings.TrimSpace(input)
			}

			if appName == "" {
				return fmt.Errorf("app name cannot be empty")
			}

			appDir, err := initializeApp(appName)
			if err != nil {
				return fmt.Errorf("failed to initialize app: %w", err)
			}

			logger.Infof("App '%s' initialized successfully!", appName)
			logger.Infof("App directory: %s", appDir)
			logger.Info("")
			logger.Info("Next steps:")
			logger.Info("  cd %s", appName)
			logger.Info("  geelato model create User        # Create your first model")
			logger.Info("  geelato api create getUserList  # Create your first API")
			logger.Info("  geelato app push                # Push to cloud")

			return nil
		},
	}

	return cmd
}

func initializeApp(appName string) (string, error) {
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

	now := time.Now().Format(time.RFC3339)

	files := map[string]string{
		"geelato.json": generateGeelatoJSON(appName, now),
		"README.md":    generateREADME(appName),
		".gitignore":   generateGitignore(),
	}

	for path, content := range files {
		filePath := filepath.Join(appDir, path)
		if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
			os.RemoveAll(appDir)
			return "", fmt.Errorf("failed to create file %s: %w", path, err)
		}
	}

	dirs := []string{
		"doc",
		"example/user-management-app/meta/User",
		"example/user-management-app/meta/Department",
		"example/user-management-app/api/user",
		"example/user-management-app/page/user",
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

	exampleFiles := map[string]string{
		"example/user-management-app/.gitignore":                generateExampleGitignore(),
		"example/user-management-app/README.md":                 generateExampleREADME(),
		"example/user-management-app/geelato.json":              generateExampleGeelatoJSON(),
		"example/user-management-app/api/user/getList.api.js":   generateExampleAPIGetList(),
		"example/user-management-app/api/user/getDetail.api.js": generateExampleAPIGetDetail(),
		"example/user-management-app/api/user/saveUser.api.js":  generateExampleAPISaveUser(),
	}

	for path, content := range exampleFiles {
		filePath := filepath.Join(appDir, path)
		if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
			os.RemoveAll(appDir)
			return "", fmt.Errorf("failed to create example file %s: %w", path, err)
		}
	}

	exampleMetaFiles := map[string]string{
		"example/user-management-app/meta/User/User.define.json":                  generateExampleUserDefine(),
		"example/user-management-app/meta/User/User.columns.json":                 generateExampleUserColumns(),
		"example/user-management-app/meta/User/User.check.json":                   generateExampleUserCheck(),
		"example/user-management-app/meta/User/User.fk.json":                      generateExampleUserFK(),
		"example/user-management-app/meta/User/User.default.view.sql":             generateExampleUserView("User"),
		"example/user-management-app/meta/Department/Department.define.json":      generateExampleDeptDefine(),
		"example/user-management-app/meta/Department/Department.columns.json":     generateExampleDeptColumns(),
		"example/user-management-app/meta/Department/Department.check.json":       generateExampleDeptCheck(),
		"example/user-management-app/meta/Department/Department.fk.json":          generateExampleDeptFK(),
		"example/user-management-app/meta/Department/Department.default.view.sql": generateExampleUserView("Department"),
	}

	for path, content := range exampleMetaFiles {
		filePath := filepath.Join(appDir, path)
		if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
			os.RemoveAll(appDir)
			return "", fmt.Errorf("failed to create example meta file %s: %w", path, err)
		}
	}

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

func downloadTemplatesZip(appDir string) error {
	url := fmt.Sprintf("%s/geelato-cli/templates.zip", githubRawURL)

	logger.Info("Downloading templates...")

	client := &http.Client{Timeout: 60 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return fmt.Errorf("failed to download templates: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("HTTP %d", resp.StatusCode)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response: %w", err)
	}

	zipReader, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		return fmt.Errorf("failed to read zip: %w", err)
	}

	for _, file := range zipReader.File {
		if filepath.Ext(file.Name) != ".md" {
			continue
		}

		destPath := filepath.Join(appDir, "doc", filepath.Base(file.Name))
		if err := os.MkdirAll(filepath.Dir(destPath), 0755); err != nil {
			return fmt.Errorf("failed to create directory: %w", err)
		}

		destFile, err := os.Create(destPath)
		if err != nil {
			return fmt.Errorf("failed to create file: %w", err)
		}

		srcFile, err := file.Open()
		if err != nil {
			destFile.Close()
			return fmt.Errorf("failed to open zip entry: %w", err)
		}

		_, err = io.Copy(destFile, srcFile)
		srcFile.Close()
		destFile.Close()

		if err != nil {
			return fmt.Errorf("failed to write file: %w", err)
		}

		logger.Debug("Extracted %s", filepath.Base(destPath))
	}

	return nil
}

func writeFallbackDoc(destPath, filename string) error {
	return nil
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

func copyFile(src, dst string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	_, err = io.Copy(dstFile, srcFile)
	return err
}

func generateGeelatoJSON(appName, now string) string {
	return `{
  "meta": {
    "version": "1.0.0",
    "appId": "app_` + appName + `",
    "name": "` + appName + `",
    "description": "A Geelato application",
    "createdAt": "` + now + `"
  },
  "config": {
    "api": {
      "url": "",
      "timeout": 30
    },
    "sync": {
      "autoPush": false,
      "autoPull": false
    }
  }
}`
}

func generateREADME(appName string) string {
	return `# ` + appName + `

This is a new application created with Geelato CLI.

## Directory Structure

` + appName + `/
|-- doc/                    # Documentation
|   |-- entity-model-spec.md
|   |-- api-script-guide.md
|   |-- workflow-guide.md
|   |-- mcp-guide.md
|-- example/               # Example applications
|   |-- user-management-app/
|-- meta/                  # Model definitions
|-- api/                   # API scripts
|-- page/                  # Page definitions
|-- workflow/              # Workflow definitions

## Quick Start

### 1. Create a Model

    geelato model create User

### 2. Create an API

    geelato api create getUserList

### 3. Push to Cloud

    geelato app push

## Documentation

- [Entity Model Specification](./doc/实体模型定义规范.md)
- [API Script Guide](./doc/API脚本编写指引和规范.md)
- [Workflow Guide](./doc/工作流流程编排指引和规范.md)
- [MCP Guide](./doc/MCP配置指引.md)

## Examples

See [example/user-management-app](./example/user-management-app) for complete application examples.
`
}

func generateGitignore() string {
	return `# Geelato CLI
.geelato/sync-state.json
.geelato/cache/

# Node.js
node_modules/
npm-debug.log*
yarn-debug.log*
yarn-error.log*

# IDE
.idea/
.vscode/
*.swp
*.swo
*~

# OS
.DS_Store
Thumbs.db

# Logs
logs/
*.log

# Build
dist/
build/
`
}

func generateExampleGitignore() string {
	return `# Geelato CLI
.geelato/

# IDE
.idea/
.vscode/
`
}

func generateExampleREADME() string {
	return `# user-management-app

User management example application.

## Models

- User - User entity
- Department - Department entity

## APIs

- getUserList - Get user list
- getDetail - Get user detail
- saveUser - Save user
`
}

func generateExampleGeelatoJSON() string {
	return `{
  "meta": {
    "version": "1.0.0",
    "appId": "app_user_mgmt_001",
    "name": "user-management-app",
    "description": "User management example application"
  }
}
`
}

func generateExampleUserDefine() string {
	return `{
  "meta": {
    "version": "1.0.0",
    "createdAt": "2024-01-15T10:00:00Z"
  },
  "table": {
    "id": "tbl_user_001",
    "title": "User",
    "entityName": "User",
    "tableName": "platform_user",
    "tableSchema": "platform",
    "tableType": "entity",
    "tableComment": "System user table"
  }
}
`
}

func generateExampleUserColumns() string {
	return `{
  "meta": {
    "version": "1.0.0",
    "tableId": "tbl_user_001"
  },
  "columns": [
    {
      "id": "col_user_id",
      "columnName": "id",
      "dataType": "bigint",
      "isPrimaryKey": true,
      "isNullable": false,
      "comment": "Primary key"
    },
    {
      "id": "col_user_name",
      "columnName": "name",
      "dataType": "varchar",
      "length": 100,
      "isNullable": false,
      "comment": "User name"
    },
    {
      "id": "col_user_login_name",
      "columnName": "login_name",
      "dataType": "varchar",
      "length": 50,
      "isNullable": false,
      "comment": "Login name"
    },
    {
      "id": "col_user_email",
      "columnName": "email",
      "dataType": "varchar",
      "length": 200,
      "isNullable": true,
      "comment": "Email address"
    },
    {
      "id": "col_user_phone",
      "columnName": "phone",
      "dataType": "varchar",
      "length": 20,
      "isNullable": true,
      "comment": "Phone number"
    },
    {
      "id": "col_user_status",
      "columnName": "status",
      "dataType": "int",
      "isNullable": false,
      "defaultValue": "1",
      "comment": "Status: 0-disabled, 1-enabled"
    },
    {
      "id": "col_user_dept_id",
      "columnName": "dept_id",
      "dataType": "bigint",
      "isNullable": true,
      "comment": "Department ID"
    },
    {
      "id": "col_user_created_at",
      "columnName": "created_at",
      "dataType": "datetime",
      "isNullable": false,
      "defaultValue": "CURRENT_TIMESTAMP",
      "comment": "Creation time"
    },
    {
      "id": "col_user_updated_at",
      "columnName": "updated_at",
      "dataType": "datetime",
      "isNullable": false,
      "defaultValue": "CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP",
      "comment": "Update time"
    }
  ]
}
`
}

func generateExampleUserCheck() string {
	return `{
  "meta": {
    "version": "1.0.0",
    "tableId": "tbl_user_001"
  },
  "checks": []
}
`
}

func generateExampleUserFK() string {
	return `{
  "meta": {
    "version": "1.0.0",
    "tableId": "tbl_user_001"
  },
  "foreignKeys": [
    {
      "id": "fk_user_dept",
      "tableId": "tbl_user_001",
      "tableName": "platform_user",
      "columnName": "dept_id",
      "foreignTable": "platform_department",
      "foreignColumn": "id",
      "onDelete": "SET_NULL",
      "onUpdate": "CASCADE"
    }
  ]
}
`
}

func generateExampleUserView(entityName string) string {
	return "-- @meta\nSELECT * FROM platform_" + strings.ToLower(entityName) + " WHERE deleted_at IS NULL\n"
}

func generateExampleDeptDefine() string {
	return `{
  "meta": {
    "version": "1.0.0",
    "createdAt": "2024-01-15T10:00:00Z"
  },
  "table": {
    "id": "tbl_dept_001",
    "title": "Department",
    "entityName": "Department",
    "tableName": "platform_department",
    "tableSchema": "platform",
    "tableType": "entity",
    "tableComment": "Department table"
  }
}
`
}

func generateExampleDeptColumns() string {
	return `{
  "meta": {
    "version": "1.0.0",
    "tableId": "tbl_dept_001"
  },
  "columns": [
    {
      "id": "col_dept_id",
      "columnName": "id",
      "dataType": "bigint",
      "isPrimaryKey": true,
      "isNullable": false,
      "comment": "Primary key"
    },
    {
      "id": "col_dept_name",
      "columnName": "name",
      "dataType": "varchar",
      "length": 100,
      "isNullable": false,
      "comment": "Department name"
    },
    {
      "id": "col_dept_code",
      "columnName": "code",
      "dataType": "varchar",
      "length": 50,
      "isNullable": false,
      "comment": "Department code"
    },
    {
      "id": "col_dept_parent_id",
      "columnName": "parent_id",
      "dataType": "bigint",
      "isNullable": true,
      "comment": "Parent department ID"
    },
    {
      "id": "col_dept_created_at",
      "columnName": "created_at",
      "dataType": "datetime",
      "isNullable": false,
      "defaultValue": "CURRENT_TIMESTAMP",
      "comment": "Creation time"
    }
  ]
}
`
}

func generateExampleDeptCheck() string {
	return `{
  "meta": {
    "version": "1.0.0",
    "tableId": "tbl_dept_001"
  },
  "checks": [
    {
      "id": "chk_dept_code",
      "expression": "code REGEXP '^[A-Z0-9_]+$'",
      "description": "Department code must contain only uppercase letters, numbers and underscores"
    }
  ]
}
`
}

func generateExampleDeptFK() string {
	return `{
  "meta": {
    "version": "1.0.0",
    "tableId": "tbl_dept_001"
  },
  "foreignKeys": [
    {
      "id": "fk_dept_parent",
      "tableId": "tbl_dept_001",
      "tableName": "platform_department",
      "columnName": "parent_id",
      "foreignTable": "platform_department",
      "foreignColumn": "id",
      "onDelete": "SET_NULL",
      "onUpdate": "CASCADE"
    }
  ]
}
`
}

func generateExampleAPIGetList() string {
	return `/**
 * @api
 * @name getList
 * @path /api/user/getList
 * @method POST
 * @description Get user list
 * @group user
 * @version 1.0.0
 */

// @param
// name: pageNum
// type: Integer
// required: true
// default: 1
// description: Page number

// @param
// name: pageSize
// type: Integer
// required: true
// default: 10
// description: Page size

// @return
// type: PageResult
// description: Paginated user list

(function() {
    var pageNum = parseInt($params.pageNum || 1);
    var pageSize = parseInt($params.pageSize || 10);

    var countResult = $db.query("SELECT COUNT(*) as total FROM platform_user");
    var total = countResult[0].total;

    var offset = (pageNum - 1) * pageSize;
    var list = $db.query(
        "SELECT id, name, login_name, email, status, created_at FROM platform_user ORDER BY created_at DESC LIMIT ? OFFSET ?",
        [pageSize, offset]
    );

    return {
        code: 200,
        message: "success",
        data: {
            list: list,
            total: total,
            pageNum: pageNum,
            pageSize: pageSize,
            pages: Math.ceil(total / pageSize)
        }
    };
})();
`
}

func generateExampleAPIGetDetail() string {
	return `/**
 * @api
 * @name getDetail
 * @path /api/user/getDetail
 * @method POST
 * @description Get user detail
 * @group user
 * @version 1.0.0
 */

// @param
// name: id
// type: Integer
// required: true
// description: User ID

(function() {
    var id = parseInt($params.id);
    if (!id) {
        return { code: 400, message: "User ID is required", data: null };
    }

    var result = $db.query("SELECT * FROM platform_user WHERE id = ?", [id]);
    if (result.length === 0) {
        return { code: 404, message: "User not found", data: null };
    }

    return { code: 200, message: "success", data: result[0] };
})();
`
}

func generateExampleAPISaveUser() string {
	return `/**
 * @api
 * @name saveUser
 * @path /api/user/saveUser
 * @method POST
 * @description Save user
 * @group user
 * @version 1.0.0
 */

// @param
// name: id
// type: Integer
// required: false
// description: User ID (empty for new user)

// @param
// name: name
// type: String
// required: true
// description: User name

(function() {
    var id = $params.id;
    var name = $params.name;
    var loginName = $params.loginName;

    if (!name || !loginName) {
        return { code: 400, message: "Name and login name are required", data: null };
    }

    if (id) {
        $db.execute("UPDATE platform_user SET name = ?, login_name = ? WHERE id = ?", [name, loginName, parseInt(id)]);
    } else {
        $db.execute("INSERT INTO platform_user (name, login_name) VALUES (?, ?)", [name, loginName]);
    }

    return { code: 200, message: "success", data: { success: true } };
})();
`
}
