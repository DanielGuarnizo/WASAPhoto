package api

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"net/http"

	"errors"
	"database/sql"
)

type JSONErrorMsg struct {
	Message string `json:"message"`
}


// add a like in the post of another user
// ps are the parameters passed by the URL path, the package httprouter will retrieve those values in the path
// ctx is a contect object, is a way of passing data accross API boundaries 
func (rt *_router) likePhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Specify the content-type the server will return to the client
	w.Header().Set("Content-Type", "application/json")

	
	// Create a variable of type Like to parse the data passed in the request body
	var l Like
	
	// Read and parse the JSON data from the request body into a Like object.
	// The .Decode method parses the data into the object's memory address &l.
	err := json.NewDecoder(r.Body).Decode(&l)
	if err != nil {
		// Handle error (e.g., invalid JSON format)
		rt.baseLogger.WithError(err).Warning("invalid JSON format")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(JSONErrorMsg{Message: "bad request body"})
		return
	}
	
	// defer closing the request body
    defer r.Body.Close()
	userid := l.User_ID
	// Authentication 
	authorized := Authentication(w,r,userid)
	if authorized == false{
		return 
	}
	
	// Generate a new unique ID for the like
	likeid := generateUniqueID()
	newlike := Like{
		Post_ID: l.Post_ID,
		Like_ID: likeid,
		User_ID: l.User_ID,
	}

	// Save the new like to the database
	err = rt.db.SetLike(newlike.LikeToDataBase())
	if err != nil {
		rt.baseLogger.WithError(err).Warning("Error saving like into database")
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(JSONErrorMsg{Message: "Internal server error"})
		return
	}

	// Encode the newLike to a variable before writing it to the response writer
	response, err := json.Marshal(newlike)
	if err != nil {
		rt.baseLogger.WithError(err).Warning("Error encoding like for response")
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(JSONErrorMsg{Message: "Internal server error"})
		return
	}

	// Write the encoded response to the response writer
	w.WriteHeader(http.StatusCreated)
	_, err = w.Write(response)
	if err != nil {
		rt.baseLogger.WithError(err).Warning("Error writing response")
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(JSONErrorMsg{Message: "Internal server error"})
		return
	}
}



func (rt *_router) unlikePhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params){
	w.Header().Set("Content-Type", "application/json")

	likeid := ps.ByName("likeid")

	// to handle the error if the request wa made in a wrong way
	if likeid == "" {
		// Handle the case when "likeid" is not present in the request.
		rt.baseLogger.Warning("the likeid in the path is empty")
		w.WriteHeader(http.StatusBadRequest)
		return 
	}

	dblike, err := rt.db.GetLike(likeid)
	if err != nil {
		rt.baseLogger.WithError(err).Warning("like not found in the database")
		w.WriteHeader(http.StatusNotFound)
		_ = json.NewEncoder(w).Encode(JSONErrorMsg{Message: "comment not found in the database"})
		return
	}
	var l Like
	l.LikeFromDataBase(dblike)
	userid := l.User_ID

	// Authentication 
	authorized := Authentication(w,r,userid)
	if authorized == false{
		return 
	}


	// Remove from the database the like given the likeid 
	err = rt.db.RemoveLike(likeid)
	if err != nil {
		// check it the err is of the same type of sql.ErrNoRows
		if errors.Is(err, sql.ErrNoRows) {
			// Resource not found
			rt.baseLogger.WithError(err).Warning("like not found in the database")
			w.WriteHeader(http.StatusNotFound)
			_ = json.NewEncoder(w).Encode(JSONErrorMsg{Message: "like not found in the database"})
			return
		}

		// Handle other errors
		rt.baseLogger.WithError(err).Warning("Error removing like from the database")
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(JSONErrorMsg{Message: "Internal server error"})
		return
	}

	w.WriteHeader(http.StatusNoContent)

}

