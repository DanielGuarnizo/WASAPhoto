package database

func (db *appdbimpl) SetFollow(username string, searchUsername string, userid string) error {
	_, err := db.c.Exec(`INSERT INTO followees (follower, followed, user_id) VALUES (?,?,?)`, username, searchUsername, userid)
	if err != nil {
		return err
	}
	return nil
}

func (db *appdbimpl) RemoveFollow(follower string, followed string) error {
	_, err := db.c.Exec(`DELETE FROM followees WHERE follower = ? AND followed = ?`, follower, followed)
	if err != nil {
		return err
	}
	return nil
}

func (db *appdbimpl) GetNumberOfFollowing(username string) (int, error) {
	var count int
	err := db.c.QueryRow(`SELECT COUNT(*) FROM followees WHERE  follower= ?`, username).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (db *appdbimpl) GetNumberOfFollowers(username string) (int, error) {
	var count int
	err := db.c.QueryRow(`SELECT COUNT(*) FROM followees WHERE  followed= ?`, username).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (db *appdbimpl) GetFollowing(username string) ([]string, error) {
	query := "SELECT followed FROM followees WHERE follower = ?"
	rows, err := db.c.Query(query, username)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var following []string

	for rows.Next() {
		var followee string
		if err := rows.Scan(&followee); err != nil {
			return nil, err
		}
		following = append(following, followee)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return following, nil
}

func (db *appdbimpl) GetFollowingN(username string) ([]string, error) {
	query := "SELECT followed FROM followees WHERE follower = ?"
	rows, err := db.c.Query(query, username)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var following []string

	for rows.Next() {
		var followee string
		if err := rows.Scan(&followee); err != nil {
			return nil, err
		}
		following = append(following, followee)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return following, nil
}

func (db *appdbimpl) GetFollowers(username string) ([]string, error) {
	query := "SELECT follower FROM followees WHERE followed = ?"
	rows, err := db.c.Query(query, username)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var followers []string

	for rows.Next() {
		var follower string
		if err := rows.Scan(&follower); err != nil {
			return nil, err
		}
		followers = append(followers, follower)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return followers, nil
}
