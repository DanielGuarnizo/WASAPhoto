package api

import (
	"net/http"
)

// Handler returns an instance of httprouter.Router that handle APIs registered here
func (rt *_router) Handler() http.Handler {
	// Register routes, this is done to each requet we specify in the API file
	rt.router.POST("/session", rt.wrap(rt.doLogin))
	rt.router.POST("/likes/", rt.wrap(rt.likePhoto))
	rt.router.DELETE("/likes/{likeid}", rt.wrap(rt.unlikePhoto))
	rt.router.POST("/comments/", rt.wrap(rt.commentPhoto))
	rt.router.DELETE("/comments/{commentid}", rt.wrap(rt.uncommentPhoto))
	rt.router.PATCH("/users/{userid}", rt.wrap(rt.setMyUserName))
	rt.router.POST("/users/{userid}/posts", rt.wrap(rt.uploadPhoto))
	rt.router.DELETE("/users/{userid}/posts/{postid}", rt.wrap(rt.deletePhoto))
	rt.router.POST("/users/{userid}/following", rt.wrap(rt.followUser))
	rt.router.DELETE("/users/{userid}/following/{followedid}", rt.wrap(rt.unfollowUser))
	rt.router.POST("/users/{userid}/muted/", rt.wrap(rt.banUser))
	rt.router.DELETE("/users/{userid}/muted/", rt.wrap(rt.unbanUser))
	rt.router.GET("/users/{username}/profile", rt.wrap(rt.getUserProfile))
	rt.router.GET("/users/{userid}/stream", rt.wrap(rt.getMyStream))

	// Special routes
	rt.router.GET("/liveness", rt.liveness)

	return rt.router
}
