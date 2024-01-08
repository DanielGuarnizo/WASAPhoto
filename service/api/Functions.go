package api

// import (
// 	"math/rand"
// 	"time"
// )

// // Function to generate a unique ID
// func generateUniqueID() string {
// 	rand.Seed(time.Now().UnixNano())

// 	const charset = "0123456789abcdefABCDEF"
// 	const length = 32

// 	b := make([]byte, length)
// 	for i := range b {
// 		b[i] = charset[rand.Intn(len(charset))]
// 	}

// 	return string(b)
// }

import (
	"encoding/json"
	"github.com/google/uuid"
	"net/http"
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
