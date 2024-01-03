package api

import (
	"net/http"
)




// Handler returns an instance of httprouter.Router that handle APIs registered here
func (rt *_router) Handler() http.Handler {
	// Register routes, this is done to each requet we specify in the API file 
	rt.router.POST("/session", rt.doLogin)
	rt.router.POST("/likes/", rt.likePhoto)
	rt.router.DELETE("/likes/{likeid}", rt.unlikePhoto)
	rt.router.POST("/comments/", rt.commentPhoto)
	rt.router.DELETE("/comments/{commentid}",rt.uncommentPhoto)
	rt.router.PATCH("/users/{userid}", rt.setMyUserName)
	rt.router.POST("/users/{userid}/posts", rt.uploadPhoto)
	rt.router.DELETE("/users/{userid}/posts/{postid}", rt.deletePhoto)
	rt.router.POST("/users/{userid}/following", rt.followUser)
	rt.router.DELETE("/users/{userid}/following/{followedid}", rt.unfollowUser)
	rt.router.POST("/users/{userid}/muted/", rt.bandUser)
	rt.router.DELETE("/users/{userid}/muted/", rt.unbandUser)
	rt.router.GET("/users/{username}/profile", rt.getUserProfile)
	
	
	//rt.router.GET("/", rt.getHelloWorld)
	//rt.router.GET("/context", rt.wrap(rt.getContextReply))

	// Special routes
	//rt.router.GET("/liveness", rt.liveness)

	return rt.router
}
