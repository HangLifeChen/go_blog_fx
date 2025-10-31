package utils

import (
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
)

// upload image
func UploadImage(handler *multipart.FileHeader, saveDir string, maxFileSize int64, allowedExtensions map[string]bool) (string, error) {
	saveDir = filepath.Join(path.Join(GetRootDir(), "uploads"), saveDir)
	// limit file size
	if handler.Size > maxFileSize {
		return "", fmt.Errorf("file size exceeds the limit of %s", FormatFileSize(maxFileSize))
	}

	if len(allowedExtensions) > 0 {
		ext := filepath.Ext(handler.Filename)
		if _, ok := allowedExtensions[ext]; !ok {
			return "", errors.New("unsupported file type")
		}
	}

	file, _ := handler.Open()
	defer file.Close()

	if err := os.MkdirAll(saveDir, os.ModePerm); err != nil {
		return "", err
	}

	// generate unique file name
	fileName := GenerateUniqueFileName(handler.Filename)
	savePath := filepath.Join(saveDir, fileName)

	// create destination file
	dst, err := os.Create(savePath)
	if err != nil {
		return "", err
	}
	defer dst.Close()

	// copy uploaded file content to destination file
	if _, err := io.Copy(dst, file); err != nil {
		return "", err
	}

	return savePath, nil
}

// get root directory
func GetRootDir() string {
	baseDir, _ := os.Getwd()
	return filepath.Dir(filepath.Dir(baseDir))
}

// format file size, automatically convert to B, KB, MB, GB, or TB format
func FormatFileSize(size int64) string {
	const (
		KB = 1024
		MB = KB * 1024
		GB = MB * 1024
		TB = GB * 1024
	)

	switch {
	case size >= TB:
		return fmt.Sprintf("%.2f TB", float64(size)/float64(TB))
	case size >= GB:
		return fmt.Sprintf("%.2f GB", float64(size)/float64(GB))
	case size >= MB:
		return fmt.Sprintf("%.2f MB", float64(size)/float64(MB))
	case size >= KB:
		return fmt.Sprintf("%.2f KB", float64(size)/float64(KB))
	default:
		return fmt.Sprintf("%d B", size)
	}
}

// generate unique file name
func GenerateUniqueFileName(originalName string) string {
	// get file extension
	ext := filepath.Ext(originalName)

	// generate UUID
	uniqueID := uuid.New().String()

	// timestamp (accurate to milliseconds)
	timestamp := time.Now().UnixMilli()

	// combine file name: UUID_timestamp.extension
	return fmt.Sprintf("%s_%d%s", uniqueID, timestamp, ext)
}

func GetDirAndParentPath(path string) (string, string) {
	cleanPath := strings.TrimRight(path, "/") // 去除末尾的斜杠

	dirName := filepath.Base(cleanPath)
	parentPath := filepath.Dir(cleanPath)

	if parentPath == "." {
		parentPath = "/"
	}

	return dirName, parentPath
}

// dir, parent := GetDirAndParentPath("/a/b/c/")
// dir = "c"
// parent = "/a/b"
