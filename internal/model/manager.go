package model

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/geelato/cli/pkg/crypto"
)

type Manager struct {
	cwd string
}

func NewManager(cwd string) *Manager {
	return &Manager{cwd: cwd}
}

func (m *Manager) LoadModel(entityName string) (*Model, error) {
	entityDir := filepath.Join(m.cwd, "meta", entityName)
	definePath := filepath.Join(entityDir, entityName+".define.json")
	columnsPath := filepath.Join(entityDir, entityName+".columns.json")

	var model Model
	model.Name = entityName

	defineData, err := os.ReadFile(definePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read define file: %w", err)
	}

	var define struct {
		Table struct {
			ID         string `json:"id"`
			Title      string `json:"title"`
			EntityName string `json:"entityName"`
			TableName  string `json:"tableName"`
			TableType  string `json:"tableType"`
			Comment    string `json:"comment"`
		} `json:"table"`
	}

	if err := json.Unmarshal(defineData, &define); err != nil {
		return nil, fmt.Errorf("failed to parse define file: %w", err)
	}

	model.Table = define.Table.ID
	model.Description = define.Table.Comment

	columnsData, err := os.ReadFile(columnsPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read columns file: %w", err)
	}

	var columnsDef struct {
		Columns []ColumnDefinition `json:"columns"`
	}

	if err := json.Unmarshal(columnsData, &columnsDef); err != nil {
		return nil, fmt.Errorf("failed to parse columns file: %w", err)
	}

	var fields []Field
	for _, col := range columnsDef.Columns {
		fields = append(fields, Field{
			ID:         col.ID,
			Name:       col.ColumnName,
			Type:       col.DataType,
			Length:     col.Length,
			Nullable:   col.IsNullable,
			PrimaryKey: col.IsPrimaryKey,
			Default:    col.DefaultValue,
			Comment:    col.Comment,
		})
	}
	model.Fields = fields
	return &model, nil
}

func (m *Manager) AddField(entityName string, field FieldDefinition) error {
	entityDir := filepath.Join(m.cwd, "meta", entityName)
	columnsPath := filepath.Join(entityDir, entityName+".columns.json")

	columnsData, err := os.ReadFile(columnsPath)
	if err != nil {
		return fmt.Errorf("failed to read columns file: %w", err)
	}

	var columnsDef struct {
		Meta    MetaInfo     `json:"meta"`
		Columns []ColumnDefinition `json:"columns"`
	}

	if err := json.Unmarshal(columnsData, &columnsDef); err != nil {
		return fmt.Errorf("failed to parse columns file: %w", err)
	}

	if field.ID == "" {
		field.ID = fmt.Sprintf("col_%s_%s", strings.ToLower(entityName), field.ColumnName)
	}

	columnsDef.Columns = append(columnsDef.Columns, ColumnDefinition{
		ID:           field.ID,
		ColumnName:   field.ColumnName,
		DataType:     field.DataType,
		Length:       field.Length,
		Precision:    field.Precision,
		Scale:        field.Scale,
		IsPrimaryKey: field.IsPrimaryKey,
		IsNullable:   field.Nullable,
		DefaultValue: field.DefaultValue,
		Comment:      field.Comment,
	})

	output, err := json.MarshalIndent(columnsDef, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal columns: %w", err)
	}

	if err := os.WriteFile(columnsPath, output, 0644); err != nil {
		return fmt.Errorf("failed to write columns file: %w", err)
	}

	return nil
}

func (m *Manager) AddView(entityName string, view ViewDefinition) error {
	entityDir := filepath.Join(m.cwd, "meta", entityName)
	viewPath := filepath.Join(entityDir, view.Name+".view.sql")

	if view.Content == "" {
		view.Content = fmt.Sprintf(`-- @view %s
-- Description: %s
SELECT * FROM %s WHERE deleted_at IS NULL`,
			view.Name, view.Description, "platform_"+strings.ToLower(entityName))
	}

	if err := os.WriteFile(viewPath, []byte(view.Content), 0644); err != nil {
		return fmt.Errorf("failed to create view file: %w", err)
	}

	return nil
}

func (m *Manager) AddCheck(entityName string, check CheckDefinition) error {
	entityDir := filepath.Join(m.cwd, "meta", entityName)
	checkPath := filepath.Join(entityDir, entityName+".check.json")

	checkData, err := os.ReadFile(checkPath)
	if err != nil {
		return fmt.Errorf("failed to read check file: %w", err)
	}

	var checkDef struct {
		Meta   MetaInfo       `json:"meta"`
		Checks []CheckConstraint `json:"checks"`
	}

	if err := json.Unmarshal(checkData, &checkDef); err != nil {
		return fmt.Errorf("failed to parse check file: %w", err)
	}

	if check.ID == "" {
		check.ID = fmt.Sprintf("chk_%s_%d", strings.ToLower(entityName), len(checkDef.Checks)+1)
	}

	checkDef.Checks = append(checkDef.Checks, CheckConstraint{
		ID:         check.ID,
		Expression: check.Expression,
		Comment:    check.Comment,
	})

	output, err := json.MarshalIndent(checkDef, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal checks: %w", err)
	}

	if err := os.WriteFile(checkPath, output, 0644); err != nil {
		return fmt.Errorf("failed to write check file: %w", err)
	}

	return nil
}

func (m *Manager) AddForeignKey(entityName string, fk ForeignKeyDefinition) error {
	entityDir := filepath.Join(m.cwd, "meta", entityName)
	fkPath := filepath.Join(entityDir, entityName+".fk.json")

	fkData, err := os.ReadFile(fkPath)
	if err != nil {
		return fmt.Errorf("failed to read fk file: %w", err)
	}

	var fkDef struct {
		Meta       MetaInfo     `json:"meta"`
		ForeignKeys []ForeignKey `json:"foreignKeys"`
	}

	if err := json.Unmarshal(fkData, &fkDef); err != nil {
		return fmt.Errorf("failed to parse fk file: %w", err)
	}

	if fk.ID == "" {
		fk.ID = fmt.Sprintf("fk_%s_%s", strings.ToLower(entityName), fk.ColumnName)
	}

	fkDef.ForeignKeys = append(fkDef.ForeignKeys, ForeignKey{
		ID:            fk.ID,
		TableName:     "platform_" + strings.ToLower(entityName),
		ColumnName:    fk.ColumnName,
		ForeignTable:  fk.ForeignTable,
		ForeignColumn: fk.ForeignColumn,
		OnDelete:      fk.OnDelete,
		OnUpdate:      fk.OnUpdate,
	})

	output, err := json.MarshalIndent(fkDef, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal foreign keys: %w", err)
	}

	if err := os.WriteFile(fkPath, output, 0644); err != nil {
		return fmt.Errorf("failed to write fk file: %w", err)
	}

	return nil
}

func (m *Manager) AddPermission(entityName string, perm PermissionDefinition) error {
	entityDir := filepath.Join(m.cwd, "meta", entityName)
	permPath := filepath.Join(entityDir, entityName+".perm.json")

	permContent := fmt.Sprintf(`{
  "meta": {
    "version": "1.0.0",
    "entityName": "%s",
    "createdAt": "%s"
  },
  "permissions": [
    {
      "action": "%s",
      "role": "%s",
      "condition": "%s"
    }
  ]
}`, entityName, time.Now().Format(time.RFC3339), perm.Action, perm.Role, perm.Condition)

	if err := os.WriteFile(permPath, []byte(permContent), 0644); err != nil {
		return fmt.Errorf("failed to create permission file: %w", err)
	}

	return nil
}

func (m *Manager) ListFields(entityName string) ([]ColumnDefinition, error) {
	entityDir := filepath.Join(m.cwd, "meta", entityName)
	columnsPath := filepath.Join(entityDir, entityName+".columns.json")

	columnsData, err := os.ReadFile(columnsPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read columns file: %w", err)
	}

	var columnsDef struct {
		Columns []ColumnDefinition `json:"columns"`
	}

	if err := json.Unmarshal(columnsData, &columnsDef); err != nil {
		return nil, fmt.Errorf("failed to parse columns file: %w", err)
	}

	return columnsDef.Columns, nil
}

func (m *Manager) ListViews(entityName string) ([]string, error) {
	entityDir := filepath.Join(m.cwd, "meta", entityName)
	var views []string

	entries, err := os.ReadDir(entityDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read entity directory: %w", err)
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := entry.Name()
		if strings.HasSuffix(name, ".view.sql") {
			viewName := strings.TrimSuffix(name, ".view.sql")
			views = append(views, viewName)
		}
	}

	return views, nil
}

type ViewDefinition struct {
	Name        string
	Description string
	Content     string
}

type CheckDefinition struct {
	ID         string
	Expression string
	Comment    string
}

type ForeignKeyDefinition struct {
	ID            string
	ColumnName    string
	ForeignTable  string
	ForeignColumn string
	OnDelete      string
	OnUpdate      string
}

type PermissionDefinition struct {
	Action    string
	Role      string
	Condition string
}

type FieldDefinition struct {
	ID           string
	ColumnName   string
	DataType     string
	Length       int
	Precision    int
	Scale        int
	IsPrimaryKey bool
	Nullable     bool
	DefaultValue string
	Comment      string
}

type MetaInfo struct {
	Version string `json:"version"`
	TableID string `json:"tableId,omitempty"`
}

func GenerateColumnID(entityName, columnName string) string {
	return fmt.Sprintf("col_%s_%s_%s", strings.ToLower(entityName), columnName, crypto.MD5String([]byte(time.Now().Format("20060102150405")))[:8])
}
