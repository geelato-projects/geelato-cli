package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/geelato/cli/cmd/initializer"
	"github.com/geelato/cli/pkg/logger"
	"github.com/geelato/cli/pkg/prompt"
	"github.com/spf13/cobra"
)

func NewModelCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "model",
		Short: "模型管理",
		Long:  `管理数据模型，包括创建模型、添加字段、视图、检查约束、权限等操作`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}

	cmd.AddCommand(modelCreateCmd)
	cmd.AddCommand(NewModelListCmd())
	cmd.AddCommand(NewModelAddCmd())

	return cmd
}

var modelCreateCmd = &cobra.Command{
	Use:   "create <model-name>",
	Short: "创建模型",
	Long: `创建一个新的数据模型

示例:
  geelato model create User
  geelato model create Product`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		var modelName string
		if len(args) > 0 {
			modelName = args[0]
		} else {
			input, err := prompt.Input("Enter model name (e.g., User):")
			if err != nil {
				return fmt.Errorf("failed to input model name: %w", err)
			}
			modelName = strings.TrimSpace(input)
		}

		if modelName == "" {
			return fmt.Errorf("model name cannot be empty")
		}

		if err := createModel(modelName); err != nil {
			return fmt.Errorf("failed to create model: %w", err)
		}

		logger.Infof("Model '%s' created successfully!", modelName)
		logger.Info("")
		logger.Info("Created files:")
		logger.Info("  meta/%s/%s.define.json", modelName, modelName)
		logger.Info("  meta/%s/%s.columns.json", modelName, modelName)
		logger.Info("  meta/%s/%s.check.json", modelName, modelName)
		logger.Info("  meta/%s/%s.fk.json", modelName, modelName)
		logger.Info("  meta/%s/%s.default.view.sql", modelName, modelName)

		return nil
	},
}

func createModel(modelName string) error {
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %w", err)
	}

	geelatoPath := filepath.Join(cwd, "geelato.json")
	if _, err := os.Stat(geelatoPath); os.IsNotExist(err) {
		return fmt.Errorf("current directory is not a valid Geelato application")
	}

	// 读取 geelato.json 获取 appId
	appId, err := getAppIdFromGeelatoJSON(geelatoPath)
	if err != nil {
		return fmt.Errorf("failed to read appId from geelato.json: %w", err)
	}

	entityName := strings.Title(modelName)
	tableName := "platform_" + strings.ToLower(modelName)
	now := time.Now().Format(time.RFC3339)
	tableID := "tbl_" + strings.ToLower(entityName)

	entityDir := filepath.Join(cwd, "meta", entityName)
	if err := os.MkdirAll(entityDir, 0755); err != nil {
		return fmt.Errorf("failed to create entity directory: %w", err)
	}

	return initializer.CreateModelFiles(entityDir, entityName, tableName, tableID, now, appId)
}

// getAppIdFromGeelatoJSON 从 geelato.json 文件中读取 appId
func getAppIdFromGeelatoJSON(filePath string) (string, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to read geelato.json: %w", err)
	}

	var config struct {
		Meta struct {
			AppID string `json:"appId"`
		} `json:"meta"`
	}

	if err := json.Unmarshal(data, &config); err != nil {
		return "", fmt.Errorf("failed to parse geelato.json: %w", err)
	}

	return config.Meta.AppID, nil
}

func NewModelListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "列出所有模型",
		Long:  `列出当前应用中的所有模型`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runModelList()
		},
	}

	cmd.Flags().Bool("json", false, "JSON格式输出")

	return cmd
}

func runModelList() error {
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %w", err)
	}

	metaDir := filepath.Join(cwd, "meta")
	entries, err := os.ReadDir(metaDir)
	if err != nil {
		return fmt.Errorf("failed to read meta directory: %w", err)
	}

	logger.Info("Models in current application:")
	logger.Info("================================")

	var models []string
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		entityName := entry.Name()
		definePath := filepath.Join(metaDir, entityName, entityName+".define.json")
		if _, err := os.Stat(definePath); err == nil {
			models = append(models, entityName)
			logger.Info("  - %s", entityName)
		}
	}

	logger.Info("")
	logger.Infof("Total: %d models", len(models))

	return nil
}

func NewModelAddCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add",
		Short: "向模型添加组件",
		Long: `向已存在的模型添加字段、视图、检查约束或权限

可添加的组件:
  field       - 添加字段到模型
  view        - 添加视图定义
  check       - 添加检查约束
  permission  - 添加权限配置`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}

	cmd.AddCommand(NewAddFieldSubCmd())
	cmd.AddCommand(NewAddViewSubCmd())
	cmd.AddCommand(NewAddCheckSubCmd())
	cmd.AddCommand(NewAddPermissionSubCmd())

	return cmd
}

func NewAddFieldSubCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "field <entity-name> <name:type[:length]> [comment]",
		Short: "添加字段",
		Long: `向指定模型添加新字段

字段类型: string, int, bigint, decimal, datetime, boolean, text

示例:
  geelato model add field User name:string:50
  geelato model add field User age:int
  geelato model add field User email:string:100:邮箱地址`,
		Args: cobra.MinimumNArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			entityName := args[0]
			fieldSpec := args[1]
			return runAddField(entityName, fieldSpec)
		},
	}

	cmd.Flags().Bool("pk", false, "设为主键")
	cmd.Flags().Bool("required", false, "设为必填")
	cmd.Flags().String("default", "", "默认值")

	return cmd
}

func runAddField(entityName, fieldSpec string) error {
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %w", err)
	}

	entityDir := filepath.Join(cwd, "meta", entityName)
	if _, err := os.Stat(entityDir); os.IsNotExist(err) {
		return fmt.Errorf("entity '%s' does not exist", entityName)
	}

	columnsPath := filepath.Join(entityDir, entityName+".columns.json")
	if _, err := os.Stat(columnsPath); os.IsNotExist(err) {
		return fmt.Errorf("columns file does not exist for entity '%s'", entityName)
	}

	geelatoPath := filepath.Join(cwd, "geelato.json")
	appId, err := getAppIdFromGeelatoJSON(geelatoPath)
	if err != nil {
		return fmt.Errorf("failed to read appId from geelato.json: %w", err)
	}

	// Parse field spec: name:type[:length]
	parts := strings.Split(fieldSpec, ":")
	if len(parts) < 2 {
		return fmt.Errorf("invalid field spec format, expected name:type[:length]")
	}

	fieldName := parts[0]
	dataType := parts[1]
	length := 0
	comment := ""
	if len(parts) >= 3 {
		fmt.Sscanf(parts[2], "%d", &length)
	}
	if len(parts) >= 4 {
		comment = parts[3]
	}

	// Determine column type based on data type
	columnType := getColumnType(dataType, length)
	dateTimePrecision := ""
	if dataType == "datetime" {
		dateTimePrecision = "0"
	}

	// Prepare template data
	tableName := "platform_" + strings.ToLower(entityName)
	tableID := "tbl_" + strings.ToLower(entityName)
	ordinalPosition := 1

	// Read existing columns to get ordinal position
	columnsData, err := os.ReadFile(columnsPath)
	if err != nil {
		return fmt.Errorf("failed to read columns file: %w", err)
	}

	var columnsDef struct {
		Meta    map[string]interface{}   `json:"meta"`
		Columns []map[string]interface{} `json:"columns"`
	}

	if err := json.Unmarshal(columnsData, &columnsDef); err != nil {
		return fmt.Errorf("failed to parse columns file: %w", err)
	}
	ordinalPosition = len(columnsDef.Columns) + 1

	// Render column using template
	tm := initializer.NewTemplateManager()
	columnData := initializer.ColumnTemplateData{
		EntityName:        entityName,
		EntityNameLower:   strings.ToLower(entityName),
		FieldName:         fieldName,
		FieldNameLower:    strings.ToLower(fieldName),
		AppID:             appId,
		TableID:           tableID,
		TableName:         tableName,
		ColumnName:        stringsToSnakeCase(fieldName),
		ColumnType:        columnType,
		DataType:          dataType,
		Length:            length,
		DateTimePrecision: dateTimePrecision,
		OrdinalPosition:   ordinalPosition,
		Comment:           comment,
	}

	columnContent, err := tm.RenderColumnTemplate("templates/meta/simple/column.json.tmpl", columnData)
	if err != nil {
		return fmt.Errorf("failed to render column template: %w", err)
	}

	// Parse rendered column JSON
	var newColumn map[string]interface{}
	if err := json.Unmarshal([]byte(columnContent), &newColumn); err != nil {
		return fmt.Errorf("failed to parse rendered column: %w", err)
	}

	columnsDef.Columns = append(columnsDef.Columns, newColumn)

	// Write back
	output, err := json.MarshalIndent(columnsDef, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal columns: %w", err)
	}

	if err := os.WriteFile(columnsPath, output, 0644); err != nil {
		return fmt.Errorf("failed to write columns file: %w", err)
	}

	logger.Infof("Field '%s' added to entity '%s' successfully!", fieldName, entityName)
	return nil
}

// getColumnType returns the column type string based on data type and length
func getColumnType(dataType string, length int) string {
	switch dataType {
	case "string":
		if length > 0 {
			return fmt.Sprintf("varchar(%d)", length)
		}
		return "varchar(255)"
	case "int":
		return "int"
	case "bigint":
		return "bigint"
	case "decimal":
		if length > 0 {
			return fmt.Sprintf("decimal(%d,2)", length)
		}
		return "decimal(10,2)"
	case "datetime":
		return "datetime"
	case "boolean":
		return "tinyint(1)"
	case "text":
		return "text"
	default:
		return "varchar(255)"
	}
}

// stringsToSnakeCase converts CamelCase to snake_case
func stringsToSnakeCase(s string) string {
	var result []byte
	for i := 0; i < len(s); i++ {
		if s[i] >= 'A' && s[i] <= 'Z' {
			if i > 0 {
				result = append(result, '_')
			}
			result = append(result, s[i]+32)
		} else {
			result = append(result, s[i])
		}
	}
	return string(result)
}

func NewAddViewSubCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "view <entity-name> <view-name>",
		Short: "添加视图",
		Long:  `向指定模型添加新视图`,
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			entityName := args[0]
			viewName := args[1]
			return runAddView(entityName, viewName)
		},
	}
}

func runAddView(entityName, viewName string) error {
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %w", err)
	}

	entityDir := filepath.Join(cwd, "meta", entityName)
	if _, err := os.Stat(entityDir); os.IsNotExist(err) {
		return fmt.Errorf("entity '%s' does not exist", entityName)
	}

	viewPath := filepath.Join(entityDir, entityName+"."+viewName+".view.sql")
	if _, err := os.Stat(viewPath); err == nil {
		return fmt.Errorf("view '%s' already exists", viewName)
	}

	content := fmt.Sprintf("-- @meta\nSELECT * FROM platform_%s WHERE deleted_at IS NULL", strings.ToLower(entityName))
	if err := os.WriteFile(viewPath, []byte(content), 0644); err != nil {
		return fmt.Errorf("failed to create view file: %w", err)
	}

	logger.Infof("View '%s' added to entity '%s' successfully!", viewName, entityName)
	return nil
}

func NewAddCheckSubCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "check <entity-name> <expression> [description]",
		Short: "添加检查约束",
		Long:  `向指定模型添加检查约束`,
		Args:  cobra.RangeArgs(2, 3),
		RunE: func(cmd *cobra.Command, args []string) error {
			entityName := args[0]
			expression := args[1]
			description := ""
			if len(args) > 2 {
				description = args[2]
			}
			return runAddCheck(entityName, expression, description)
		},
	}
}

func runAddCheck(entityName, expression, description string) error {
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %w", err)
	}

	entityDir := filepath.Join(cwd, "meta", entityName)
	if _, err := os.Stat(entityDir); os.IsNotExist(err) {
		return fmt.Errorf("entity '%s' does not exist", entityName)
	}

	checkPath := filepath.Join(entityDir, entityName+".check.json")
	data, err := os.ReadFile(checkPath)
	if err != nil {
		return fmt.Errorf("failed to read check file: %w", err)
	}

	var checkDef struct {
		Meta   map[string]interface{}   `json:"meta"`
		Checks []map[string]interface{} `json:"checks"`
	}

	if err := json.Unmarshal(data, &checkDef); err != nil {
		return fmt.Errorf("failed to parse check file: %w", err)
	}

	newCheck := map[string]interface{}{
		"id":         fmt.Sprintf("chk_%s_%d", strings.ToLower(entityName), len(checkDef.Checks)+1),
		"expression": expression,
		"comment":    description,
	}

	checkDef.Checks = append(checkDef.Checks, newCheck)

	output, err := json.MarshalIndent(checkDef, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal checks: %w", err)
	}

	if err := os.WriteFile(checkPath, output, 0644); err != nil {
		return fmt.Errorf("failed to write check file: %w", err)
	}

	logger.Infof("Check constraint added to entity '%s' successfully!", entityName)
	return nil
}

func NewAddPermissionSubCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "permission <entity-name> <operation> <role>",
		Short: "添加权限",
		Long:  `向指定模型添加权限配置`,
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			entityName := args[0]
			operation := args[1]
			role := args[2]
			return runAddPermission(entityName, operation, role)
		},
	}
}

func runAddPermission(entityName, operation, role string) error {
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %w", err)
	}

	entityDir := filepath.Join(cwd, "meta", entityName)
	if _, err := os.Stat(entityDir); os.IsNotExist(err) {
		return fmt.Errorf("entity '%s' does not exist", entityName)
	}

	permissionPath := filepath.Join(entityDir, entityName+".permission.json")

	var permissionDef struct {
		Meta        map[string]interface{}   `json:"meta"`
		Permissions []map[string]interface{} `json:"permissions"`
	}

	if _, err := os.Stat(permissionPath); err == nil {
		data, err := os.ReadFile(permissionPath)
		if err != nil {
			return fmt.Errorf("failed to read permission file: %w", err)
		}
		if err := json.Unmarshal(data, &permissionDef); err != nil {
			return fmt.Errorf("failed to parse permission file: %w", err)
		}
	} else {
		permissionDef = struct {
			Meta        map[string]interface{}   `json:"meta"`
			Permissions []map[string]interface{} `json:"permissions"`
		}{
			Meta: map[string]interface{}{
				"version":   "1.0.0",
				"tableName": "platform_" + strings.ToLower(entityName),
			},
			Permissions: []map[string]interface{}{},
		}
	}

	newPermission := map[string]interface{}{
		"id":        fmt.Sprintf("perm_%s_%s_%s", strings.ToLower(entityName), operation, strings.ToLower(role)),
		"operation": operation,
		"role":      role,
		"allow":     true,
	}

	permissionDef.Permissions = append(permissionDef.Permissions, newPermission)

	output, err := json.MarshalIndent(permissionDef, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal permissions: %w", err)
	}

	if err := os.WriteFile(permissionPath, output, 0644); err != nil {
		return fmt.Errorf("failed to write permission file: %w", err)
	}

	logger.Infof("Permission added to entity '%s' successfully!", entityName)
	return nil
}
