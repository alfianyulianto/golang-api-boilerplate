package storage

import (
	"fmt"
	"github.com/google/uuid"
	"mime/multipart"
	"os"
	"path/filepath"
)

type LocalStorage struct {
}

func NewLocalStorage() *LocalStorage {
	return &LocalStorage{}
}

func (l *LocalStorage) UploadFile(fileHeader *multipart.FileHeader, path string) (string, error) {
	file, err := fileHeader.Open()
	if err != nil {
		return "", err
	}
	defer file.Close()

	destDir := filepath.Join("./uploads", path)
	if err = os.MkdirAll(destDir, os.ModePerm); err != nil {
		return "", err
	}

	fileName := uuid.NewString() + filepath.Ext(fileHeader.Filename)
	destPath := filepath.Join(destDir, fileName)
	out, err := os.Create(destPath)
	if err != nil {
		return "", err
	}
	defer out.Close()

	if _, err = out.ReadFrom(file); err != nil {
		return "", err
	}

	return fmt.Sprintf("uploads/%s/%s", path, fileName), nil
}

func (l *LocalStorage) DeleteFile(path string) error {
	return os.Remove(path)
}
