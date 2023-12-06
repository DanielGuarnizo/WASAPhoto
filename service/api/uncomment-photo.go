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

func (rt *_router) uncommentPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	commentid, err := strconv.Atoi(ps.Byname("commentid"))



	// Iterate through Comments to find the index of the comment with the matching ID
	var foundIndex := -1
	// key, value := iterated over an array of elements 
	for i, comment := range Comments {
		if comment.ID == commentid {
			foundIndex = i
			break
		}
	}

	// If the comments with the matching ID is found, remove it from Comments
	if foundIndex != -1 {
		Comments = append(Comments[:foundIndex], Comments[foundIndex+1:]...)
		w.WriteHeader(http.StatusNoContent)
		return 
	} else {
		// If the comment with the matching ID is not found, return a 404 Not Found status
		w.WriteHeader(http.StatusNotFound)
		return 
	}
}