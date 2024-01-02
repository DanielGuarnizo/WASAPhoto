package database 

func (db *appdbimpl) UploadPhoto(post Post) (error) {
	_, err := db.c.Exec(`INSERT INTO posts (user_id, post_id, uploaded, image, numberOfComments, numberOfLikes) VALUES (?,?,?,?,?,?)`, post.User_ID, post.Post_ID, post.Uploaded, post.Image, post.NumberOfComments, post.NumberOfLikes)
	if err != nil {
		return err
	}
	return nil 
}

func (db *appdbimpl) DeletePhoto(postid string) (error) {
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
			&post.Image.Image, // Assuming Image is a string column in the database
			&post.NumberOfComments,
			&post.NumberOfLikes,
		)
		if err != nil {
			continue
		}

		// Populate the Comments and Likes lists
		post.Comments, _ = db.GetComments(post.Post_ID)
		post.Likes, _ = db.GetLikes(post.Post_ID)

		posts = append(posts, post)
	}

	// Check for errors from iterating over rows
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return posts, nil
}