package api 

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
	"encoding/json"
	"strconv"
	"math/rand"
	"time"
)
// ID represents a general ID format
type ID string

// Comment represents a comment on a post
type Comment struct {
	ID     ID    `json:"id"`
	UserID ID    `json:"userId"`
	Body   string `json:"body"`
}

// Like represents a like on a post
type Like struct {
	ID     ID `json:"id"`
	UserID ID `json:"userId"`
}

// Image represents an image of a post
type Image struct {
	Image string `json:"image"`
}

// Post represents a user's post on their profile
type Post struct {
	ID              ID       `json:"id"`
	Uploaded        string   `json:"uploaded"`
	Image           Image    `json:"image"`
	Comments        []Comment `json:"comments"`
	NumberOfComments int      `json:"numberOfComments"`
	Likes           []Like   `json:"likes"`
	NumberOfLikes   int      `json:"numberOfLikes"`
}

// Profile represents a user's profile on the app
type Profile struct {
	ID             ID      `json:"id"`
	Username       string  `json:"username"`
	Photos         []Post  `json:"photos"`
	NumberOfPosts  int     `json:"numberOfPosts"`
	UserFollowers  int     `json:"userFollowers"`
	UserFollowing  int     `json:"userFollowing"`
}

// Example Profile 1
profile1 := Profile{
	ID:            "1a2b3c4d5e6f7A8B9C0D1E2F3A4B5C6D7E8F9",
	Username:      "user1",
	Photos:        []Post{},
	NumberOfPosts: 0,
	UserFollowers: 2,
	UserFollowing: 1985,
}

// Example Profile 2
profile2 := Profile{
	ID:            "9e8d7c6b5a4A3B2C1D0E9F8A7B6C5D4E3F2A1B0",
	Username:      "user2",
	Photos:        []Post{},
	NumberOfPosts: 0,
	UserFollowers: 5,
	UserFollowing: 100,
}


// Example Post
post := Post{
	ID:              "1a2b3c4d5e6f7A8B9C0D1E2F3A4B5C6D7E8F9",
	Uploaded:        "2023-01-01T12:34:56",
	Image:           Image{Image: "base64encodedimage"},
	Comments:        []Comment{},
	NumberOfComments: 0,
	Likes:           []Like{},
	NumberOfLikes:   0,
}

profile1.Photos = append(profile.Photos, post)

var (
	// Mutex to synchronize access to profiles slice
	profilesMutex sync.Mutex

	// Slice to store user profiles (simulating an in-memory database)
	profiles []Profile
)

// add the example profiles to the slice profiles 
profiles = []Profile{profile1, profile2}

func (rt *_router) uploadPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	userid, err := strconv.Atoi(ps.ByName("userid"))

	if err != nil {
		w.WriteHeader(StatusBadRequest)
		return 
	}

	profilesMutex.Lock()
	defer profilesMutex.Unlock()

	// find the user profile with the given userid in path 
	var userindex int
	found := false
	for i, profile := range profiles {
		if profile.ID == userid {
			userindex = i 
			found = true 
			break
		}
	}

	// check if the profile was faound or not 
	if !found {
		// User not found, return 404
		w.WriteHeader(http.StatusNotFound)
		return
	}

	var newPost Post 
	err := json.NewDecoder(r.Body).Decode(&newPost)

	if err != nil {
		w:WriteHeader(StatusBadRequest)
		return 
	}
	newID := generateUniqueID()
	newPost.ID = newID

	// add the new post ot the profile which perfomr the action 
	profiles[userindex].Photos = append(profiles[userindex].Photos, newPost)

	// Respond with the updated new post 
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(newPost)
	if err != nil {
		// Error encoding response, log the error
		fmt.Printf("Error encoding response: %v\n", err)
		return
	}
}