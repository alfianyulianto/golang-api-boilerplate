package storage

import "mime/multipart"

type StorageProvider interface {
	UploadFile(file *multipart.FileHeader, path string) (string, error)
	DeleteFile(path string) error
}
