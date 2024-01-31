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
)

type Like struct {
	Post_ID string `json:"post_id"`
	Like_ID string `json:"like_id"`
	User_ID string `json:"user_id"` // Change from userId to UserID
}

type Comment struct {
	Post_ID    string `json:"post_id"`
	Comment_ID string `json:"comment_id"`
	User_ID    string `json:"user_id"` // Change from userId to UserID
	Body       string `json:"body"`
}

// type Image struct {
// 	Image string `json:"image"`
// }

type Post struct {
	User_ID          string    `json:"user_id"`
	Post_ID          string    `json:"post_id"`
	Uploaded         string    `json:"uploaded"`
	Image            string    `json:"image"`
	Comments         []Comment `json:"comments"`
	NumberOfComments int       `json:"numberOfComments"`
	Likes            []Like    `json:"likes"`
	NumberOfLikes    int       `json:"numberOfLikes"`
}

type User struct {
	User_ID  string `json:"user_id"`
	Username string `json:"username"`
}

type Profile struct {
	User          User   `json:"user"`
	Photos        []Post `json:"photos"`
	NumberOfPosts int    `json:"numberOfPosts"`
	UserFollowers int    `json:"userFollowers"`
	UserFollowing int    `json:"userFollowing"`
}

// AppDatabase is the high-level interface for the DB
type AppDatabase interface {
	Ping() error

	// like methods
	GetLike(string) (Like, error)
	SetLike(Like) error
	RemoveLike(string) error

	// Commemt methods
	GetComment(string) (Comment, error)
	SetComment(Comment) error
	RemoveComment(string) error

	// User methods
	SetUsername(string, string) (User, error)
	GetUserByName(string) (User, error)
	GetName(string) (string, error)
	CreateUser(string, string) error

	// photo methods
	UploadPhoto(Post) error
	DeletePhoto(string) error
	GetUserIDForPost(string) (string, error)
	GetPhotos(string) ([]Post, error)
	GetLastPosts([]string) ([]Post, error)

	// follow methods
	SetFollow(string, string) error
	RemoveFollow(string, string) error
	GetNumberOfFollowers(string) (int, error)
	GetNumberOfFollowing(string) (int, error)
	GetFollowing(string) ([]string, error)

	// band methods
	GetBans(string) ([]string, error)
	BandUser(string, string, string) error
	UnbandUser(string, string) error

	// the name of the parameter is up to the implementation of the interface
	// we only define the type of the parameter the method should have
}

type appdbimpl struct {
	// sql.DB normally indicates a connection to a SQL database
	c *sql.DB
}

// New returns a new instance of AppDatabase based on the SQLite connection `db`.
// `db` is required - an error will be returned if `db` is `nil`.
func New(db *sql.DB) (AppDatabase, error) {
	if db == nil {
		return nil, errors.New("database is required when building an AppDatabase")
	}

	// Enable foreign keys
	_, err := db.Exec("PRAGMA foreign_keys = ON")
	if err != nil {
		return nil, err
	}
	//fmt.Println("fuori il if del database")

	// Check if table exists. If not, the database is empty, and we need to create the structure
	// var tableName string
	// err = db.QueryRow(`SELECT username FROM sqlite_master WHERE type='table' AND name='users';`).Scan(&tableName)
	var create = true
	// if errors.Is(err, sql.ErrNoRows) {
	if create {

		//fmt.Println("enter to create the tables ")

		// here we will go to define the structure of the tables in my database
		usersTable := `CREATE TABLE IF NOT EXISTS users (
			user_id string NOT NULL,
			username string UNIQUE NOT NULL,
			PRIMARY KEY (user_id)
		);`

		// we use type TEXT in the definition of the table given that the data type
		// is capable of storing large amounts of characters
		postsTable := `CREATE TABLE IF NOT EXISTS posts (
			user_id string NOT NULL,
			post_id string NOT NULL,
			uploaded DATETIME,
			image TEXT,
			numberOfComments INTEGER,
			numberOfLikes INTEGER,
			FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE ON UPDATE CASCADE,
			PRIMARY KEY (post_id)
		);`

		likesTable := `CREATE TABLE IF NOT EXISTS likes (
			post_id string NOT NULL,
			like_id string NOT NULL, 
			user_id string NOT NULL,
			FOREIGN KEY (post_id) REFERENCES posts(post_id) ON DELETE CASCADE ON UPDATE CASCADE,
			FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE ON UPDATE CASCADE,
			PRIMARY KEY (like_id)
		);`

		commentsTable := `CREATE TABLE IF NOT EXISTS comments (
			post_id string NOT NULL,
			comment_id string NOT NULL,
			user_id string NOT NULL,
			body TEXT,
			PRIMARY KEY (comment_id),
			FOREIGN KEY (post_id) REFERENCES posts(post_id) ON DELETE CASCADE ON UPDATE CASCADE,
			FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE ON UPDATE CASCADE
		);`

		bansTable := `CREATE TABLE IF NOT EXISTS bans (
			banisher string NOT NULL,
			banished string NOT NULL,
			user_id string NOT NULL,
			FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE ON UPDATE CASCADE
		);`

		followersTable := `CREATE TABLE IF NOT EXISTS followers (
			follower string NOT NULL,
			followed string NOT NULL,
			FOREIGN KEY (follower) REFERENCES users(user_id) ON DELETE CASCADE ON UPDATE CASCADE,
			FOREIGN KEY (followed) REFERENCES users(user_id) ON DELETE CASCADE ON UPDATE CASCADE
		);`

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
