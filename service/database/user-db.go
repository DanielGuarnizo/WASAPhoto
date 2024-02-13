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
	var user User
	_, err := db.c.Exec(`UPDATE users SET username = ? WHERE user_id = ? `, newname, userid)
	if err != nil {
		return user, err
	}
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
		return nil, err
	}
	return UserID, nil
}

func (db *appdbimpl) GetUserIdPost(post_id string) (string, error) {
	var UserID string
	err := db.c.QueryRow("SELECT user_id FROM posts WHERE post_id=?", post_id).Scan(&UserID)
	if err != nil {
		return nil, err
	}
	return UserID, nil
}
