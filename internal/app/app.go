package app

import (
	"time"
)

type Application struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Path        string    `json:"path"`
	Template    string    `json:"template"`
	Author      string    `json:"author"`
	Version     string    `json:"version"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

type ApplicationInfo struct {
	Name        string    `json:"name"`
	Path        string    `json:"path"`
	Description string    `json:"description"`
	Template    string    `json:"template"`
	Version     string    `json:"version"`
	LastModTime time.Time `json:"lastModTime"`
}
