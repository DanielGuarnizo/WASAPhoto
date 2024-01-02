package api 

import (
	"net/http"
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"errors"
	"database/sql"
)

// // Example Profile 1
// profile1 := Profile{
// 	ID:            "1a2b3c4d5e6f7A8B9C0D1E2F3A4B5C6D7E8F9",
// 	Username:      "user1",
// 	Photos:        []Post{},
// 	NumberOfPosts: 0,
// 	UserFollowers: 2,
// 	UserFollowing: 1985,
// }

// // Example Profile 2
// profile2 := Profile{
// 	ID:            "9e8d7c6b5a4A3B2C1D0E9F8A7B6C5D4E3F2A1B0",
// 	Username:      "user2",
// 	Photos:        []Post{},
// 	NumberOfPosts: 0,
// 	UserFollowers: 5,
// 	UserFollowing: 100,
// }


// // Example Post
// post := Post{
// 	ID:              "1a2b3c4d5e6f7A8B9C0D1E2F3A4B5C6D7E8F9",
// 	Uploaded:        "2023-01-01T12:34:56",
// 	Image:           Image{Image: "base64encodedimage"},
// 	Comments:        []Comment{},
// 	NumberOfComments: 0,
// 	Likes:           []Like{},
// 	NumberOfLikes:   0,
// }

// profile1.Photos = append(profile.Photos, post)

// var (
// 	// Mutex to synchronize access to profiles slice
// 	profilesMutex sync.Mutex

// 	// Slice to store user profiles (simulating an in-memory database)
// 	profiles []Profile
// )

// // add the example profiles to the slice profiles 
// profiles = []Profile{profile1, profile2}


func (rt *_router) setMyUserName(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")

	// get from the path the userid and handle error 
	userid := ps.ByName("userid")
	if userid == "" {
		rt.baseLogger.Warning("the userid in the path is empty")
		w.WriteHeader(http.StatusBadRequest)
		return 
	}

	// read and parse the json data from the request body into an username upadate object 
	type UsernameUpdate struct {
		NewUsername string `json:"newUsername"`
	}
	var update UsernameUpdate
	err := json.NewDecoder(r.Body).Decode(&update)
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

func (rt *_router) getUserProfile(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")
	// retrieve the esential information to perfomr the operation
	userid := ps.ByName("userid")
	username := r.URL.Query().Get("username")

	// Authentication 
	authorized := Authentication(w,r,userid)
	if authorized == false{
		return 
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
	//var photos []Post
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

		count1, err := rt.db.GetNumberOfFollowers(ReqUser.User_ID)
		if err != nil {
			rt.baseLogger.WithError(err).Warning("Error getting number of followers from database")
			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(JSONErrorMsg{Message: "Internal server error"})
			return
		}
		profile.UserFollowers = count1

		count2, err := rt.db.GetNumberOfFollowing(ReqUser.User_ID) 
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
