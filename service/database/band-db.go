package database

import (
	"database/sql"
	"errors"
)

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
			// continue
			return nil, err
		}
		UserList = append(UserList, banished)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return UserList, nil
}

func (db *appdbimpl) BanUser(banisher string, banished string, userid string) error {
	_, err := db.c.Exec(`INSERT INTO bans (banisher, banished, user_id) VALUES (?,?,?)`, banisher, banished, userid)
	if err != nil {
		return err
	}

	// Remove followed from banished to banisher
	_, err = db.c.Exec(`DELETE FROM followees WHERE followed = ? AND follower = ? `, banisher, banished)
	if !errors.Is(err, sql.ErrNoRows) {
		return err
	}
	if err != nil {
		return err
	}

	// Remove followed from banisher to banished
	_, err = db.c.Exec(`DELETE FROM followees WHERE followed = ? AND follower = ? `, banished, banisher)
	if !errors.Is(err, sql.ErrNoRows) {
		return err
	}
	if err != nil {
		return err
	}

	return nil
}

func (db *appdbimpl) UnbanUser(banisher string, banished string) error {
	_, err := db.c.Exec(`DELETE FROM bans WHERE banisher = ? AND banished = ? `, banisher, banished)
	if err != nil {
		return err
	}
	return nil
}
