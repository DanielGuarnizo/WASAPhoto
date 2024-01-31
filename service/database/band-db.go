package database

func (db *appdbimpl) GetBans(username string) ([]string, error) {
	var UserList []string

	query := `SELECT banished FROM bans WHERE banisher = ?`
	rows, err := db.c.Query(query, username)
	if err != nil {
		return UserList, err
	}
	defer rows.Close()

	// retrieve wors from wuery and put then in a list of strings
	for rows.Next() {
		var banished string
		if err := rows.Scan(&banished); err != nil {
			continue
		}
		UserList = append(UserList, banished)
	}
	return UserList, nil
}

func (db *appdbimpl) BandUser(banisher string, banished string, userid string) error {
	_, err := db.c.Exec(`INSERT INTO bans (banisher, banished, user_id) VALUES (?,?,?)`, banisher, banished, userid)
	if err != nil {
		return err
	}
	return nil
}

func (db *appdbimpl) UnbandUser(banisher string, banished string) error {
	_, err := db.c.Exec(`DELETE FROM bans WHERE banisher = ? AND banished = ? `, banisher, banished)
	if err != nil {
		return err
	}
	return nil
}
