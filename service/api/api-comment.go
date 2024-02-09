package api

import (
	"WASAPhoto/service/api/reqcontext"
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (rt *_router) commentPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	w.Header().Set("content-type", "application/json")

	// Get needed information to perfomr the operation
	userid := ps.ByName("userid")
	name, _ := rt.db.GetName(userid)
	id := r.Header.Get("Authorization")

	// Authentication
	is_valid, err := rt.db.Validate(name, id)
	if is_valid == false {
		w.WriteHeader(http.StatusUnauthorized)

		// You can include additional information in the response body if needed
		response := map[string]string{
			"error":   "UnauthorizedError",
			"message": "Authentication information is missing or invalid",
		}

		// Convert the response to JSON and write it to the response body
		_ = json.NewEncoder(w).Encode(response)
		return
	}

	var comment Comment

	// fetch the comment pass in the request body
	err = json.NewDecoder(r.Body).Decode(&comment)
	if err != nil {
		// this error happend when the request has a bad json structure
		rt.baseLogger.WithError(err).Warning("Invalid JSON format")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(JSONErrorMsg{Message: "Bad request body"})
		return
	}

	// defer closing the request body
	defer r.Body.Close()

	// Generated a new unique ID for the comment
	commentid := generateUniqueID()
	newComment := Comment{
		Post_ID:    comment.Post_ID,
		Comment_ID: commentid,
		Commenter:  comment.Commenter,
		User_ID:    comment.User_ID,
		Body:       comment.Body,
	}

	// Save the new comment in the database
	err = rt.db.SetComment(newComment.Post_ID, newComment.Comment_ID, newComment.Commenter, newComment.User_ID, newComment.Body)
	if err != nil {
		rt.baseLogger.WithError(err).Warning("Error saving comment into database")
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(JSONErrorMsg{Message: "Internal server error"})
		return
	}

	// Encode the comment_id in to a variable before writing it to the response writer
	response, err := json.Marshal(commentid)
	if err != nil {
		rt.baseLogger.WithError(err).Warning("Error encoding UserList for response")
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(JSONErrorMsg{Message: "Internal server error"})
		return
	}

	// Write the encoded response to the response writer
	w.WriteHeader(http.StatusCreated)
	// rt.baseLogger.Warning(UserList)
	_, err = w.Write(response)
	if err != nil {
		rt.baseLogger.WithError(err).Warning("Error writing response")
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(JSONErrorMsg{Message: "Internal server error"})
		return
	}
}

func (rt *_router) uncommentPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	w.Header().Set("content-type", "application/json")

	// Get needed information to perfomr the operation
	userid := ps.ByName("userid")
	commentid := ps.ByName("commentid")
	name, _ := rt.db.GetName(userid)
	id := r.Header.Get("Authorization")

	// Authentication
	is_valid, err := rt.db.Validate(name, id)
	if is_valid == false {
		w.WriteHeader(http.StatusUnauthorized)

		// You can include additional information in the response body if needed
		response := map[string]string{
			"error":   "UnauthorizedError",
			"message": "Authentication information is missing or invalid",
		}

		// Convert the response to JSON and write it to the response body
		_ = json.NewEncoder(w).Encode(response)
		return
	}

	// Remove from the database the comment given the commentid
	err = rt.db.RemoveComment(commentid)
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

func (rt *_router) getComments(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	w.Header().Set("Content-Type", "application/json")

	// Get from the path the needed data
	userid := ps.ByName("userid")
	postid := ps.ByName("postid")
	if postid == "" || userid == "" {
		// Handle the case when "likeid" is not present in the request.
		rt.baseLogger.Warning("the postid ot the userid in the path is empty")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Get the list of likes of a post from the database
	dbComments, err := rt.db.GetComments(postid)
	if err != nil {
		rt.baseLogger.WithError(err).Warning("Error getting likes from database")
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(JSONErrorMsg{Message: "Internal server error"})
		return
	}

	CommentsList := GetCommentsFromDatabase(dbComments)

	// Encode the CommentList in to a variable before writing it to the response writer
	response, err := json.Marshal(CommentsList)
	if err != nil {
		rt.baseLogger.WithError(err).Warning("Error encoding comment list for response")
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
