package api

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

type Like struct {
	id     int
	UserID int // Change from userId to UserID
}

type JSONErrorMsg struct {
	Message string
}

var Likes = []Like{
	Like{
		id:     1,
		UserID: 1993238,
	},
	Like{
		id:     2,
		UserID: 1984033,
	},
}

// add a like in the post of another user
func (rt *_router) likePhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// ps are the parameters passed by the URL path, the package httprouter will retrieve those values in the path

	// this specifies the content-type the server will return to the client
	w.Header().Set("Content-Type", "application/json")

	// create a variable of type like in which we will parse the data passed in the request body
	var like Like

	// read and parse the JSON data from the request body into a Like object.
	// the .Decode method parses the data retrieved in the object memory address &like
	err := json.NewDecoder(r.Body).Decode(&like)

	if err != nil {
		// Handle error (e.g., invalid JSON format)
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	new_id := len(Likes)

	var new_like = Like{
		id:     new_id,
		UserID: like.UserID, // Change from userid to UserID
	}

	Likes = append(Likes, new_like)

	// Encode the new_like to a variable before writing it to the response writer
	response, err := json.Marshal(new_like)
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

	// ... (rest of your code if any)
}
