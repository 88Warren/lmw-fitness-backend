package database

import (
	"fmt"
	"log"
	"os"
)

func ReadHTMLFile(filename string) (string, error) {
	filePath := "database/content/blog/" + filename
	content, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to read file %s: %w", filePath, err)
	}
	log.Printf("Successfully read HTML file: %s", filePath)
	return string(content), nil
}
