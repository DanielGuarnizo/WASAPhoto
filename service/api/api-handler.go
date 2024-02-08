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
	rt.router.PUT("/users/:userid/posts/:postid/likes", rt.wrap(rt.likePhoto))
	// ? CHECKED
	rt.router.GET("/users/:userid/posts/:postid/likes", rt.wrap(rt.getLikers))
	// ? CHECKED
	rt.router.DELETE("/users/:userid/posts/:postid/likes", rt.wrap(rt.unlikePhoto))

	// COMMENT OPERATIONS
	// ? CHECKED
	rt.router.PUT("/comments", rt.wrap(rt.commentPhoto))
	// ? CHECKED
	rt.router.DELETE("/users/:userid/posts/:postid/comments/:commentid", rt.wrap(rt.uncommentPhoto))
	// ? CHECKED
	rt.router.GET("/users/:userid/posts/:postid/comments", rt.wrap(rt.getComments))

	// make post into user profile
	// ? CHECKED
	rt.router.PUT("/users/:userid/posts", rt.wrap(rt.uploadPhoto))
	// ? CHECKED
	rt.router.DELETE("/users/:userid/posts/:postid", rt.wrap(rt.deletePhoto))

	// follow functions
	// ? CHECKED
	rt.router.PUT("/users/:userid/follows/:username", rt.wrap(rt.followUser))
	// ? CHECKED
	rt.router.DELETE("/users/:userid/follows/:username", rt.wrap(rt.unfollowUser))
	// ? CHECKED
	rt.router.GET("/users/:userid/following", rt.wrap(rt.getFollowing))
	// ? CHECKED
	rt.router.GET("/users/:userid/followers", rt.wrap(rt.getFollowers))

	// Bans funcitons
	// ? CHECKED
	rt.router.GET("/users/:userid/bans/:username", rt.wrap(rt.getUserBans))
	// ? CHECKED
	rt.router.PUT("/users/:userid/bans/:username", rt.wrap(rt.banUser))
	// ? CHECKED
	rt.router.DELETE("/users/:userid/bans/:username", rt.wrap(rt.unbanUser))

	// main functions
	// ? CHECKED
	rt.router.GET("/users/:userid/profile", rt.wrap(rt.getUserProfile))

	rt.router.GET("/users/:userid/stream", rt.wrap(rt.getMyStream))

	// Special routes
	rt.router.GET("/liveness", rt.liveness)

	return rt.router
}
