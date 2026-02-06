package crypto

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"hash"
	"io"
	"os"
)

type HashType string

const (
	MD5    HashType = "md5"
	SHA1   HashType = "sha1"
	SHA256 HashType = "sha256"
)

func HashString(data []byte, hashType HashType) string {
	var h hash.Hash
	switch hashType {
	case MD5:
		h = md5.New()
	case SHA1:
		h = sha1.New()
	case SHA256:
		h = sha256.New()
	default:
		h = sha256.New()
	}

	h.Write(data)
	return hex.EncodeToString(h.Sum(nil))
}

func HashFile(path string, hashType HashType) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", fmt.Errorf("打开文件失败: %w", err)
	}
	defer file.Close()

	return HashReader(file, hashType)
}

func HashReader(reader io.Reader, hashType HashType) (string, error) {
	var h hash.Hash
	switch hashType {
	case MD5:
		h = md5.New()
	case SHA1:
		h = sha1.New()
	case SHA256:
		h = sha256.New()
	default:
		h = sha256.New()
	}

	if _, err := io.Copy(h, reader); err != nil {
		return "", fmt.Errorf("计算哈希失败: %w", err)
	}

	return hex.EncodeToString(h.Sum(nil)), nil
}

func MD5String(data []byte) string {
	return HashString(data, MD5)
}

func MD5File(path string) (string, error) {
	return HashFile(path, MD5)
}

func SHA1String(data []byte) string {
	return HashString(data, SHA1)
}

func SHA1File(path string) (string, error) {
	return HashFile(path, SHA1)
}

func SHA256String(data []byte) string {
	return HashString(data, SHA256)
}

func SHA256File(path string) (string, error) {
	return HashFile(path, SHA256)
}
