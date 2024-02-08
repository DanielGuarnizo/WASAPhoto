package database

import (
	"database/sql"
	"io/ioutil"
)

func (db *appdbimpl) UploadPhoto(post Post) error {
	_, err := db.c.Exec(`INSERT INTO posts (user_id, post_id, uploaded, image, numberOfComments, numberOfLikes) VALUES (?,?,?,?,?,?)`, post.User_ID, post.Post_ID, post.Uploaded, post.Image, post.NumberOfComments, post.NumberOfLikes)
	if err != nil {
		return err
	}
	return nil
}

func (db *appdbimpl) DeletePhoto(postid string) error {
	_, err := db.c.Exec(`DELETE FROM posts WHERE post_id = ?`, postid)
	if err != nil {
		return err
	}
	return nil
}

func (db *appdbimpl) GetPhotos(userid string) ([]Post, error) {
	// Execute the query to get posts for a given userid

	query := "SELECT * FROM posts WHERE user_id = ?"
	rows, err := db.c.Query(query, userid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Iterate through the result set and populate the Post slice
	var posts []Post
	for rows.Next() {
		var post Post
		err := rows.Scan(
			&post.User_ID,
			&post.Post_ID,
			&post.Uploaded,
			&post.Image, // Assuming Image is a string column in the database
			&post.NumberOfComments,
			&post.NumberOfLikes,
		)
		if err != nil {
			continue
		}

		// Populate the Comments and Likes lists
		post.Comments, _ = db.GetComments(post.Post_ID)
		post.Likes, _ = db.GetLikes(post.Post_ID)
		post.NumberOfComments = len(post.Comments)
		post.NumberOfLikes = len(post.Likes)
		post.Image, err = loadImageFromFileSystem(post.Image)
		if err != nil {
			return nil, err
		}

		posts = append(posts, post)
	}

	// Check for errors from iterating over rows
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return posts, nil
}

func (db *appdbimpl) GetLastPosts(usernames []string) ([]Post, error) {

	var posts []Post

	for _, username := range usernames {
		// get userid given the username
		var userid string
		err := db.c.QueryRow("SELECT user_id FROM users WHERE username=?", username).Scan(&userid)

		query := `SELECT * FROM posts WHERE user_id = ? ORDER BY uploaded DESC LIMIT 1`
		rows, err := db.c.Query(query, userid)
		if err != nil {
			return nil, err
		}
		defer rows.Close()
		for rows.Next() {
			var post Post
			err := rows.Scan(
				&post.User_ID,
				&post.Post_ID,
				&post.Uploaded,
				&post.Image, // Assuming Image is a string column in the database
				&post.NumberOfComments,
				&post.NumberOfLikes,
			)
			if err != nil {
				continue
			}

			// Populate the Comments and Likes lists
			post.Comments, _ = db.GetComments(post.Post_ID)
			post.Likes, _ = db.GetLikes(post.Post_ID)
			post.NumberOfComments = len(post.Comments)
			post.NumberOfLikes = len(post.Likes)
			post.Image, err = loadImageFromFileSystem(post.Image)
			if err != nil {
				return nil, err
			}

			posts = append(posts, post)
		}
		if err := rows.Err(); err != nil {
			return nil, err
		}
	}

	return posts, nil
	// return lastPosts, nil
}
func (db *appdbimpl) GetUserIDForPost(postID string) (string, error) {
	var User_ID string

	query := "SELECT user_id FROM posts WHERE post_id = ?"
	err := db.c.QueryRow(query, postID).Scan(&User_ID)

	if err != nil {
		if err == sql.ErrNoRows {
			// No post found for the given post ID
			return "", err
		}
		return "", err
	}

	return User_ID, nil
}

func loadImageFromFileSystem(imagePath string) (string, error) {
	// Read the image file
	imageBytes, err := ioutil.ReadFile(imagePath)
	if err != nil {
		return "", err
	}

	// Convert image bytes to string
	imageString := string(imageBytes)

	return imageString, nil
}
