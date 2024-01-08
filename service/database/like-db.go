package database

import ()

func (db *appdbimpl) GetLike(likeid string) (like Like, err error) {
	err = db.c.QueryRow(`SELECT * FROM likes WHERE like_id = ?`, likeid).Scan(&like.Post_ID, &like.Like_ID, &like.User_ID)
	if err != nil {
		return Like{}, err
	}
	return like, nil
}

func (db *appdbimpl) SetLike(like Like) error {

	_, err := db.c.Exec(`INSERT INTO likes (post_id, like_id, user_id) VALUES (?,?,?)`, like.Post_ID, like.Like_ID, like.User_ID)
	if err != nil {
		return err
	}
	return nil
}

func (db *appdbimpl) RemoveLike(likeid string) error {
	_, err := db.c.Exec(`DELETE FROM likes WHERE like_id = ?`, likeid)
	if err != nil {
		return err
	}
	return nil
}

func (db *appdbimpl) GetLikes(postid string) ([]Like, error) {
	// Execute the query to get likes for a given postID
	query := "SELECT * FROM likes WHERE post_id = ?"
	rows, err := db.c.Query(query, postid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Iterate through the result set and populate the Like slice
	var likes []Like
	for rows.Next() {
		var like Like
		err := rows.Scan(
			&like.Post_ID,
			&like.Like_ID,
			&like.User_ID,
		)
		if err != nil {
			continue
		}
		likes = append(likes, like)
	}

	// Check for errors from iterating over rows
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return likes, nil
}
