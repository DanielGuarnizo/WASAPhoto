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

	// Get needed information to perfomr the operation
	userid := ps.ByName("userid")
	if userid == "" {
		// Handle the case when "likeid" is not present in the request.
		rt.baseLogger.Warning("the userid in the path is empty")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Get name that will be used to the authetication
	username, err := rt.db.GetName(userid)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error("error getting username")
		return
	}

	// Authentication
	id := r.Header.Get("Authorization")
	is_valid, err := rt.db.Validate(username, id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error("error in Validate")
		return
	}
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

	// read and parse the json data from the request body into an username upadate object
	type Body struct {
		SearchUsername string `json:"searchUsername"`
	}
	var input Body
	err = json.NewDecoder(r.Body).Decode(&input)
	if err != nil || input.SearchUsername == "" {
		// Handle error (e.g., invalid JSON format)
		rt.baseLogger.WithError(err).Warning("invalid JSON format")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(JSONErrorMsg{Message: "bad request body"})
		return
	}

	// defer closing the request body
	defer r.Body.Close()

	// Save into the database the relation between the follower(userid) and the
	// followed(property of the object pass into the request body)
	err = rt.db.SetFollow(username, input.SearchUsername, userid)
	if err != nil {
		rt.baseLogger.WithError(err).Warning("Error saving follow into database")
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(JSONErrorMsg{Message: "Internal server error"})
		return
	}

	// retrieve user_id to return
	UserID, err := rt.db.GetUserID(input.SearchUsername)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error("error getting user_id")
		return
	}

	// Encode the UserID in to a variable before writing it to the response writer
	response, err := json.Marshal(UserID)
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

func (rt *_router) unfollowUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	w.Header().Set("Content-Type", "application/json")

	// Get userid from the path and hanlde error
	userid := ps.ByName("userid")
	followed := ps.ByName("followed")
	if userid == "" || followed == "" {
		// Handle the case when "likeid" is not present in the request.
		rt.baseLogger.Warning("the userid or followed in the path is empty")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Get name that will be used to the authetication
	follower, err := rt.db.GetName(userid)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error("error getting username")
		return
	}

	// Authentication
	id := r.Header.Get("Authorization")
	is_valid, err := rt.db.Validate(follower, id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error("error in Validate")
		return
	}
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

	// Remove from the database the follow given the userid and the followed name
	err = rt.db.RemoveFollow(follower, followed)
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

	// Get userid from the path and hanlde error
	userid := ps.ByName("userid")
	if userid == "" {
		// Handle the case when "likeid" is not present in the request.
		rt.baseLogger.Warning("the userid in the path is empty")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Get name that will be used to the authetication
	username, err := rt.db.GetName(userid)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error("error getting username")
		return
	}

	// Authentication
	id := r.Header.Get("Authorization")
	is_valid, err := rt.db.Validate(username, id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error("error in Validate")
		return
	}
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

	// Get userid from the path and hanlde error
	userid := ps.ByName("userid")
	if userid == "" {
		// Handle the case when "likeid" is not present in the request.
		rt.baseLogger.Warning("the userid in the path is empty")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Retrieve usename of the user by using userid
	username, err := rt.db.GetName(userid)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error("error getting username")
		return
	}

	// Authentication
	id := r.Header.Get("Authorization")
	is_valid, err := rt.db.Validate(username, id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error("error in Validate")
		return
	}
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
