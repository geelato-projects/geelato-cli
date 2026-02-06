package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/geelato/cli/internal/model"
	"github.com/geelato/cli/pkg/logger"
	"github.com/geelato/cli/pkg/prompt"
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

	entityName := strings.Title(modelName)
	tableName := "platform_" + strings.ToLower(modelName)
	now := time.Now().Format(time.RFC3339)
	tableId := "tbl_" + strings.ToLower(entityName)

	entityDir := filepath.Join(cwd, "meta", entityName)
	if err := os.MkdirAll(entityDir, 0755); err != nil {
		return fmt.Errorf("failed to create entity directory: %w", err)
	}

	defineContent := fmt.Sprintf(`{
  "meta": {
    "version": "1.0.0",
    "createdAt": "%s"
  },
  "table": {
    "id": "%s",
    "title": "%s",
    "entityName": "%s",
    "tableName": "%s",
    "tableSchema": "platform",
    "tableType": "entity",
    "tableComment": "%s entity"
  }
}`, now, tableId, modelName, entityName, tableName, entityName)

	definePath := filepath.Join(entityDir, entityName+".define.json")
	if err := os.WriteFile(definePath, []byte(defineContent), 0644); err != nil {
		return fmt.Errorf("failed to create define file: %w", err)
	}

	columnsContent := fmt.Sprintf(`{
  "meta": {
    "version": "1.0.0",
    "tableId": "%s"
  },
  "columns": [
    {
      "id": "col_%s_id",
      "columnName": "id",
      "dataType": "bigint",
      "isPrimaryKey": true,
      "isNullable": false,
      "comment": "Primary key"
    },
    {
      "id": "col_%s_created_at",
      "columnName": "created_at",
      "dataType": "datetime",
      "isNullable": false,
      "defaultValue": "CURRENT_TIMESTAMP",
      "comment": "Creation time"
    },
    {
      "id": "col_%s_updated_at",
      "columnName": "updated_at",
      "dataType": "datetime",
      "isNullable": false,
      "defaultValue": "CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP",
      "comment": "Update time"
    }
  ]
}`, tableId, strings.ToLower(entityName), strings.ToLower(entityName), strings.ToLower(entityName))

	columnsPath := filepath.Join(entityDir, entityName+".columns.json")
	if err := os.WriteFile(columnsPath, []byte(columnsContent), 0644); err != nil {
		return fmt.Errorf("failed to create columns file: %w", err)
	}

	checkContent := fmt.Sprintf(`{
  "meta": {
    "version": "1.0.0",
    "tableId": "%s"
  },
  "checks": []
}`, tableId)

	checkPath := filepath.Join(entityDir, entityName+".check.json")
	if err := os.WriteFile(checkPath, []byte(checkContent), 0644); err != nil {
		return fmt.Errorf("failed to create check file: %w", err)
	}

	fkContent := fmt.Sprintf(`{
  "meta": {
    "version": "1.0.0",
    "tableId": "%s"
  },
  "foreignKeys": []
}`, tableId)

	fkPath := filepath.Join(entityDir, entityName+".fk.json")
	if err := os.WriteFile(fkPath, []byte(fkContent), 0644); err != nil {
		return fmt.Errorf("failed to create fk file: %w", err)
	}

	viewContent := fmt.Sprintf(`-- @meta
SELECT * FROM %s WHERE deleted_at IS NULL`, tableName)

	viewPath := filepath.Join(entityDir, entityName+".default.view.sql")
	if err := os.WriteFile(viewPath, []byte(viewContent), 0644); err != nil {
		return fmt.Errorf("failed to create view file: %w", err)
	}

	return nil
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
		return fmt.Errorf("model '%s' does not exist", entityName)
	}

	parts := strings.Split(fieldSpec, ":")
	if len(parts) < 2 {
		return fmt.Errorf("invalid field spec: %s (format: name:type[:length])", fieldSpec)
	}

	columnName := parts[0]
	dataType := strings.ToLower(parts[1])

	length := 0
	if len(parts) >= 3 {
		fmt.Sscanf(parts[2], "%d", &length)
	}

	comment := ""
	if len(parts) >= 4 {
		comment = parts[3]
	}
	if comment == "" {
		promptStr := fmt.Sprintf("Enter comment for field '%s':", columnName)
		comment, _ = prompt.Input(promptStr)
	}

	mgr := model.NewManager(cwd)

	field := model.FieldDefinition{
		ColumnName: columnName,
		DataType:  dataType,
		Length:    length,
		Comment:   comment,
		Nullable:  true,
	}

	if err := mgr.AddField(entityName, field); err != nil {
		return fmt.Errorf("failed to add field: %w", err)
	}

	logger.Success("Field added to model '%s':", entityName)
	logger.Infof("  Name: %s", columnName)
	logger.Infof("  Type: %s", dataType)
	if length > 0 {
		logger.Infof("  Length: %d", length)
	}
	logger.Infof("  Comment: %s", comment)

	return nil
}

func NewAddViewSubCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "view <entity-name> <view-name>",
		Short: "添加视图",
		Long: `向指定模型添加视图定义

示例:
  geelato model add view User active-users`,
		Args: cobra.MinimumNArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			entityName := args[0]
			viewName := args[1]
			return runAddView(entityName, viewName)
		},
	}

	return cmd
}

func runAddView(entityName, viewName string) error {
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %w", err)
	}

	entityDir := filepath.Join(cwd, "meta", entityName)
	if _, err := os.Stat(entityDir); os.IsNotExist(err) {
		return fmt.Errorf("model '%s' does not exist", entityName)
	}

	mgr := model.NewManager(cwd)

	view := model.ViewDefinition{
		Name:        viewName,
		Description: viewName,
	}

	if err := mgr.AddView(entityName, view); err != nil {
		return fmt.Errorf("failed to add view: %w", err)
	}

	logger.Success("View '%s' added to model '%s'", viewName, entityName)
	return nil
}

func NewAddCheckSubCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "check <entity-name> <expression>",
		Short: "添加检查约束",
		Long: `向指定模型添加检查约束

示例:
  geelato model add check User "status IN (0,1)"`,
		Args: cobra.MinimumNArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			entityName := args[0]
			expression := args[1]
			return runAddCheck(entityName, expression)
		},
	}

	return cmd
}

func runAddCheck(entityName, expression string) error {
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %w", err)
	}

	entityDir := filepath.Join(cwd, "meta", entityName)
	if _, err := os.Stat(entityDir); os.IsNotExist(err) {
		return fmt.Errorf("model '%s' does not exist", entityName)
	}

	mgr := model.NewManager(cwd)

	check := model.CheckDefinition{
		Expression: expression,
	}

	if err := mgr.AddCheck(entityName, check); err != nil {
		return fmt.Errorf("failed to add check: %w", err)
	}

	logger.Success("Check constraint added to model '%s'", entityName)
	logger.Infof("  Expression: %s", expression)
	return nil
}

func NewAddPermissionSubCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "permission <entity-name> <action> <role>",
		Short: "添加权限",
		Long: `向指定模型添加权限配置

示例:
  geelato model add permission User read admin`,
		Args: cobra.MinimumNArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			entityName := args[0]
			action := args[1]
			role := args[2]
			return runAddPermission(entityName, action, role)
		},
	}

	return cmd
}

func runAddPermission(entityName, action, role string) error {
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %w", err)
	}

	entityDir := filepath.Join(cwd, "meta", entityName)
	if _, err := os.Stat(entityDir); os.IsNotExist(err) {
		return fmt.Errorf("model '%s' does not exist", entityName)
	}

	mgr := model.NewManager(cwd)

	perm := model.PermissionDefinition{
		Action: action,
		Role:   role,
	}

	if err := mgr.AddPermission(entityName, perm); err != nil {
		return fmt.Errorf("failed to add permission: %w", err)
	}

	logger.Success("Permission added to model '%s'", entityName)
	logger.Infof("  Action: %s", action)
	logger.Infof("  Role: %s", role)
	return nil
}
