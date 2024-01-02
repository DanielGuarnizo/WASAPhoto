package api

import (
	"net/http"
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"errors"
	"database/sql"
)



// var Comments = []Comment{
// 	Comment{
// 		ID: 	"1a2b3c4d5e6f7A8B9C0D1E2F3A4B5C6",
// 		UserID: "1a2b3c4d5e6f7A8B9C0D1E2F3A4B5C6",
// 		Body: 	"this is the first comment",
// 	},
// 	Comment{
// 		ID: 	"9F8E7D6C5B4A3F2E1D0C9B8A7F6E5D4",
// 		UserID: "9F8E7D6C5B4A3F2E1D0C9B8A7F6E5D4",
// 		Body: 	"this is a nother comment",
// 	},
// }

func (rt *_router) commentPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("content-type", "application/json")

	var comment Comment 

	// fetch the comment pass in the request body
	err := json.NewDecoder(r.Body).Decode(&comment)
	if err != nil {
		// this error happend when the request has a bad json structure
		rt.baseLogger.WithError(err).Warning("Invalid JSON format")
		w.WriteHeader(http.StatusBadRequest)
		_= json.NewEncoder(w).Encode(JSONErrorMsg{Message: "Bad request body"})
		return 
	}

	// defer closing the request body
    defer r.Body.Close()

	// Generated a new unique ID for the comment 
	commentid := generateUniqueID()
	newComment := Comment{
		Post_ID: comment.Post_ID,
		Comment_ID: commentid,
		User_ID: comment.User_ID,
		Body: comment.Body,
	}

	// Save the new comment in the database 
	err = rt.db.SetComment(newComment.CommentToDataBase())
	if err != nil {
		rt.baseLogger.WithError(err).Warning("Error saving comment into database")
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(JSONErrorMsg{Message: "Internal server error"})
		return
	}

	// Encode the newComment to a variable before writing it to the response writer
	response, err := json.Marshal(newComment)
	if err != nil {
		rt.baseLogger.WithError(err).Warning("Comment returned an error in encoding")
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(JSONErrorMsg{Message: "Internal server error"})
		return
	}

	// Write the encoded response to the response writer
	w.WriteHeader(http.StatusCreated)
	_, err = w.Write(response)
	if err != nil {
		rt.baseLogger.WithError(err).Warning("Comment returned an error in writing response")
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(JSONErrorMsg{Message: "Internal server error"})
		return
	}

}


func (rt *_router) uncommentPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("content-type", "application/json")
	
	commentid := ps.ByName("commentid")
	if commentid == "" {
		rt.baseLogger.Warning("The commentid in the path is empty")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Remove from the database the comment given the commentid 
	err := rt.db.RemoveComment(commentid)
	if err != nil {
		// check it the err is od the same type of sql.ErrNoRows
		if errors.Is(err, sql.ErrNoRows) {
			// Resource not found
			rt.baseLogger.WithError(err).Warning("comment not found in the database")
			w.WriteHeader(http.StatusNotFound)
			_ = json.NewEncoder(w).Encode(JSONErrorMsg{Message: "Comment not found in the database"})
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

