package database

// GetName is an example that shows you how to query data
func (db *appdbimpl) GetName(userid string) (string, error) {
	var name string
	err := db.c.QueryRow("SELECT username FROM users WHERE user_id=?", userid).Scan(&name)
	return name, err
}
