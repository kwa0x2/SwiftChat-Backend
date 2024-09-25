package utils

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func LoadTemplate(templateName string, email string) (string, error) {
	_, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("error getting current directory: %w", err)
	}

	templatePath := filepath.Join("email_templates", templateName+".html")

	htmlContent, err := ioutil.ReadFile(templatePath)
	if err != nil {
		return "", fmt.Errorf("error reading HTML file: %w", err)
	}

	htmlStr := string(htmlContent)

	htmlStr = strings.ReplaceAll(htmlStr, "[email]", email)

	return htmlStr, nil
}
