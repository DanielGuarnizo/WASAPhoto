package api 

import (
	"WASAPhoto/service/database"
)

type Like struct {
	Post_ID string `json:"post_id"`
	Like_ID string `json:"like_id"`
	User_ID string `json:"user_id"` // Change from userId to UserID
}

func (l *Like) LikeToDataBase() database.Like {
	return database.Like {
		Post_ID: l.Post_ID,
		User_ID: l.User_ID,
	}
}

func (l *Like) LikeFromDataBase(like database.Like) {
	l.Post_ID = like.Post_ID
	l.Like_ID = like.Like_ID
	l.User_ID = like.User_ID
}

// Comment structure and function to pass object to and from different packages 
type Comment struct {
	Post_ID     string `json:"post_id"`
	Comment_ID 	string `json:"comment_id"`
	User_ID 	string `json:"user_id"` // Change from userId to UserID
	Body   		string `json:"body"`
}

func (c *Comment) CommentToDataBase() database.Comment {
	return database.Comment {
		Post_ID: c.Post_ID, 
		Comment_ID: c.Comment_ID, 
		User_ID: c.User_ID,
		Body: c.Body,
	}
}

func (c *Comment) CommentFromDataBase(comment database.Comment) {
	c.Post_ID = comment.Post_ID
	c.Comment_ID = comment.Comment_ID
	c.User_ID = comment.User_ID
	c.Body = comment.Body
}

//Image represents an image of a post
type Image struct {
	Image string `json:"image"`
}

// Post represents a user's post on their profile
type Post struct {	
	User_ID         	string   	`json:"user_id"`
	Post_ID 			string   	`json:"post_id"`
	Uploaded        	string   	`json:"uploaded"`
	Image           	Image    	`json:"image"`
	Comments        	[]Comment	`json:"comments"`
	NumberOfComments	int      	`json:"numberOfComments"`
	Likes           	[]Like   	`json:"likes"`
	NumberOfLikes   	int      	`json:"numberOfLikes"`
}




// Profile structure and function to pass object to and from different packages 
type User struct {
	User_ID        string  `json:"user_id"`
	Username       string  `json:"username"`
}

type Profile struct {
	User 			User   `json:"user"`
	Photos         []Post  `json:"photos"`
	NumberOfPosts  int     `json:"numberOfPosts"`
	UserFollowers  int     `json:"userFollowers"`
	UserFollowing  int     `json:"userFollowing"`
}

func (u *User) UserFromDataBase(user database.User) {
	u.User_ID = user.User_ID
	u.Username = user.Username
}

// function to pass object to and from different packages 

func (p *Post) PostToDataBase() database.Post {
	return database.Post{
		User_ID: p.User_ID,
		Post_ID: p.Post_ID,
		Uploaded: p.Uploaded,
		Image: database.Image{
			Image: p.Image.Image, // Access the Image field within p.Image
		},
	}
}

func (p *Post) PostFromDataBase(dbPost database.Post) {
	p.User_ID = dbPost.User_ID
	p.Post_ID = dbPost.Post_ID
	p.Uploaded = dbPost.Uploaded

	// Convert the database.Image to api.Image
	p.Image.Image = dbPost.Image.Image

	// Convert the list of database.Comment to api.Comment
	for _, dbComment := range dbPost.Comments {
		apiComment := Comment{}
		apiComment.CommentFromDataBase(dbComment)
		p.Comments = append(p.Comments, apiComment)
	}

	// Convert the list of database.Like to api.Like
	for _, dbLike := range dbPost.Likes {
		apiLike := Like{}
		apiLike.LikeFromDataBase(dbLike)
		p.Likes = append(p.Likes, apiLike)
	}
}


// GetPhotosFromDatabase converts a list of database.Post objects to a list of api.Post objects
func GetPhotosFromDatabase(dbPhotos []database.Post) []Post {
	var apiPhotos []Post

	for _, dbPhoto := range dbPhotos {
		apiPhoto := Post{}
		apiPhoto.PostFromDataBase(dbPhoto)

		// Convert comments from database format to API format
		apiPhoto.Comments = make([]Comment, len(dbPhoto.Comments))
		for i, dbComment := range dbPhoto.Comments {
			apiComment := Comment{}
			apiComment.CommentFromDataBase(dbComment)
			apiPhoto.Comments[i] = apiComment
		}

		// Convert likes from database format to API format
		apiPhoto.Likes = make([]Like, len(dbPhoto.Likes))
		for i, dbLike := range dbPhoto.Likes {
			apiLike := Like{}
			apiLike.LikeFromDataBase(dbLike)
			apiPhoto.Likes[i] = apiLike
		}

		// Add the converted api.Post to the list
		apiPhotos = append(apiPhotos, apiPhoto)
	}

	return apiPhotos
}

// Following structure and fucntion to pass pbjects to and from different packages
type FollowRequest struct {
    FollowedID       string `json:"followingID"`
    FollowingUsername string `json:"followingUsername"`
}

// Band structure 

type MuteRequest struct {
	MutedUserID string 
}

// Stream 

type Stream struct {
	Photos []Post `json:"photos"`
}