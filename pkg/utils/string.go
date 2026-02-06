package utils

import (
	"crypto/rand"
	"encoding/hex"
	"strings"
	"unicode"
)

func TrimSpace(s string) string {
	return strings.TrimSpace(s)
}

func ToLower(s string) string {
	return strings.ToLower(s)
}

func ToUpper(s string) string {
	return strings.ToUpper(s)
}

func ToTitle(s string) string {
	return strings.ToTitle(s)
}

func Contains(s, substr string) bool {
	return strings.Contains(s, substr)
}

func HasPrefix(s, prefix string) bool {
	return strings.HasPrefix(s, prefix)
}

func HasSuffix(s, suffix string) bool {
	return strings.HasSuffix(s, suffix)
}

func Split(s, sep string) []string {
	if s == "" {
		return []string{}
	}
	return strings.Split(s, sep)
}

func Join(strs []string, sep string) string {
	return strings.Join(strs, sep)
}

func Replace(s, old, new string, n int) string {
	return strings.Replace(s, old, new, n)
}

func ReplaceAll(s, old, new string) string {
	return strings.ReplaceAll(s, old, new)
}

func Trim(s, cutset string) string {
	return strings.Trim(s, cutset)
}

func TrimLeft(s, cutset string) string {
	return strings.TrimLeft(s, cutset)
}

func TrimRight(s, cutset string) string {
	return strings.TrimRight(s, cutset)
}

func TrimPrefix(s, prefix string) string {
	return strings.TrimPrefix(s, prefix)
}

func TrimSuffix(s, suffix string) string {
	return strings.TrimSuffix(s, suffix)
}

func Fields(s string) []string {
	return strings.Fields(s)
}

func Capitalize(s string) string {
	if s == "" {
		return s
	}
	runes := []rune(s)
	runes[0] = unicode.ToUpper(runes[0])
	return string(runes)
}

func Uncapitalize(s string) string {
	if s == "" {
		return s
	}
	runes := []rune(s)
	runes[0] = unicode.ToLower(runes[0])
	return string(runes)
}

func SnakeCase(s string) string {
	var result strings.Builder
	for i, r := range s {
		if unicode.IsUpper(r) {
			if i > 0 {
				result.WriteRune('_')
			}
			result.WriteRune(unicode.ToLower(r))
		} else {
			result.WriteRune(r)
		}
	}
	return result.String()
}

func CamelCase(s string) string {
	var result strings.Builder
	upperNext := true
	for _, r := range s {
		if r == '_' || r == '-' || r == ' ' {
			upperNext = true
			continue
		}
		if upperNext {
			result.WriteRune(unicode.ToUpper(r))
			upperNext = false
		} else {
			result.WriteRune(unicode.ToLower(r))
		}
	}
	return result.String()
}

func RandomHex(n int) (string, error) {
	bytes := make([]byte, n)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
