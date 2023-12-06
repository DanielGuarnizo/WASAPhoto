package api

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
	"encoding/json"
	"strconv"
	"math/rand"
	"time"
)

type Comment struct {
	ID     string `json:"id"`
	UserID string `json:"userId"` // Change from userId to UserID
	Body   string `json:"body"`
}

var Comments = []Comment{
	Comment{
		ID: 	"1a2b3c4d5e6f7A8B9C0D1E2F3A4B5C6",
		UserID: "1a2b3c4d5e6f7A8B9C0D1E2F3A4B5C6",
		Body: 	"this is the first comment",
	},
	Comment{
		ID: 	"9F8E7D6C5B4A3F2E1D0C9B8A7F6E5D4",
		UserID: "9F8E7D6C5B4A3F2E1D0C9B8A7F6E5D4",
		Body: 	"this is a nother comment",
	},
}

func (rt *_router) commentPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("content-type", "application/json")

	var comment Comment 

	err := json.NewDecoder(r.Body).Decode(&comment)

	if err != nil {
		// this error happend when the request has a bad json structure 
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return 
	}

	newID := generateUniqueID()

	newComment := Comment{
		ID: comment.ID,
		UserID: comment.UserID,
		Body: comment.Body,
	}

	Comments = append(Comment, newComment)

	// Encode the newComment to a variable before writing it to the response writer
	response, err := json.Marshal(newComment)

	if err != nil {
		rt.baseLogger.WithError(err).Warning("Comment returned an error in encoding")
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(JSONErrorMsg{Message: "internal server error or bad request body"})
		return
	}

	// Write the encoded response to the response writer
	_, err = w.Write(response)
	if err != nil {
		rt.baseLogger.WithError(err).Warning("Comment returned an error in writing response")
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
