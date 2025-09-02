package utils

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func SaveFile(file *multipart.FileHeader, destinationDir, fileType string, maxSizeMB int64) (string, error) {
	allowedTypes := map[string][]string{
		"image": {"image/jpeg", "image/png", "image/gif", "image/webp"},
		"video": {"video/mp4", "video/mov"},
		"pdf":   {"application/pdf"},
	}

	if types, ok := allowedTypes[fileType]; ok {
		isAllowed := false
		contentType := file.Header.Get("Content-Type")
		for _, t := range types {
			if contentType == t {
				isAllowed = true
				break
			}
		}
		if !isAllowed {
			return "", ErrInvalidFileType
		}
	} else {
		return "", ErrInvalidFileType
	}

	if file.Size > maxSizeMB*1024*1024 {
		return "", ErrFileSizeExceeded
	}

	if _, err := os.Stat(destinationDir); os.IsNotExist(err) {
		err := os.MkdirAll(destinationDir, os.ModePerm)
		if err != nil {
			return "", fmt.Errorf("gagal membuat direktori: %v", err)
		}
	}

	extension := strings.ToLower(filepath.Ext(file.Filename))
	baseName := strings.TrimSuffix(filepath.Base(file.Filename), extension)

	baseName = sanitizeFileName(baseName)
	uniqueFileName := fmt.Sprintf("%d_%s%s", time.Now().UnixNano(), baseName, extension)
	filePath := filepath.Join(destinationDir, uniqueFileName)

	src, err := file.Open()
	if err != nil {
		return "", fmt.Errorf("gagal membuka file yang diunggah: %v", err)
	}
	defer src.Close()

	dst, err := os.Create(filePath)
	if err != nil {
		return "", fmt.Errorf("gagal membuat file tujuan: %v", err)
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		return "", fmt.Errorf("gagal menyalin file: %v", err)
	}

	return filePath, nil
}

func sanitizeFileName(filename string) string {
	filename = strings.ReplaceAll(filename, " ", "_")
	filename = strings.ReplaceAll(filename, "-", "_")
	filename = strings.ReplaceAll(filename, "(", "")
	filename = strings.ReplaceAll(filename, ")", "")

	unsafe := []string{"<", ">", ":", "\"", "|", "?", "*", "/", "\\"}
	for _, char := range unsafe {
		filename = strings.ReplaceAll(filename, char, "")
	}

	return filename
}

func DeleteFile(filePath string) error {
	if filePath == "" {
		return nil
	}

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return nil
	}

	return os.Remove(filePath)
}

func GetFileURL(filePath, baseURL string) string {
    if filePath == "" {
        return ""
    }
    
    if !strings.HasPrefix(filePath, "/") {
        filePath = "/" + filePath
    }
    
    return fmt.Sprintf("%s%s", strings.TrimSuffix(baseURL, "/"), filePath)
}
