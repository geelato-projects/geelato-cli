package model

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
)

type Parser struct {
	metaDir string
}

func NewParser() *Parser {
	return &Parser{
		metaDir: "meta",
	}
}

func (p *Parser) Parse(entityName string) (*EntityParseResult, error) {
	result := &EntityParseResult{
		EntityName: entityName,
	}

	entityDir := filepath.Join(p.metaDir, entityName)

	defineFile := filepath.Join(entityDir, entityName+".define.json")
	if data, err := os.ReadFile(defineFile); err == nil {
		if tableDef, err := ParseTableDefinition(data); err == nil {
			result.Table = tableDef
		}
	}

	columnsFile := filepath.Join(entityDir, entityName+".columns.json")
	if data, err := os.ReadFile(columnsFile); err == nil {
		if cols, err := ParseColumns(data); err == nil {
			result.Columns = cols
		}
	}

	checkFile := filepath.Join(entityDir, entityName+".check.json")
	if data, err := os.ReadFile(checkFile); err == nil {
		if checks, err := ParseCheckConstraints(data); err == nil {
			result.Checks = checks
		}
	}

	fkFile := filepath.Join(entityDir, entityName+".fk.json")
	if data, err := os.ReadFile(fkFile); err == nil {
		if fks, err := ParseForeignKeys(data); err == nil {
			result.ForeignKeys = fks
		}
	}

	return result, nil
}

func ParseTableDefinition(data []byte) (*TableDefinition, error) {
	var result struct {
		Meta struct {
			Version  string `json:"version"`
			CreatedAt string `json:"createdAt"`
		} `json:"meta"`
		Table TableDefinition `json:"table"`
	}

	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}

	return &result.Table, nil
}

func ParseColumns(data []byte) ([]ColumnDefinition, error) {
	var result struct {
		Meta struct {
			Version  string `json:"version"`
			TableID  string `json:"tableId"`
		} `json:"meta"`
		Columns []ColumnDefinition `json:"columns"`
	}

	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}

	return result.Columns, nil
}

func ParseCheckConstraints(data []byte) ([]CheckConstraint, error) {
	var result struct {
		Meta   struct {
			Version string `json:"version"`
			TableID string `json:"tableId"`
		} `json:"meta"`
		Checks []CheckConstraint `json:"checks"`
	}

	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}

	return result.Checks, nil
}

func ParseForeignKeys(data []byte) ([]ForeignKey, error) {
	var result struct {
		Meta struct {
			Version string `json:"version"`
			TableID string `json:"tableId"`
		} `json:"meta"`
		ForeignKeys []ForeignKey `json:"foreignKeys"`
	}

	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}

	return result.ForeignKeys, nil
}

type EntityParseResult struct {
	EntityName   string
	Table        *TableDefinition
	Columns      []ColumnDefinition
	Checks       []CheckConstraint
	ForeignKeys  []ForeignKey
}

func (p *Parser) ParseAll() ([]EntityParseResult, error) {
	var results []EntityParseResult

	entries, err := os.ReadDir(p.metaDir)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		result, err := p.Parse(entry.Name())
		if err != nil {
			continue
		}

		results = append(results, *result)
	}

	return results, nil
}

func ValidateFieldType(fieldType string) bool {
	validTypes := []string{
		"bigint", "int", "integer", "smallint", "tinyint",
		"varchar", "char", "text", "nvarchar", "ntext",
		"decimal", "numeric", "float", "double", "real",
		"datetime", "date", "time", "timestamp",
		"boolean", "bool",
		"json", "jsonb",
		"binary", "varbinary", "blob",
	}

	fieldType = strings.ToLower(fieldType)
	for _, vt := range validTypes {
		if fieldType == vt {
			return true
		}
	}

	return false
}
