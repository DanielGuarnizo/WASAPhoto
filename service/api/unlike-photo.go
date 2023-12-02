package api

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"strconv"
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


func (rt *_router) unlikePhoto(w http.ResponseWriter, r *http.Request, ps httprouter.parameters){
	likeid, err = strconv.Atoi(ps.Byname("likeid"))

	// to handle the error if the request wa made in a wrong way
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return 
	}

	// Iterate through Likes to find the index of the like with the matching ID
	var foundIndex := -1
	// key, value := iterated over an array of elements 
	for i, like := range Likes {
		if like.ID == likeid {
			foundIndex = i
			break
		}
	}

	// If the like with the matching ID is found, remove it from Likes
	if foundIndex != -1 {
		Likes = append(Likes[:foundIndex], Likes[foundIndex+1:]...)
		w.WriteHeader(http.StatusNoContent)
	} else {
		// If the like with the matching ID is not found, return a 404 Not Found status
		w.WriteHeader(http.StatusNotFound)
	}
}