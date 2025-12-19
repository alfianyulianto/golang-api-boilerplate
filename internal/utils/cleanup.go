package utils

import (
	"github.com/alfianyulianto/pds-service/pkg/storage"
)

func CleanUpFilesOnFail(storage storage.StorageProvider, success *bool, paths ...string) {
	if !*success {
		for _, path := range paths {
			if path != "" {
				_ = storage.DeleteFile(path)
			}
		}
	}
}
