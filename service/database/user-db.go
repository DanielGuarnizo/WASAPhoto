package database

func (db *appdbimpl) Validate(username string, id string) (bool, error) {
	var count int
	var is_valid bool
	//  checking the exixtece of a user with the given username and userid and counting the occurences
	err := db.c.QueryRow(`SELECT COUNT(*) FROM users WHERE user_id= ? AND username = ?`, id, username).Scan(&count)
	if err != nil {
		return false, err
	}
	is_valid = (count == 1)

	return is_valid, nil
}

func (db *appdbimpl) SetUsername(newname string, userid string) (User, error) {
	var oldUsername string
	var user User

	// Seect old name
	err := db.c.QueryRow(`SELECT username FROM users WHERE user_id = ? `, userid).Scan(&oldUsername)
	if err != nil {
		return user, err
	}
	// Update the new name
	_, err = db.c.Exec(`UPDATE users SET username = ? WHERE user_id = ? `, newname, userid)
	if err != nil {
		return user, err
	}
	// Update folowees table with the new username
	_, err = db.c.Exec("UPDATE followees SET follower = ? WHERE follower = ?", newname, oldUsername)
	if err != nil {
		return user, err
	}
	_, err = db.c.Exec("UPDATE followees SET followed = ? WHERE followed = ?", newname, oldUsername)
	if err != nil {
		return user, err
	}

	// Update bans table with the new user name
	_, err = db.c.Exec("UPDATE bans SET banisher = ? WHERE banisher = ?", newname, oldUsername)
	if err != nil {
		return user, err
	}
	_, err = db.c.Exec("UPDATE bans SET banished = ? WHERE banished = ?", newname, oldUsername)
	if err != nil {
		return user, err
	}

	// Update the likes table
	_, err = db.c.Exec("UPDATE likes SET liker = ? WHERE liker = ?", newname, oldUsername)
	if err != nil {
		return user, err
	}

	// return the user with the new username
	err = db.c.QueryRow(`SELECT * FROM users WHERE user_id = ?`, userid).Scan(&user.User_ID, &user.Username)
	if err != nil {
		return user, err
	}

	return user, nil
}

func (db *appdbimpl) GetUserByName(username string) (User, error) {
	var user User
	err := db.c.QueryRow(`SELECT * FROM users WHERE username = ?`, username).Scan(&user.User_ID, &user.Username)
	if err != nil {
		return User{}, err
	}
	return user, nil
}

func (db *appdbimpl) CreateUser(newUserid string, username string) error {
	_, err := db.c.Exec(`INSERT INTO users (user_id,username) VALUES (?,?)`, newUserid, username)
	if err != nil {
		return err
	}
	return nil
}

func (db *appdbimpl) GetUserID(username string) (string, error) {
	var UserID string
	err := db.c.QueryRow("SELECT user_id FROM users WHERE username=?", username).Scan(&UserID)
	if err != nil {
		return UserID, err
	}
	return UserID, nil
}

func (db *appdbimpl) GetUserIdPost(post_id string) (string, error) {
	var UserID string
	err := db.c.QueryRow("SELECT user_id FROM posts WHERE post_id=?", post_id).Scan(&UserID)
	if err != nil {
		return UserID, err
	}
	return UserID, nil
}
