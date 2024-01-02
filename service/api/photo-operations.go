package api 

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
	"encoding/json"
	"errors"
	"database/sql"
)
// // ID represents a general ID format
// type ID string

// // Comment represents a comment on a post
// type Comment struct {
// 	ID     ID    `json:"id"`
// 	UserID ID    `json:"userId"`
// 	Body   string `json:"body"`
// }

// // Like represents a like on a post
// type Like struct {
// 	ID     ID `json:"id"`
// 	UserID ID `json:"userId"`
// }

// // Image represents an image of a post
// type Image struct {
// 	Image string `json:"image"`
// }

// // Post represents a user's post on their profile
// type Post struct {
// 	ID              ID       `json:"id"`
// 	Uploaded        string   `json:"uploaded"`
// 	Image           Image    `json:"image"`
// 	Comments        []Comment `json:"comments"`
// 	NumberOfComments int      `json:"numberOfComments"`
// 	Likes           []Like   `json:"likes"`
// 	NumberOfLikes   int      `json:"numberOfLikes"`
// }

// // Profile represents a user's profile on the app
// type Profile struct {
// 	ID             ID      `json:"id"`
// 	Username       string  `json:"username"`
// 	Photos         []Post  `json:"photos"`
// 	NumberOfPosts  int     `json:"numberOfPosts"`
// 	UserFollowers  int     `json:"userFollowers"`
// 	UserFollowing  int     `json:"userFollowing"`
// }

// // Example Profile 1
// profile1 := Profile{
// 	ID:            "1a2b3c4d5e6f7A8B9C0D1E2F3A4B5C6D",
// 	Username:      "user1",
// 	Photos:        []Post{},
// 	NumberOfPosts: 0,
// 	UserFollowers: 2,
// 	UserFollowing: 1985,
// }

// // Example Profile 2
// profile2 := Profile{
// 	ID:            "9e8d7c6b5a4A3B2C1D0E9F8A7B6C5D4E3F",
// 	Username:      "user2",
// 	Photos:        []Post{},
// 	NumberOfPosts: 0,
// 	UserFollowers: 5,
// 	UserFollowing: 100,
// }


// // Example Post
// post := Post{
// 	ID:              "1a2b3c4d5e6f7A8B9C0D1E2F3A4B5C6D",
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

func (rt *_router) uploadPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")
	
	// get from the path the userid and handle error 
	userid := ps.ByName("userid")
	if userid == "" {
		rt.baseLogger.Warning("the userid in the path is empty")
		w.WriteHeader(http.StatusBadRequest)
		return 
	}

	// create a variable post to parse the data passed in the request body 
	var p Post

	// Read and parse the JSON data from the request body into a Post object.
	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		rt.baseLogger.WithError(err).Warning("Invalid JSON format")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(JSONErrorMsg{Message: "bad request body"})
		return 
	}
	
	// defer closing the request body
    defer r.Body.Close()

	// Generate a new unique ID for the like 
	postid := generateUniqueID()
	newpost := Post{
		User_ID: userid,
		Post_ID: postid,
		Uploaded: p.Uploaded,
		Image: p.Image,
	}

	// Save the new post in the database 
	err = rt.db.UploadPhoto(newpost.PostToDataBase())
	if err != nil {
		rt.baseLogger.WithError(err).Warning("Error saving post into database")
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(JSONErrorMsg{Message: "Internal server error"})
		return
	}

	// Encode the newpost to a variable before writing it to the response writer
	response, err := json.Marshal(newpost)
	if err != nil {
		rt.baseLogger.WithError(err).Warning("Error encoding post for response")
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

func (rt *_router) deletePhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")

	postid := ps.ByName("postid")
	if postid == "" {
		// Handle the case when "likeid" is not present in the request.
		rt.baseLogger.Warning("the postid in the path is empty")
		w.WriteHeader(http.StatusBadRequest)
		return 
	}

	// Remove from the database the post given the postid 
	err := rt.db.DeletePhoto(postid)
	if err != nil {
		// check it the err is of the same type of sql.ErrNoRows
		if errors.Is(err, sql.ErrNoRows) {
			// Resource not found
			rt.baseLogger.WithError(err).Warning("Post not found in the database")
			w.WriteHeader(http.StatusNotFound)
			_ = json.NewEncoder(w).Encode(JSONErrorMsg{Message: "Post not found in the database"})
			return
		}

		// Handle other errors
		rt.baseLogger.WithError(err).Warning("Error removing post from the database")
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(JSONErrorMsg{Message: "Internal server error"})
		return
	}

	w.WriteHeader(http.StatusNoContent)
}