package database

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
	_, err := db.c.Exec(`INSERT INTO users (?,?) VALUES (?,?)`, newUserid, username)
	if err != nil {
		return err
	}
	return nil
}
