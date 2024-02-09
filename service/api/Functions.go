package api

import (
	"os"
	"path/filepath"

	"github.com/google/uuid"

	"fmt"
	"io/ioutil"
)

// Function to generate a unique ID using UUID
func generateUniqueID() string {
	id := uuid.New()
	return id.String()
}

func saveImageToFileSystem(postID string, image string) (string, error) {
	// Get the current working directory
	currentDir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	// Navigate up to the root directory
	rootDir := filepath.Join(currentDir, "..", "..")

	// Create the directory if it doesn't exist
	photoDir := filepath.Join(rootDir, "service", "api", "photos")
	if err := os.MkdirAll(photoDir, 0755); err != nil {
		return "", err
	}

	// Save the image file
	imagePath := filepath.Join(photoDir, fmt.Sprintf("%s.jpg", postID))
	err = ioutil.WriteFile(imagePath, []byte(image), 0644)
	if err != nil {
		return "", err
	}

	return imagePath, nil
}
