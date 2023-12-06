package api

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"math/rand"
	"time"
)

type Like struct {
	ID     string `json:"id"`
	UserID string `json:"userId"` // Change from userId to UserID
}

type JSONErrorMsg struct {
	Message string `json:"message"`
}

var Likes = []Like{
	Like{
		ID: "1a2b3c4d5e6f7A8B9C0D1E2F3A4B5C6",
		UserID: "1a2b3c4d5e6f7A8B9C0D1E2F3A4B5C6",
	},
	Like{
		ID: "9F8E7D6C5B4A3F2E1D0C9B8A7F6E5D4",
		UserID: "9F8E7D6C5B4A3F2E1D0C9B8A7F6E5D4",
	},
}

// add a like in the post of another user
func (rt *_router) likePhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// ps are the parameters passed by the URL path, the package httprouter will retrieve those values in the path

	// this specifies the content-type the server will return to the client
	w.Header().Set("Content-Type", "application/json")

	// create a variable of type like in which we will parse the data passed in the request body
	var like Like

	// read and parse the JSON data from the request body into an Identifier object.
	// the .Decode method parses the data retrieved in the object memory address &like
	err := json.NewDecoder(r.Body).Decode(&like)

	if err != nil {
		// Handle error (e.g., invalid JSON format)
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	newID := generateUniqueID()

	newLike := Like{
		ID:     newID,
		UserID: like.UserID,
	}

	Likes = append(Likes, newLike)

	// Encode the newLike to a variable before writing it to the response writer
	response, err := json.Marshal(newLike)
	if err != nil {
		rt.baseLogger.WithError(err).Warning("like returned an error in encoding")
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(JSONErrorMsg{Message: "internal server error or bad request body"})
		return
	}

	// Write the encoded response to the response writer
	_, err = w.Write(response)
	if err != nil {
		rt.baseLogger.WithError(err).Warning("like returned an error in writing response")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

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

