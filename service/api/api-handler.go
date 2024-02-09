package api

import (
	"net/http"
)

// Handler returns an instance of httprouter.Router that handle APIs registered here
func (rt *_router) Handler() http.Handler {

	// USER OPERATIONS
	// ? CHECKED
	rt.router.POST("/session", rt.wrap(rt.doLogin))
	// ? CHECKED
	rt.router.PUT("/users/:userid", rt.wrap(rt.setMyUserName))

	// LIKES OPERATIONS
	// ? CHECKED
	rt.router.POST("/users/:userid/posts/:postid/likes", rt.wrap(rt.likePhoto))
	// ? CHECKED
	rt.router.GET("/users/:userid/posts/:postid/likes", rt.wrap(rt.getLikers))
	// ? CHECKED
	rt.router.DELETE("/users/:userid/posts/:postid/likes/:likerUsername", rt.wrap(rt.unlikePhoto))

	// COMMENT OPERATIONS
	// ? CHECKED
	rt.router.POST("/users/:userid/posts/:postid/comments", rt.wrap(rt.commentPhoto))
	// ? CHECKED
	rt.router.DELETE("/users/:userid/posts/:postid/comments/:commentid", rt.wrap(rt.uncommentPhoto))
	// ? CHECKED
	rt.router.GET("/users/:userid/posts/:postid/comments", rt.wrap(rt.getComments))

	// make post into user profile
	// ? CHECKED
	rt.router.POST("/users/:userid/posts", rt.wrap(rt.uploadPhoto))
	// ? CHECKED
	rt.router.DELETE("/users/:userid/posts/:postid", rt.wrap(rt.deletePhoto))

	// follow functions
	// ? CHECKED
	rt.router.POST("/users/:userid/following", rt.wrap(rt.followUser))
	// ? CHECKED
	rt.router.DELETE("/users/:userid/following/:followed", rt.wrap(rt.unfollowUser))
	// ? CHECKED
	rt.router.GET("/users/:userid/following", rt.wrap(rt.getFollowing))
	// ? CHECKED
	rt.router.GET("/users/:userid/followers", rt.wrap(rt.getFollowers))

	// Bans funcitons
	// ? CHECKED
	rt.router.POST("/users/:userid/bans", rt.wrap(rt.banUser))
	// ? CHECKED
	rt.router.DELETE("/users/:userid/bans/:banished", rt.wrap(rt.unbanUser))
	// ? CHECKED
	rt.router.GET("/users/:userid/bans", rt.wrap(rt.getUserBans))

	// main functions
	// ? CHECKED
	rt.router.GET("/users/:userid/profile", rt.wrap(rt.getUserProfile))
	// ? CHECKED
	rt.router.GET("/users/:userid/stream", rt.wrap(rt.getMyStream))

	// Special routes
	rt.router.GET("/liveness", rt.liveness)

	return rt.router
}
