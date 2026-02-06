package app

import (
	"os"
	"path/filepath"
	"strings"
)

type ValidationResult struct {
	Valid     bool
	Models    int
	APIs      int
	Workflows int
	Errors    []string
}

type Validator struct {
	cwd    string
	errors []string
}

func NewValidator() *Validator {
	return &Validator{
		errors: []string{},
	}
}

func (v *Validator) Validate(cwd string) (*ValidationResult, error) {
	result := &ValidationResult{
		Valid:  true,
		Errors: []string{},
	}
	v.cwd = cwd
	v.errors = []string{}

	result.Models = v.validateDir("meta", []string{".json"})
	result.APIs = v.validateDir("api", []string{".js"})
	result.Workflows = v.validateDir("workflow", []string{".xml", ".bpmn"})

	result.Errors = v.errors
	if len(result.Errors) > 0 {
		result.Valid = false
	}

	return result, nil
}

func (v *Validator) validateDir(name string, extensions []string) int {
	count := 0
	dirPath := filepath.Join(v.cwd, name)

	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		v.errors = append(v.errors, "Required directory missing: "+name+"/")
		return 0
	}

	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}

		ext := filepath.Ext(path)
		for _, e := range extensions {
			if ext == e {
				count++
				break
			}
		}
		return nil
	})

	if err != nil {
		v.errors = append(v.errors, "Error scanning "+name+"/: "+err.Error())
		return 0
	}

	return count
}

func isIgnoredDir(dirName string) bool {
	ignored := []string{"node_modules", ".git", "vendor", "__pycache__", ".idea", ".vscode"}
	for _, ig := range ignored {
		if dirName == ig {
			return true
		}
	}
	return false
}

func isIgnoredFile(fileName string) bool {
	ignoredPrefixes := []string{".", "_"}
	for _, prefix := range ignoredPrefixes {
		if strings.HasPrefix(fileName, prefix) {
			return true
		}
	}
	return false
}
