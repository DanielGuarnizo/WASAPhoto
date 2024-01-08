package database

import ()

func (db *appdbimpl) GetComment(commentid string) (Comment, error) {
	var comment Comment
	err := db.c.QueryRow(`SELECT * FROM comments WHERE comment_id = ?`, commentid).Scan(&comment.Post_ID, &comment.Comment_ID, &comment.User_ID)
	if err != nil {
		return Comment{}, err
	}
	return comment, nil
}

func (db *appdbimpl) SetComment(comment Comment) error {
	_, err := db.c.Exec(`INSERT INTO comments (post_id, comment_id, user_id, body)`, comment.Post_ID, comment.Comment_ID, comment.User_ID, comment.Body)
	if err != nil {
		return err
	}
	return nil
}

func (db *appdbimpl) RemoveComment(commentid string) error {
	_, err := db.c.Exec(`DELETE FROM comments WHERE 	comment_id = ? `, commentid)
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
			&comment.User_ID,
			&comment.Body,
		)
		if err != nil {

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
