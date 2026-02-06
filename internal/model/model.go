package model

import (
	"time"
)

type Model struct {
	Name        string    `json:"name"`
	Table       string    `json:"table"`
	Description string    `json:"description"`
	Fields      []Field   `json:"fields"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

type Field struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Type       string `json:"type"`
	Length     int    `json:"length,omitempty"`
	Required   bool   `json:"required,omitempty"`
	Nullable   bool   `json:"nullable,omitempty"`
	PrimaryKey bool   `json:"primaryKey,omitempty"`
	Unique     bool   `json:"unique,omitempty"`
	Default    string `json:"default,omitempty"`
	Comment    string `json:"comment,omitempty"`
	Reference  string `json:"reference,omitempty"`
}

type TableDefinition struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	EntityName  string    `json:"entityName"`
	TableName   string    `json:"tableName"`
	TableSchema string    `json:"tableSchema"`
	TableType   string    `json:"tableType"`
	Comment     string    `json:"comment"`
}

type ColumnDefinition struct {
	ID           string `json:"id"`
	ColumnName   string `json:"columnName"`
	DataType     string `json:"dataType"`
	Length       int    `json:"length,omitempty"`
	Precision    int    `json:"precision,omitempty"`
	Scale       int    `json:"scale,omitempty"`
	IsPrimaryKey bool   `json:"isPrimaryKey"`
	IsNullable  bool   `json:"isNullable"`
	DefaultValue string `json:"defaultValue,omitempty"`
	Comment     string `json:"comment,omitempty"`
}

type CheckConstraint struct {
	ID         string `json:"id"`
	Expression string `json:"expression"`
	Comment    string `json:"comment,omitempty"`
}

type ForeignKey struct {
	ID           string `json:"id"`
	TableName    string `json:"tableName"`
	ColumnName   string `json:"columnName"`
	ForeignTable string `json:"foreignTable"`
	ForeignColumn string `json:"foreignColumn"`
	OnDelete     string `json:"onDelete,omitempty"`
	OnUpdate     string `json:"onUpdate,omitempty"`
}
