package database

func (db *appdbimpl) SetComment(postid string, commentid string, commenter string, userid string, body string) error {
	_, err := db.c.Exec(`INSERT INTO comments (post_id, comment_id, commenter, user_id, body) VALUES (?,?,?,?,?)`, postid, commentid, commenter, userid, body)
	if err != nil {
		return err
	}
	return nil
}

func (db *appdbimpl) RemoveComment(commentid string) error {
	_, err := db.c.Exec(`DELETE FROM comments WHERE comment_id = ? `, commentid)
	if err != nil {
		return err
	}
	return nil
}

func (db *appdbimpl) GetComments(postid string) ([]Comment, error) {
	// Execute the query to get comments for a given postID
	query := "SELECT * FROM comments WHERE post_id = ?"
	rows, err := db.c.Query(query, postid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Iterate through the result set and populate the Comment slice
	var comments []Comment
	for rows.Next() {
		var comment Comment
		err := rows.Scan(
			&comment.Post_ID,
			&comment.Comment_ID,
			&comment.Commenter,
			&comment.User_ID,
			&comment.Body,
		)
		if err != nil {
			// return nil, err
			continue
		}
		comments = append(comments, comment)
	}

	// Check for errors from iterating over rows
	if err := rows.Err(); err != nil {
		return nil, err

	}

	return comments, nil
}
