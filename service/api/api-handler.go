package api

import (
	"net/http"
)

// Handler returns an instance of httprouter.Router that handle APIs registered here
func (rt *_router) Handler() http.Handler {

	// USER OPERATIONS
	//? CHECKED
	rt.router.POST("/session", rt.wrap(rt.doLogin))
	//? CHECKED
	rt.router.PUT("/users/:userid", rt.wrap(rt.setMyUserName))

	// LIKES AND COMMENT OPERATIONS
	rt.router.POST("/likes/", rt.wrap(rt.likePhoto))
	rt.router.DELETE("/likes/:likeid", rt.wrap(rt.unlikePhoto))
	rt.router.POST("/comments/", rt.wrap(rt.commentPhoto))
	rt.router.DELETE("/comments/:commentid", rt.wrap(rt.uncommentPhoto))

	// make post into user profile
	//? CHECKED
	rt.router.PUT("/users/:userid/posts", rt.wrap(rt.uploadPhoto))
	//? CHECKED
	rt.router.DELETE("/users/:userid/posts/:postid", rt.wrap(rt.deletePhoto))

	// follow functions
	rt.router.POST("/users/:userid/following", rt.wrap(rt.followUser))
	rt.router.DELETE("/users/:userid/following/:followedid", rt.wrap(rt.unfollowUser))

	// Bans funcitons
	//? CHECKED
	rt.router.GET("/users/:userid/bans/:username", rt.wrap(rt.getUserBans))
	//? CHECKED
	rt.router.PUT("/users/:userid/bans/:username", rt.wrap(rt.banUser))
	//? CHECKED
	rt.router.DELETE("/users/:userid/bans/:username", rt.wrap(rt.unbanUser))

	// main functions
	//? CHECKED
	rt.router.GET("/users/:userid/profile", rt.wrap(rt.getUserProfile))

	rt.router.GET("/users/:userid/stream", rt.wrap(rt.getMyStream))

	// Special routes
	rt.router.GET("/liveness", rt.liveness)

	// rt.router.POST("/session", rt.wrap(rt.doLogin))
	// rt.router.POST("/likes/", rt.wrap(rt.likePhoto))
	// rt.router.DELETE("/likes/:likeid", rt.wrap(rt.unlikePhoto))
	// rt.router.POST("/comments/", rt.wrap(rt.commentPhoto))
	// rt.router.DELETE("/comments/:commentid", rt.wrap(rt.uncommentPhoto))
	// rt.router.PATCH("/users/:id", rt.wrap(rt.setMyUserName))

	// rt.router.PUT("/users/:id/posts/", rt.wrap(rt.uploadPhoto))
	// rt.baseLogger.Warning("Registering route: /users/:userid/posts")

	// rt.router.DELETE("/users/:id/posts/:postid", rt.wrap(rt.deletePhoto))
	// rt.router.POST("/users/:id/following", rt.wrap(rt.followUser))
	// rt.router.DELETE("/users/:id/following/:followedid", rt.wrap(rt.unfollowUser))

	// // user operations
	// rt.router.GET("/users/:id/bans/:username", rt.wrap(rt.getUserBans))
	// rt.router.POST("/users/:id/muted/", rt.wrap(rt.banUser))
	// rt.router.DELETE("/users/:id/muted/", rt.wrap(rt.unbanUser))

	// rt.router.GET("/users/:id/profile", rt.wrap(rt.getUserProfile))

	// rt.router.GET("/users/:id/stream", rt.wrap(rt.getMyStream))

	// // Special routes
	// rt.router.GET("/liveness", rt.liveness)

	return rt.router
}
