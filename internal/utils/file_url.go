package utils

import (
	"fmt"
	"github.com/spf13/viper"
)

func BuildFileURL(config *viper.Viper, filePath string) *string {
	if filePath == "" {
		return nil
	}

	port := config.GetInt("app.port")
	baseURL := config.GetString("app.base_url")
	if baseURL[len(baseURL)-1] == '/' {
		baseURL = baseURL[:len(baseURL)-1]
	}

	if len(filePath) > 0 && filePath[0] == '/' {
		filePath = filePath[1:]
	}

	var url string
	if config.GetString("app.env") == "production" {
		url = fmt.Sprintf("%s/%s", baseURL, filePath)
	} else {
		url = fmt.Sprintf("%s:%d/%s", baseURL, port, filePath)
	}

	return &url
}
