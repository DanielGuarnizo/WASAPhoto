/*
Package database is the middleware between the app database and the code. All data (de)serialization (save/load) from a
persistent database are handled here. Database specific logic should never escape this package.

To use this package you need to apply migrations to the database if needed/wanted, connect to it (using the database
data source name from config), and then initialize an instance of AppDatabase from the DB connection.

For example, this code adds a parameter in `webapi` executable for the database data source name (add it to the
main.WebAPIConfiguration structure):

	DB struct {
		Filename string `conf:""`
	}

This is an example on how to migrate the DB and connect to it:

	// Start Database
	logger.Println("initializing database support")
	db, err := sql.Open("sqlite3", "./foo.db")
	if err != nil {
		logger.WithError(err).Error("error opening SQLite DB")
		return fmt.Errorf("opening SQLite: %w", err)
	}
	defer func() {
		logger.Debug("database stopping")
		_ = db.Close()
	}()

Then you can initialize the AppDatabase and pass it to the api package.
*/
package database

import (
	"database/sql"
	"errors"
	"fmt"
	"WASAPhoto/service/api/structure"
)

// AppDatabase is the high level interface for the DB
type AppDatabase interface {
	// GetName() (string, error)
	// SetName(name string) error
	// Ping() error
	AddLike(structure.Like) (structure.Like, error)
		// the name of the parameter is up to the implementation of the interface 
		// we only define the type of the parameter the method should have 
	
}

type appdbimpl struct {
	// sql.DB normally indicates a connection to a sql database 
	c *sql.DB
}

// New returns a new instance of AppDatabase based on the SQLite connection `db`.
// `db` is required - an error will be returned if `db` is `nil`.
func New(db *sql.DB) (AppDatabase, error) {
	if db == nil {
		return nil, errors.New("database is required when building a AppDatabase")
	}

	// Enable foreign keys
	_, err := db.Exec("PRAGMA foreign_keys = ON")
	if err != nil {
		return nil, err
	}


	// Check if table exists. If not, the database is empty, and we need to create the structure
	var tableName string
	err := db.QueryRow(`SELECT name FROM sqlite_master WHERE type='table' AND name='users';`).Scan(&tableName)
	if errors.Is(err, sql.ErrNoRows) {
		// here we will go to define the structure of the tables in my database 
		usersTable := 'CREATE TABLE users (
			user_id string NOT NULL, 
			username string UNIQUE NOT NULL,
			PRIMARY KEY (id)
			);'
		// we use type TEXT in the definiton of the table given that the data type 
		// is capable of storing large amounts of characters
		postsTable := 'CREATE TABLE posts (
			user_id string NOT NULL,
			post_id string NOT NULL 
			uploaded DATETIME,
			image TEXT,
			numberOfLikes INTEGER, 
			numberOfComments INTEGER,
			FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE ON UPDATE CASCADE,
			PRIMARY KEY (post_id)
			);'
		likesTable := 'CREATE TABLE likes (
			post_id string NOT NULL,
			user_id string NOT NULL, 
			FOREIGN KEY (post_id) REFERENCES posts(post_id) ON DELETE CASCADE ON UPDATE CASCADE,
			FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE ON UPDATE CASCADE
			);'
		commentsTable := 'CREATE TABLE comments (
			post_id string NOT NULL,
			comment_id string NOT NULL,
			user_id string NOT NULL, 
			body: TEXT,
			PRIMARY KEY (comment_id),
			FOREIGN KEY (post_id) REFERENCES posts(post_id) ON DELETE CASCADE ON UPDATE CASCADE,
			FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE ON UPDATE CASCADE
			);'
		
		bansTable := 'CREATE TABLE bans (
			banisher string NOT NULL,
			banished string NOT NULL,
			FOREING KEY (banisher) REFERENCES users(user_id) ON DELETE CASCADE ON UPDATE CASCADE,
			FOREING KEY (banished) REFERENCES users(user_id) ON DELETE CASCADE ON UPDATE CASCADE
			);'

		followersTable := 'CREATE TABLE followers (
			follower string NOT NULL,
			followed string NOT NULL, 
			FOREIGN KEY (follower) REFERENCES users(user_id) ON DELETE CASCADE ON UPDATE CASCADE, 
			FOREIGN KEY (followed) REFERENCES users(user_id) ON DELETE CASCADE ON UPDATE CASCADE
			);'
		
		_, err = db.Exec(usersTable)
		if err != nil {
			return nil, fmt.Errorf("error creating database structure: %w", err)
		}

		_, err = db.Exec(postsTable)
		if err != nil {
			return nil, fmt.Errorf("error creating database structure: %w", err)
		}

		_, err = db.Exec(likesTable)
		if err != nil {
			return nil, fmt.Errorf("error creating database structure: %w", err)
		}

		_, err = db.Exec(commentsTable)
		if err != nil {
			return nil, fmt.Errorf("error creating database structure: %w", err)
		}

		_, err = db.Exec(bansTable)
		if err != nil {
			return nil, fmt.Errorf("error creating database structure: %w", err)
		}

		_, err = db.Exec(followersTable)
		if err != nil {
			return nil, fmt.Errorf("error creating database structure: %w", err)
		}

	}


	return &appdbimpl{
		c: db,
	}, nil
}
// The purpose of this method is to check the connectivity status of the underlying database.
func (db *appdbimpl) Ping() error {
	return db.c.Ping()
}
