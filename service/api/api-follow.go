package api

import (
	"WASAPhoto/service/api/reqcontext"
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (rt *_router) followUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	w.Header().Set("Content-Type", "application/json")

	userid := ps.ByName("userid")
	searchUsername := ps.ByName("username")
	// to handle the error if the request wa made in a wrong way
	if userid == "" || searchUsername == "" {
		// Handle the case when "likeid" is not present in the request.
		rt.baseLogger.Warning("the userid or username in the path is empty")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Authentication
	// authorized := Authentication(w, r, userid)
	// if authorized == false {
	// 	return
	// }

	// Retrieve usename of the user by using userid
	username, err := rt.db.GetName(userid)

	// Save into the database the relation between the follower(userid) and the
	// followed(property of the object pass into the request body)
	err = rt.db.SetFollow(username, searchUsername, userid)
	if err != nil {
		rt.baseLogger.WithError(err).Warning("Error saving follow into database")
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(JSONErrorMsg{Message: "Internal server error"})
		return
	}

	// Respond with a 204 status code (success, no content)
	w.WriteHeader(http.StatusNoContent)
}

func (rt *_router) unfollowUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	w.Header().Set("Content-Type", "application/json")

	// load parameter from the path
	userid := ps.ByName("userid")
	searchUsername := ps.ByName("username")

	// to handle the error if the request wa made in a wrong way
	if userid == "" || searchUsername == "" {
		// Handle the case when "likeid" is not present in the request.
		rt.baseLogger.Warning("the userid or username in the path is empty")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Authentication
	// authorized := Authentication(w, r, userid)
	// if authorized == false {
	// 	return
	// }

	// Remove from the database the follow given the userid and the followedid
	err := rt.db.RemoveFollow(userid, searchUsername)
	if err != nil {
		// check it the err is of the same type of sql.ErrNoRows
		if errors.Is(err, sql.ErrNoRows) {
			// Resource not found
			rt.baseLogger.WithError(err).Warning("follow not found in the database")
			w.WriteHeader(http.StatusNotFound)
			_ = json.NewEncoder(w).Encode(JSONErrorMsg{Message: "follow not found in the database"})
			return
		}

		// Handle other errors
		rt.baseLogger.WithError(err).Warning("Error removing follow from the database")
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(JSONErrorMsg{Message: "Internal server error"})
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (rt *_router) getFollowing(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	w.Header().Set("Content-Type", "application/json")
	userid := ps.ByName("userid")

	// Retrieve usename of the user by using userid
	username, err := rt.db.GetName(userid)
	// Authentication

	// authorized := Authentication(w, r, userid)
	// if authorized == false {
	// 	return
	// }

	// Get the list of following from the database
	UserList, err := rt.db.GetFollowingN(username)
	if err != nil {
		rt.baseLogger.WithError(err).Warning("Error retreaving user list of following ")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Encode the list of users  in to a variable before writing it to the response writer
	response, err := json.Marshal(UserList)
	if err != nil {
		rt.baseLogger.WithError(err).Warning("Error encoding UserList for response")
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(JSONErrorMsg{Message: "Internal server error"})
		return
	}

	// Write the encoded response to the response writer
	w.WriteHeader(http.StatusOK)
	// rt.baseLogger.Warning(UserList)
	_, err = w.Write(response)
	if err != nil {
		rt.baseLogger.WithError(err).Warning("Error writing response")
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(JSONErrorMsg{Message: "Internal server error"})
		return
	}

}

func (rt *_router) getFollowers(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	w.Header().Set("Content-Type", "application/json")
	userid := ps.ByName("userid")
	// Retrieve usename of the user by using userid
	username, err := rt.db.GetName(userid)

	// Authentication

	// authorized := Authentication(w, r, userid)
	// if authorized == false {
	// 	return
	// }

	// Get the list of following from the database
	UserList, err := rt.db.GetFollowers(username)
	if err != nil {
		rt.baseLogger.WithError(err).Warning("Error retreaving user list of followers")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Encode the list of users  in to a variable before writing it to the response writer
	response, err := json.Marshal(UserList)
	if err != nil {
		rt.baseLogger.WithError(err).Warning("Error encoding UserList for response")
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(JSONErrorMsg{Message: "Internal server error"})
		return
	}

	// Write the encoded response to the response writer
	w.WriteHeader(http.StatusOK)
	// rt.baseLogger.Warning(UserList)
	_, err = w.Write(response)
	if err != nil {
		rt.baseLogger.WithError(err).Warning("Error writing response")
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(JSONErrorMsg{Message: "Internal server error"})
		return
	}

}
