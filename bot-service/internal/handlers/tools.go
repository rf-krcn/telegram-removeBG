package handlers

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func addSuffixBeforeExtension(filename string, suffix string) string {
	ext := filepath.Ext(filename)
	nameWithoutExt := strings.TrimSuffix(filename, ext)
	newFilename := nameWithoutExt + suffix + ext
	return newFilename
}

func isImage(filename string) bool {
	imageExtensions := []string{".jpg", ".jpeg", ".png", ".gif", ".bmp"}

	// Convert the filename extension to lowercase for case-insensitive comparison
	ext := strings.ToLower(filepath.Ext(filename))

	for _, imageExt := range imageExtensions {
		if ext == imageExt {
			return true
		}
	}

	return false
}

func saveFile(data []byte) (string, error) {
	// Generate a unique filename based on the current timestamp
	uniqueFilename := generateUniqueFilename()

	// Save the file
	err := ioutil.WriteFile(uniqueFilename, data, 0644)
	if err != nil {
		return "", err
	}

	return uniqueFilename, nil
}

func generateUniqueFilename() string {
	// Use current timestamp to create a unique filename
	timestamp := time.Now().UnixNano()
	filename := fmt.Sprintf("file_%d", timestamp)

	// Get the current working directory
	wd, err := os.Getwd()
	if err != nil {
		panic(err) // Handle the error appropriately in a real application
	}

	// Combine the working directory and the generated filename
	uniqueFilename := filepath.Join(wd, filename)

	return uniqueFilename
}
