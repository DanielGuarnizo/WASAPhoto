package api

import (
	"WASAPhoto/service/api/reqcontext"
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (rt *_router) setMyUserName(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
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

	// read and parse the json data from the request body into an username upadate object
	type UsernameUpdate struct {
		NewUsername string `json:"newUsername"`
	}
	var update UsernameUpdate
	err = json.NewDecoder(r.Body).Decode(&update)
	if err != nil {
		// Handle error (e.g., invalid JSON format)
		rt.baseLogger.WithError(err).Warning("invalid JSON format")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(JSONErrorMsg{Message: "bad request body"})
		return
	}

	// defer closing the request body
	defer r.Body.Close()

	// set the new name the user pass in the request body in the database
	newname := update.NewUsername
	dbuser, err := rt.db.SetUsername(newname, userid)
	if err != nil {

		if errors.Is(err, sql.ErrNoRows) {
			// Resource not found
			rt.baseLogger.WithError(err).Warning("user not found in the database")
			w.WriteHeader(http.StatusNotFound)
			_ = json.NewEncoder(w).Encode(JSONErrorMsg{Message: "user not found in the database"})
			return
		}

		rt.baseLogger.WithError(err).Warning("Error updating new name into database")
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(JSONErrorMsg{Message: "Internal server error"})
		return
	}

	// Return the updated profile
	var updatedUser User
	updatedUser.UserFromDataBase(dbuser)

	// Serialize the updated user as JSON and write it to the response
	response, err := json.Marshal(updatedUser)
	if err != nil {
		rt.baseLogger.WithError(err).Warning("Error encoding JSON response")
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(JSONErrorMsg{Message: "Internal server error"})
		return
	}

	w.WriteHeader(http.StatusOK)
	_, err = w.Write(response)
	if err != nil {
		rt.baseLogger.WithError(err).Warning("Error writing response")
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(JSONErrorMsg{Message: "Internal server error"})
		return
	}
}

func (rt *_router) getUserProfile(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	w.Header().Set("Content-Type", "application/json")

	// Get userid and username from the path,query  and hanlde error
	userid := ps.ByName("userid")
	username := r.URL.Query().Get("username")
	if userid == "" {
		// Handle the case when "likeid" is not present in the request.
		rt.baseLogger.Warning("the userid in the path is empty or the username in query is empty")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Get name that will be used to the authetication
	name, err := rt.db.GetName(userid)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error("error getting username")
		return
	}

	// Authentication
	id := r.Header.Get("Authorization")
	is_valid, err := rt.db.Validate(name, id)
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
	}

	// Get the user from the database given the username
	var ReqUser User
	dbReqUser, err := rt.db.GetUserByName(username)
	if err != nil {
		rt.baseLogger.WithError(err).Warning("Error checking if user exists")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	ReqUser.UserFromDataBase(dbReqUser)

	// retrive all the data needed for fetch profile given the request user
	var profile Profile
	// var photos []Post
	dbPhotos, err := rt.db.GetPhotos(ReqUser.User_ID)
	if err != nil {
		rt.baseLogger.WithError(err).Warning("Error getting photos from database")
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(JSONErrorMsg{Message: "Internal server error"})
		return
	}
	// Convert the list of database.Post to a list of api.Post
	apiPhotos := GetPhotosFromDatabase(dbPhotos)

	// fetch profile
	{
		profile.User = ReqUser
		profile.Photos = apiPhotos
		profile.NumberOfPosts = len(apiPhotos)

		count1, err := rt.db.GetNumberOfFollowers(ReqUser.Username)
		if err != nil {
			rt.baseLogger.WithError(err).Warning("Error getting number of followers from database")
			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(JSONErrorMsg{Message: "Internal server error"})
			return
		}
		profile.UserFollowers = count1

		count2, err := rt.db.GetNumberOfFollowing(ReqUser.Username)
		if err != nil {
			rt.baseLogger.WithError(err).Warning("Error getting number of following from database")
			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(JSONErrorMsg{Message: "Internal server error"})
			return
		}
		profile.UserFollowing = count2
	}

	// Encode the profile to a variable before writing it to the response writer
	response, err := json.Marshal(profile)
	if err != nil {
		rt.baseLogger.WithError(err).Warning("Error encoding profile for response")
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

// /users/{userid}/stream
func (rt *_router) getMyStream(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
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

	// get the list of all the followees of the user given it's userid
	followees, err := rt.db.GetFollowing(username)
	if err != nil {
		rt.baseLogger.WithError(err).Warning("Error retrieving followwes from  database")
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(JSONErrorMsg{Message: "Internal server error"})
		return
	}
	// rt.baseLogger.Warning(followees)

	// get the list of lastest post of the following users given the user_id of each of them
	// then it returns a list of post and we conver from database.Post type to api.Post type
	dblastPosts, err := rt.db.GetLastPosts(followees)
	if err != nil {
		rt.baseLogger.WithError(err).Warning("Error getting photos from database")
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(JSONErrorMsg{Message: "Internal server error"})
		return
	}

	lastPosts := GetPhotosFromDatabase(dblastPosts)

	// we constructe the stream
	var stream Stream
	stream.Photos = lastPosts

	// Encode the stream in to a variable before writing it to the response writer
	response, err := json.Marshal(stream)
	if err != nil {
		rt.baseLogger.WithError(err).Warning("Error encoding stream for response")
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
