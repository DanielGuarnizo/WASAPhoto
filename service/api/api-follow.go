package api

import (
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"WASAPhoto/service/api/reqcontext"
)

func (rt *_router) followUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	w.Header().Set("Content-Type", "application/json")

	userid := ps.ByName("userid")
	// to handle the error if the request wa made in a wrong way
	if userid == "" {
		// Handle the case when "likeid" is not present in the request.
		rt.baseLogger.Warning("the userid in the path is empty")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Authentication
	authorized := Authentication(w, r, userid)
	if authorized == false {
		return
	}

	// Read and parse the JSON data from the request body into a FollowRequest object
	var followReq FollowRequest
	err := json.NewDecoder(r.Body).Decode(&followReq)
	if err != nil {
		rt.baseLogger.WithError(err).Warning("invalid JSON format")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(JSONErrorMsg{Message: "bad request body"})
		return
	}
	// defer closing the request body
	defer r.Body.Close()

	// Save into the database the relation between the follower(userid) and the
	//followed(property of the object pass into the request body)
	err = rt.db.SetFollow(userid, followReq.FollowedID)
	if err != nil {
		rt.baseLogger.WithError(err).Warning("Error saving follow into database")
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(JSONErrorMsg{Message: "Internal server error"})
		return
	}

	// Encode the follow (which is the same of the request body) to a variable before writing it to the response writer
	response, err := json.Marshal(followReq)
	if err != nil {
		rt.baseLogger.WithError(err).Warning("Error encoding follow for response")
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

func (rt *_router) unfollowUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	w.Header().Set("Content-Type", "application/json")

	// load parameter from the path
	userid := ps.ByName("userid")
	followedid := ps.ByName("followedid")
	// to handle the error if the request wa made in a wrong way
	if userid == "" || followedid == "" {
		// Handle the case when "likeid" is not present in the request.
		rt.baseLogger.Warning("the userid or followedid in the path is empty")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// Authentication
	authorized := Authentication(w, r, userid)
	if authorized == false {
		return
	}

	// Remove from the database the follow given the userid and the followedid
	err := rt.db.RemoveFollow(userid, followedid)
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
