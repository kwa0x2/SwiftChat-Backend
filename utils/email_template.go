package utils

import (
	"fmt"
	"github.com/getsentry/sentry-go"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// region "LoadTemplate" loads an HTML email template, replaces the placeholder with the provided email, and returns the final HTML string.
func LoadTemplate(templateName string, email string) (string, error) {
	// Get the current working directory
	_, err := os.Getwd()
	if err != nil {
		sentry.CaptureException(err)
		return "", fmt.Errorf("error getting current directory: %w", err) // Return an error if fetching the current directory fails
	}

	// Define the path to the template file
	templatePath := filepath.Join("email_templates", templateName+".html")

	// Read the content of the HTML template file
	htmlContent, readErr := ioutil.ReadFile(templatePath)
	if readErr != nil {
		sentry.CaptureException(readErr)
		return "", fmt.Errorf("error reading HTML file: %w", readErr) // Return an error if reading the file fails
	}

	// Convert the content to a string
	htmlStr := string(htmlContent)

	// Replace the placeholder [email] with the provided email address
	htmlStr = strings.ReplaceAll(htmlStr, "[email]", email)

	return htmlStr, nil // Return the final HTML string
}

// endregion
