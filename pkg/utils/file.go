package utils

import (
	"os"
	"path/filepath"
	"strings"
)

func Exists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

func IsFile(path string) bool {
	info, err := os.Stat(path)
	return err == nil && !info.IsDir()
}

func IsDir(path string) bool {
	info, err := os.Stat(path)
	return err == nil && info.IsDir()
}

func ReadFile(path string) ([]byte, error) {
	return os.ReadFile(path)
}

func WriteFile(path string, data []byte, perm os.FileMode) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	return os.WriteFile(path, data, perm)
}

func AppendFile(path string, data []byte, perm os.FileMode) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, perm)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.Write(data)
	return err
}

func CopyFile(src, dst string) error {
	data, err := ReadFile(src)
	if err != nil {
		return err
	}
	return WriteFile(dst, data, 0644)
}

func RemoveFile(path string) error {
	return os.Remove(path)
}

func RemoveDir(path string) error {
	return os.RemoveAll(path)
}

func ListDir(path string) ([]string, error) {
	entries, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	var names []string
	for _, entry := range entries {
		names = append(names, entry.Name())
	}
	return names, nil
}

func ListFiles(path string, extensions ...string) ([]string, error) {
	var files []string
	err := filepath.Walk(path, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		if len(extensions) > 0 {
			for _, ext := range extensions {
				if strings.HasSuffix(filePath, ext) {
					files = append(files, filePath)
					break
				}
			}
		} else {
			files = append(files, filePath)
		}

		return nil
	})

	return files, err
}

func ListDirs(path string) ([]string, error) {
	var dirs []string
	err := filepath.Walk(path, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() && filePath != path {
			dirs = append(dirs, filePath)
		}

		return nil
	})

	return dirs, err
}

func Getwd() (string, error) {
	return os.Getwd()
}

func Chdir(path string) error {
	return os.Chdir(path)
}

func HomeDir() (string, error) {
	return os.UserHomeDir()
}

func TempDir() string {
	return os.TempDir()
}
