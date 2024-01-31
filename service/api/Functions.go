package api

import (
	"encoding/json"
	"net/http"
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

func Authentication(w http.ResponseWriter, r *http.Request, reqUserid string) bool {
	userToken := r.Header.Get("Authorization")
	if userToken != reqUserid || userToken == "" {
		handleUnauthorizedError(w)
		return false
	} else {
		return true
	}
}

func handleUnauthorizedError(w http.ResponseWriter) {
	w.WriteHeader(http.StatusUnauthorized)

	// You can include additional information in the response body if needed
	response := map[string]string{
		"error":   "UnauthorizedError",
		"message": "Authentication information is missing or invalid",
	}

	// Convert the response to JSON and write it to the response body
	_ = json.NewEncoder(w).Encode(response)
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
