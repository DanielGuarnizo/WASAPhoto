package database

func (db *appdbimpl) SetLike(postid string, liker string, userid string) error {

	_, err := db.c.Exec(`INSERT INTO likes (post_id, liker, user_id) VALUES (?,?,?)`, postid, liker, userid)
	if err != nil {
		return err
	}
	return nil
}

func (db *appdbimpl) RemoveLike(userid string, postid string) error {
	_, err := db.c.Exec(`DELETE FROM likes WHERE user_id = ? AND post_id = ?`, userid, postid)
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
			&like.Liker,
			&like.User_ID,
		)
		if err != nil {
			return nil, err
		}
		likes = append(likes, like)
	}

	// Check for errors from iterating over rows
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return likes, nil
}

func (db *appdbimpl) GetLikers(postid string) ([]string, error) {
	var UserList []string

	query := `SELECT liker FROM likes WHERE post_id = ?`
	rows, err := db.c.Query(query, postid)
	if err != nil {
		return UserList, err
	}
	defer rows.Close()

	// retrieve wors from wuery and put then in a list of strings
	for rows.Next() {
		var liker string
		if err := rows.Scan(&liker); err != nil {
			// return nil, err
			continue
		}
		UserList = append(UserList, liker)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return UserList, nil
}
