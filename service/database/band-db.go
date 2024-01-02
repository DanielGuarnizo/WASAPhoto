package database 

func (db *appdbimpl) BandUser(userid string, mutedUserid string) (error) {
	_, err := db.c.Exec(`INSERT INTO bans (banisher, banished) VALUES (?,?)`, userid, mutedUserid)
	if err != nil {
		return err
	}
	return nil 
}

func (db *appdbimpl) UnbandUser(userid string, mutedUserid string) (error) {
	_, err := db.c.Exec(`DELETE FROM bans WHERE banisher = ? AND banished = ? `, userid, mutedUserid)
	if err != nil {
		return err
	}
	return nil
}