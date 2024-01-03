package database 

func (db *appdbimpl) SetFollow(userid string, followedid string) (error) {
    _, err := db.c.Exec(`INSERT INTO followers (follower, followed) VALUES (?,?)`, userid, followedid)
    if err != nil {
        return err
    }
    return nil 
}

func (db *appdbimpl) RemoveFollow(userid string, followedid string) (error) {
	_, err := db.c.Exec(`DELETE FROM followers WHERE follower = ? AND followed = ?`, userid, followedid)
	if err != nil {
		return err 
	}
	return nil 
}

func (db *appdbimpl) GetNumberOfFollowing(userid string) (int, error) {
	var count int
	err := db.c.QueryRow(`SELECT COUNT(*) FROM followers WHERE  follower= ?`, userid).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}


func (db *appdbimpl) GetNumberOfFollowers(userid string) (int, error) {
	var count int
	err := db.c.QueryRow(`SELECT COUNT(*) FROM followers WHERE  followed= ?`, userid).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (db *appdbimpl) GetFollowing(userid string) ([]string, error) {
	query := "SELECT * FROM followees WHERE follower = ?"
	rows, err := db.c.Query(query, userid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var followees []string

	for rows.Next() {
		var followee string
		if err := rows.Scan(&followee); err != nil {
			return nil, err
		}
		followees = append(followees, followee)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return followees, nil
}
